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

type Pokemon struct {
	Name string `json:"name"`
	Url  string `json:"url"`
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

func getPokemonsInLocationArea(endpoint string) ([]Pokemon, error) {
	// define response struct
	var pokemonEncounters struct {
		Encounters []struct {
			Pokemon Pokemon `json:"pokemon"`
		} `json:"pokemon_encounters"`
	}
	pokemons := []Pokemon{}

	res, err := http.Get(endpoint)
	if err != nil {
		return pokemons, fmt.Errorf("error: failed getting response %w", err)
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return pokemons, fmt.Errorf("error: reading body from response: %w", err)
	}
	if res.StatusCode > 299 {
		return pokemons, fmt.Errorf("error: response failed with status code: %d", res.StatusCode)
	}
	if err := json.Unmarshal(body, &pokemonEncounters); err != nil {
		return pokemons, fmt.Errorf("error: unmarshal operation failed: %w", err)
	}
	for _, pokemon := range pokemonEncounters.Encounters {
		pokemons = append(pokemons, pokemon.Pokemon)
	}
	return pokemons, nil
}
