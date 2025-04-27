package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"time"
)

const (
	CMD_MAP     string = "map"
	CMD_MAPB    string = "mapb"
	CMD_HELP    string = "help"
	CMD_EXIT    string = "exit"
	CMD_EXPLORE string = "explore"
	CMD_CATCH   string = "catch"
	CMD_INSPECT string = "inspect"
	CMD_POKEDEX string = "pokedex"
)

type Config struct {
	Next     string
	Previous string
	Params   []string
}

type Command struct {
	Name        string
	Description string
	Config      *Config
	Command     func(*Config) error
}

func getRegistry() map[string]Command {
	return map[string]Command{
		CMD_POKEDEX: {
			Name:        "Pokedex",
			Description: "Show all Pokémon from the Pokedex.",
			Config:      &Config{},
			Command:     commandPokedex,
		},
		CMD_INSPECT: {
			Name:        "Inspect",
			Description: "Inspect a Pokémon from the Pokedex.",
			Config:      &Config{},
			Command:     commandInspect,
		},
		CMD_CATCH: {
			Name:        "Catch",
			Description: "Try catch a Pokémon.",
			Config: &Config{
				Next: "https://pokeapi.co/api/v2/pokemon/",
			},
			Command: commandCatch,
		},
		CMD_EXPLORE: {
			Name:        "Explore",
			Description: "Shows the names of all the Pokémons located in an area in the Pokemon world.",
			Config: &Config{
				Next: "https://pokeapi.co/api/v2/location-area/",
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
			Name:        "map back",
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

func commandExit(config *Config) error {
	defer os.Exit(0)
	fmt.Println("Closing the Pokedex... Goodbye!")
	return nil
}

func commandHelp(config *Config) error {
	registry := getRegistry()
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("\nusage: <command>")
	fmt.Printf("\nThese are common Pokedex commands used in various situations:\n\n")
	for _, data := range registry {
		fmt.Printf("    %s \t%s\n", data.Name, data.Description)
	}
	return nil
}

var mapConfig = Config{
	Next:     "https://pokeapi.co/api/v2/location-area/?offset=0&limit=20",
	Previous: "",
}

func commandMapForward(config *Config) error {
	if config.Next == "" {
		return fmt.Errorf("error: cant't map forward")
	}
	return Map(config, config.Next, CMD_MAP)
}

func commandMapBack(config *Config) error {
	if config.Previous == "" {
		return fmt.Errorf("error: cant't map back")
	}
	return Map(config, config.Previous, CMD_MAPB)
}

func Map(config *Config, url string, cmd string) error {
	var pokeLocationArea PokeLocationArea
	cachedEntry, ok := PokeCache.Get(cmd)
	if ok {
		if err := json.Unmarshal(cachedEntry.Val, &pokeLocationArea); err != nil {
			return fmt.Errorf("error: unmarshal operation failed from cached entry: %w", err)
		}
	} else {
		p, err := getLocationAreas(url)
		if err != nil {
			return fmt.Errorf("error: failed getting location areas (%w)", err)
		}
		// cache data
		data, err := json.Marshal(p)
		if err != nil {
			return fmt.Errorf("error: marshal operation failed: %w", err)
		}
		PokeCache.Add(cmd, data)

		pokeLocationArea = p
		// update config's pagination
		config.Next = pokeLocationArea.Next
		config.Previous = pokeLocationArea.Previous
	}
	// Print results
	for _, result := range pokeLocationArea.Results {
		fmt.Println(result.Name)
	}
	return nil
}

func commandExplore(config *Config) error {
	var pokemons []Pokemon

	fullCommand := CMD_EXPLORE + config.Params[0]
	cachedEntry, ok := PokeCache.Get(fullCommand)
	if ok {
		if err := json.Unmarshal(cachedEntry.Val, &pokemons); err != nil {
			return fmt.Errorf("error: unmarshal operation failed from cached entry: %w", err)
		}
	} else {
		fullUrl := config.Next + config.Params[0]
		p, err := getPokemonsInLocationArea(fullUrl)
		if err != nil {
			return fmt.Errorf("error: failed getting pokemons in location area (%w)", err)
		}
		// cache data
		pokemonsJsonData, err := json.Marshal(p)
		if err != nil {
			return fmt.Errorf("error: marshal operation failed: %w", err)
		}
		PokeCache.Add(fullCommand, pokemonsJsonData)

		pokemons = p
	}
	// Print results
	for _, pokemon := range pokemons {
		fmt.Println(pokemon.Name)
	}
	return nil
}

func commandCatch(config *Config) error {
	var pokemon Pokemon

	fullCommand := CMD_EXPLORE + config.Params[0]
	cachedEntry, ok := PokeCache.Get(fullCommand)
	if ok {
		if err := json.Unmarshal(cachedEntry.Val, &pokemon); err != nil {
			return fmt.Errorf("error: unmarshal operation failed from cached entry: %w", err)
		}
	} else {
		fullUrl := config.Next + config.Params[0]
		p, err := getPokemon(fullUrl)
		if err != nil {
			return fmt.Errorf("error: failed getting pokemons in location area (%w)", err)
		}
		// cache data
		pokemonJsonData, err := json.Marshal(p)
		if err != nil {
			return fmt.Errorf("error: marshal operation failed: %w", err)
		}
		PokeCache.Add(fullCommand, pokemonJsonData)

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
		if _, ok := UserPokedex.Get(pokemon.Name); !ok {
			UserPokedex.Add(pokemon)
		}
	} else {
		fmt.Printf("%s escaped!\n", pokemon.Name)
	}
	return nil
}

func commandInspect(config *Config) error {
	pokedexEntry, ok := UserPokedex.Get(config.Params[0])
	if !ok {
		fmt.Println("you have not caught that pokemon")
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

func commandPokedex(config *Config) error {
	pokemonNames := UserPokedex.GetAll()
	if len(pokemonNames) == 0 {
		fmt.Println("your Pokedex is empty... Try catch some Pokémons first!")
	} else {
		fmt.Println("Your Pokedex:")
		for _, name := range pokemonNames {
			fmt.Printf("  - %s\n", name)
		}
	}
	return nil
}
