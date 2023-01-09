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
	fmt.Println("Inside main number of goroutines are: ", runtime.NumGoroutine())
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
	quit := make(chan int)
	filesCh := make(chan string)
	jobsCh := make(chan string)

	wg := sync.WaitGroup{}

	go func() {
		defer close(filesCh)
		for _, filename := range filenames {
			filesCh <- filename
		}
	}()
	go func() {
		for {
			select {
			case fn := <-filesCh:
				if fn == "" {
					log.Println("closing job channel positive")
					close(jobsCh)
					return
				}
				jobsCh <- fn
			case <-quit:
				log.Println("closing job channel negative")
				close(jobsCh)
				return
			}
		}
	}()
	for i := 0; i < runtime.NumCPU(); i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case filename := <-jobsCh:
					if filename == "" {
						log.Println("Job channel closed returning goroutine")
						return
					}
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
			}
		}()
	}
	for i := 0; i < len(filenames); i++ {
		select {
		case err := <-errCh:
			close(quit)
			go func() {
				for {
					select {
					case <-resCh:
						continue
					case <-errCh:
						continue
					}
				}
			}()
			wg.Wait()
			close(resCh)
			close(errCh)
			return err
		case data := <-resCh:
			consolidate = append(consolidate, data...)
		}
	}
	wg.Wait()
	_, err := fmt.Fprintln(w, opFunc.operation(consolidate))
	return err
}
