package handlers

import (
	"log"
	"net/http"
	"regexp"
	"strconv"

	"github.com/dev4ndy/products/data"
)

type Products struct {
	l *log.Logger
}

func NewProducts(l *log.Logger) *Products {
	return &Products{l}
}

func (p *Products) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	isProductsResource, err := regexp.MatchString(`/products/*`, r.URL.Path)
	if err != nil || !isProductsResource {
		http.Error(rw, "router needs to be implemented", http.StatusBadRequest)
		return
	}
	switch r.Method {
	case http.MethodGet:
		p.getProducts(rw, r)
		break
	case http.MethodPost:
		p.addProduct(rw, r)
		break
	case http.MethodPut:
		p.l.Printf("%s", r.URL.Path)
		regex := regexp.MustCompile(`/([0-9]+)`)
		groups := regex.FindAllStringSubmatch(r.URL.Path, -1)
		p.l.Printf("%v", groups)
		if len(groups) != 1 {
			http.Error(rw, "Invalid URI", http.StatusBadRequest)
			return
		}
		if len(groups[0]) != 2 {
			http.Error(rw, "Invalid URI", http.StatusBadRequest)
			return
		}
		idString := groups[0][1]
		id, err := strconv.Atoi(idString)

		if err != nil {
			http.Error(rw, "Invalid URI", http.StatusBadRequest)
			return
		}
		p.updateProduct(id, rw, r)
		break
	default:
		rw.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (p *Products) getProducts(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle GET Products")
	productList := data.GetProducts()
	err := productList.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to marshal json", http.StatusInternalServerError)
	}
}

func (p *Products) addProduct(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle POST Products")
	product := &data.Product{}
	err := product.FromJSON(r.Body)
	if err != nil {
		http.Error(rw, "Unable to unmarshal json", http.StatusBadRequest)
	}

	p.l.Printf("Product: %#v", product)

	data.AddProduct(product)
}

func (p *Products) updateProduct(id int, rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle PUT Products")
	product := &data.Product{}
	err := product.FromJSON(r.Body)
	if err != nil {
		http.Error(rw, "Unable to unmarshal json", http.StatusBadRequest)
	}
	err = data.UpdateProduct(id, product)
	if err == data.ErrProductNotFound {
		http.Error(rw, "Product not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(rw, "Produt not found", http.StatusInternalServerError)
		return
	}
}
