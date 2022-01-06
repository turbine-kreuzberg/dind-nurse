package main

import (
	"context"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"syscall"

	nurse "github.com/turbine-kreuzberg/dind-nurse/pkg"
	"github.com/urfave/cli/v2"
)

func main() {
	err := run()
	if err != nil {
		log.Fatalf("nurse failed: %v", err)
	}
}

func run() error {
	app := &cli.App{
		Name:  "Nurse",
		Usage: "Keeps docker in docker healthy.",
		Commands: []*cli.Command{
			{
				Name:  "server",
				Usage: "Start the server.",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "addr", Value: ":2375", Usage: "Address to run service on."},
					&cli.StringFlag{Name: "target", Value: "http://127.0.0.1:12375", Usage: "Docker daemon to forward requests to."},
					&cli.IntFlag{Name: "dind-memory-limit", Value: 200 * 1024 * 1024, Usage: "Restart memory watermark for Docker daemon."},
					&cli.IntFlag{Name: "parallel-request-limit", Value: 4, Usage: "Maximum of request to process in parallel."},
				},
				Action: server,
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		return err
	}

	return nil
}

func server(c *cli.Context) error {
	log.Println("set up service")

	targetURL, err := url.Parse(c.String("target"))
	if err != nil {
		return err
	}

	svc := nurse.NewService(targetURL, c.Int("dind-memory-limit"), c.Int("parallel-request-limit"))

	server := httpServer(svc, c.String("addr"))

	log.Println("starting server")

	go func() {
		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	// awaitShutdown
	log.Println("running")

	stop := make(chan os.Signal, 2)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	// shutdown
	log.Println("shutting down")

	ctx, cancel := context.WithTimeout(context.Background(), nurse.MaxExecutionTime)
	defer cancel()

	err = server.Shutdown(ctx)
	if err != nil {
		return err
	}

	log.Println("shutdown complete")

	return nil
}

func httpServer(h http.Handler, addr string) *http.Server {
	httpServer := &http.Server{
		ReadTimeout:  nurse.MaxExecutionTime,
		WriteTimeout: nurse.MaxExecutionTime,
	}
	httpServer.Addr = addr
	httpServer.Handler = h

	return httpServer
}
