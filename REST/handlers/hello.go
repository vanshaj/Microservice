package handlers

import (
	"log"
	"net/http"
	"regexp"
	"strconv"

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
	if req.Method == http.MethodGet {
		h.getProducts(resp, req)
		return
	}
	if req.Method == http.MethodPost {
		h.addProduct(resp, req)
		return
	}
	if req.Method == http.MethodPut {
		// expect the id in the URI
		r := regexp.MustCompile(`/([0-9]+)`)
		g := r.FindAllStringSubmatch(req.URL.Path, -1)
		if len(g) != 1 {
			h.l.Println("Invalid URI more than one id")
			http.Error(resp, "Invalid URI", http.StatusBadRequest)
			return
		}
		if len(g[0]) != 2 {
			h.l.Println("Invalid URI more than 1 capture group")
			http.Error(resp, "Invalid URI", http.StatusBadRequest)
			return
		}
		idString := g[0][1]
		id, err := strconv.Atoi(idString)
		if err != nil {
			http.Error(resp, "Id doesnot contain integer", http.StatusBadRequest)
		}
		h.updateProducts(id, resp, req)
		return
	}
	resp.WriteHeader(http.StatusMethodNotAllowed)
}

func (p *Hello) getProducts(resp http.ResponseWriter, req *http.Request) {
	productList := data.GetProducts()
	err := productList.ToJSON(resp)
	if err != nil {
		http.Error(resp, "Unable to marshal", http.StatusInternalServerError)
	}
}

func (p *Hello) addProduct(resp http.ResponseWriter, req *http.Request) {
	prod := &data.Product{}
	err := prod.FromJSON(req.Body)
	if err != nil {
		http.Error(resp, "Unable to unmarshal json", http.StatusBadRequest)
	}
	data.AddProduct(prod)
	p.l.Printf("Prod: %#v", prod)
}

func (h *Hello) updateProducts(id int, resp http.ResponseWriter, req *http.Request) {
	prod := &data.Product{}
	err := prod.FromJSON(req.Body)
	if err != nil {
		http.Error(resp, "Unable to unmarshal json", http.StatusBadRequest)
		return
	}
	h.l.Println("Id to be updated ", id)
	err = data.UpdateProduct(id, prod)
	if err == data.ErrProductNotFound {
		http.Error(resp, "Product not found", http.StatusNotFound)
		return
	} else if err != nil {
		h.l.Println(err)
		http.Error(resp, "Not found", http.StatusInternalServerError)
		return
	}
}
