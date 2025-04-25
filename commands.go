package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	Next     string
	Previous string
}

type Command struct {
	Name        string
	Description string
	Config      *Config
	Command     func(*Config) error
}

func getRegistry() map[string]Command {
	return map[string]Command{
		"map": {
			Name:        "map",
			Description: "Shows the names of the next 20 location areas in the Pokemon world.",
			Config:      &mapConfig,
			Command:     commandMapForward,
		},
		"mapb": {
			Name:        "map back",
			Description: "Shows the names of the previous 20 location areas in the Pokemon world.",
			Config:      &mapConfig,
			Command:     commandMapBack,
		},
		"help": {
			Name:        "help",
			Description: "Shows the list of commands",
			Command:     commandHelp,
		},
		"exit": {
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
	return Map(config, config.Next, "map")
}

func commandMapBack(config *Config) error {
	if config.Previous == "" {
		return fmt.Errorf("error: cant't map back")
	}
	return Map(config, config.Previous, "mapb")
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
