package command

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"sync"
	"time"

	"github.com/charlesaraya/pokedex-go/pokeapi"
	"github.com/charlesaraya/pokedex-go/saveload"
	"github.com/charlesaraya/pokedex-go/terminal"
)

const (
	CMD_MAP         string = "map"
	CMD_MAPB        string = "mapb"
	CMD_HELP        string = "help"
	CMD_EXIT        string = "exit"
	CMD_EXPLORE     string = "explore"
	CMD_CATCH       string = "catch"
	CMD_INSPECT     string = "inspect"
	CMD_POKEDEX     string = "pokedex"
	CMD_SAVE        string = "save"
	CMD_LOAD        string = "load"
	CMD_WHEREAMI    string = "whereami"
	FLAG_WHEREAMI_R string = "-r"
	FLAG_WHEREAMI_L string = "-l"
	CMD_VISIT       string = "visit"
	CMD_ENCOUNTER   string = "encounter"
)

type Config struct {
	Next     string
	Previous string
	Params   []string
}

type Flag struct {
	Name        string
	Description string
}

type Command struct {
	Name        string
	Description string
	Flags       []Flag
	Config      *Config
	Command     func(*Config, *Cache) error
}

type Cache struct {
	CachedEntries map[string]*CacheEntry
	Pokedex       *pokeapi.Pokedex
	Mu            sync.RWMutex
}

type CacheEntry struct {
	CreatedAt time.Time
	Val       []byte
}

func NewCache(interval time.Duration) *Cache {
	var cache *Cache = &Cache{
		CachedEntries: make(map[string]*CacheEntry),
	}
	go cache.reapLoop(interval)
	return cache
}

func (c *Cache) reapLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	// run every tick interval
	for range ticker.C {
		for k, v := range c.CachedEntries {
			// clean cached entries older than interval
			if time.Since(v.CreatedAt) > interval {
				delete(c.CachedEntries, k)
			}
		}
	}
}

func (c *Cache) Add(key string, val []byte) {
	c.Mu.Lock()
	c.CachedEntries[key] = &CacheEntry{
		CreatedAt: time.Now(),
		Val:       val,
	}
	defer c.Mu.Unlock()
}

func (c *Cache) Get(key string) (*CacheEntry, bool) {
	c.Mu.RLock()
	cachedEntry, ok := c.CachedEntries[key]
	defer c.Mu.RUnlock()
	return cachedEntry, ok
}

func GetRegistry() map[string]Command {
	mapConfig := Config{
		Next: pokeapi.ENDPOINT_LOCATION_AREA + pokeapi.PAGINATION,
	}
	return map[string]Command{
		CMD_ENCOUNTER: {
			Name:        "encounter",
			Description: "Triggers a random Pokémon encounter in the currently visited area.",
			Config: &Config{
				Next: pokeapi.ENDPOINT_LOCATION_AREA,
			},
			Command: commandEncounter,
		},
		CMD_VISIT: {
			Name:        "visit",
			Description: "Visits a location area.",
			Config: &Config{
				Next: pokeapi.ENDPOINT_LOCATION_AREA,
			},
			Command: commandVisit,
		},
		CMD_WHEREAMI: {
			Name:        "whereami",
			Description: "Shows the player's current location area.",
			Flags: []Flag{
				{
					Name:        "-r",
					Description: "Shows the current region instead.",
				},
				{
					Name:        "-l",
					Description: "Shows the current location instead.",
				},
			},
			Config: &Config{
				Params: []string{},
			},
			Command: commandWhereAmI,
		},
		CMD_LOAD: {
			Name:        "load",
			Description: "Loads previous game.",
			Config:      &Config{},
			Command:     commandLoad,
		},
		CMD_SAVE: {
			Name:        "save",
			Description: "Save current game.",
			Config:      &Config{},
			Command:     commandSave,
		},
		CMD_POKEDEX: {
			Name:        "pokedex",
			Description: "Show all Pokémon from the Pokedex.",
			Config:      &Config{},
			Command:     commandPokedex,
		},
		CMD_INSPECT: {
			Name:        "inspect",
			Description: "Inspect a Pokémon from the Pokedex.",
			Config:      &Config{},
			Command:     commandInspect,
		},
		CMD_CATCH: {
			Name:        "catch",
			Description: "Try catch a Pokémon.",
			Config: &Config{
				Next: pokeapi.ENDPOINT_POKEMON,
			},
			Command: commandCatch,
		},
		CMD_EXPLORE: {
			Name:        "explore",
			Description: "Shows the names of all the Pokémons located in an area in the Pokemon world.",
			Config: &Config{
				Next: pokeapi.ENDPOINT_LOCATION_AREA,
			},
			Command: commandExplore,
		},
		CMD_MAP: {
			Name:        "map",
			Description: "Shows the names of the next 20 location areas in the Pokemon world.",
			Config:      &mapConfig,
			Command:     commandMapForward,
		},
		CMD_MAPB: {
			Name:        "mapb",
			Description: "Shows the names of the previous 20 location areas in the Pokemon world.",
			Config:      &mapConfig,
			Command:     commandMapBack,
		},
		CMD_HELP: {
			Name:        "help",
			Description: "Shows the list of commands",
			Command:     commandHelp,
		},
		CMD_EXIT: {
			Name:        "exit",
			Description: "Exit the Pokedex CLI",
			Command:     commandExit,
		},
	}
}

