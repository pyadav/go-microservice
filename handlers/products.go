package handlers

import (
	"log"
	"net/http"
	"regexp"
	"strconv"

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

// ServeHTTP is a main entry point for the handler and stisfies the http.Handler
// interface
func (h *Products) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		h.getProducts(w, r)
		return
	}

	if r.Method == http.MethodPost {
		h.addProducts(w, r)
		return
	}

	if r.Method == http.MethodPut {
		// except id in the URI
		regex := regexp.MustCompile(`/([0-9]+)`)
		matches := regex.FindAllStringSubmatch(r.URL.Path, -1)

		if len(matches) != 1 {
			http.Error(w, "Invalid URI", http.StatusBadRequest)
			return
		}

		matchString := matches[0][1]
		id, err := strconv.Atoi(matchString)
		if err != nil {
			http.Error(w, "Invalid URI", http.StatusBadRequest)
			return
		}

		h.updateProduct(id, w, r)
	}

	// catch all
	w.WriteHeader(http.StatusMethodNotAllowed)
}

// getProducts return products form the data stores
func (h *Products) getProducts(w http.ResponseWriter, r *http.Request) {
	productsList := data.GetProducts()

	if err := productsList.ToJSON(w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// addProducts add product into the data stores
func (h *Products) addProducts(w http.ResponseWriter, r *http.Request) {
	h.logger.Println("Add a product")

	product := &data.Product{}
	if err := product.FromJSON(r.Body); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	data.AddProduct(product)
}

// updateProduct will update info of a product
func (h *Products) updateProduct(id int, w http.ResponseWriter, r *http.Request) {
	h.logger.Println("update a product")

	product := &data.Product{}
	if err := product.FromJSON(r.Body); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	err := data.UpdateProduct(id, product)
	if err == data.ErrProductNotFound {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
