package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

func main() {
	fmt.Println(countScanner(os.Stdin))
}

func countScanner(r io.Reader) int {
	var count int
	scanner := bufio.NewScanner(r)
	scanner.Split(bufio.ScanWords)

	for scanner.Scan() {
		count++
	}

	return count
}
