package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type LocationArea struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type PokeLocationArea struct {
	Count    int            `json:"count"`
	Next     string         `json:"next"`
	Previous string         `json:"previous"`
	Results  []LocationArea `json:"results"`
}

func getLocationAreas(endpoint string) (PokeLocationArea, error) {
	locationArea := PokeLocationArea{}

	res, err := http.Get(endpoint)
	if err != nil {
		return locationArea, fmt.Errorf("error: failed getting response %w", err)
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return locationArea, fmt.Errorf("error: reading body from response: %w", err)
	}
	if res.StatusCode > 299 {
		return locationArea, fmt.Errorf("error: response failed with status code: %d", res.StatusCode)
	}
	if err := json.Unmarshal(body, &locationArea); err != nil {
		return locationArea, fmt.Errorf("error: unmarshal operation failed: %w", err)
	}
	return locationArea, nil
}
