package main

import (
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
	config      *Config
	Command     func(*Config) error
}

func getRegistry() map[string]Command {
	return map[string]Command{
		"map": {
			Name:        "map",
			Description: "Shows the names of 20 location areas in the Pokemon world.",
			config:      &mapConfig,
			Command:     commandMap,
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
		fmt.Printf("    %s\t%s\n", data.Name, data.Description)
	}
	return nil
}

var mapConfig = Config{
	Next:     "https://pokeapi.co/api/v2/location-area/?offset=0&limit=20",
	Previous: "",
}

func commandMap(config *Config) error {
	pokeLocationArea, err := getLocationAreas(config.Next)
	if err != nil {
		return err
	}
	// update config's pagination
	config.Next = pokeLocationArea.Next
	config.Previous = pokeLocationArea.Previous
	// Print results
	for _, result := range pokeLocationArea.Results {
		fmt.Println(result.Name)
	}
	return nil
}
