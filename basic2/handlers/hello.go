package handlers

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
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
	h.l.Println("Hello World")
	//Read data from user
	data, err := ioutil.ReadAll(req.Body)
	if err != nil {
		resp.WriteHeader(http.StatusBadRequest)
		resp.Write([]byte("Oops didn't work"))
		// Or we can use http.Error(resp, "Opps didn't work", http.StatusBadRequest)
		return
	}
	h.l.Printf("Data %s", data)
	//Write data back to user
	fmt.Fprintf(resp, "hello I receive your data as %s\n", data) // http.ResponseWriter is of type Writer interface and Fprintf() takes a writer interface to write to it
}
