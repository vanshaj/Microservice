package main

import (
	"strings"
	"testing"
)

func TestCountWords(t *testing.T) {
	tests := []struct {
		Name           string
		Content        string
		ParseLines     bool
		ExpectedOutput int
	}{
		{"ParseWords", "word1 word2 word3", false, 3},
		{"ParseLines", "word1\n word2 word3\n word4", true, 3},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			res := countScanner(strings.NewReader(tt.Content), tt.ParseLines)
			if res != tt.ExpectedOutput {
				t.Errorf("Expected %d, got %d ", tt.ExpectedOutput, res)
			}
		})
	}
}
