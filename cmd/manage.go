package cmd

import (
	"context"
	"fmt"

	"github.com/cedana/cedana/pkg/api/services"
	"github.com/cedana/cedana/pkg/api/services/task"
	"github.com/cedana/cedana/pkg/utils"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"google.golang.org/grpc/status"
)

var manageCmd = &cobra.Command{
	Use:   "manage",
	Short: "Start managing a process or a container",
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

var manageProcessCmd = &cobra.Command{
	Use:   "process",
	Short: "Manage a process",
	Args:  cobra.ArbitraryArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		cts := cmd.Context().Value(utils.CtsKey).(*services.ServiceClient)

		pid, _ := cmd.Flags().GetInt(pidFlag)
		jid, _ := cmd.Flags().GetString(idFlag)
		gpuEnabled, _ := cmd.Flags().GetBool(gpuEnabledFlag)

		manageArgs := &task.ManageArgs{
			JID: jid,
			PID: int32(pid),
			GPU: gpuEnabled,
		}

		resp, err := cts.Manage(ctx, manageArgs)
		if err != nil {
			st, ok := status.FromError(err)
			if ok {
				log.Error().Err(st.Err()).Msg("manage process failed")
			} else {
				log.Error().Err(err).Msg("manage process failed")
			}
			return err
		}
		log.Info().Msgf("Managing process: %v", resp)

		return nil
	},
}

var manageRuncCmd = &cobra.Command{
	Use:   "runc",
	Short: "Manage a runc container",
	Args:  cobra.ArbitraryArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		cts := cmd.Context().Value(utils.CtsKey).(*services.ServiceClient)

		id, _ := cmd.Flags().GetString(idFlag)
		root, _ := cmd.Flags().GetString(rootFlag)
		gpuEnabled, _ := cmd.Flags().GetBool(gpuEnabledFlag)

		manageArgs := &task.RuncManageArgs{
			ContainerID: id,
			Root:        root,
			GPU:         gpuEnabled,
		}

		resp, err := cts.RuncManage(ctx, manageArgs)
		if err != nil {
			st, ok := status.FromError(err)
			if ok {
				log.Error().Err(st.Err()).Msg("manage runc container failed")
			} else {
				log.Error().Err(err).Msg("manage runc container failed")
			}
			return err
		}
		log.Info().Msgf("Managing runc container: %v", resp)

		return nil
	},
}

func init() {
	// process
	manageProcessCmd.Flags().StringP(idFlag, "i", "", "job id to use")
	manageProcessCmd.Flags().Int(pidFlag, 0, "pid")
	manageProcessCmd.MarkFlagRequired(pidFlag)
	manageProcessCmd.Flags().BoolP(gpuEnabledFlag, "g", false, "runc root")
	manageCmd.AddCommand(manageProcessCmd)

	// runc
	manageRuncCmd.Flags().StringP(idFlag, "i", "", "container id")
	manageRuncCmd.MarkFlagRequired(idFlag)
	manageRuncCmd.Flags().StringP(rootFlag, "r", "default", "runc root")
	manageRuncCmd.Flags().BoolP(gpuEnabledFlag, "g", false, "runc root")
	manageCmd.AddCommand(manageRuncCmd)

	rootCmd.AddCommand(manageCmd)
}
