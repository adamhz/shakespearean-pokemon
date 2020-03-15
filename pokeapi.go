package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/pkg/errors"
)

const (
	urlPokeAPISpecies = "https://pokeapi.co/api/v2/pokemon-species/%s"
)

// Language represents the language used to describe the pokemon, e.g. "en".
type Language struct {
	Name string `json:"name"`
}

// FlavorTextEntry is a descritpion of the pokemon. There can be many.
type FlavorTextEntry struct {
	FlavorText string   `json:"flavor_text"`
	Language   Language `json:"language"`
}

// PokemonSpecies represents the pokemon species data model returned by the PokeAPI
// see: https://pokeapi.co/docs/v2.html/#pokemon-species
// note: we only specify the fields we need
type PokemonSpecies struct {
	FlavorTextEntries []FlavorTextEntry `json:"flavor_text_entries"`
}

// PokeAPIClt implements DescriptionGetter
type PokeAPIClt struct{}

// GetDescription fetches the description from the PokeAPI
func (p *PokeAPIClt) GetDescription(pokemon string) (string, error) {
	clt := http.Client{}
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf(urlPokeAPISpecies, pokemon), nil)
	if err != nil {
		return "", err
	}
	res, err := clt.Do(req)
	if err != nil {
		return "", errors.Wrapf(err, "failed to make request to PokeAPI")
	}

	if res.StatusCode == http.StatusNotFound {
		return "", errors.New(fmt.Sprintf("unable to find pokemon: %s", pokemon))
	}
	if res.StatusCode >= 300 {
		return "", errors.New("error finding pokemon")
	}
	return getDescription(res.Body)
}

func getDescription(body io.Reader) (string, error) {
	var pokemonSpecies PokemonSpecies
	err := json.NewDecoder(body).Decode(&pokemonSpecies)
	if err != nil {
		return "", err
	}

	for _, flavorTextEntry := range pokemonSpecies.FlavorTextEntries {
		if flavorTextEntry.Language.Name == "en" {

			text := strings.Join(strings.Fields(strings.TrimSpace(flavorTextEntry.FlavorText)), " ")
			return text, nil
		}
	}
	return "", errors.New("no english description")
}
