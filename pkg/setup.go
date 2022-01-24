package nurse

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
)

func Setup(ctx context.Context) error {
	err := setupBuildxBuilder(ctx)
	if err != nil {
		return err
	}

	return nil
}

func setupBuildxBuilder(ctx context.Context) error {
	cmd := exec.CommandContext(
		ctx,
		"sh",
		"-c",
		"docker buildx create --use --bootstrap --name ci-builder",
	)
	cmd.Stderr = os.Stderr
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, "DOCKER_HOST=tcp://127.0.0.1:12375")

	log.Println(cmd.String())

	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("create buildx builder: %v: %s", err, string(output))
	}

	return nil
}
