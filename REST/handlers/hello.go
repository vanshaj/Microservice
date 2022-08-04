package handlers

import (
	"log"
	"net/http"

	"github.com/vanshaj/Microservice/REST/data"
)

// Inorder for us to create a handler we have to implement ServerHTTP() method
type Hello struct {
	l *log.Logger
}

//Dependency Injection we can pass any other logger and it will create the object of Hello and return us the same
//Now we write the log depending upon where do we want to write
func NewHello(l *log.Logger) *Hello {
	return &Hello{l}
}

func (h *Hello) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	productList := data.GetProducts()
	err := productList.ToJSON(resp)
	if err != nil {
		http.Error(resp, "Unable to marshal", http.StatusInternalServerError)
	}
}
