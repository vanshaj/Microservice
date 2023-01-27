package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	directory := flag.String("dir", "", "directory to run pipeline on")
	flag.Parse()
	err := run(*directory, os.Stdout)
	if err != nil {
		log.Fatal(err)
	}
}

func run(dir string, out io.Writer) error {
	sig := make(chan os.Signal, 1)
	errCh := make(chan error)
	done := make(chan struct{})

	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	pipelines := make([]executor, 4)
	pipelines[0] = newStep(
		"Build",
		"go",
		"Go Build: Success",
		dir,
		[]string{"build", "."})
	pipelines[1] = newStep(
		"Test",
		"go",
		"Go test: Success",
		dir,
		[]string{"test", "-v", "./..."})
	pipelines[2] = newExceptionStep(
		"Format",
		"gofmt",
		"Go format: Success",
		dir,
		[]string{"-l", "."})
	pipelines[3] = newTimeoutStep(
		"git push",
		"git",
		"Git Push: Success",
		dir,
		[]string{"push", "origin", "master"},
		10*time.Second)
	go func() {
		for _, pipeline := range pipelines {
			msg, err := pipeline.execute()
			if err != nil {
				errCh <- err
				return
			}
			fmt.Fprintln(out, msg)
		}
		close(done)
	}()
	for {
		select {
		case rec := <-sig:
			signal.Stop(sig)
			return fmt.Errorf("%s: Existing %w", rec, ErrSignal)
		case err := <-errCh:
			return err
		case <-done:
			return nil
		}
	}
}
