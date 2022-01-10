package nurse

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func AvoidOOM(ctx context.Context, memoryLimit int) error {
	pid, err := PIDofDinD(ctx)
	if err != nil {
		return err
	}

	memory, err := MemoryOfPID(ctx, pid)
	if err != nil {
		return err
	}

	log.Printf("dockerd memory usage %d Byte", memory)

	if memory > memoryLimit {
		err = RestartDockerd(ctx)
		if err != nil {
			return err
		}

		err = AwaitDocker(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

func PIDofDinD(ctx context.Context) (int, error) {
	cmd := exec.CommandContext(
		ctx,
		"sh",
		"-c",
		"pgrep dockerd",
	)
	cmd.Stderr = os.Stderr

	output, err := cmd.Output()
	if err != nil {
		return 0, fmt.Errorf("find pid: %v", err)
	}

	pid, err := strconv.Atoi(strings.TrimSpace(string(output)))
	if err != nil {
		return 0, fmt.Errorf("parse pid: %v", err)
	}

	return pid, nil
}

func MemoryOfPID(ctx context.Context, pid int) (int, error) {
	cmd := exec.CommandContext(
		ctx,
		"sh",
		"-c",
		fmt.Sprintf("grep ^VmRSS /proc/%d/status | awk '{print $2}'", pid),
	)
	cmd.Stderr = os.Stderr

	output, err := cmd.Output()
	if err != nil {
		return 0, fmt.Errorf("read memory limit: %v", err)
	}

	memory, err := strconv.Atoi(strings.TrimSpace(string(output)))
	if err != nil {
		return 0, fmt.Errorf("parse memory limit: %v", err)
	}

	memory = memory * 1024

	return memory, nil
}

func RestartDockerd(ctx context.Context) error {
	log.Printf("running killall dockerd")

	cmd := exec.CommandContext(
		ctx,
		"killall",
		"dockerd",
	)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("killall dockerd: %v", err)
	}

	return nil
}

func AwaitDocker(ctx context.Context) error {
	cmd := exec.CommandContext(
		ctx,
		"sh",
		"-c",
		"until docker version; do sleep 1; done",
	)
	cmd.Stderr = os.Stderr
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, "DOCKER_HOST=tcp://127.0.0.1:12375")

	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("await docker startup: %v: %s", err, string(output))
	}

	return nil
}
