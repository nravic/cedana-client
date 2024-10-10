package api

import (
	"context"
	"encoding/json"

	"github.com/cedana/cedana/pkg/api/services/task"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *service) JobDump(ctx context.Context, args *task.JobDumpArgs) (*task.JobDumpResp, error) {
	res := &task.JobDumpResp{}

	state, err := s.getState(ctx, args.JID)
	if err != nil {
		err = status.Error(codes.NotFound, err.Error())
		return nil, err
	}

	// Check if normal process or container
	if state.ContainerID == "" {
		dumpResp, err := s.Dump(ctx, &task.DumpArgs{
			JID:      args.JID,
			Type:     args.Type,
			Stream:   args.Stream,
			Dir:      args.Dir,
			CriuOpts: args.CriuOpts,
		})
		if err != nil {
			return nil, err
		}
		res.State = dumpResp.State
		res.DumpStats = dumpResp.DumpStats
		res.CheckpointID = dumpResp.CheckpointID
		res.UploadID = dumpResp.UploadID
		res.Message = dumpResp.Message
	} else {
		// Runc
	}

	return res, nil
}

func (s *service) JobQuery(ctx context.Context, args *task.JobQueryArgs) (*task.JobQueryResp, error) {
	res := &task.JobQueryResp{}

	if len(args.JIDs) > 0 {
		for _, jid := range args.JIDs {
			state, err := s.getState(ctx, jid)
			if err != nil {
				return nil, status.Error(codes.NotFound, "job not found")
			}
			if state != nil {
				res.Processes = append(res.Processes, state)
			}
		}
	} else {
		pidSet := make(map[int32]bool)
		for _, pid := range args.PIDs {
			pidSet[pid] = true
		}

		list, err := s.db.ListJobs(ctx)
		if err != nil {
			return nil, status.Error(codes.Internal, "failed to retrieve jobs from database")
		}
		for _, job := range list {
			state := task.ProcessState{}
			err = json.Unmarshal(job.State, &state)
			if err != nil {
				return nil, status.Error(codes.Internal, "failed to unmarshal state")
			}
			if len(pidSet) > 0 && !pidSet[state.PID] {
				continue
			}
			res.Processes = append(res.Processes, &state)
		}
	}

	return res, nil
}