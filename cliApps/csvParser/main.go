package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
)

func main() {
	op := flag.String("op", "sum", "operation to be performed")
	col := flag.Int("col", 1, "column to be performed on")
	flag.Parse()
	err := run(flag.Args(), *op, *col, os.Stdout)
	if err != nil {
		log.Fatal(err)
	}
}

func run(filenames []string, operation string, column int, w io.Writer) error {
	var opFunc stat
	switch operation {
	case "sum":
		opFunc = statFunc(sum)
	case "avg":
		opFunc = statFunc(avg)
	default:
		return ErrInvalidOperation
	}
	consolidate := make([]float64, 0)
	for _, filename := range filenames {
		fileInfo, err := os.Stat(filename)
		if err != nil {
			return err
		}
		if fileInfo.IsDir() {
			return ErrNoFiles
		}
		f, err := os.Open(filename)
		defer f.Close()
		if err != nil {
			return err
		}
		data, err := csvToFloat(f, column)
		if err != nil {
			return err
		}
		consolidate = append(consolidate, data...)
	}
	_, err := fmt.Fprintln(w, opFunc.operation(consolidate))
	return err
}
