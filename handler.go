package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

// DescriptionGetter is our type for getting pokemon description
type DescriptionGetter interface {
	GetDescription(pokemon string) (string, error)
}

// ShakespeareConverter is our type for convert text to shakespearean style
type ShakespeareConverter interface {
	ConvertText(text string) (string, error)
}

// PokeAPIClt implements DescriptionGetter
type PokeAPIClt struct{}

// GetDescription fetches the description from the PokeAPI
// see: https://pokeapi.co/docs/v2.html/#pokemon-species
func (p *PokeAPIClt) GetDescription(pokemon string) (string, error) {
	return "foo", nil
}

// ShakespeareTranslatorClt implements ShakespeareConverter
type ShakespeareTranslatorClt struct{}

// ConvertText converts the input text to Shakespearean style using the Fun Translations API.
// see: https://funtranslations.com/api/shakespeare
func (s *ShakespeareTranslatorClt) ConvertText(text string) (string, error) {
	return "bar", nil
}

// Handler handles http requests
type Handler struct {
	d DescriptionGetter
	s ShakespeareConverter
}

// HandleGetPokemon handles requests to `/pokemon/<pokemon name>`.
// Fetches pokemon description and converts it Shakespearean style.
func (h *Handler) HandleGetPokemon(w http.ResponseWriter, r *http.Request) {
	pathParams := mux.Vars(r)
	name := pathParams["name"]
	log.Printf("GetPokemon: %s", name)

	text, err := h.d.GetDescription(name)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		log.Printf("error fetching description: %s", err)
		return
	}

	description, err := h.s.ConvertText(text)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("error converting text: %s", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	res := &GetPokemonResponse{
		Name:        name,
		Description: description,
	}
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("error marshalling response: %+v", err)
	}
}
