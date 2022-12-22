package main

import (
	"compress/gzip"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
)

func filterout(path string, extension string, size int64, info fs.FileInfo) bool {
	if info.IsDir() || filepath.Ext(path) != extension {
		log.Printf("%s is a directory\n", path)
		return false
	}
	if size > info.Size() {
		log.Printf("%s is a file but size is less\n", path)
		return false
	}
	log.Printf("%s is a file and size is greater\n", path)
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

func archiveFile(dest string, root string, path string) error {
	info, err := os.Stat(dest)
	if err != nil {
		return err
	}
	if !info.IsDir() {
		return fmt.Errorf("destination path is not  a directory")
	}
	relDir, err := filepath.Rel(root, filepath.Dir(path))
	if err != nil {
		return err
	}
	destfile := fmt.Sprintf("%s.gz", filepath.Base(dest))
	targetPath := filepath.Join(dest, relDir, destfile)
	if err = os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
		return err
	}
	f, err := os.OpenFile(targetPath, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	in, err := os.Open(path)
	if err != nil {
		return err
	}
	defer in.Close()
	gzipW := gzip.NewWriter(f)
	gzipW.Name = filepath.Base(path)

	if _, err = io.Copy(gzipW, in); err != nil {
		return err
	}
	if err := gzipW.Close(); err != nil {
		return err
	}
	return nil
}
