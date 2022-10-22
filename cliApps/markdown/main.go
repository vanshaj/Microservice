package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday/v2"
)

const (
	header = `<!DOCTYPE html>
	 <html> 
	 <head> 
	 <meta http-equiv="content-type" content="text/html; charset=utf-8"> <title>Markdown Preview Tool</title> 
	 </head> 
	 <body>`

	footer = `</body> </html>`
)

func main() {
	file := flag.String("file", "", "file name to parse")
	flag.Parse()

	if *file == "" {
		fmt.Errorf("File name to parse is empty")
		os.Exit(1)
	}

	err := run(*file)
	if err != nil {
		fmt.Errorf(err.Error())
		os.Exit(1)
	}
}

func run(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	input, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}
	htmlData := parseContent(input)
	outName := fmt.Sprintf("%s.html", filepath.Base(filename))
	fmt.Println(outName)
	return saveHTML(outName, htmlData)
}

func parseContent(input []byte) []byte {
	output := blackfriday.Run(input)
	body := bluemonday.UGCPolicy().SanitizeBytes(output)

	var htmlOutput bytes.Buffer
	htmlOutput.WriteString(header)
	htmlOutput.Write(body)
	htmlOutput.WriteString(footer)
	return htmlOutput.Bytes()
}

func saveHTML(filename string, htmlData []byte) error {
	return os.WriteFile(filename, htmlData, 0644)
}
