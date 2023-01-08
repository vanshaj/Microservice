package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"sync"
)

func main() {
	op := flag.String("op", "sum", "operation to be performed")
	col := flag.Int("col", 1, "column to be performed on")
	dir := flag.Bool("dir", false, "is directory passed")
	flag.Parse()
	err := run(flag.Args(), *op, *col, *dir, os.Stdout)
	if err != nil {
		log.Fatal(err)
	}
}

func run(filenames []string, operation string, column int, dir bool, w io.Writer) error {
	var opFunc stat
	if len(filenames) == 0 {
		return ErrNoFiles
	}
	if column < 1 {
		return ErrInvalidColumn
	}
	switch operation {
	case "sum":
		opFunc = statFunc(sum)
	case "avg":
		opFunc = statFunc(avg)
	default:
		return ErrInvalidOperation
	}
	consolidate := make([]float64, 0)
	if dir {
		return nil
	}
	resCh := make(chan []float64)
	errCh := make(chan error)
	doneCh := make(chan struct{})
	filesCh := make(chan string)

	wg := sync.WaitGroup{}

	go func() {
		defer close(filesCh)
		for _, filename := range filenames {
			filesCh <- filename
		}
	}()

	for i := 0; i < runtime.NumCPU(); i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for filename := range filesCh {
				//fileInfo, err := os.Stat(filename)
				//if err != nil {
				//	errCh <- err
				//}
				//if fileInfo.IsDir() {
				//	errCh <- ErrNoFiles
				//}
				f, err := os.Open(filename)
				defer f.Close()
				if err != nil {
					errCh <- err
				}
				data, err := csvToFloat(f, column)
				if err != nil {
					errCh <- err
				}
				resCh <- data
			}
		}()
	}

	go func() {
		wg.Wait()
		close(doneCh)
	}()
	for {
		select {
		case err := <-errCh:
			return err
		case data := <-resCh:
			consolidate = append(consolidate, data...)
		case <-doneCh:
			_, err := fmt.Fprintln(w, opFunc.operation(consolidate))
			return err
		}
	}
}
