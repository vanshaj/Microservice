package handler

import (
	"fmt"
	"net/http"

	"github.com/vanshaj/Microservice/Udemy/BasicWebApp/pkg/config"
	"github.com/vanshaj/Microservice/Udemy/BasicWebApp/pkg/render"
)

const basePath = "/home/vanshaj/Projects/Golang/Microservice/Udemy/BasicWebApp"

var Repo *Repository

type Repository struct {
	App *config.AppConfig
}

func NewRepo(app *config.AppConfig) *Repository {
	return &Repository{
		App: app,
	}
}

func NewHandlers(r *Repository) {
	Repo = r
}

func (r *Repository) Home(w http.ResponseWriter, req *http.Request) {
	if r.App.UseCache == false {
		cache, err := render.CreateTemplateCache()
		if err != nil {
			fmt.Fprintln(w, "unable to create cache, reason ", err)
			return
		}
		r.App.TemplateCache = cache
	}
	render.RenderTemplate(w, basePath+"/templates/home.page.tmpl")
}

func (r *Repository) About(w http.ResponseWriter, req *http.Request) {
	render.RenderTemplate(w, basePath+"/templates/about.page.tmpl")
}
