package main

import (
	"bytes"
	"encoding/json"
	"testing"

	. "github.com/onsi/gomega"
)

func TestGetDescription(t *testing.T) {
	tests := []struct {
		name                string
		pokemonSpecies      PokemonSpecies
		shouldWork          bool
		expectedDescription string
	}{
		{
			name: "base",
			pokemonSpecies: PokemonSpecies{
				FlavorTextEntries: []FlavorTextEntry{
					{
						FlavorText: "pokemon",
						Language: Language{
							Name: "en",
						},
					},
				},
			},
			shouldWork:          true,
			expectedDescription: "pokemon",
		},
		{
			name: "unsanitized input",
			pokemonSpecies: PokemonSpecies{
				FlavorTextEntries: []FlavorTextEntry{
					{
						FlavorText: "pokemon\ncan \tfly",
						Language: Language{
							Name: "en",
						},
					},
				},
			},
			shouldWork:          true,
			expectedDescription: "pokemon can fly",
		},
		{
			name: "no english description",
			pokemonSpecies: PokemonSpecies{
				FlavorTextEntries: []FlavorTextEntry{
					{
						FlavorText: "pokemon",
						Language: Language{
							Name: "fr",
						},
					},
				},
			},
			shouldWork: false,
		},
		{
			name: "no flavor entries",
			pokemonSpecies: PokemonSpecies{
				FlavorTextEntries: []FlavorTextEntry{},
			},
			shouldWork: false,
		},
		{
			name: "no flavor entries",
			pokemonSpecies: PokemonSpecies{
				FlavorTextEntries: []FlavorTextEntry{},
			},
			shouldWork: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := NewGomegaWithT(t)

			body, err := json.Marshal(tt.pokemonSpecies)
			g.Expect(err).ToNot(HaveOccurred())

			res, err := getDescription(bytes.NewReader(body))
			if !tt.shouldWork {
				g.Expect(err).To(HaveOccurred())
			} else {
				g.Expect(res).To(Equal(tt.expectedDescription))
			}
		})
	}
}
