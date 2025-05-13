package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/charlesaraya/pokedex-go/internal/pokedex"
)

const (
	ENDPOINT_POKEMON       string = "https://pokeapi.co/api/v2/pokemon/"
	ENDPOINT_LOCATION_AREA string = "https://pokeapi.co/api/v2/location-area/"
	ENDPOINT_LOCATION      string = "https://pokeapi.co/api/v2/location/"
	PAGINATION             string = "?offset=0&limit=20"
)

type NamedResource struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type LocationArea struct {
	Name     string        `json:"name"`
	URL      string        `json:"url"`
	Location NamedResource `json:"location"`
}

type LocationAreas struct {
	Count    int             `json:"count"`
	Next     string          `json:"next"`
	Previous string          `json:"previous"`
	Results  []NamedResource `json:"results"`
}

type Location struct {
	Name   string          `json:"name"`
	URL    string          `json:"url"`
	Areas  []NamedResource `json:"areas"`
	Region NamedResource   `json:"region"`
}

type PokemonEncounter struct {
	Name   string
	Chance int
}

func GetLocationArea(endpoint string) (LocationArea, error) {
	locationArea := LocationArea{}

	res, err := http.Get(endpoint)
	if err != nil {
		return locationArea, fmt.Errorf("failed to get response %w", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return locationArea, fmt.Errorf("failed to read the response body: %w", err)
	}
	if res.StatusCode > 299 {
		return locationArea, fmt.Errorf("failed response with status code: %d", res.StatusCode)
	}
	if err := json.Unmarshal(body, &locationArea); err != nil {
		return locationArea, fmt.Errorf("failed to unmarshal location area: %w", err)
	}
	return locationArea, nil
}

func GetLocationAreas(endpoint string) (LocationAreas, error) {
	locationArea := LocationAreas{}

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

func GetPokemonEncounters(endpoint string) ([]PokemonEncounter, error) {
	var pokemonEncounters struct {
		Encounters []struct {
			Pokemon struct {
				Name string `json:"name"`
			} `json:"pokemon"`
			VersionDetails []struct {
				EncounterDetails []struct {
					MaxLevel int `json:"max_level"`
					MinLevel int `json:"min_level"`
				} `json:"encounter_details"`
				Chance int `json:"max_chance"`
			} `json:"version_details"`
		} `json:"pokemon_encounters"`
	}
	encounters := []PokemonEncounter{}

	res, err := http.Get(endpoint)
	if err != nil {
		return encounters, fmt.Errorf("error: failed getting response %w", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return encounters, fmt.Errorf("error: reading body from response: %w", err)
	}
	if res.StatusCode > 299 {
		return encounters, fmt.Errorf("error: response failed with status code: %d", res.StatusCode)
	}
	if err := json.Unmarshal(body, &pokemonEncounters); err != nil {
		return encounters, fmt.Errorf("error: unmarshal operation failed: %w", err)
	}
	for _, e := range pokemonEncounters.Encounters {
		encounter := PokemonEncounter{
			Name:   e.Pokemon.Name,
			Chance: e.VersionDetails[0].Chance,
		}
		encounters = append(encounters, encounter)
	}
	return encounters, nil
}

func GetPokemonsInLocationArea(endpoint string) ([]pokedex.Pokemon, error) {
	// define response struct
	var pokemonEncounters struct {
		Encounters []struct {
			Pokemon pokedex.Pokemon `json:"pokemon"`
		} `json:"pokemon_encounters"`
	}
	pokemons := []pokedex.Pokemon{}

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

func GetPokemon(endpoint string) (pokedex.Pokemon, error) {
	pokemon := pokedex.Pokemon{}

	res, err := http.Get(endpoint)
	if err != nil {
		return pokemon, fmt.Errorf("error: failed getting response %w", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return pokemon, fmt.Errorf("error: reading body from response: %w", err)
	}
	if res.StatusCode > 299 {
		return pokemon, fmt.Errorf("error: response failed with status code: %d", res.StatusCode)
	}
	if err := json.Unmarshal(body, &pokemon); err != nil {
		return pokemon, fmt.Errorf("error: unmarshal operation failed: %w", err)
	}
	return pokemon, nil
}
