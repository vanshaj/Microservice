package main

import (
	"strings"
	"testing"
)

func TestCountWords(t *testing.T) {
	b := strings.NewReader("word1 word2\n")
	exp := 2
	res := countScanner(b)
	if res != exp {
		t.Errorf("Expected %d got %d", exp, res)
	}
}
