package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// Handler responds to an HTTP reqest
//type Handler interface {
// ServeHTTP(ResponseWriter, *Request) this should write reply headers and data to Resposne Writer and then return
//}

func main() {
	//Greedy matching anything without other matches will execute this
	//This automatically converts our function to a handler type and then register it to a default serve mux.
	// Here server is a default handler and that default handler is a http.ServeMux
	http.HandleFunc("/", func(resp http.ResponseWriter, req *http.Request) {
		log.Println("Hello World")

		//Read data from user
		data, err := ioutil.ReadAll(req.Body)
		if err != nil {
			resp.WriteHeader(http.StatusBadRequest)
			resp.Write([]byte("Oops didn't work"))
			// Or we can use http.Error(resp, "Opps didn't work", http.StatusBadRequest)
			return
		}

		log.Printf("Data %s", data)

		//Write data back to user
		fmt.Fprintf(resp, "hello I receive your data as %s\n", data) // http.ResponseWriter is of type Writer interface and Fprintf() takes a writer interface to write to it

	})
	http.HandleFunc("/bye", func(resp http.ResponseWriter, req *http.Request) {
		log.Println("bye")
	})

	http.ListenAndServe(":8082", nil) // Currently uses default serve mux , here it will check which of the handlers are registered and call ServeHTTP methods of the handlers
}
