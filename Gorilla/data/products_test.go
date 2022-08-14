package data

import "testing"

func TestCheckValidationFunc(t *testing.T) {
	p := &Product{
		Name:  "mobile",
		Price: 100,
		SKU:   "axx-bxa-cxd",
	}
	err := p.Validate()
	if err != nil {
		t.Fatal(err)
	}
}
