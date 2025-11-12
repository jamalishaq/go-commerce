package handlers

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"

	"github.com/jamalishaq/go-commerce/product-api/data"
)

// Products is a http.Handler
type Products struct {
	l *log.Logger
}

// NewProducts creates a products handler with the given logger
func NewProducts(l *log.Logger) *Products {
	return &Products{l}
}

// ServeHTTP is the main entry point for the handler and staisfies the http.Handler
// interface
func (p *Products) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	// handle the request for a list of products
	if r.Method == http.MethodGet {
		p.getProducts(rw, r)
		return
	}

	if r.Method == http.MethodPost {
		p.addProduct(rw, r)
		return
	}

	if r.Method == http.MethodPut {
		p.l.Println("PUT", r.URL.Path)
		id, err := getURLID(r.URL.Path, p)
		if err != nil {
			http.Error(rw, "Product not found", http.StatusNotFound)
			return
		}

		p.updateProducts(id, rw, r)
		return
	}

	if r.Method == http.MethodDelete {
		id, err := getURLID(r.URL.Path, p)
		if err != nil {
			http.Error(rw, "Invalid URL", http.StatusNotFound)
			return
		}

		p.deleteProduct(id, rw)
		return
	}
	// catch all
	// if no method is satisfied return an error
	rw.WriteHeader(http.StatusMethodNotAllowed)
}

// getProducts returns the products from the data store
func (p *Products) getProducts(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle GET Products")

	// fetch the products from the datastore
	lp := data.GetProducts()

	// serialize the list to JSON
	err := lp.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to marshal json", http.StatusInternalServerError)
	}
}

func (p *Products) addProduct(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle POST Product")

	prod := &data.Product{}

	err := prod.FromJSON(r.Body)
	if err != nil {
		http.Error(rw, "Unable to unmarshal json", http.StatusBadRequest)
	}

	data.AddProduct(prod)
}

func (p *Products) updateProducts(id int, rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle PUT Product")

	prod := &data.Product{}

	err := prod.FromJSON(r.Body)
	if err != nil {
		http.Error(rw, "Unable to unmarshal json", http.StatusBadRequest)
	}

	err = data.UpdateProduct(id, prod)
	if err == data.ErrProductNotFound {
		http.Error(rw, "Product not found", http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(rw, "Product not found", http.StatusInternalServerError)
		return
	}
}

func (p *Products) deleteProduct(id int, rw http.ResponseWriter) {
	p.l.Println("Handle DELETE Product")
	err := data.DeleteProduct(id)
	if err != nil {
		http.Error(rw, "Unable to delete product", http.StatusInternalServerError)
		return
	}
}

func getURLID(url string, p *Products) (int, error) {
	// expect the id in the URI
	reg := regexp.MustCompile(`/([0-9]+)`)
	g := reg.FindAllStringSubmatch(url, -1)

	if len(g) != 1 {
		p.l.Println("Invalid URI more than one id")
		return -1, fmt.Errorf("invalid url")
	}

	if len(g[0]) != 2 {
		p.l.Println("Invalid URI more than one capture group")
		return -1, fmt.Errorf("invalid url")
	}

	idString := g[0][1]
	id, err := strconv.Atoi(idString)
	if err != nil {
		p.l.Println("Invalid URI unable to convert to numer", idString)
		return -1, fmt.Errorf("invalid url")
	}

	return id, nil
}
