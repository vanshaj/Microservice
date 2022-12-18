package main

import (
	"log"
	"net/http"
	"os"
	"testing"
)

// Will get executed before each test
func TestMain(m *testing.M) {
	log.Println("Before test")
	os.Exit(m.Run())
}

//Place to create your custom objects, structs to be used
type customHandler struct{}

func (c *customHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {}
