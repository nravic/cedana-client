package cmd

// This file contains all the restore-related commands when starting `cedana restore ...`

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"time"

	task "buf.build/gen/go/cedana/task/protocolbuffers/go"
	"github.com/cedana/cedana/pkg/api"
	"github.com/cedana/cedana/pkg/api/services"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"google.golang.org/grpc/status"

	"github.com/cedana/cedana/pkg/utils"
	"github.com/mdlayher/vsock"
)

var restoreCmd = &cobra.Command{
	Use:   "restore",
	Short: "Manually restore a process or container from a checkpoint located at input path: [process, runc (container), containerd (container)]",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		port, _ := cmd.Flags().GetUint32(portFlag)
		cts, err := services.NewClient(port)
		if err != nil {
			return fmt.Errorf("Error creating client: %v", err)
		}
		ctx := context.WithValue(cmd.Context(), utils.CtsKey, cts)
		cmd.SetContext(ctx)
		return nil
	},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		cts := cmd.Context().Value(utils.CtsKey).(*services.ServiceClient)
		cts.Close()
	},
}

var restoreProcessCmd = &cobra.Command{
	Use:   "process",
	Short: "Manually restore a process from a checkpoint located at input path",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		cts := cmd.Context().Value(utils.CtsKey).(*services.ServiceClient)

		path := args[0]
		tcpEstablished, _ := cmd.Flags().GetBool(tcpEstablishedFlag)
		tcpClose, _ := cmd.Flags().GetBool(tcpCloseFlag)
		stream, _ := cmd.Flags().GetInt32(streamFlag)
		bucket, _ := cmd.Flags().GetString(bucketFlag)
		if stream > 0 {
			if _, err := exec.LookPath("cedana-image-streamer"); err != nil {
				log.Error().Msgf("Cannot find cedana-image-streamer in PATH")
				return err
			}
      var err error
			if bucket != "" {
				ctx, err = awsSetup(bucket, ctx, false)
				if err != nil {
					log.Error().Msgf("Error setting up AWS bucket for direct remoting")
					return err
				}
			}
		} else if bucket != "" {
			return fmt.Errorf("Dump to AWS S3 bucket only possible with --stream")
		}
		restoreArgs := task.RestoreArgs{
			CheckpointID:   "Not implemented",
			CheckpointPath: path,
			Stream:         stream,
			Bucket:         bucket,
			CriuOpts: &task.CriuOpts{
				TcpEstablished: tcpEstablished,
				TcpClose:       tcpClose,
			},
		}

		resp, err := cts.Restore(ctx, &restoreArgs)
		if err != nil {
			st, ok := status.FromError(err)
			if ok {
				log.Error().Str("message", st.Message()).Str("code", st.Code().String()).Msgf("Failed")
			} else {
				log.Error().Err(err).Msgf("Failed")
			}
			return err
		}
		log.Info().Str("message", resp.Message).Int32("PID", resp.State.PID).Interface("stats", resp.RestoreStats).Msgf("Success")

		return nil
	},
}

