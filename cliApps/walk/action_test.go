package main

import (
	"os"
	"testing"
)

func TestFilterOut(t *testing.T) {
	testCases := []struct {
		Name      string
		Path      string
		Extension string
		Size      int64
		Expected  bool
	}{
		{"DirectoryMatch", "./testdata", "", 0, false},
		{"FilterExtensionMatch", "./testdata/hello.txt", ".txt", 0, true},
		{"SizeMatch", "./testdata/hello.txt", ".txt", 30, true},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			info, err := os.Stat(tc.Path)
			if err != nil {
				t.Fatal(err)
			}
			f := filterout(tc.Path, tc.Extension, tc.Size, info)
			if f != tc.Expected {
				t.Errorf("Expected %t got %t\n", tc.Expected, f)
			}
		})
	}
}
