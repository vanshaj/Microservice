package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestSayHelloBack(t *testing.T) {
	data := "hello"
	req := httptest.NewRequest(http.MethodGet, "/", strings.NewReader(data))
	w := httptest.NewRecorder()
	sayHelloBack(w, req)
	res := w.Result()
	defer res.Body.Close()
	dataGot, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Errorf("expected error to be nil got %v\n", err)
	}
	expected := fmt.Sprintf("hello I receive your data as %s\n", data)

	if string(dataGot) != expected {
		t.Errorf("expected \"%s\" but got \"%s\"", expected, string(dataGot))
	}

}
