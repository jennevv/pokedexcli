package pokeapi

import (
	"encoding/json"
	"io"
	"net/http"
	"time"
)

type Response struct {
	Count    int      `json:"count"`
	Next     string   `json:"next"`
	Previous string   `json:"previous"`
	Results  []Result `json:"results"`
}

type Result struct {
	Name string `json:"name"`
}

type PokeClient struct {
	httpClient *http.Client
}

func NewClient() *PokeClient {
	return &PokeClient{
		&http.Client{Timeout: 10 * time.Second},
	}
}

func (p *PokeClient) Get(url string) ([]byte, error) {
	res, err := p.httpClient.Get(url)
	if err != nil {
		return []byte{}, err
	}

	body, err := io.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		return []byte{}, err
	}

	return body, nil
}

func (p *PokeClient) Unmarshal(body []byte) (Response, error) {
	var response Response

	err := json.Unmarshal(body, &response)
	if err != nil {
		return Response{}, err
	}

	return response, nil
}
