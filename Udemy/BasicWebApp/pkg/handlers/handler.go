package handler

import (
	"net/http"

	render "github.com/vanshaj/Microservice/Udemy/BasicWebApp/pkg/renders"
)

const basePath = "/home/vanshaj/Projects/Golang/Microservice/Udemy/BasicWebApp"

func Home(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, basePath+"/templates/home.page.tmpl")
}

func About(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, basePath+"/templates/about.page.tmpl")
}