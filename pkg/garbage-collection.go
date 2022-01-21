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

func CollectGargabe(ctx context.Context, dockerPath string, upperDiskUsageLimit, lowerDiskUsageLimit int) error {
	keepDays := 28

	inUse, err := diskUsage(ctx, dockerPath)
	if err != nil {
		return err
	}

	log.Printf("usage of %s is %d%%", dockerPath, inUse)

	if inUse < upperDiskUsageLimit {
		return nil
	}

	err = dockerSystemVolumes(ctx)
	if err != nil {
		return err
	}

	for {
		if keepDays < 0 {
			return fmt.Errorf("failed to reduce disk usage")
		}

		err = dockerSystemPrune(ctx, keepDays)
		if err != nil {
			return err
		}

		err = dockerBuildxPrune(ctx, keepDays)
		if err != nil {
			return err
		}

		keepDays--

		inUse, err := diskUsage(ctx, dockerPath)
		if err != nil {
			return err
		}

		if inUse < lowerDiskUsageLimit {
			return nil
		}
	}
}

// diskUsage returns the diskUsage in percent
func diskUsage(ctx context.Context, path string) (int, error) {
	cmd := exec.CommandContext(
		ctx,
		"sh",
		"-c",
		fmt.Sprintf("df %s | tail -n-1 | awk '{print $5+0}'", path),
	)

	cmd.Stderr = os.Stderr

	output, err := cmd.Output()
	if err != nil {
		return 0, fmt.Errorf("execute disk free: %v: %s", err, string(output))
	}

	percent, err := strconv.Atoi(strings.TrimSpace(string(output)))
	if err != nil {
		return 0, fmt.Errorf("parse free disk: %v", err)
	}

	return percent, nil
}

func dockerSystemVolumes(ctx context.Context) error {
	cmd := exec.CommandContext(
		ctx,
		"sh",
		"-c",
		"docker system prune -f --volumes",
	)
	cmd.Stderr = os.Stderr
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, "DOCKER_HOST=tcp://127.0.0.1:12375")

	log.Println(cmd.String())

	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("prune docker volumes: %v: %s", err, string(output))
	}

	return nil
}

func dockerSystemPrune(ctx context.Context, keepDays int) error {
	cmd := exec.CommandContext(
		ctx,
		"sh",
		"-c",
		fmt.Sprintf("docker system prune --all --force --filter until=%dh", keepDays*24),
	)
	cmd.Stderr = os.Stderr
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, "DOCKER_HOST=tcp://127.0.0.1:12375")

	log.Println(cmd.String())

	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("docker system prune: %v: %s", err, string(output))
	}

	return nil
}

func dockerBuildxPrune(ctx context.Context, keepDays int) error {
	cmd := exec.CommandContext(
		ctx,
		"sh",
		"-c",
		fmt.Sprintf("docker buildx prune --all --force --filter until=%dh", keepDays*24),
	)
	cmd.Stderr = os.Stderr
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, "DOCKER_HOST=tcp://127.0.0.1:12375")

	log.Println(cmd.String())

	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("docker system prune: %v: %s", err, string(output))
	}

	return nil
}
