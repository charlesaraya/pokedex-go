package main

import (
	"fmt"
	"os"
	"time"

	"github.com/charlesaraya/pokedex-go/cache"
	"github.com/charlesaraya/pokedex-go/terminal"
)

var duration, _ = time.ParseDuration("5s")
var PokeCache = cache.NewCache(duration)
var UserPokedex = NewPokedex()

func main() {
	err := terminal.EnableRawMode()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer terminal.DisableRawMode()

	commandHistory, historyIdx := terminal.InitCommandHistory()
	inputBuffer, cursor := terminal.InitBuffer()
	buf := make([]byte, 3)

	terminal.RedrawLine(inputBuffer, cursor)
	for {
		if _, err := os.Stdin.Read(buf); err != nil {
			fmt.Println("Error reading:", err)
			break
		}
		switch terminal.GetKey(buf) {
		case terminal.KEY_UP:
			if historyIdx > 0 {
				if historyIdx == len(commandHistory) && len(inputBuffer) > 0 {
					commandHistory = append(commandHistory, string(inputBuffer))
				} else if len(inputBuffer) > 0 {
					commandHistory[historyIdx] = string(inputBuffer)
				}
				historyIdx--
				terminal.UpdateBuffer(commandHistory[historyIdx], &inputBuffer, &cursor)
				terminal.RedrawLine(inputBuffer, cursor)
			}
		case terminal.KEY_DOWN:
			if historyIdx < len(commandHistory) {
				commandHistory[historyIdx] = string(inputBuffer)
				historyIdx++
				if historyIdx == len(commandHistory) {
					terminal.ResetBuffer(&inputBuffer, &cursor)
				} else {
					terminal.UpdateBuffer(commandHistory[historyIdx], &inputBuffer, &cursor)
				}
				terminal.RedrawLine(inputBuffer, cursor)
			}
		case terminal.KEY_RIGHT:
			terminal.MoveCursorRight(&cursor, len(inputBuffer))
		case terminal.KEY_LEFT:
			terminal.MoveCursorLeft(&cursor)
		case terminal.KEY_BACKSPACE:
			if ok := terminal.DeleteFromBuffer(&inputBuffer, &cursor); ok {
				terminal.RedrawLine(inputBuffer, cursor)
			}
		case terminal.KEY_PRINTABLE:
			terminal.AddToBuffer(rune(buf[0]), &inputBuffer, &cursor)
			terminal.RedrawLine(inputBuffer, cursor)
		case terminal.KEY_ENTER:
			if len(inputBuffer) > 0 {
				fmt.Println() // move to next line
				// Check if the user entered a valid command
				fullCommand := terminal.CleanInput(string(inputBuffer))
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
				terminal.AddCommand(string(inputBuffer), &commandHistory, &historyIdx)
				terminal.ResetBuffer(&inputBuffer, &cursor)
				terminal.RedrawLine(inputBuffer, cursor)
				continue
			}
			historyIdx = len(commandHistory)
			fmt.Println()
			terminal.RedrawLine(inputBuffer, cursor)
		}
	}

}
