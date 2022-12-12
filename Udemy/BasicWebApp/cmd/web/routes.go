package main

import (
	"net/http"

	chi "github.com/go-chi/chi/v5"
	handler "github.com/vanshaj/Microservice/Udemy/BasicWebApp/pkg/handlers"
)

func routes(repo *handler.Handler) http.Handler {
	mux := chi.NewRouter()
	mux.Use(SetCookie)
	mux.Use(GetIP)
	mux.Use(SessionLoad)
	mux.Get("/", repo.Home)
	mux.Get("/about", repo.About)
	mux.Get("/data", repo.Data)
	fileServer := http.FileServer(http.Dir("./static"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))
	return mux
}
