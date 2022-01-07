package nurse

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func CollectGargabe(dockerPath string, usageLimit int) error {
	keepDays := 28

	inUse, err := diskUsage(dockerPath)
	if err != nil {
		return err
	}

	log.Printf("usage of %s is %d%%", dockerPath, inUse)

	if inUse >= usageLimit {
		err := dockerSystemVolumes()
		if err != nil {
			return err
		}
	}

	for {
		inUse, err := diskUsage(dockerPath)
		if err != nil {
			return err
		}

		if inUse < usageLimit {
			return nil
		}

		if keepDays < 0 {
			return fmt.Errorf("failed to reduce disk usage")
		}

		err = dockerSystemPrune(keepDays)
		if err != nil {
			return err
		}

		keepDays--
	}
}

// diskUsage returns the diskUsage in percent
func diskUsage(path string) (int, error) {
	cmd := exec.Command(
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

func dockerSystemPrune(keepDays int) error {
	log.Printf("docker system prune -a -f --filter until=%dh", keepDays*24)

	cmd := exec.Command(
		"sh",
		"-c",
		fmt.Sprintf("docker system prune -a -f --filter until=%dh", keepDays*24),
	)
	cmd.Stderr = os.Stderr
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, "DOCKER_HOST=tcp://127.0.0.1:12375")

	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("docker system prune: %v: %s", err, string(output))
	}

	return nil
}

func dockerSystemVolumes() error {
	log.Printf("docker system prune -f --volumes")

	cmd := exec.Command(
		"sh",
		"-c",
		"docker system prune -f --volumes",
	)
	cmd.Stderr = os.Stderr
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, "DOCKER_HOST=tcp://127.0.0.1:12375")

	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("prune docker volumes: %v: %s", err, string(output))
	}

	return nil
}
