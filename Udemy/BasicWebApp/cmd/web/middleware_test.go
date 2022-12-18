package main

import (
	"fmt"
	"net/http"
	"testing"
)

func TestNoSurf(t *testing.T) {
	myHandler := &customHandler{}
	h := NoSurf(myHandler)
	switch v := h.(type) {
	case http.Handler:
		//do nothing
	default:
		t.Error(fmt.Sprintln("Type is not http.handler but it is ", v))
	}
}