func commandExit(config *Config, c *Cache) error {
	defer os.Exit(0)
	fmt.Println("Closing the Pokedex... Goodbye!")
	return nil
}

func commandHelp(config *Config, c *Cache) error {
	registry := GetRegistry()
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("\nusage: <command>")
	fmt.Printf("\nThese are common Pokedex commands used in various situations:\n\n")
	for _, data := range registry {
		fmt.Printf("    %s \t%s\n", data.Name, data.Description)
		for _, flag := range data.Flags {
			fmt.Printf("    \t\t%s \t%s\n", flag.Name, flag.Description)
		}
	}
	return nil
}

func commandMapForward(config *Config, c *Cache) error {
	if config.Next == "" {
		return fmt.Errorf("error: cant't map forward")
	}
	return Map(config, config.Next, CMD_MAP, c)
}

func commandMapBack(config *Config, c *Cache) error {
	if config.Previous == "" {
		return fmt.Errorf("error: cant't map back")
	}
	return Map(config, config.Previous, CMD_MAPB, c)
}

func Map(config *Config, url string, cmd string, c *Cache) error {
	var pokeLocationArea pokeapi.LocationAreas
	cachedEntry, ok := c.Get(cmd)
	if ok {
		if err := json.Unmarshal(cachedEntry.Val, &pokeLocationArea); err != nil {
			return fmt.Errorf("error: unmarshal operation failed from cached entry: %w", err)
		}
	} else {
		p, err := pokeapi.GetLocationAreas(url)
		if err != nil {
			return fmt.Errorf("error: failed getting location areas (%w)", err)
		}
		// cache data
		data, err := json.Marshal(p)
		if err != nil {
			return fmt.Errorf("error: marshal operation failed: %w", err)
		}
		c.Add(cmd, data)

		pokeLocationArea = p
		// update config's pagination
		config.Next = pokeLocationArea.Next
		config.Previous = pokeLocationArea.Previous
	}
	// Print results
	names := make([]string, len(pokeLocationArea.Results))
	for i, result := range pokeLocationArea.Results {
		names[i] = result.Name
	}
	terminal.PrettyPrint(names)
	return nil
}

