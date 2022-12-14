package main

import (
	"flag"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
)

type comfig struct {
	ext         string
	size        int64
	list        bool
	del         bool
	wLog        io.Writer
	archiveDest string
}

func main() {
	ext := flag.String("ext", "txt", "extension of file")
	size := flag.Int64("size", 0, "min size for file")
	list := flag.Bool("list", true, "list the files")
	root := flag.String("root", "/tmp", "root the files")
	del := flag.Bool("del", false, "delete the files")
	logF := flag.String("logF", "", "log file to log deleted file names")
	archiveDestination := flag.String("arch", "", "destination to archive deleted files")
	flag.Parse()
	var (
		delf *os.File
		err  error
	)
	if *logF == "" {
		delf = os.Stdout
	} else {
		delf, err = os.OpenFile(*logF, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0644)
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}
		defer delf.Close()
	}
	c := comfig{
		ext:         *ext,
		size:        *size,
		list:        *list,
		del:         *del,
		wLog:        delf,
		archiveDest: *archiveDestination,
	}
	err = run(*root, os.Stdout, c)
	if err != nil {
		log.Fatal(err)
	}
}

func run(root string, w io.Writer, c comfig) error {
	delLogger := log.New(c.wLog, "Deleted: ", log.LstdFlags)
	err := filepath.Walk(root,
		func(path string, info fs.FileInfo, err error) error {
			info, err = os.Stat(path)
			if err != nil {
				return err
			}
			log.Println("path travelling is ", path)
			observed := filterout(path, c.ext, c.size, info)
			if !observed {
				return nil
			}
			if c.list {
				listfile(w, path)
			}
			if c.archiveDest != "" {
				if err = archiveFile(c.archiveDest, root, path); err != nil {
					return err
				}
			}
			if c.del {
				err := deletefile(path, delLogger)
				if err != nil {
					return err
				}
			}
			return nil
		})
	return err
}
