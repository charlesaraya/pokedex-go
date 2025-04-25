package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/charlesaraya/pokedex-go/internal"
)

func cleanInput(text string) []string {
	if len(text) == 0 {
		return []string{}
	}
	lowercased := strings.ToLower(text)
	fields := strings.Fields(lowercased)
	return fields
}

var duration, _ = time.ParseDuration("5s")
var PokeCache = internal.NewCache(duration)

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
				if err := data.Command(data.Config); err != nil {
					fmt.Printf("Error: %s command produced an error (%s)\n", command, err)
				}
				commandEntered = true
			}
		}
		if !commandEntered {
			fmt.Println("Unknown command")
		}
	}

}
