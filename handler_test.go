package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"
)

const (
	testDescription          = "gotta catch em all"
	testConvertedDescription = "to be, or not to be"
)

type mockDescriptionGetter struct {
	pokemon string
	exists  bool
	g       *GomegaWithT
}

func (m *mockDescriptionGetter) GetDescription(pokemon string) (string, error) {
	m.g.Expect(pokemon).To(Equal(m.pokemon))
	if m.exists {
		return testDescription, nil
	}
	return "", errors.New("not found")
}

type mockShakespeareConverter struct {
	shouldConvert bool
	g             *GomegaWithT
}

func (m *mockShakespeareConverter) ConvertText(text string) (string, error) {
	m.g.Expect(text).To(Equal(testDescription))
	if m.shouldConvert {
		return testConvertedDescription, nil
	}
	return "", errors.New("failed to convert")
}

func TestGetPokemonHandler(t *testing.T) {
	tests := []struct {
		name          string
		pokemon       string
		exists        bool
		shouldConvert bool
	}{
		{
			name:    "base",
			pokemon: "charizard",
			exists:  true,
		},
		{
			name:    "missing pokemon",
			pokemon: "potato",
			exists:  false,
		},
		{
			name:          "unable to shakespearean-ise",
			pokemon:       "true",
			exists:        true,
			shouldConvert: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := NewGomegaWithT(t)
			// initialise handler with mocked external calls
			handler := Handler{
				d: &mockDescriptionGetter{g: g, pokemon: tt.pokemon, exists: tt.exists},
				s: &mockShakespeareConverter{g: g, shouldConvert: tt.shouldConvert},
			}

			// set up router
			url := fmt.Sprintf("http://host.com/pokemon/%s", tt.pokemon)
			req := httptest.NewRequest("GET", url, nil)
			router := mux.NewRouter()
			router.HandleFunc("/pokemon/{name}", handler.HandleGetPokemon)

			// make test call and record response
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			resp := w.Result()

			// assertions

			// if pokemon does not exists, check if we return 404
			if !tt.exists {
				g.Expect(resp.StatusCode).To(Equal(http.StatusNotFound))
				return
			}

			// if shakespeare translator fails, return 500
			if !tt.shouldConvert {
				g.Expect(resp.StatusCode).To(Equal(http.StatusInternalServerError))
				return
			}

			var response GetPokemonResponse
			err := json.NewDecoder(resp.Body).Decode(&response)
			g.Expect(err).ToNot(HaveOccurred())

			g.Expect(resp.StatusCode).To(Equal(http.StatusOK))
			g.Expect(resp.Header.Get("Content-Type")).To(Equal("application/json"))
			g.Expect(response.Description).To(Equal(testConvertedDescription))
			g.Expect(response.Name).To(Equal(tt.pokemon))

		})
	}
}
