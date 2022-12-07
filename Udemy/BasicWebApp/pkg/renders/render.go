package render

import (
	"fmt"
	"html/template"
	"net/http"
)

func RenderTemplate(w http.ResponseWriter, path string) {
	fmt.Println("LOGGER - filePath", path)
	parseTemplate, _ := template.ParseFiles(path, "./templates/base.page.tmpl")
	err := parseTemplate.Execute(w, nil)
	if err != nil {
		fmt.Fprintln(w, "Unable to render template")
		return
	}
}