func commandExplore(config *Config, c *Cache) error {
	var pokemons []pokeapi.Pokemon
	var locationAreaName string
	if len(config.Params) == 0 {
		locationAreaName = c.Pokedex.CurrentLocation.LocationArea
	} else {
		locationAreaName = config.Params[0]
	}
	fullCommand := CMD_EXPLORE + locationAreaName
	cachedEntry, ok := c.Get(fullCommand)
	if ok {
		if err := json.Unmarshal(cachedEntry.Val, &pokemons); err != nil {
			return fmt.Errorf("error: unmarshal operation failed from cached entry: %w", err)
		}
	} else {
		fullUrl := config.Next + locationAreaName
		p, err := pokeapi.GetPokemonsInLocationArea(fullUrl)
		if err != nil {
			return fmt.Errorf("error: failed getting pokemons in location area (%w)", err)
		}
		// cache data
		pokemonsJsonData, err := json.Marshal(p)
		if err != nil {
			return fmt.Errorf("error: marshal operation failed: %w", err)
		}
		c.Add(fullCommand, pokemonsJsonData)

		pokemons = p
	}
	// Print results
	names := make([]string, len(pokemons))
	for i, pokemon := range pokemons {
		names[i] = pokemon.Name
	}
	terminal.PrettyPrint(names)
	return nil
}

func commandCatch(config *Config, c *Cache) error {
	var pokemon pokeapi.Pokemon

	var pokemonName string
	if len(config.Params) == 0 {
		cachedEntry, ok := c.Get(CMD_ENCOUNTER)
		if ok {
			pokemonName = string(cachedEntry.Val)
		} else {
			fmt.Print("Nothing to catch!\n")
			return nil
		}
	} else {
		pokemonName = config.Params[0]
	}
	fullCommand := CMD_EXPLORE + pokemonName

	cachedEntry, ok := c.Get(fullCommand)
	if ok {
		if err := json.Unmarshal(cachedEntry.Val, &pokemon); err != nil {
			return fmt.Errorf("error: unmarshal operation failed from cached entry: %w", err)
		}
	} else {
		fullUrl := config.Next + pokemonName
		p, err := pokeapi.GetPokemon(fullUrl)
		if err != nil {
			return fmt.Errorf("error: failed getting pokemons in location area (%w)", err)
		}
		// cache data
		pokemonJsonData, err := json.Marshal(p)
		if err != nil {
			return fmt.Errorf("error: marshal operation failed: %w", err)
		}
		c.Add(fullCommand, pokemonJsonData)

		pokemon = p
	}
	fmt.Printf("Throwing a Pokeball at %s!", pokemon.Name)
	// We generate ellipsis every sec to add excitement
	duration, _ := time.ParseDuration("1s")
	ticker := time.NewTicker(duration)
	defer ticker.Stop()
	done := make(chan bool)
	tickerCount := 3
	go func() {
		for range ticker.C {
			fmt.Printf(".")
			<-ticker.C
			tickerCount--
			if tickerCount == 0 {
				done <- true
				return
			}
		}
	}()
	<-done
	fmt.Println()
	// TODO: Add a formula that generates a decent Capture Probability
	if rand.Float64() > 0.5 {
		fmt.Printf("%s was caught!\n", pokemon.Name)
		if _, ok := c.Pokedex.Get(pokemon.Name); !ok {
			c.Pokedex.Add(pokemon)
		}
	} else {
		fmt.Printf("%s escaped!\n", pokemon.Name)
	}
	return nil
}

func commandInspect(config *Config, c *Cache) error {
	pokedexEntry, ok := c.Pokedex.Get(config.Params[0])
	if !ok {
		fmt.Println("You have not caught that pokemon")
		return nil
	}
	fmt.Printf("Name: %s\n", pokedexEntry.Pokemon.Name)
	fmt.Printf("Height: %v\n", pokedexEntry.Pokemon.Height)
	fmt.Printf("Weight: %v\n", pokedexEntry.Pokemon.Weight)
	fmt.Printf("Stats:\n")
	for _, stat := range pokedexEntry.Pokemon.Stats {
		fmt.Printf("  -%s: %v\n", stat.Stat.Name, stat.Base)
	}
	fmt.Printf("Types:\n")
	for _, pokemonType := range pokedexEntry.Pokemon.Types {
		fmt.Printf("  - %s\n", pokemonType.Type.Name)
	}
	return nil
}

