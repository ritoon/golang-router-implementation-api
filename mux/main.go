package main

import (
	//"compress/gzip"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

type Product struct {
	Code string
	Name string
}

var store = map[string]*Product{}

// simulate an Id in a database
var ID = 3

func init() {
	store["1"] = &Product{"1", "Gopher"}
	store["2"] = &Product{"2", "Ice Cream"}
	store["3"] = &Product{"3", "Trip to L.A."}
}

func main() {
	fmt.Println("****** Test for a product ******")
	// Router creation
	r := mux.NewRouter()
	// subrouter creation
	s := r.Host("localhost").Subrouter()
	// allow calls without a slash
	s.StrictSlash(true)
	// url rooting for retriving a product post
	s.HandleFunc("/products", PostProduct).Methods("POST")
	// url rooting for retriving a product list with get
	s.HandleFunc("/products", GetAllProducts).Methods("GET")
	// url rooting for one product with its ID
	s.HandleFunc("/products/{id:[0-9]+}", GetProduct).Methods("GET")
	// url rooting for deleting a product
	s.HandleFunc("/products/delete/{id:[0-9]+}", DeleteProduct)
	//
	// finally attach http handler to the router
	http.Handle("/", r)
	// log infos
	log.Println("Listening on port 8080 -> go to http://localhost:8080")
	// listen on port 8080
	http.ListenAndServe(":8080", nil)
}

func GetProduct(w http.ResponseWriter, r *http.Request) {
	log.Println("call : GetProduct")
	vars := mux.Vars(r)
	id := vars["id"]
	p := store[id]
	renderJson(w, p)
}

func GetAllProducts(w http.ResponseWriter, r *http.Request) {
	log.Println("call : GetAllProducts")

	msg := make([]Product, len(store))
	i := 0
	for _, product := range store {
		msg[i] = *product
		i++
	}
	renderJson(w, msg)
}
func DeleteProduct(w http.ResponseWriter, r *http.Request) {
	log.Println("call : DeleteProduct")
	vars := mux.Vars(r)
	id := vars["id"]
	p := store[id]
	delete(store, id)
	msg := "All good : " + p.Name + " has been deleted "
	renderJson(w, msg)
}

func PostProduct(w http.ResponseWriter, r *http.Request) {
	log.Println("call : PostProduct")
	p := Product{}
	p.Name = r.FormValue("name")
	var msg string
	ID++
	store[p.Name] = &Product{strconv.Itoa(ID), p.Name}
	msg = "All good : " + p.Name + " has been added "
	renderJson(w, msg)
}

// renderJson let you send some json response
func renderJson(w http.ResponseWriter, msg interface{}) {
	log.Println("call : renderJson")

	// CHOICE 1 : classy way to give a json response
	w.Header().Add("Accept-Charset", "utf-8")
	// encode in json the message
	log.Println(msg)
	b, err := json.Marshal(msg)
	log.Println(string(b))
	if err != nil {
		// if any error to encode send it to the client
		// you could also log it on the backend
		fmt.Fprint(w, "error")
	}
	fmt.Fprint(w, string(b))

	// // CHOICE 2 : gzip way that compress the result
	// w.Header().Add("Accept-Charset", "utf-8")
	// w.Header().Add("Content-Type", "application/json")
	// w.Header().Set("Content-Encoding", "gzip")

	// // Gzip data
	// gz := gzip.NewWriter(w)
	// json.NewEncoder(gz).Encode(msg)
	// gz.Close()
}
