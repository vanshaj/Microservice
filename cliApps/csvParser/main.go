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
	// Use of intermediate channel jobsCh
	// The closing of this intermediate channel will cause the worker function to retrieve empty filenames in case
	go func() {
		for {
			select {
			case fn := <-filesCh:
				// When filesCh is closed then case will read the default string value from the channel
				// This code block will get execute when all the filenames have been processed, so we now close the jobsCh
				if fn == "" {
					log.Println("closing job channel positive")
					close(jobsCh)
					return
				}
				// send filename to jobsCh , positive scenario
				jobsCh <- fn
			case <-quit:
				// This quit channed will get called when there is an error in the processing of file names
				log.Println("closing job channel negative")
				close(jobsCh)
				return
			}
		}
	}()
	// Worker function
	for i := 0; i < runtime.NumCPU(); i++ {
		wg.Add(1)
		go func() {
			// Once return is called from for block we will call wg.Done()
			defer wg.Done()
			for {
				select {
				case filename := <-jobsCh:
					// When job channel is closed filename will be empty string thus we will return all the goroutines
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
	// Keep on listening, [till len of filenames]
	for i := 0; i < len(filenames); i++ {
		select {
		case err := <-errCh:
			// if we have an error we don't want to perform anythin so we close the quit channel which will triger the case of line 74 case <- quit and this close the jobCh
			close(quit)
			go func() {
				// We will return all goroutines because jobCh is closed and filename in line 92 will be empty but what about some goroutines who are  in progress
				// For those we will just listen to the noise and continue in backgound
				for {
					select {
					case <-resCh:
						continue
					case <-errCh:
						continue
					}
				}
			}()
			// Wait for all worker goroutines to return
			wg.Wait()
			close(resCh)
			close(errCh)
			return err
		case data := <-resCh:
			consolidate = append(consolidate, data...)
		}
	}
	// Wait for all worker goroutines to complete the task and return
	wg.Wait()
	_, err := fmt.Fprintln(w, opFunc.operation(consolidate))
	return err
}
