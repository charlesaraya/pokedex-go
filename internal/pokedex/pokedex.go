package pokedex

import (
	"sync"
	"time"
)

const (
	STARTING_REGION        string = "kanto"
	STARTING_LOCATION      string = "pallet-town"
	STARTING_LOCATION_AREA string = "pallet-town-area"
)

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

type PlayerLocation struct {
	Region       string
	Location     string
	LocationArea string
}

type Pokedex struct {
	PokedexEntries  map[string]*PokedexEntry
	CurrentLocation PlayerLocation
	Mu              sync.RWMutex
}

func NewPokedex() *Pokedex {
	var pokedex *Pokedex = &Pokedex{
		PokedexEntries: make(map[string]*PokedexEntry),
		CurrentLocation: PlayerLocation{
			Region:       STARTING_REGION,
			Location:     STARTING_LOCATION,
			LocationArea: STARTING_LOCATION_AREA,
		},
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