var restoreJobCmd = &cobra.Command{
	Use:   "job",
	Short: "Manually restore a previously dumped process or container from an input id",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		cts := cmd.Context().Value(utils.CtsKey).(*services.ServiceClient)

		jid := args[0]
		tcpEstablished, _ := cmd.Flags().GetBool(tcpEstablishedFlag)
		tcpCloseFlag, _ := cmd.Flags().GetBool(tcpCloseFlag)
		root, _ := cmd.Flags().GetString(rootFlag)
		stream, _ := cmd.Flags().GetInt32(streamFlag)
		bundle, err := cmd.Flags().GetString(bundleFlag)
		consoleSocket, err := cmd.Flags().GetString(consoleSocketFlag)
		detach, err := cmd.Flags().GetBool(detachFlag)
		img, err := cmd.Flags().GetString(imgFlag)
		bucket, err := cmd.Flags().GetString(bucketFlag)
		if stream > 0 {
			if _, err := exec.LookPath("cedana-image-streamer"); err != nil {
				log.Error().Msgf("Cannot find cedana-image-streamer in PATH")
				return err
			}
      var err error
			if bucket != "" {
				ctx, err = awsSetup(bucket, ctx, false)
				if err != nil {
					log.Error().Msgf("Error setting up AWS bucket for direct remoting")
					return err
				}
			}
		} else if bucket != "" {
			return fmt.Errorf("Dump to AWS S3 bucket only possible with --stream")
		}
		restoreArgs := &task.JobRestoreArgs{
			JID:            jid,
			Stream:         stream,
			Bucket:         bucket,
			CheckpointPath: img,
			CriuOpts: &task.CriuOpts{
				TcpEstablished: tcpEstablished,
				TcpClose:       tcpCloseFlag,
			},
			RuncOpts: &task.RuncOpts{
				Root:          getRuncRootPath(root),
				Bundle:        bundle,
				ConsoleSocket: consoleSocket,
				Detach:        detach,
			},
		}

		attach, _ := cmd.Flags().GetBool(attachFlag)
		if attach {
			stream, err := cts.JobRestoreAttach(ctx, &task.AttachArgs{Args: &task.AttachArgs_JobRestoreArgs{JobRestoreArgs: restoreArgs}})
			if err != nil {
				st, ok := status.FromError(err)
				if ok {
					log.Error().Err(st.Err()).Msg("restore failed")
				} else {
					log.Error().Err(err).Msg("restore failed")
				}
				return err
			}

			// Handler stdout, stderr
			exitCode := make(chan int)
			go func() {
				for {
					resp, err := stream.Recv()
					if err != nil {
						log.Error().Err(err).Msg("stream ended")
						exitCode <- 1
						return
					}
					if resp.Stdout != "" {
						fmt.Print(resp.Stdout)
					} else if resp.Stderr != "" {
						fmt.Fprint(os.Stderr, resp.Stderr)
					} else {
						exitCode <- int(resp.GetExitCode())
						return
					}
				}
			}()

			// Handle stdin
			go func() {
				scanner := bufio.NewScanner(os.Stdin)
				for scanner.Scan() {
					if err := stream.Send(&task.AttachArgs{Stdin: scanner.Text() + "\n"}); err != nil {
						log.Error().Err(err).Msg("error sending stdin")
						return
					}
				}
			}()

			os.Exit(<-exitCode)

			// TODO: Add signal handling properly
		}

		resp, err := cts.JobRestore(ctx, restoreArgs)
		if err != nil {
			st, ok := status.FromError(err)
			if ok {
				log.Error().Str("message", st.Message()).Str("code", st.Code().String()).Msgf("Failed")
			} else {
				log.Error().Err(err).Msgf("Failed")
			}
			return err
		}
		log.Info().Str("message", resp.Message).Int32("PID", resp.State.PID).Interface("stats", resp.RestoreStats).Msgf("Success")

		return nil
	},
}

