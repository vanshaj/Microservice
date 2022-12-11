package main

import (
	"bytes"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

const (
	inputFile  = "./testdata/test1.md"
	outputFile = "test1.md.html"
	goldenFile = "./testdata/test1.md.html"
)

func TestParseContent(t *testing.T) {
	testCase := []struct {
		name string
	}{
		{"normal markdown test"},
	}
	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			input, err := os.Open(inputFile)
			if err != nil {
				t.FailNow()
			}
			defer input.Close()
			inputData, err := ioutil.ReadAll(input)
			if err != nil {
				t.Fatal(err)
			}
			actualOutput := parseContent(inputData)
			expectedOutput, err := ioutil.ReadFile(goldenFile)

			if !bytes.Equal(expectedOutput, actualOutput) {
				t.Error("markdown fails")
			}
		})
	}
}

func TestMain(t *testing.T) {
	// tt := []struct{

	// }
	var mockStdOut bytes.Buffer
	if err := run(inputFile, &mockStdOut, true); err != nil {
		t.Fatal(err)
	}
	resultFile := strings.TrimSpace(mockStdOut.String())
	actualData, err := ioutil.ReadFile(resultFile)
	if err != nil {
		t.Fatal(err)
	}
	expectedData, err := ioutil.ReadFile(goldenFile)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(expectedData, actualData) {
		t.Error("Result does not match")
	}
	os.Remove(resultFile)
}
