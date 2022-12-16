package main

import (
	"net/http"

	chi "github.com/go-chi/chi/v5"
	handler "github.com/vanshaj/Microservice/Udemy/BasicWebApp/internal/handlers"
)

func routes(repo *handler.Handler) http.Handler {
	mux := chi.NewRouter()
	mux.Use(SetCookie)
	mux.Use(GetIP)
	mux.Use(SessionLoad)
	mux.Use(NoSurf)
	mux.Get("/", repo.Home)
	mux.Get("/about", repo.About)
	mux.Get("/make-reservation", repo.Reservation)
	mux.Get("/generals", repo.Generals)
	mux.Get("/majors", repo.Majors)
	mux.Get("/search-availability", repo.Availability)
	mux.Post("/search-availability", repo.PostAvailability)
	mux.Get("/contact", repo.Contact)
	fileServer := http.FileServer(http.Dir("./static"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))
	return mux
}
