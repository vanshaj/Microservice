package main

import (
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"log"
	"testing"
)

func TestRoutes(t *testing.T) {
	var app *app.Config
	mux := routes(&app)
	switch v := mux.(type) {
	case *chi.Mux:
		//
	default:
		t.Error(fmt.Sprintln("type is incorrect"))
	}
}
