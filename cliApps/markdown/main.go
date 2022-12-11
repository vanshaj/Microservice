package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"runtime"

	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday/v2"
)

const (
	header = `<!DOCTYPE html>
	 <html> 
	 <head> 
	 <meta http-equiv="content-type" content="text/html; charset=utf-8"> <title>Markdown Preview</title> 
	 </head> 
	 <body>`

	footer = `</body> </html>`
)

func main() {
	file := flag.String("file", "", "file name to parse")
	skipPreview := flag.Bool("skip", false, "skip preview")
	flag.Parse()

	if *file == "" {
		fmt.Errorf("File name to parse is empty")
		os.Exit(1)
	}

	err := run(*file, os.Stdout, *skipPreview)
	if err != nil {
		fmt.Errorf(err.Error())
		os.Exit(1)
	}
}

func run(filename string, out io.Writer, skipPreview bool) error {
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

	temp, err := ioutil.TempFile("", "mdp*.html")
	if err != nil {
		return err
	}
	if err := temp.Close(); err != nil {
		return err
	}
	outName := temp.Name()
	fmt.Fprint(out, outName)
	if err := saveHTML(outName, htmlData); err != nil {
		return err
	}
	if skipPreview {
		return nil
	}
	return preview(outName)
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

func preview(fileName string) error {
	cName := ""
	var cParams []string
	switch runtime.GOOS {
	case "windows":
		cName = "cmd.exe"
		cParams = []string{"/C", "start"}
	default:
		return fmt.Errorf("OS not supported")
	}
	cParams = append(cParams, fileName)
	cPath, err := exec.LookPath(cName)
	if err != nil {
		return err
	}
	return exec.Command(cPath, cParams...).Run()
}
