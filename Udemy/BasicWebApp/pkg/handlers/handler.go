package handler

import (
	"net/http"

	"github.com/vanshaj/Microservice/Udemy/BasicWebApp/pkg/config"
	"github.com/vanshaj/Microservice/Udemy/BasicWebApp/pkg/models"
	"github.com/vanshaj/Microservice/Udemy/BasicWebApp/pkg/render"
)

const basePath = "/home/vanshaj/Projects/Golang/Microservice/Udemy/BasicWebApp"

var repo *Handler

type Handler struct {
	App *config.AppConfig
}

func NewHandler(app *config.AppConfig) *Handler {
	return &Handler{
		App: app,
	}
}

func (r *Handler) Home(w http.ResponseWriter, req *http.Request) {
	render.RenderTemplate(w, basePath+"/templates/home.page.tmpl", nil)
}

func (r *Handler) About(w http.ResponseWriter, req *http.Request) {
	render.RenderTemplate(w, basePath+"/templates/about.page.tmpl", nil)
}
func (r *Handler) Data(w http.ResponseWriter, req *http.Request) {
	stringMap := map[string]string{
		"test": "hello",
	}
	td := &models.TemplateData{
		StringMap: stringMap,
	}
	render.RenderTemplate(w, basePath+"/templates/data.page.tmpl", td)
}
