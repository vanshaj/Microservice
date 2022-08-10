package handlers

import (
	"io/ioutil"
	"log"
	"net/http/httptest"
	"os"
	"testing"
)

func TestBye(t *testing.T) {
	req := httptest.NewRequest("GET", "localhost:8082/bye", nil)
	res := httptest.NewRecorder()
	l := log.New(os.Stdout, "Loggin: ", log.LstdFlags)
	b := NewBye(l)
	b.ServeHTTP(res, req)

	result := res.Result()
	defer result.Body.Close()
	data, err := ioutil.ReadAll(result.Body)
	if err != nil {
		t.Errorf("Unable to read data from the output")
	}

	expected := "bye from server"
	if expected != string(data) {
		t.Errorf("Expected %s, got %s", expected, data)
	}

}
