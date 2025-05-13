package main

import (
	"fmt"
	"os"
	"time"

	"github.com/charlesaraya/pokedex-go/internal/cache"
	"github.com/charlesaraya/pokedex-go/internal/commands"
	"github.com/charlesaraya/pokedex-go/internal/pokedex"
	"github.com/charlesaraya/pokedex-go/internal/terminal"
)

func main() {
	var duration, _ = time.ParseDuration("5s")
	var cache = cache.NewCache(duration)
	cache.Pokedex = pokedex.NewPokedex()

	err := terminal.EnableRawMode()
	if err != nil {
		fmt.Println("Error: failed to terminal enable raw mode", err)
		return
	}
	defer terminal.DisableRawMode()

	commandHistory, historyIdx := terminal.InitCommandHistory()
	inputBuffer, cursor := terminal.InitBuffer()
	buf := make([]byte, 3)
	registry := commands.GetRegistry()

	for {
		terminal.RedrawLine(inputBuffer, cursor)
		if _, err := os.Stdin.Read(buf); err != nil {
			fmt.Println("Error: failed reading from input buffer:", err)
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
			}
		case terminal.KEY_RIGHT:
			terminal.MoveCursorRight(&cursor, len(inputBuffer))
		case terminal.KEY_LEFT:
			terminal.MoveCursorLeft(&cursor)
		case terminal.KEY_BACKSPACE:
			terminal.DeleteFromBuffer(&inputBuffer, &cursor)
		case terminal.KEY_PRINTABLE:
			terminal.AddToBuffer(rune(buf[0]), &inputBuffer, &cursor)
		case terminal.KEY_ENTER:
			if len(inputBuffer) > 0 {
				fmt.Println() // move to next line
				// Check if the user entered a valid command
				fullCommand := terminal.CleanInput(string(inputBuffer))
				Cmd, ok := registry[fullCommand[0]]
				if !ok {
					fmt.Printf("Error: unknown command %q\n", fullCommand[0])
				} else {
					if len(fullCommand) > 1 {
						Cmd.Config.Params = fullCommand[1:]
					}
					if err := Cmd.Command(Cmd.Config, cache); err != nil {
						fmt.Printf("Error: %s command produced an error: %s\n", Cmd.Name, err)
					}
					terminal.AddCommand(string(inputBuffer), &commandHistory, &historyIdx)
				}
			}
			historyIdx = len(commandHistory)
			terminal.ResetBuffer(&inputBuffer, &cursor)
		}
	}
}
