package quotes

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

// seed is used to create a pseudo-random number
var seed = rand.NewSource(time.Now().UnixNano())

// API holds all the information to connect to an external api as a source for quotes.
type API struct {
	url    string
	key    string
	client http.Client
}

// fullQuoteResponse is the full json response from the api. Only used to take the key information keyed as 'docs'.
type fullQuoteResponse struct {
	Data []Quote `json:"docs"`
}

// fullCharacterResponse is the full json response from the api. Only used to take the key information keyed as 'docs'.
type fullCharacterResponse struct {
	Data []Character `json:"docs"`
}

// Quote is a single quote with character id.
type Quote struct {
	Text          string `json:"dialog"`
	CharacterID   string `json:"character"`
	CharacterName string `json:"name"`
}

// Character is a single character with character id.
type Character struct {
	CharacterID string `json:"_id"`
	Name        string `json:"name"`
}

// NewAPIClient is used to build a new http client to connect to the api.
func NewAPIClient(url, apiKey string) API {
	client := http.Client{
		Timeout: 5 * time.Second}

	api := API{
		url:    url,
		key:    apiKey,
		client: client,
	}

	return api
}

// NewQuote returns a new randomly selected quote from the api
func (a API) NewQuote() (Quote, error) {
	quote := Quote{}
	quotes, err := a.GetQuotes()
	if err != nil {
		return quote, err
	}
	characters, err := a.GetCharacters()
	if err != nil {
		return quote, err
	}

	rand.New(seed)
	index := rand.Intn(len(quotes))

	quote = quotes[index]

	quote.CharacterName, err = matchCharacter(quotes[index], characters)
	if err != nil {
		return quote, err
	}

	return quote, nil
}

// GetQuotes grabs all quotes from the lord of the rings api.
func (a API) GetQuotes() ([]Quote, error) {
	quotes := []Quote{}

	res, err := a.newRequest("/quote")
	if err != nil {
		return quotes, err
	}

	defer res.Body.Close()

	var jsonResp fullQuoteResponse

	err = json.NewDecoder(res.Body).Decode(&jsonResp)
	if err != nil {
		return quotes, err
	}

	quotes = jsonResp.Data

	return quotes, nil
}

// GetCharacters grabs all the character names and ids and returns a map between them.
func (a API) GetCharacters() (map[string]string, error) {
	characterMap := make(map[string]string)

	res, err := a.newRequest("/character")
	if err != nil {
		return characterMap, err
	}

	defer res.Body.Close()

	var jsonResp fullCharacterResponse

	err = json.NewDecoder(res.Body).Decode(&jsonResp)
	if err != nil {
		return characterMap, err
	}

	for _, character := range jsonResp.Data {
		characterMap[character.CharacterID] = character.Name
	}

	return characterMap, nil
}

func (a API) newRequest(endpoint string) (*http.Response, error) {
	res := &http.Response{}

	fullURL := fmt.Sprintf("%s%s", a.url, endpoint)

	req, err := http.NewRequest(http.MethodGet, fullURL, nil)
	if err != nil {
		return res, err
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", a.key))

	res, err = a.client.Do(req)
	if err != nil {
		return res, err
	}

	return res, nil
}

func matchCharacter(quote Quote, characterMap map[string]string) (string, error) {
	name, found := characterMap[quote.CharacterID]
	if !found {
		return "", errors.New("could not find character name")
	}

	return name, nil
}
