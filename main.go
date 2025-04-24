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

var Registry = map[string]Command{
	"exit": {
		Name:        "exit",
		Description: "Exit the Pokedex CLI",
		Command:     commandExit,
	},
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

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		// get input from user
		scanner.Scan()
		userInput := cleanInput(scanner.Text())

		// Check if the user entered a command
		for command, data := range Registry {
			if command == userInput[0] {
				if err := data.Command(); err != nil {
					fmt.Printf("Error: %s command produced an error", command)
				}
			} else {
				fmt.Println("Unknown command")
			}
		}
	}

}
