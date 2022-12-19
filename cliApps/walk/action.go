package main

import (
	"fmt"
	"io"
	"io/fs"
	"path/filepath"
)

func filterout(path string, extension string, size int64, info fs.FileInfo) bool {
	if info.IsDir() || filepath.Ext(path) != extension {
		return false
	}
	if size > info.Size() {
		return false
	}
	return true
}

func listfile(w io.Writer, path string) {
	fmt.Fprintf(w, "%s\n", path)
}
