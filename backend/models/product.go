package models

import (
	"fmt"
)

type Product struct {
	Name    string `json:"name"`
	Barcode string `json:"barcode"`
}

func NewProduct(Name string, Barcode string) Product {
	e := Product{Name, Barcode}
	return e
}

func (e Product) LeavesRemaining() {
	fmt.Printf("%s %s\n", e.Name, e.Barcode)
}
