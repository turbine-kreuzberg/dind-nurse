package nurse

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"
)

func Setup(ctx context.Context) error {
	err := setupBuildxBuilder(ctx)
	if err != nil {
		return err
	}

	err = awaitDocker(ctx)
	if err != nil {
		return err
	}

	testBuild(ctx)

	testBuildx(ctx)

	return nil
}

func setupBuildxBuilder(ctx context.Context) error {
	log.Println("set up buildx builder")

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

func awaitDocker(ctx context.Context) error {
	log.Println("wait for docker startup")

	ctx, cancel := context.WithTimeout(ctx, time.Minute)
	defer cancel()

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
		return fmt.Errorf("await docker: %v: %s", err, string(output))
	}

	return nil
}

func testBuild(ctx context.Context) {
	log.Println("test docker build")

	cmd := exec.CommandContext(
		ctx,
		"sh",
		"-c",
		`mkdir -p /test-build
echo "FROM hello-world" > /test-build/Dockerfile
docker build /test-build`,
	)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, "DOCKER_HOST=tcp://127.0.0.1:12375")

	log.Println(cmd.String())

	_ = cmd.Run() // ignore failures to allow startup of nurse to keep docker healyhy
}

func testBuildx(ctx context.Context) {
	log.Println("test buildx build")

	cmd := exec.CommandContext(
		ctx,
		"sh",
		"-c",
		`mkdir -p /test-build
echo "FROM hello-world" > /test-build/Dockerfile
docker buildx build /test-build`,
	)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, "DOCKER_HOST=tcp://127.0.0.1:12375")

	log.Println(cmd.String())

	_ = cmd.Run() // ignore failures to allow startup of nurse to keep docker healyhy
}
