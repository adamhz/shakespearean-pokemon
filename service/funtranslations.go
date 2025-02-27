package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"github.com/pkg/errors"
)

const (
	urlFuntranslations = "https://api.funtranslations.com/translate/shakespeare.json"
)

// TranslateResponse represents the response returned from the Shakespeare translator API.
// see: https://funtranslations.com/api/shakespeare
type TranslateResponse struct {
	Contents *Contents `json:"contents"`
}

// Contents represents the translation if succesful.
type Contents struct {
	Translated  string `json:"translated"`
	Text        string `json:"text"`
	Translation string `json:"translation"`
}

// ShakespeareTranslatorClt implements ShakespeareConverter
type ShakespeareTranslatorClt struct{}

// ConvertText converts the input text to Shakespearean style using the Fun Translations API.
// see: https://funtranslations.com/api/shakespeare
func (s *ShakespeareTranslatorClt) ConvertText(text string) (string, error) {
	log.Printf("fetch translation for text: %s", text)

	clt := http.Client{}
	req, err := http.NewRequest(http.MethodPost, urlFuntranslations, nil)
	if err != nil {
		return "", err
	}

	// set query param
	params := url.Values{}
	params.Add("text", text)
	req.URL.RawQuery = params.Encode()

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	res, err := clt.Do(req)
	if err != nil {
		return "", errors.Wrapf(err, "failed to make request")
	}

	if res.StatusCode >= 300 {
		bs, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return "", errors.Wrapf(err, fmt.Sprintf("error reading response body: %+v", res))
		}
		return "", errors.New(fmt.Sprintf("error translating using Fun Translation API: %s", bs))
	}

	var translateResponse TranslateResponse
	err = json.NewDecoder(res.Body).Decode(&translateResponse)
	if err != nil {
		log.Printf("response.Body: %s", res.Body)
		return "", errors.Wrapf(err, "failed to unmarshal response body")
	}
	if translateResponse.Contents == nil {
		return "", errors.New("failed to translate")
	}

	return translateResponse.Contents.Translated, nil
}
