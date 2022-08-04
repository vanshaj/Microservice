package data

import (
	"encoding/json"
	"fmt"
	"io"
)

// Product defines the structure for an API product
type Product struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float32 `json:"price"`
	SKU         string  `json:"sku"`
}

type Products []Product

func (p *Products) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w) //encoder takes writer interface to write to after encoding to json
	return e.Encode(p)      // encode methods takes the struct and write the json conversion to the writer received above
}

func (p *Product) FromJSON(r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(p)
}

func GetProducts() Products {
	return productList
}

func AddProduct(p *Product) error {
	p.ID = getNextID()
	productList = append(productList, *p)
	return nil
}

var ErrProductNotFound = fmt.Errorf("Product not found")

func UpdateProduct(id int, updateProd *Product) error {
	for index, prod := range productList {
		if prod.ID == id {
			productList[index] = *updateProd
			return nil
		}
	}
	return ErrProductNotFound
}

func getNextID() int {
	lenP := productList[len(productList)-1].ID
	lenP += 1
	return lenP
}

var productList = []Product{
	Product{
		ID:          1,
		Name:        "C",
		Description: "C",
		Price:       12.34,
		SKU:         "abc323",
	},
	Product{
		ID:          2,
		Name:        "b",
		Description: "b",
		Price:       2.29,
		SKU:         "fjd34",
	},
}
