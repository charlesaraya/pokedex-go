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
	res, err := http.Get(endpoint)
	if err != nil {
		return PokeLocationArea{}, fmt.Errorf("error: failed getting response %w", err)
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return PokeLocationArea{}, fmt.Errorf("error: reading body from response: %w", err)
	}
	if res.StatusCode > 299 {
		return PokeLocationArea{}, fmt.Errorf("error: response failed with status code: %d", res.StatusCode)
	}
	var locationAreaObject PokeLocationArea
	if err := json.Unmarshal(body, &locationAreaObject); err != nil {
		return PokeLocationArea{}, fmt.Errorf("error: unmarshal operation failed: %w", err)
	}
	return locationAreaObject, nil
}
