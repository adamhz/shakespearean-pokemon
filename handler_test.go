package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	. "github.com/onsi/gomega"
)

const (
	testDescription          = "gotta catch em all"
	testConvertedDescription = "to be, or not to be"
)

type mockDescriptionGetter struct {
	pokemon string
	g       *GomegaWithT
}

func (m *mockDescriptionGetter) GetDescription(pokemon string) (string, error) {
	m.g.Expect(pokemon).To(Equal(m.pokemon))
	return testDescription, nil
}

type mockShakespeareConverter struct {
	g *GomegaWithT
}

func (m *mockShakespeareConverter) ConvertText(text string) (string, error) {
	m.g.Expect(text).To(Equal(testDescription))
	return testConvertedDescription, nil
}

func TestGetPokemonHandler(t *testing.T) {
	tests := []struct {
		name    string
		pokemon string
		exists  bool
	}{
		{
			name:    "base",
			pokemon: "potato",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := NewGomegaWithT(t)
			// initialise handler with mocked external calls
			handler := Handler{
				d: &mockDescriptionGetter{g: g, pokemon: tt.pokemon},
				s: &mockShakespeareConverter{g: g},
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

			var response GetPokemonResponse
			err := json.NewDecoder(resp.Body).Decode(&response)
			g.Expect(err).ToNot(HaveOccurred())

			// assertions
			g.Expect(resp.StatusCode).To(Equal(http.StatusOK))
			g.Expect(resp.Header.Get("Content-Type")).To(Equal("application/json"))
			g.Expect(response.Description).To(Equal(testConvertedDescription))
			g.Expect(response.Name).To(Equal(tt.pokemon))
		})
	}
}
