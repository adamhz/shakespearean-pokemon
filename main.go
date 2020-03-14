package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// GetPokemonResponse is returned by the our API
type GetPokemonResponse struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func HandleGetPokemon(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	res := &GetPokemonResponse{
		Name:        "foo",
		Description: "bar",
	}
	err := json.NewEncoder(w).Encode(res)
	if err != nil {
		log.Printf("error marshalling response: %+v", err)
	}
}

func main() {
	// configure server to serve on port 3000
	addr := fmt.Sprintf("0.0.0.0:%d", 3000)

	// set up route /pokemon/<pokemon name>
	router := mux.NewRouter()
	router.HandleFunc("/pokemon/{name}", HandleGetPokemon)

	// serve
	log.Printf("starting server on: %s", addr)
	err := http.ListenAndServe(addr, router)
	if err != nil {
		panic(err)
	}
}
