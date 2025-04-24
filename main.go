package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func cleanInput(text string) []string {
	if len(text) == 0 {
		return []string{}
	}
	lowercased := strings.ToLower(text)
	fields := strings.Fields(lowercased)
	return fields
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
