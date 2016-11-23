package main

import (
	"fmt"
	"path"
	"os"
	"strings"
	"os/exec"
)

const mountCmd = "mount"

func initDriver([]string) (FlexOutput, error) {
	return FlexOutput{FlexStatusSuccess, "no initialisation needed", ""}, nil
}

func attach(args []string) (FlexOutput, error) {
	if len(args) < 3 {
		return FlexOutput{FlexStatusFailure, "invalid arguments specified. expected json options.", ""}, fmt.Errorf("invalid arguments specified, expected 'attach [json options]'")
	}

	opts := args[2]
	options, err := parseOptions(opts)

	if err != nil {
		return FlexOutput{FlexStatusFailure, fmt.Sprintf("failed to parse options: %s", err.Error()), ""}, err
	}

	if len(options.VolumeID) == 0 {
		return FlexOutput{FlexStatusFailure, fmt.Sprintf("volumeID must be specified"), ""}, fmt.Errorf("a volume ID must be provided to guarantee persistence")
	}

	// TODO: should we have a default here?
	if len(options.Path) == 0 {
		return FlexOutput{FlexStatusFailure, fmt.Sprintf("a base path must be specified"), ""}, fmt.Errorf("a base storage directory path must be specified")
	}

	if options.Permissions == 0 {
		// default permissions to 0755, as hostPath does
		options.Permissions = 0755
	}

	dir := path.Join(options.Path, options.VolumeID)

	err = os.MkdirAll(dir, options.Permissions)

	if err != nil {
		return FlexOutput{FlexStatusFailure, fmt.Sprintf("error provisioning directory: %s"), ""}, err
	}

	return FlexOutput{FlexStatusSuccess, "", dir}, nil
}

func detach(args []string) (FlexOutput, error) {
	if len(args) < 3 {
		return FlexOutput{FlexStatusFailure, "no device to detach specified", ""}, fmt.Errorf("no device to detach specified")
	}

	target := args[2]

	err := os.RemoveAll(target)

	if err != nil {
		return FlexOutput{FlexStatusFailure, fmt.Sprintf("failed to cleanup volume: %s", err.Error()), ""}, nil
	}

	return FlexOutput{FlexStatusSuccess, "", ""}, nil
}

func mount(args []string) (FlexOutput, error) {
	if len(args) < 5 {
		return FlexOutput{FlexStatusFailure, "invalid arguments specified, expected 'mount [target dir] [mount device] [json options]'", ""},
			fmt.Errorf("invalid arguments specified, expected 'mount [target dir] [mount device] [json options]'")
	}

	target := args[2]
	source := args[3]
	opts := args[4]

	options, err := parseOptions(opts)
	if err != nil {
		return FlexOutput{FlexStatusFailure, fmt.Sprintf("failed to parse options: %s", err.Error()), ""}, err
	}

	if options.ReadWrite == "" {
		options.ReadWrite = "rw"
	}

	mountArgs := makeMountArgs(source, target, "", []string{"bind", options.ReadWrite})

	command := exec.Command(mountCmd, mountArgs...)
	output, err := command.CombinedOutput()

	if err != nil {
		return FlexOutput{FlexStatusFailure, fmt.Sprintf("mount failed: %v\nMounting command: %s\nMounting arguments: %s %s %v\nOutput: %s\n",
			err, mountCmd, source, target, options, string(output)), ""}, err
	}

	return FlexOutput{FlexStatusSuccess, "", ""},  nil
}

func unmount(args []string) (FlexOutput, error) {
	if len(args) < 3 {
		return FlexOutput{FlexStatusFailure, "no mount path specified for unmount", ""}, fmt.Errorf("no mount path specified for unmount")
	}

	target := args[2]

	command := exec.Command("umount", target)

	output, err := command.CombinedOutput()
	if err != nil {
		return FlexOutput{FlexStatusFailure, fmt.Sprintf("Unmount failed: %v\nUnmounting arguments: %s\nOutput: %s\n", err, target, string(output)), ""}, err
	}

	return FlexOutput{FlexStatusSuccess, "", ""}, nil
}

// makeMountArgs makes the arguments to the mount(8) command.
func makeMountArgs(source, target, fstype string, options []string) []string {
	// Build mount command as follows:
	//   mount [-t $fstype] [-o $options] [$source] $target
	mountArgs := []string{}
	if len(fstype) > 0 {
		mountArgs = append(mountArgs, "-t", fstype)
	}
	if len(options) > 0 {
		mountArgs = append(mountArgs, "-o", strings.Join(options, ","))
	}
	if len(source) > 0 {
		mountArgs = append(mountArgs, source)
	}
	mountArgs = append(mountArgs, target)

	return mountArgs
}