package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
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
	for _, pipeline := range pipelines {
		msg, err := pipeline.execute()
		if err != nil {
			return err
		}
		fmt.Fprintln(out, msg)
	}
	return nil
}
