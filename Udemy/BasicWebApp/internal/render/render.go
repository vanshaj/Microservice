package render

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"

	"github.com/vanshaj/Microservice/Udemy/BasicWebApp/internal/config"
	"github.com/vanshaj/Microservice/Udemy/BasicWebApp/internal/models"
)

var app *config.AppConfig

func NewTemplate(a *config.AppConfig) {
	app = a
}

func RenderTemplate(w http.ResponseWriter, path string, td *models.TemplateData) {
	myCache := app.TemplateCache
	if app.UseCache == false {
		cache, err := CreateTemplateCache()
		if err != nil {
			fmt.Fprintln(w, "unable to create cache, reason ", err)
			return
		}
		app.TemplateCache = cache
	}
	basePath := filepath.Base(path)
	v, ok := myCache[basePath]
	if !ok {
		fmt.Fprintln(w, "Unable to find template ", basePath)
		return
	}
	buf := &bytes.Buffer{}
	var err error
	if td == nil {
		err = v.Execute(buf, nil)
	} else {
		fmt.Println("Template data is ", td.CSRFToken)
		err = v.Execute(buf, td)
	}
	if err != nil {
		fmt.Fprintln(w, "Unable to execute ", basePath, " reason ", err)
		return
	}
	buf.WriteTo(w)
}

func CreateTemplateCache() (map[string]*template.Template, error) {
	myCache := map[string]*template.Template{}
	files, err := filepath.Glob("./templates/*.page.tmpl")
	if err != nil {
		return myCache, err
	}
	for _, file := range files {
		name := filepath.Base(file)
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
