package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
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
	pipelines := make([]*step, 2)
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
	for _, pipeline := range pipelines {
		msg, err := pipeline.execute()
		if err != nil {
			return err
		}
		fmt.Fprintln(out, msg)
	}
	return nil
}
