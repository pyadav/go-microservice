package data

import "testing"

func TestCheckValidations(t *testing.T) {
	product := &Product{
		Name:  "TestProduct",
		Price: 1.00,
		SKU:   "abc-abc-abc",
	}

	err := product.Validate()
	if err != nil {
		t.Fatal(err)
	}
}