var restoreKataCmd = &cobra.Command{
	Use:   "kata",
	Short: "Manually restore a workload in the kata-vm [vm-name] from a directory [-d]",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		vm := args[0]

		port, _ := cmd.Flags().GetUint32(portFlag)
		cts, err := services.NewVSockClient(vm, port)
		if err != nil {
			log.Error().Msgf("Error creating client: %v", err)
			return err
		}
		defer cts.Close()

		path, _ := cmd.Flags().GetString(dirFlag)
		tcpEstablished, _ := cmd.Flags().GetBool(tcpEstablishedFlag)
		restoreArgs := task.RestoreArgs{
			CheckpointID:   vm,
			CheckpointPath: "/tmp/dmp.tar",
			CriuOpts: &task.CriuOpts{
				TcpEstablished: tcpEstablished,
			},
		}

		go func() {
			time.Sleep(1 * time.Second)

			// extract cid from the process tree on host
			cid, err := utils.ExtractCID(vm)
			if err != nil {
				return
			}

			conn, err := vsock.Dial(cid, api.KATA_TAR_FILE_RECEIVER_PORT, nil)
			if err != nil {
				return
			}
			defer conn.Close()

			// Open the file
			file, err := os.Open(path)
			if err != nil {
				return
			}
			defer file.Close()

			buffer := make([]byte, 1024)

			// Read from file and send over VSOCK connection
			for {
				bytesRead, err := file.Read(buffer)
				if err != nil {
					if err == io.EOF {
						break
					}
					return
				}

				_, err = conn.Write(buffer[:bytesRead])
				if err != nil {
					return
				}
			}
		}()

		resp, err := cts.KataRestore(ctx, &restoreArgs)
		if err != nil {
			st, ok := status.FromError(err)
			if ok {
				log.Error().Msgf("Restore task failed: %v, %v: %v", st.Code(), st.Message(), st.Details())
			} else {
				log.Error().Msgf("Restore task failed: %v", err)
			}
			return err
		}
		log.Info().Msgf("Response: %v", resp.Message)

		return nil
	},
}

var containerdRestoreCmd = &cobra.Command{
	Use:   "containerd",
	Short: "Manually restore a running container to a directory",
	Args:  cobra.ArbitraryArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		cts := cmd.Context().Value(utils.CtsKey).(*services.ServiceClient)

		ref, _ := cmd.Flags().GetString(imgFlag)
		id, _ := cmd.Flags().GetString(idFlag)
		restoreArgs := &task.ContainerdRestoreArgs{
			ImgPath:     ref,
			ContainerID: id,
		}

		resp, err := cts.ContainerdRestore(ctx, restoreArgs)
		if err != nil {
			st, ok := status.FromError(err)
			if ok {
				log.Error().Msgf("Restore task failed: %v, %v", st.Message(), st.Code())
			} else {
				log.Error().Msgf("Restore task failed: %v", err)
			}
			return err
		}
		log.Info().Msgf("Response: %v", resp.Message)

		return nil
	},
}

var runcRestoreCmd = &cobra.Command{
	Use:   "runc",
	Short: "Manually restore a running runc container to a directory",
	Args:  cobra.ArbitraryArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		cts := cmd.Context().Value(utils.CtsKey).(*services.ServiceClient)

		root, err := cmd.Flags().GetString(rootFlag)
		bundle, err := cmd.Flags().GetString(bundleFlag)
		consoleSocket, err := cmd.Flags().GetString(consoleSocketFlag)
		detach, err := cmd.Flags().GetBool(detachFlag)
		netPid, err := cmd.Flags().GetInt32(netPidFlag)
		tcpEstablished, _ := cmd.Flags().GetBool(tcpEstablishedFlag)
		tcpClose, _ := cmd.Flags().GetBool(tcpCloseFlag)
		img, _ := cmd.Flags().GetString(imgFlag)
		id, _ := cmd.Flags().GetString(idFlag)
		fileLocks, _ := cmd.Flags().GetBool(fileLocksFlag)

		opts := &task.RuncOpts{
			Root:          getRuncRootPath(root),
			Bundle:        bundle,
			ConsoleSocket: consoleSocket,
			Detach:        detach,
			NetPid:        netPid,
			ContainerID:   id,
		}

		restoreArgs := &task.RuncRestoreArgs{
			ImagePath: img,
			Opts:      opts,
			CriuOpts: &task.CriuOpts{
				FileLocks:      fileLocks,
				TcpEstablished: tcpEstablished,
				TcpClose:       tcpClose,
			},
		}

		resp, err := cts.RuncRestore(ctx, restoreArgs)
		if err != nil {
			st, ok := status.FromError(err)
			if ok {
				log.Error().Str("message", st.Message()).Str("code", st.Code().String()).Msgf("Failed")
			} else {
				log.Error().Err(err).Msgf("Failed")
			}
			return err
		}
		log.Info().Str("message", resp.Message).Interface("stats", resp.RestoreStats).Msgf("Success")

		return nil
	},
}

