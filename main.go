package main

import (
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
var UserPokedex = NewPokedex()

func main() {
	err := enableRawMode()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer disableRawMode()

	commandHistory, historyIdx := initCommandHistory()
	inputBuffer, cursor := initBuffer()
	buf := make([]byte, 3)

	redrawLine(inputBuffer, cursor)
	for {
		if _, err := os.Stdin.Read(buf); err != nil {
			fmt.Println("Error reading:", err)
			break
		}
		switch getKey(buf) {
		case KEY_UP:
			if historyIdx > 0 {
				if historyIdx == len(commandHistory) && len(inputBuffer) > 0 {
					commandHistory = append(commandHistory, string(inputBuffer))
				} else if len(inputBuffer) > 0 {
					commandHistory[historyIdx] = string(inputBuffer)
				}
				historyIdx--
				updateBuffer(commandHistory[historyIdx], &inputBuffer, &cursor)
				redrawLine(inputBuffer, cursor)
			}
		case KEY_DOWN:
			if historyIdx < len(commandHistory) {
				commandHistory[historyIdx] = string(inputBuffer)
				historyIdx++
				if historyIdx == len(commandHistory) {
					resetBuffer(&inputBuffer, &cursor)
				} else {
					updateBuffer(commandHistory[historyIdx], &inputBuffer, &cursor)
				}
				redrawLine(inputBuffer, cursor)
			}
		case KEY_RIGHT:
			moveCursorRight(&cursor, len(inputBuffer))
		case KEY_LEFT:
			moveCursorLeft(&cursor)
		case KEY_BACKSPACE:
			if ok := deleteFromBuffer(&inputBuffer, &cursor); ok {
				redrawLine(inputBuffer, cursor)
			}
		case KEY_PRINTABLE:
			addToBuffer(rune(buf[0]), &inputBuffer, &cursor)
			redrawLine(inputBuffer, cursor)
		case KEY_ENTER:
			if len(inputBuffer) > 0 {
				fmt.Println() // move to next line
				// Check if the user entered a valid command
				fullCommand := cleanInput(string(inputBuffer))
				registry := getRegistry()
				isKnownCommand := false
				for command, data := range registry {
					if command == fullCommand[0] {
						if len(fullCommand) > 1 {
							data.Config.Params = fullCommand[1:]
						}
						if err := data.Command(data.Config); err != nil {
							fmt.Printf("Error: %s command produced an error (%s)\n", command, err)
						}
						isKnownCommand = true
					}
				}
				if !isKnownCommand {
					fmt.Println("Unknown command")
				}
				addCommand(string(inputBuffer), &commandHistory, &historyIdx)
				resetBuffer(&inputBuffer, &cursor)
				redrawLine(inputBuffer, cursor)
				continue
			}
			historyIdx = len(commandHistory)
			fmt.Println()
			redrawLine(inputBuffer, cursor)
		}
	}

}
