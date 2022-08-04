package handlers

import (
	"log"
	"net/http"
	"time"
)

type bye struct {
	l *log.Logger
}

func NewBye(l *log.Logger) *bye {
	return &bye{l}
}

func (h *bye) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	/*
		ctx, cancel := context.WithTimeout(req.Context(), 1*time.Second)
		go func() {
			<-ctx.Done()
			cancel()
		}()
		doesnot work*/
	h.l.Println("Entered Bye")
	time.Sleep(3 * time.Second)
	resp.Write([]byte("bye from server"))
}
