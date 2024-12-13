package gpu

// Internal definitions for the GPU controller

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"sync"
	"syscall"

	"buf.build/gen/go/cedana/cedana-gpu/grpc/go/gpu/gpugrpc"
	"buf.build/gen/go/cedana/cedana-gpu/protocolbuffers/go/gpu"
	"github.com/cedana/cedana/pkg/utils"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

const (
	CONTROLLER_HOST               = "localhost"
	CONTROLLER_LOG_PATH_FORMATTER = "/tmp/cedana-gpu-controller-%s.log"
	CONTROLLER_LOG_FLAGS          = os.O_CREATE | os.O_WRONLY | os.O_TRUNC
	CONTROLLER_LOG_PERMS          = 0644

	// Signal sent to job when GPU controller exits prematurely. The intercepted job
	// is guaranteed to exit upon receiving this signal, and prints to stderr
	// about the GPU controller's failure.
	CONTROLLER_PREMATURE_EXIT_SIGNAL = syscall.SIGUSR1
)

type Controller struct {
	ErrBuf *bytes.Buffer

	*exec.Cmd
	gpugrpc.ControllerClient
	*grpc.ClientConn
}

type Controllers struct {
	sync.Map
}

func (m *Controllers) Get(jid string) *Controller {
	c, ok := m.Load(jid)
	if !ok {
		return nil
	}
	return c.(*Controller)
}

func spawnController(ctx context.Context, lifetime context.Context, wg *sync.WaitGroup, binary string, jid string) (*Controller, chan int, error) {
	port, err := utils.GetFreePort()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get free port: %w", err)
	}

	controller := &Controller{
		ErrBuf: &bytes.Buffer{},
		Cmd:    exec.CommandContext(lifetime, binary, jid, "--port", strconv.Itoa(port)),
	}

	controller.Stderr = controller.ErrBuf
	controller.Stdout = nil // TODO: capture controller logs
	controller.SysProcAttr = &syscall.SysProcAttr{
		Pdeathsig: syscall.SIGTERM,
	}
	controller.Cancel = func() error { return controller.Cmd.Process.Signal(syscall.SIGTERM) } // NO SIGKILL!!!

	err = controller.Start()
	if err != nil {
		return nil, nil, fmt.Errorf(
			"failed to start GPU controller: %w",
			utils.GRPCErrorShort(err, controller.ErrBuf.String()),
		)
	}

	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	conn, err := grpc.NewClient(fmt.Sprintf("%s:%d", CONTROLLER_HOST, port), opts...)
	if err != nil {
		controller.Process.Signal(syscall.SIGTERM)
		return nil, nil, fmt.Errorf(
			"failed to create GPU controller client: %w",
			utils.GRPCErrorShort(err, controller.ErrBuf.String()),
		)
	}
	controller.ClientConn = conn
	controller.ControllerClient = gpugrpc.NewControllerClient(conn)

	// Cleanup controller on exit, and signal job of its exit

	exited := make(chan int, 1)
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer close(exited)
		defer conn.Close()

		err := controller.Wait()
		if err != nil {
			log.Trace().Err(err).Msg("GPU controller Wait()")
		}
		log.Debug().Int("code", controller.ProcessState.ExitCode()).Msg("GPU controller exited")
	}()

	log.Debug().Int("port", port).Msg("waiting for GPU controller...")

	err = controller.WaitForHealthCheck(ctx, wg)
	if err != nil {
		controller.Process.Signal(syscall.SIGTERM)
		conn.Close()
		return nil, nil, err
	}

	return controller, exited, nil
}

// Health checks the GPU controller, blocking on connection until ready.
// This can be used as a proxy to wait for the controller to be ready.
func (controller *Controller) WaitForHealthCheck(ctx context.Context, wg *sync.WaitGroup) error {
	waitCtx, cancel := context.WithTimeout(ctx, HEALTH_TIMEOUT)
	defer cancel()

	// Wait for early controller exit, and cancel the blocking health check
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-utils.WaitForPidCtx(waitCtx, uint32(controller.Process.Pid))
		cancel()
	}()

	resp, err := controller.HealthCheck(waitCtx, &gpu.HealthCheckRequest{}, grpc.WaitForReady(true))
	if resp != nil {
		log.Debug().
			Int32("devices", resp.DeviceCount).
			Str("version", resp.Version).
			Int32("driver", resp.GetAvailableAPIs().GetDriverVersion()).
			Msg("GPU health check")
	}
	if err != nil || !resp.Success {
		controller.Process.Signal(syscall.SIGTERM)
		controller.Close()
		if err == nil {
			err = status.Errorf(codes.FailedPrecondition, "GPU health check failed")
			controller.ErrBuf.WriteString("GPU health check failed")
		}
		return utils.GRPCErrorShort(err, controller.ErrBuf.String())
	}
	return nil
}