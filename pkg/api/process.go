package api

// Implements the task service functions for processes

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"

	task "buf.build/gen/go/cedana/task/protocolbuffers/go"
	"github.com/cedana/cedana/pkg/utils"
	"github.com/rs/xid"
	"github.com/shirou/gopsutil/v3/process"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/rs/zerolog/log"
)

const (
	OUTPUT_FILE_PATH  string      = "/var/log/cedana-output.log"
	OUTPUT_FILE_PERMS os.FileMode = 0o777
	OUTPUT_FILE_FLAGS int         = os.O_WRONLY | os.O_CREATE | os.O_TRUNC
)

func (s *service) Start(ctx context.Context, args *task.StartArgs) (*task.StartResp, error) {
	return s.startHelper(ctx, args, nil)
}

func (s *service) StartAttach(stream grpc.BidiStreamingServer[task.AttachArgs, task.AttachResp]) error {
	in, err := stream.Recv()
	if err != nil {
		return err
	}
	args := in.GetStartArgs()

	_, err = s.startHelper(stream.Context(), args, stream)
	return err
}

func (s *service) Manage(ctx context.Context, args *task.ManageArgs) (*task.ManageResp, error) {
	state := &task.ProcessState{}
	var err error
	if args.JID == "" {
		args.JID = xid.New().String()
	} else {
		// Check if job ID is already in use
		state, err = s.getState(ctx, args.JID)
		if state != nil {
			err = status.Error(codes.AlreadyExists, "job ID already exists")
			return nil, err
		}
	}

	// Check if process PID already running as a managed job
	queryResp, err := s.JobQuery(ctx, &task.JobQueryArgs{PIDs: []int32{args.PID}})
	if queryResp != nil && len(queryResp.Processes) > 0 {
		if utils.PidExists(uint32(args.PID)) {
			err = status.Error(codes.AlreadyExists, "PID already running as a managed job")
			return nil, err
		}
	}

	exitCode := utils.WaitForPid(args.PID)

	state, err = s.generateState(ctx, args.PID)
	if err != nil {
		return nil, status.Error(codes.NotFound, "is the process running?")
	}
	state.JID = args.JID
	state.JobState = task.JobState_JOB_RUNNING
	state.GPU = args.GPU

	if args.GPU {
		log.Info().Msg("GPU support requested, assuming process was already started with LD_PRELOAD")
		if args.GPU {
			err = s.StartGPUController(ctx, args.UID, args.GID, args.Groups, args.JID)
			if err != nil {
				return nil, fmt.Errorf("failed to start GPU controller: %v", err)
			}
		}
	}

	err = s.updateState(ctx, state.JID, state)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to update state after manage")
	}

	// Wait for server shutdown to gracefully exit, since job is now managed
	// Also wait for process exit, to update state
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		select {
		case <-s.serverCtx.Done():
			log.Info().Str("JID", state.JID).Int32("PID", state.PID).Msg("server shutting down, killing process")
			err := syscall.Kill(int(state.PID), syscall.SIGKILL)
			if err != nil {
				log.Error().Err(err).Str("JID", state.JID).Int32("PID", state.PID).Msg("failed to kill process. already dead?")
			}
		case <-exitCode:
			log.Info().Str("JID", state.JID).Int32("PID", state.PID).Msg("process exited")
		}
		state, err = s.getState(context.WithoutCancel(ctx), state.JID)
		if err != nil {
			log.Error().Err(err).Msg("failed to get latest state, DB might be inconsistent")
		}
		state.JobState = task.JobState_JOB_DONE
		err := s.updateState(context.WithoutCancel(ctx), state.JID, state)
		if err != nil {
			log.Error().Err(err).Msg("failed to update state after done")
		}
		if s.GetGPUController(state.JID) != nil {
			err = s.StopGPUController(state.JID)
			if err != nil {
				log.Error().Err(err).Msg("failed to stop GPU controller")
			}
		}
	}()

	// Clean up GPU controller and also handle premature exit
	if s.GetGPUController(state.JID) != nil {
		s.wg.Add(1)
		go func() {
			defer s.wg.Done()
			s.WaitGPUController(state.JID)

			// Should kill process if still running since GPU controller might have exited prematurely
			syscall.Kill(int(state.PID), syscall.SIGKILL)
		}()
	}

	return &task.ManageResp{Message: "success", State: state}, nil
}

