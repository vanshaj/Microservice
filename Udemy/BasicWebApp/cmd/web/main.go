package main

import (
	"log"
	"net/http"
	"time"

	scs "github.com/alexedwards/scs/v2"
	"github.com/vanshaj/Microservice/Udemy/BasicWebApp/internal/config"
	handler "github.com/vanshaj/Microservice/Udemy/BasicWebApp/internal/handlers"
	"github.com/vanshaj/Microservice/Udemy/BasicWebApp/internal/render"
)

var app *config.AppConfig

func main() {
	app = &config.AppConfig{
		UseCache:     false,
		InProduction: false,
	}
	session := scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction
	app.SessionManager = session
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
