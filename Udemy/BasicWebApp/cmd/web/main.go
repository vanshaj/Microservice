package main

import (
	"net/http"

	handler "github.com/vanshaj/Microservice/Udemy/BasicWebApp/pkg/handlers"
)

func main() {
	http.HandleFunc("/home", handler.Home)
	http.HandleFunc("/about", handler.About)
	_ = http.ListenAndServe(":8080", nil)

}