func (s *service) Dump(ctx context.Context, args *task.DumpArgs) (*task.DumpResp, error) {
	dumpStats := task.DumpStats{
		DumpType: task.DumpType_PROCESS,
	}
	ctx = context.WithValue(ctx, utils.DumpStatsKey, &dumpStats)

	if args.Dir == "" {
		args.Dir = viper.GetString("shared_storage.dump_storage_dir")
		if args.Dir == "" {
			return nil, status.Error(codes.InvalidArgument, "dump storage dir not provided/found in config")
		}
	}

	state := &task.ProcessState{}
	pid := args.PID

	var err error
	if args.JID != "" { // if managed job
		state, err = s.getState(ctx, args.JID)
		if err != nil {
			return nil, status.Error(codes.NotFound, "job ID not found: "+err.Error())
		}
		if state.GPU && s.gpuEnabled == false {
			return nil, status.Error(codes.FailedPrecondition, "GPU support is not enabled in daemon")
		}
		pid = state.PID
	}

	state, err = s.generateState(ctx, pid)
	if err != nil {
		err = status.Error(codes.Internal, err.Error())
		return nil, err
	}

	err = s.dump(ctx, state, args)
	if err != nil {
		st := status.New(codes.Internal, err.Error())
		return nil, st.Err()
	}

	resp := task.DumpResp{
		Message:      fmt.Sprintf("Dumped process %d to %s", pid, args.Dir),
		CheckpointID: state.CheckpointPath, // XXX: Just return path for ID for now
	}

	// Only update state if it was a managed job
	if args.JID != "" {
		err = s.updateState(ctx, state.JID, state)
		if err != nil {
			return nil, status.Error(codes.Internal, fmt.Sprintf("failed to update state with error: %s", err.Error()))
		}
	}

	resp.State = state
	resp.DumpStats = &dumpStats

	return &resp, err
}

func (s *service) Restore(ctx context.Context, args *task.RestoreArgs) (*task.RestoreResp, error) {
	return s.restoreHelper(ctx, args, nil)
}

func (s *service) RestoreAttach(stream grpc.BidiStreamingServer[task.AttachArgs, task.AttachResp]) error {
	in, err := stream.Recv()
	if err != nil {
		return err
	}
	args := in.GetRestoreArgs()

	_, err = s.restoreHelper(stream.Context(), args, stream)
	return err
}

//////////////////////////
///// Process Utils //////
//////////////////////////

