package handlers

import (
	"context"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/vanshaj/Microservice/Gorilla/data"
)

// Inorder for us to create a handler we have to implement ServerHTTP() method
type Products struct {
	l *log.Logger
}

//Dependency Injection we can pass any other logger and it will create the object of Products and return us the same
//Now we write the log depending upon where do we want to write
func NewProducts(l *log.Logger) *Products {
	return &Products{l}
}

func (p *Products) GetProducts(resp http.ResponseWriter, req *http.Request) {
	productList := data.GetProducts()
	err := productList.ToJSON(resp)
	if err != nil {
		http.Error(resp, "Unable to marshal", http.StatusInternalServerError)
	}
}

func (p *Products) AddProduct(resp http.ResponseWriter, req *http.Request) {
	prod := req.Context().Value(KeyProduct{}).(data.Product)
	data.AddProduct(&prod)
	p.l.Printf("Prod: %#v", prod)
}

func (h *Products) UpdateProducts(resp http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(resp, "Unable to convert id", http.StatusBadRequest)
		return
	}
	prod := req.Context().Value(KeyProduct{}).(data.Product)
	err = data.UpdateProduct(id, &prod)
	if err == data.ErrProductNotFound {
		http.Error(resp, "Product not found", http.StatusNotFound)
		return
	} else if err != nil {
		h.l.Println(err)
		http.Error(resp, "Not found", http.StatusInternalServerError)
		return
	}
}

type KeyProduct struct{}

func (p *Products) MiddlewareProductHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		prod := data.Product{}
		err := prod.FromJSON(req.Body)
		if err != nil {
			p.l.Println("Error deserializing product", err)
			http.Error(resp, "Error deserializing", http.StatusBadRequest)
			return
		}
		ctx := context.WithValue(req.Context(), KeyProduct{}, prod)
		req = req.WithContext(ctx)

		next.ServeHTTP(resp, req)
	})
}
