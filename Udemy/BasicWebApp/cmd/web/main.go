package main

import (
	"log"
	"net/http"

	"github.com/vanshaj/Microservice/Udemy/BasicWebApp/pkg/config"
	handler "github.com/vanshaj/Microservice/Udemy/BasicWebApp/pkg/handlers"
	"github.com/vanshaj/Microservice/Udemy/BasicWebApp/pkg/render"
)

func main() {
	app := &config.AppConfig{
		UseCache: false,
	}
	myCache, err := render.CreateTemplateCache()

	app.TemplateCache = myCache
	if err != nil {
		log.Println(" failed creating cache ", err)
		return
	}

	render.NewTemplate(app)
	repo := handler.NewHandler(app)
	router := routes(repo)
	// http.HandleFunc("/home", repo.Home)
	// http.HandleFunc("/about", repo.About)
	// http.HandleFunc("/data", repo.Data)
	_ = http.ListenAndServe(":8080", router)

}