func (s *service) startHelper(ctx context.Context, args *task.StartArgs, stream grpc.BidiStreamingServer[task.AttachArgs, task.AttachResp]) (*task.StartResp, error) {
	if args.Task == "" {
		args.Task = viper.GetString("client.task")
	}

	state := &task.ProcessState{}

	if args.GPU {
		if s.gpuEnabled == false {
			return nil, status.Error(codes.FailedPrecondition, "GPU support is not enabled in daemon")
		}
		state.GPU = true
	}

	if args.JID == "" {
		state.JID = xid.New().String()
	} else {
		existingState, _ := s.getState(ctx, args.JID)
		if existingState != nil {
			return nil, status.Error(codes.AlreadyExists, "job ID already exists")
		}
		state.JID = args.JID
	}
	args.JID = state.JID

	pid, exitCode, err := s.run(ctx, args, stream)
	if err != nil {
		log.Error().Err(err).Msg("failed to run task")
		return nil, status.Error(codes.Internal, "failed to run task")
	}
	log.Info().Int32("PID", pid).Str("JID", state.JID).Msgf("managing process")
	state.PID = pid
	state.JobState = task.JobState_JOB_RUNNING
	err = s.updateState(ctx, state.JID, state)
	if err != nil {
		log.Error().Err(err).Msg("failed to update state after run")
		syscall.Kill(int(pid), syscall.SIGKILL) // kill cuz inconsistent state
		return nil, status.Error(codes.Internal, "failed to update state after run")
	}

	if stream != nil && exitCode != nil {
		code := <-exitCode // if streaming, wait for process to finish
		if stream != nil {
			stream.Send(&task.AttachResp{
				ExitCode: int32(code),
			})
		}
		state, err = s.getState(context.WithoutCancel(ctx), state.JID)
		if err != nil {
			log.Warn().Err(err).Msg("failed to get latest state, DB might be inconsistent")
		}
		state.JobState = task.JobState_JOB_DONE
		err = s.updateState(context.WithoutCancel(ctx), state.JID, state)
		if err != nil {
			log.Error().Err(err).Msg("failed to update state after done")
			return nil, status.Error(codes.Internal, "failed to update state after done")
		}
	} else {
		s.wg.Add(1)
		go func() {
			defer s.wg.Done()
			<-exitCode
			state, err = s.getState(context.WithoutCancel(ctx), state.JID)
			if err != nil {
				log.Warn().Err(err).Msg("failed to get latest state, DB might be inconsistent")
			}
			state.JobState = task.JobState_JOB_DONE
			err = s.updateState(context.WithoutCancel(ctx), state.JID, state)
			if err != nil {
				log.Error().Err(err).Msg("failed to update state after done")
				return
			}
		}()
	}

	return &task.StartResp{
		Message: fmt.Sprint("Job started successfully"),
		PID:     pid,
		JID:     state.JID,
	}, err
}

func (s *service) restoreHelper(ctx context.Context, args *task.RestoreArgs, stream grpc.BidiStreamingServer[task.AttachArgs, task.AttachResp]) (*task.RestoreResp, error) {
	var resp task.RestoreResp
	var pid int32
	var err error

	restoreStats := task.RestoreStats{
		DumpType: task.DumpType_PROCESS,
	}
	ctx = context.WithValue(ctx, utils.RestoreStatsKey, &restoreStats)

	if args.JID != "" {
		state, err := s.getState(ctx, args.JID)
		if err != nil {
			err = status.Error(codes.NotFound, err.Error())
			return nil, err
		}
		if state.GPU && s.gpuEnabled == false {
			return nil, status.Error(codes.FailedPrecondition, "Dump has GPU state and GPU support is not enabled in daemon")
		}
		if args.CheckpointPath == "" {
			args.CheckpointPath = state.CheckpointPath
		}
	}

	if args.CheckpointPath == "" {
		return nil, status.Error(codes.InvalidArgument, "checkpoint path cannot be empty")
	}
	stat, err := os.Stat(args.CheckpointPath)
	if os.IsNotExist(err) {
		return nil, status.Error(codes.InvalidArgument, "invalid checkpoint path: does not exist")
	}
	if (args.Stream <= 0) && (stat.IsDir() || !strings.HasSuffix(args.CheckpointPath, ".tar")) {
		return nil, status.Error(codes.InvalidArgument, "invalid checkpoint path: must be tar file")
	}
	if (args.Stream > 0) && !stat.IsDir() {
		return nil, status.Error(codes.InvalidArgument, "invalid checkpoint path: must be dir (--stream enabled)")
	}

	pid, exitCode, err := s.restore(ctx, args, stream)
	if err != nil {
		err := status.Error(codes.Internal, fmt.Sprintf("failed to restore process: %v", err))
		return nil, err
	}

	// Only update state if it was a managed job
	state := &task.ProcessState{}
	if args.JID != "" {
		state, err = s.getState(ctx, args.JID)
		if err != nil {
			log.Warn().Err(err).Msg("failed to get latest state, DB might be inconsistent")
		}
		log.Info().Int32("PID", pid).Str("JID", state.JID).Msgf("managing restored process")
		state.PID = pid
		state.JobState = task.JobState_JOB_RUNNING
		err = s.updateState(ctx, state.JID, state)
		if err != nil {
			log.Error().Err(err).Msg("failed to update state after restore")
			syscall.Kill(int(pid), syscall.SIGKILL) // kill cuz inconsistent state
			return nil, status.Error(codes.Internal, "failed to update state after restore")
		}

		if stream != nil && exitCode != nil {
			code := <-exitCode // if streaming, wait for process to finish
			if stream != nil {
				stream.Send(&task.AttachResp{
					ExitCode: int32(code),
				})
			}
			state, err = s.getState(context.WithoutCancel(ctx), state.JID)
			if err != nil {
				log.Warn().Err(err).Msg("failed to get latest state, DB might be inconsistent")
			}
			state.JobState = task.JobState_JOB_DONE
			err = s.updateState(context.WithoutCancel(ctx), state.JID, state)
			if err != nil {
				log.Error().Err(err).Msg("failed to update state after done")
				return nil, status.Error(codes.Internal, "failed to update state after done")
			}
		} else {
			s.wg.Add(1)
			go func() {
				defer s.wg.Done()
				<-exitCode
				state, err = s.getState(context.WithoutCancel(ctx), state.JID)
				if err != nil {
					log.Warn().Err(err).Msg("failed to get latest state, DB might be inconsistent")
				}
				state.JobState = task.JobState_JOB_DONE
				err = s.updateState(context.WithoutCancel(ctx), state.JID, state)
				if err != nil {
					log.Error().Err(err).Msg("failed to update state after done")
					return
				}
			}()
		}
	}

	resp = task.RestoreResp{
		Message:      fmt.Sprintf("successfully restored process: %v", pid),
		State:        state,
		RestoreStats: &restoreStats,
	}

	return &resp, nil
}

