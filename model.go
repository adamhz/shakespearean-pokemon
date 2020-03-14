package main

// GetPokemonResponse is returned by the our API
type GetPokemonResponse struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}