func init() {
	// Process
	restoreCmd.AddCommand(restoreProcessCmd)
	restoreProcessCmd.Flags().BoolP(tcpEstablishedFlag, "t", false, "restore with TCP connections established")
	restoreProcessCmd.Flags().BoolP(tcpCloseFlag, "", false, "restore with TCP connections closed")
	restoreProcessCmd.Flags().Int32P(streamFlag, "s", 0, "restore images using cedana-image-streamer")
	restoreProcessCmd.Flags().StringP(bucketFlag, "", "", "AWS S3 bucket to stream from")

	// Job
	restoreCmd.AddCommand(restoreJobCmd)
	restoreJobCmd.Flags().BoolP(tcpEstablishedFlag, "t", false, "restore with TCP connections established")
	restoreJobCmd.Flags().BoolP(tcpCloseFlag, "", false, "restore with TCP connections closed")
	restoreJobCmd.Flags().Int32P(streamFlag, "s", 0, "restore images using cedana-image-streamer")
	restoreJobCmd.Flags().StringP(bucketFlag, "", "", "AWS S3 bucket to stream from")
	restoreJobCmd.Flags().BoolP(attachFlag, "a", false, "attach stdin/stdout/stderr")
	restoreJobCmd.Flags().StringP(bundleFlag, "b", "", "(runc) bundle path")
	restoreJobCmd.Flags().StringP(consoleSocketFlag, "c", "", "(runc) console socket path")
	restoreJobCmd.Flags().BoolP(detachFlag, "e", false, "(runc) restore detached")
	restoreJobCmd.Flags().StringP(rootFlag, "r", "default", "(runc) root")
	restoreJobCmd.Flags().StringP(imgFlag, "i", "", "checkpoint image")

	// Kata
	restoreCmd.AddCommand(restoreKataCmd)
	restoreKataCmd.Flags().StringP(dirFlag, "d", "", "path of tar file (inside VM) to restore from")
	restoreKataCmd.MarkFlagRequired(dirFlag)

	// Containerd
	restoreCmd.AddCommand(containerdRestoreCmd)
	containerdRestoreCmd.Flags().String(imgFlag, "", "checkpoint image")
	containerdRestoreCmd.MarkFlagRequired(imgFlag)
	containerdRestoreCmd.Flags().StringP(idFlag, "i", "", "container id")
	containerdRestoreCmd.MarkFlagRequired(idFlag)

	// Runc
	restoreCmd.AddCommand(runcRestoreCmd)
	runcRestoreCmd.Flags().StringP(imgFlag, "", "", "checkpoint image to restore from")
	runcRestoreCmd.MarkFlagRequired(imgFlag)
	runcRestoreCmd.Flags().StringP(idFlag, "i", "", "container id")
	runcRestoreCmd.MarkFlagRequired(idFlag)
	runcRestoreCmd.Flags().StringP(bundleFlag, "b", "", "bundle path")
	runcRestoreCmd.Flags().StringP(consoleSocketFlag, "c", "", "console socket path")
	runcRestoreCmd.Flags().StringP(rootFlag, "r", "default", "runc root directory")
	runcRestoreCmd.Flags().BoolP(detachFlag, "e", false, "run runc container in detached mode")
	runcRestoreCmd.Flags().Int32P(netPidFlag, "n", 0, "provide the network pid to restore to in k3s")
	runcRestoreCmd.Flags().Bool(fileLocksFlag, false, "restore file locks")
	runcRestoreCmd.Flags().BoolP(tcpEstablishedFlag, "t", false, "tcp established")
	runcRestoreCmd.Flags().BoolP(tcpCloseFlag, "", false, "tcp close")

	rootCmd.AddCommand(restoreCmd)
}