func (s *service) run(ctx context.Context, args *task.StartArgs, stream grpc.BidiStreamingServer[task.AttachArgs, task.AttachResp]) (int32, chan int, error) {
	var pid int32
	if args.Task == "" {
		return 0, nil, fmt.Errorf("could not find task")
	}

	var err error
	if args.GPU {
		err = s.StartGPUController(ctx, args.UID, args.GID, args.Groups, args.JID)
		if err != nil {
			return 0, nil, fmt.Errorf("failed to start GPU controller: %v", err)
		}

		sharedLibPath := viper.GetString("gpu_shared_lib_path")
		if sharedLibPath == "" {
			sharedLibPath = utils.GpuSharedLibPath
		}
		if _, err := os.Stat(sharedLibPath); os.IsNotExist(err) {
			return 0, nil, fmt.Errorf("no gpu shared lib at %s", sharedLibPath)
		}
		args.Task = fmt.Sprintf("CEDANA_JID=%s LD_PRELOAD=%s %s", args.JID, sharedLibPath, args.Task)
	}

	groupsUint32 := make([]uint32, len(args.Groups))
	for i, v := range args.Groups {
		groupsUint32[i] = uint32(v)
	}
	var cmdCtx context.Context
	if stream != nil {
		cmdCtx = utils.CombineContexts(s.serverCtx, stream.Context()) // either should terminate the process
	} else {
		cmdCtx = s.serverCtx
	}
	cmd := exec.CommandContext(cmdCtx, "bash", "-c", args.Task)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setsid: true,
		Credential: &syscall.Credential{
			Uid:    uint32(args.UID),
			Gid:    uint32(args.GID),
			Groups: groupsUint32,
		},
	}
	cmd.Env = args.Env

	// working dir needs to be consistent on the checkpoint and restore side
	if args.WorkingDir != "" {
		cmd.Dir = args.WorkingDir
	}

	if stream == nil {
		if args.LogOutputFile == "" {
			args.LogOutputFile = OUTPUT_FILE_PATH
		}
		outFile, err := os.OpenFile(args.LogOutputFile, OUTPUT_FILE_FLAGS, OUTPUT_FILE_PERMS)
		defer outFile.Close()
		os.Chmod(args.LogOutputFile, OUTPUT_FILE_PERMS)
		if err != nil {
			return 0, nil, err
		}
		cmd.Stdin = nil // equivalent to /dev/null
		cmd.Stdout = outFile
		cmd.Stderr = outFile
	} else {
		stdinPipe, err := cmd.StdinPipe()
		if err != nil {
			return 0, nil, err
		}
		// Receive stdin from stream
		s.wg.Add(1)
		go func() {
			defer s.wg.Done()
			defer stdinPipe.Close()
			for {
				in, err := stream.Recv()
				if err != nil {
					log.Debug().Err(err).Msg("finished reading stdin")
					return
				}
				_, err = stdinPipe.Write([]byte(in.Stdin))
				if err != nil {
					log.Error().Err(err).Msg("failed to write to stdin")
					return
				}
			}
		}()
		// Scan stdout
		stdoutPipe, err := cmd.StdoutPipe()
		if err != nil {
			return 0, nil, err
		}
		stdoutScanner := bufio.NewScanner(stdoutPipe)
		s.wg.Add(1)
		go func() {
			defer s.wg.Done()
			defer stdoutPipe.Close()
			for stdoutScanner.Scan() {
				if err := stream.Send(&task.AttachResp{Stdout: stdoutScanner.Text() + "\n"}); err != nil {
					log.Error().Err(err).Msg("failed to send stdout")
					return
				}
			}
			if err := stdoutScanner.Err(); err != nil {
				log.Debug().Err(err).Msgf("finished reading stdout")
			}
		}()

		// Scan stdout
		stderrPipe, err := cmd.StderrPipe()
		if err != nil {
			return 0, nil, err
		}
		stderrScanner := bufio.NewScanner(stderrPipe)
		s.wg.Add(1)
		go func() {
			defer s.wg.Done()
			defer stderrPipe.Close()
			for stderrScanner.Scan() {
				if err := stream.Send(&task.AttachResp{Stderr: stderrScanner.Text() + "\n"}); err != nil {
					log.Error().Err(err).Msg("failed to send stderr")
					return
				}
			}
			if err := stderrScanner.Err(); err != nil {
				log.Debug().Err(err).Msgf("finished reading stderr")
			}
		}()
	}

	err = cmd.Start()
	if err != nil {
		return 0, nil, err
	}

	ppid := int32(os.Getpid())
	pid = int32(cmd.Process.Pid)
	closeCommonFds(ppid, pid)

	s.wg.Add(1)
	exitCode := make(chan int)
	go func() {
		defer s.wg.Done()
		err := cmd.Wait()
		if err != nil {
			log.Debug().Err(err).Msg("process Wait()")
		}
		if s.GetGPUController(args.JID) != nil {
			s.StopGPUController(args.JID)
		}
		log.Info().Int("status", cmd.ProcessState.ExitCode()).Int32("PID", pid).Msg("process exited")
		code := cmd.ProcessState.ExitCode()
		exitCode <- code
	}()

	// Clean up GPU controller and also handle premature exit
	if s.GetGPUController(args.JID) != nil {
		s.wg.Add(1)
		go func() {
			defer s.wg.Done()
			s.WaitGPUController(args.JID)

			// Should kill process if still running since GPU controller might have exited prematurely
			cmd.Process.Kill()
		}()
	}

	return pid, exitCode, err
}

func closeCommonFds(parentPID, childPID int32) error {
	parent, err := process.NewProcess(parentPID)
	if err != nil {
		return err
	}

	child, err := process.NewProcess(childPID)
	if err != nil {
		return err
	}

	parentFds, err := parent.OpenFiles()
	if err != nil {
		return err
	}

	childFds, err := child.OpenFiles()
	if err != nil {
		return err
	}

	for _, pfd := range parentFds {
		for _, cfd := range childFds {
			if pfd.Path == cfd.Path && strings.Contains(pfd.Path, ".pid") {
				// we have a match, close the FD
				err := syscall.Close(int(cfd.Fd))
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}
