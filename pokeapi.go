package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
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
	Name       string `json:"name"`
	Experience int    `json:"base_experience"`
	Url        string `json:"url"`
}

type PokedexEntry struct {
	CatchedAt time.Time
	Pokemon   Pokemon
}

type Pokedex struct {
	Pokedex map[string]*PokedexEntry
	Mu      sync.RWMutex
}

func NewPokedex() *Pokedex {
	var pokedex *Pokedex = &Pokedex{
		Pokedex: make(map[string]*PokedexEntry),
	}
	return pokedex
}

func (p *Pokedex) Add(pokemon Pokemon) {
	p.Mu.Lock()
	p.Pokedex[pokemon.Name] = &PokedexEntry{
		CatchedAt: time.Now(),
		Pokemon:   pokemon,
	}
	p.Mu.Unlock()
}

func (p *Pokedex) Get(key string) (*PokedexEntry, bool) {
	p.Mu.RLock()
	pokedexEntry, ok := p.Pokedex[key]
	p.Mu.RUnlock()
	return pokedexEntry, ok
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

func getPokemon(endpoint string) (Pokemon, error) {
	pokemon := Pokemon{}

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
