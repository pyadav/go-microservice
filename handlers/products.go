package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/pyadav/microservice/data"
)

// Products is a http.Handler
type Products struct {
	logger *log.Logger
}

// NewProducts creats a product handler with a given configuration
func NewProducts(l *log.Logger) *Products {
	return &Products{l}
}

// GetProducts return products form the data stores
func (p *Products) GetProducts(w http.ResponseWriter, r *http.Request) {
	productsList := data.GetProducts()

	if err := productsList.ToJSON(w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// AddProduct add product into the data stores
func (p *Products) AddProduct(w http.ResponseWriter, r *http.Request) {
	p.logger.Println("Add a product")

	product := r.Context().Value(KeyProduct{}).(data.Product)
	data.AddProduct(&product)
}

// UpdateProduct will update info of a product
func (p Products) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	product := r.Context().Value(KeyProduct{}).(data.Product)

	err = data.UpdateProduct(id, &product)
	if err == data.ErrProductNotFound {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// KeyProduct ...
type KeyProduct struct{}

// MiddlewareValidateProduct ...
func (p Products) MiddlewareValidateProduct(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		prod := data.Product{}

		err := prod.FromJSON(r.Body)
		if err != nil {
			p.logger.Println("[ERROR] deserializing product", err)
			http.Error(rw, "Error reading product", http.StatusBadRequest)
			return
		}

		// validate the product
		if err := prod.Validate(); err != nil {
			p.logger.Println("[ERROR] validating product", err)
			http.Error(rw, fmt.Sprintf("Error validating product: %s", err), http.StatusBadRequest)
			return
		}

		// add the product to the context
		ctx := context.WithValue(r.Context(), KeyProduct{}, prod)
		r = r.WithContext(ctx)

		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(rw, r)
	})
}
