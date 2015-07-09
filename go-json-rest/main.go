package main

import (
	"fmt"
	"github.com/ant0ine/go-json-rest/rest"
	"log"
	"net/http"
	"sync"
)

func init() {
	store["1"] = &Product{"1", "Gopher"}
	store["2"] = &Product{"2", "Ice Cream"}
	store["3"] = &Product{"3", "Trip to L.A."}
}

func main() {
	fmt.Println("****** Test for a product ******")
	api := rest.NewApi()
	api.Use(rest.DefaultDevStack...)
	router, err := rest.MakeRouter(
		rest.Get("/products", GetAllProducts),
		rest.Post("/products", PostProduct),
		rest.Get("/products/:id", GetProduct),
		rest.Delete("/products/:id", DeleteProduct),
	)
	if err != nil {
		log.Fatal(err)
	}
	api.SetApp(router)
	log.Fatal(http.ListenAndServe(":8080", api.MakeHandler()))
	log.Println("Listening on port 8080 -> go to http://localhost:8080")
}

type Product struct {
	Code string
	Name string
}

var store = map[string]*Product{}

var lock = sync.RWMutex{}

func GetProduct(w rest.ResponseWriter, r *rest.Request) {
	code := r.PathParam("id")

	lock.RLock()
	var product *Product
	if store[code] != nil {
		product = &Product{}
		*product = *store[code]
	}
	lock.RUnlock()

	if product == nil {
		rest.NotFound(w, r)
		return
	}
	w.WriteJson(product)
}

func GetAllProducts(w rest.ResponseWriter, r *rest.Request) {
	lock.RLock()
	products := make([]Product, len(store))
	i := 0
	for _, product := range store {
		products[i] = *product
		i++
	}
	lock.RUnlock()
	w.WriteJson(&products)
}

func PostProduct(w rest.ResponseWriter, r *rest.Request) {
	product := Product{}
	err := r.DecodeJsonPayload(&product)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if product.Code == "" {
		rest.Error(w, "product code required", 400)
		return
	}
	if product.Name == "" {
		rest.Error(w, "product name required", 400)
		return
	}
	lock.Lock()
	store[product.Code] = &product
	lock.Unlock()
	w.WriteJson(&product)
}

func DeleteProduct(w rest.ResponseWriter, r *rest.Request) {
	code := r.PathParam("id")

	msg := "All good : " + store[code].Name + " has been deleted "

	delete(store, code)

	w.WriteJson(&msg)
}
