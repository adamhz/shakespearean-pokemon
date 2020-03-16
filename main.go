package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/adamhz/shakespearean-pokemon/service"
)

func main() {
	// configure server to serve on port 3000
	addr := fmt.Sprintf("0.0.0.0:%d", 3000)

	// set up route `/pokemon/<pokemon name>`
	h := service.NewHandler(&service.PokeAPIClt{}, &service.ShakespeareTranslatorClt{})

	router := mux.NewRouter()
	router.HandleFunc("/pokemon/{name}", h.HandleGetPokemon)

	// serve
	log.Printf("starting server on: %s", addr)
	err := http.ListenAndServe(addr, router)
	if err != nil {
		panic(err)
	}
}
