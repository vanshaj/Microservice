package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/vanshaj/Microservice/REST/handlers"
)

// Handler responds to an HTTP reqest
//type Handler interface {
// ServeHTTP(ResponseWriter, *Request) this should write reply headers and data to Resposne Writer and then return
//}

func main() {
	// I am creating a logger which will be passed to my Handler so that handler will automatically write to this
	// If onwards I want to write to other io.Writer I will Pass other parameter and create a new logger based on that
	l := log.New(os.Stdout, "Loggin: ", log.LstdFlags)
	hh := handlers.NewHello(l)
	bb := handlers.NewBye(l)

	//Register the handler with the server
	sm := http.NewServeMux()
	sm.Handle("/", hh)
	sm.Handle("/bye", bb)

	//http.ListenAndServe(":8082", sm) // Now updated to custom mux

	//How to create your own custom server rather than using http.ListenAndServe
	//Manually create a new server
	s := &http.Server{
		Addr:         ":8082",
		Handler:      sm,
		IdleTimeout:  120,              // seconds IdleTimeout is the maximum amount of time to wait for the next request when keep-alives are enabled. If IdleTimeout is zero, the value of ReadTimeout is used. If both are zero, there is no timeout.
		ReadTimeout:  10 * time.Second, //ReadTimeout is the maximum duration for reading the entire request, including the body
		WriteTimeout: 10 * time.Second, //WriteTimeout is the maximum duration before timing out writes of the response
	}

	go func() {
		l.Println("Starting server on port 8082")
		err := s.ListenAndServe()
		if err != nil {
			l.Printf("Error starting server %s\n", err)
			os.Exit(1)
		}
	}()

	// we are trapping sigterm or interrupt and then will shutdown the server
	sigchan := make(chan os.Signal)
	signal.Notify(sigchan, os.Interrupt)
	signal.Notify(sigchan, os.Kill)

	// Block until a signal is received
	sig := <-sigchan
	log.Println("Got signal ", sig)

	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	// s.Shutdown() will wait for all the work in progress to be finished and will not accept any more requests
	s.Shutdown(ctx)
}
