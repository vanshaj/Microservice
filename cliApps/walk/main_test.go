package main

import (
	"bytes"
	"fmt"
	"os"
	"testing"
)

func TestRun(t *testing.T) {
	testCases := []struct {
		name     string
		root     string
		c        comfig
		expected string
	}{
		{"EmptyDirectory", "./testdata/empty/", comfig{".txt", 32, true, false, &bytes.Buffer{}, ""}, ""},
		{"TextFile", "./testdata", comfig{".txt", 20, true, false, &bytes.Buffer{}, ""}, "testdata/hello.txt\n"},
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

func createTempDir(t *testing.T, files map[string]int) (string, func()) {
	t.Helper()
	tempDir, err := os.MkdirTemp("", "walktest")
	if err != nil {
		t.Fatal(err)
	}
	for k, v := range files {
		for j := 1; j <= v; j++ {
			f, _ := os.CreateTemp(tempDir, fmt.Sprintf("*_file%s", k))
			_, err = f.Write([]byte("rabdom dRa"))
		}
	}
	teardown := func() {
		os.RemoveAll(tempDir)
	}
	return tempDir, teardown
}

func TestDelExtension(t *testing.T) {
	testCases := []struct {
		name        string
		cfg         comfig
		extNoDelete string
		nDelete     int
		nNoDelete   int
		exoected    string
	}{
		{"DeleteExtNoMatch", comfig{ext: ".log", del: true, wLog: &bytes.Buffer{}}, ".gz", 0, 10, ""},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var buffer bytes.Buffer
			tempDir, teardown := createTempDir(t, map[string]int{
				tc.cfg.ext:     tc.nDelete,
				tc.extNoDelete: tc.nNoDelete,
			})
			defer teardown()
			eer := run(tempDir, &buffer, tc.cfg)
			if eer != nil {
				t.Fatal(eer)
			}
			output, err := os.Open(tempDir)
			if err != nil {
				t.Error(err)
			}
			outputFiles, err := output.ReadDir(0)
			if err != nil {
				t.Error(err)
			}
			if len(outputFiles) != tc.nNoDelete {
				t.Errorf("expected no files but it contains")
			}
		})
	}

}
