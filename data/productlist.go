package data

import (
	"encoding/json"
	"fmt"
	"io"
	"time"
)

type Product struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float32 `json:"price"`
	SKU         string  `json:"sku"`
	CreatedOn   string  `json:"-"` // to omit the value when use json.marshal
	UpdatedOn   string  `json:"-"`
	DeleteOn    string  `json:"-"`
}

type Products []*Product

var ErrProductNotFound = fmt.Errorf("Product not found")

func (p *Products) ToJSON(w io.Writer) error {
	encoder := json.NewEncoder(w)
	return encoder.Encode(p)
}

func (p *Product) FromJSON(r io.Reader) error {
	decoder := json.NewDecoder(r)
	return decoder.Decode(p)
}

func GetProducts() Products {
	return productList
}

func AddProduct(p *Product) {
	p.ID = getNextID()
	productList = append(productList, p)
}

func UpdateProduct(id int, product *Product) error {
	_, index, err := getProductById(id)
	if err != nil {
		return err
	}
	product.ID = id
	productList[index] = product
	return nil
}

func getProductById(id int) (*Product, int, error) {
	for i, prod := range productList {
		if prod.ID == id {
			return prod, i, nil
		}
	}
	return nil, -1, ErrProductNotFound
}

func getNextID() int {
	lastProduct := productList[len(productList)-1]
	return lastProduct.ID + 1
}

var productList = []*Product{
	&Product{
		ID:          1,
		Name:        "Samsung Galaxy S22",
		Description: "Android Smartphone",
		Price:       1000.00,
		SKU:         "sgs22+",
		CreatedOn:   time.Now().String(),
		UpdatedOn:   time.Now().String(),
	},
	&Product{
		ID:          2,
		Name:        "Iphone 13 Pro Max",
		Description: "Apple Smartphone",
		Price:       1800.00,
		SKU:         "i13pm",
		CreatedOn:   time.Now().String(),
		UpdatedOn:   time.Now().String(),
	},
}
