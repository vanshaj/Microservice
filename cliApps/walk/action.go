package main

import (
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
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
func deletefile(path string, delLogger *log.Logger) error {
	err := os.Remove(path)
	if err != nil {
		return err
	}
	delLogger.Println("Deleting file : ", path)
	return nil
}
