package data

import (
	"encoding/json"
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

func GetProducts() Products {
	return productList
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
