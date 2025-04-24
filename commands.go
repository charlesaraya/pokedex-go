package main

import (
	"fmt"
	"os"
)

type Command struct {
	Name        string
	Description string
	Command     func() error
}

func getRegistry() map[string]Command {
	return map[string]Command{
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

func commandExit() error {
	defer os.Exit(0)
	fmt.Println("Closing the Pokedex... Goodbye!")
	return nil
}

func commandHelp() error {
	registry := getRegistry()
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("\nusage: <command>")
	fmt.Printf("\nThese are common Pokedex commands used in various situations:\n\n")
	for _, data := range registry {
		fmt.Printf("    %s\t%s\n", data.Name, data.Description)
	}
	return nil
}
