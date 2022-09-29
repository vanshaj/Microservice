package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
)

func main() {
	parseLines := flag.Bool("l", false, "count line")
	flag.Parse()
	fmt.Println(countScanner(os.Stdin, *parseLines))
}

func countScanner(r io.Reader, parseLines bool) int {
	var count int
	scanner := bufio.NewScanner(r)
	if !parseLines {
		scanner.Split(bufio.ScanWords)
	}

	for scanner.Scan() {
		count++
	}

	return count
}
