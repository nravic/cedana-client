package gpu

// Runc GPU interception utilities

import (
	"fmt"
	"strings"

	"github.com/opencontainers/runtime-spec/specs-go"
)

func AddGPUInterceptionToSpec(spec *specs.Spec, libraryPath string, jid string) error {
	// HACK: Remove nvidia prestart hook as we don't support working around it, yet
	if spec.Hooks != nil {
		for i, hook := range spec.Hooks.Prestart {
			if strings.Contains(hook.Path, "nvidia") {
				spec.Hooks.Prestart = append(spec.Hooks.Prestart[:i], spec.Hooks.Prestart[i+1:]...)
				break
			}
		}
	}

	// Remove existing /dev/shm mount if it exists
	for i, m := range spec.Mounts {
		if m.Destination == "/dev/shm" {
			spec.Mounts = append(spec.Mounts[:i], spec.Mounts[i+1:]...)
		}
	}

	spec.Mounts = append(spec.Mounts, specs.Mount{
		Destination: "/dev/shm",
		Source:      "/dev/shm",
		Type:        "bind",
		Options:     []string{"rbind", "nosuid", "nodev", "rw"},
	})

	// Mount the GPU plugin library
	spec.Mounts = append(spec.Mounts, specs.Mount{
		Destination: libraryPath,
		Source:      libraryPath,
		Type:        "bind",
		Options:     []string{"rbind", "nosuid", "nodev", "rw"},
	})

	// Set the env vars
	if spec.Process == nil {
		return fmt.Errorf("spec does not have a process")
	}
	spec.Process.Env = append(spec.Process.Env, "LD_PRELOAD="+libraryPath)
	spec.Process.Env = append(spec.Process.Env, "CEDANA_JID="+jid)

	return nil
}