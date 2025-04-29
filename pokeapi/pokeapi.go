package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

const (
	ENDPOINT_POKEMON       string = "https://pokeapi.co/api/v2/pokemon/"
	ENDPOINT_LOCATION_AREA string = "https://pokeapi.co/api/v2/location-area/"
	PAGINATION             string = "?offset=0&limit=20"
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
	Name       string        `json:"name"`
	Height     int           `json:"height"`
	Weight     int           `json:"weight"`
	Experience int           `json:"base_experience"`
	Url        string        `json:"url"`
	Stats      []PokemonStat `json:"stats"`
	Types      []PokemonType `json:"types"`
}

type PokemonStat struct {
	Stat struct {
		Name string `json:"name"`
	} `json:"stat"`
	Base int `json:"base_stat"`
}

type PokemonType struct {
	Type struct {
		Name string `json:"name"`
	} `json:"type"`
}

type PokedexEntry struct {
	CatchedAt time.Time
	Pokemon   Pokemon
}

type Pokedex struct {
	PokedexEntries map[string]*PokedexEntry
	Mu             sync.RWMutex
}

func NewPokedex() *Pokedex {
	var pokedex *Pokedex = &Pokedex{
		PokedexEntries: make(map[string]*PokedexEntry),
	}
	return pokedex
}

func (p *Pokedex) Add(pokemon Pokemon) {
	p.Mu.Lock()
	p.PokedexEntries[pokemon.Name] = &PokedexEntry{
		CatchedAt: time.Now(),
		Pokemon:   pokemon,
	}
	p.Mu.Unlock()
}

func (p *Pokedex) Get(key string) (*PokedexEntry, bool) {
	p.Mu.RLock()
	pokedexEntry, ok := p.PokedexEntries[key]
	p.Mu.RUnlock()
	return pokedexEntry, ok
}

func (p *Pokedex) GetAll() []string {
	pokemonNames := []string{}
	p.Mu.RLock()
	for key := range p.PokedexEntries {
		pokemonNames = append(pokemonNames, key)
	}
	p.Mu.RUnlock()
	return pokemonNames
}

func GetLocationAreas(endpoint string) (PokeLocationArea, error) {
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

func GetPokemonsInLocationArea(endpoint string) ([]Pokemon, error) {
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

func GetPokemon(endpoint string) (Pokemon, error) {
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
