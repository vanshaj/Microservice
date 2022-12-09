package render

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

func RenderTemplate(w http.ResponseWriter, path string) {
	fmt.Println("LOGGER - filePath", path)
	myCache, err := createCache()
	if err != nil {
		fmt.Fprintln(w, "Unable to render template, reason ", err)
		return
	}
	basePath := filepath.Base(path)
	v, ok := myCache[basePath]
	if !ok {
		fmt.Fprintln(w, "Unable to find template ", basePath)
		return
	}
	buf := &bytes.Buffer{}
	err = v.Execute(buf, nil)
	if err != nil {
		fmt.Fprintln(w, "Unable to execute ", basePath, " reason ", err)
		return
	}
	buf.WriteTo(w)
}

func createCache() (map[string]*template.Template, error) {
	myCache := map[string]*template.Template{}
	files, err := filepath.Glob("./templates/*.page.tmpl")
	if err != nil {
		return myCache, err
	}
	for _, file := range files {
		name := filepath.Base(file)
		log.Print("DEBUG: ", "file path is ", file)
		fileTemplate, err := template.New(name).ParseFiles(file)
		if err != nil {
			return myCache, err
		}
		_, err = filepath.Glob("./templates/*.layout.tmpl")
		if err != nil {
			return myCache, err
		}
		fileTemplate, err = fileTemplate.ParseGlob("./templates/*.layout.tmpl")
		if err != nil {
			return myCache, err
		}
		myCache[name] = fileTemplate
	}
	return myCache, nil
}
