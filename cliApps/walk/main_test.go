package main

import (
	"bytes"
	"testing"
)

func TestRun(t *testing.T) {
	testCases := []struct {
		name     string
		root     string
		c        comfig
		expected string
	}{
		{"EmptyDirectory", "./testdata/empty/", comfig{".txt", 32, true}, ""},
		{"TextFile", "./testdata", comfig{".txt", 20, true}, "testdata/hello.txt\n"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			output := &bytes.Buffer{}
			err := run(tc.root, output, tc.c)
			if err != nil {
				t.Fatal(err)
			}
			actual := output.String()
			if actual != tc.expected {
				t.Errorf("expected '%v', but actual is '%v' \n", tc.expected, actual)
			}
		})
	}
}
