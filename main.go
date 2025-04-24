package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
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

func cleanInput(text string) []string {
	if len(text) == 0 {
		return []string{}
	}
	lowercased := strings.ToLower(text)
	fields := strings.Fields(lowercased)
	return fields
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

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		// get input from user
		scanner.Scan()
		userInput := cleanInput(scanner.Text())

		// Check if the user entered a command
		registry := getRegistry()
		commandEntered := false
		for command, data := range registry {
			if command == userInput[0] {
				if err := data.Command(); err != nil {
					fmt.Printf("Error: %s command produced an error", command)
				}
				commandEntered = true
			}
		}
		if !commandEntered {
			fmt.Println("Unknown command")
		}
	}

}
