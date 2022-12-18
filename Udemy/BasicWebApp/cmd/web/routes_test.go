package main

import (
	"fmt"
	"github.com/go-chi/chi"
	"github.com/vanshaj/Microservice/Udemy/BasicWebApp/internal/config"
	"testing"
)

func TestRoutes(t *testing.T) {
	var app *config.AppConfig
	mux := routes(app)
	switch v := mux.(type) {
	case *chi.Mux:
		//
	default:
		t.Error(fmt.Sprintln("type is incorrect, ", v))
	}
}