func commandPokedex(config *Config, c *Cache) error {
	pokemonNames := c.Pokedex.GetAll()
	if len(pokemonNames) == 0 {
		fmt.Println("your Pokedex is empty... Try catch some Pokémons first!")
	} else {
		fmt.Println("Your Pokedex:")
		terminal.PrettyPrint(pokemonNames)
	}
	return nil
}

func commandSave(config *Config, c *Cache) error {
	if err := saveload.Save(c.Pokedex, saveload.DATA_DIR); err != nil {
		return fmt.Errorf("error saving pokedex %w", err)
	}
	return nil
}

func commandLoad(config *Config, c *Cache) error {
	pokedex, err := saveload.Load(saveload.DATA_DIR)
	if err != nil {
		return fmt.Errorf("error loading game %w", err)
	}
	c.Pokedex = pokedex
	return nil
}

func commandWhereAmI(config *Config, c *Cache) error {
	locationArea := c.Pokedex.CurrentLocation.LocationArea
	if len(config.Params) > 0 {
		flag := config.Params[0]
		switch flag {
		case FLAG_WHEREAMI_R:
			fmt.Printf("%s\n", c.Pokedex.CurrentLocation.Region)
		case FLAG_WHEREAMI_L:
			fmt.Printf("%s\n", c.Pokedex.CurrentLocation.Location)
		default:
			fmt.Printf("%s\n", locationArea)
		}
		return nil
	}
	fmt.Printf("%s\n", locationArea)
	return nil
}

func commandVisit(config *Config, c *Cache) error {
	if len(config.Params) == 0 {
		return fmt.Errorf("received no argument")
	}
	var LocationArea pokeapi.LocationArea
	locaAreaName := config.Params[0]
	fullCommand := CMD_VISIT + " " + locaAreaName
	cachedEntry, ok := c.Get(fullCommand)
	if ok {
		if err := json.Unmarshal(cachedEntry.Val, &LocationArea); err != nil {
			return fmt.Errorf("failed to unmarshal location area: %w", err)
		}
		c.Pokedex.CurrentLocation.LocationArea = LocationArea.Name
		c.Pokedex.CurrentLocation.Location = LocationArea.Location.Name
	} else {
		fullUrl := config.Next + locaAreaName
		locationArea, err := pokeapi.GetLocationArea(fullUrl)
		if err != nil {
			return fmt.Errorf("failed to retrieve location area: %w", err)
		}
		data, err := json.Marshal(locationArea)
		if err != nil {
			return fmt.Errorf("failed to marshal location area: %w", err)
		}
		c.Add(fullCommand, data)
		c.Pokedex.CurrentLocation.LocationArea = locationArea.Name
		c.Pokedex.CurrentLocation.Location = locationArea.Location.Name
	}
	return nil
}

func commandEncounter(config *Config, c *Cache) error {
	fullEndpoint := config.Next + c.Pokedex.CurrentLocation.LocationArea
	pokemonEncounters, err := pokeapi.GetPokemonEncounters(fullEndpoint)
	if err != nil {
		return fmt.Errorf("failed to get pokemon encounters: %w", err)
	}
	// roulette wheel selection
	cumulativeWeights := 0
	for _, encounter := range pokemonEncounters {
		cumulativeWeights += encounter.Chance
	}
	pick := rand.Intn(cumulativeWeights)
	cumulativeWeights = 0
	for _, encounter := range pokemonEncounters {
		if pick >= cumulativeWeights && pick <= cumulativeWeights+encounter.Chance {
			fmt.Printf("You encountered a %s!\n", encounter.Name)
			// cache encounter
			c.Add(CMD_ENCOUNTER, []byte(encounter.Name))
			break
		}
		cumulativeWeights += encounter.Chance - 1
	}
	return nil
}
