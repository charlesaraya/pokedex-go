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

	commandHistory := []string{}
	historyIdx := 0

	buf := make([]byte, 3) // Arrow keys send 3 bytes
	inputBuffer := []rune{}
	promptName := "Pokedex > "
	cursor := 0
	redrawLine(inputBuffer, cursor, promptName)
	for {
		bufLen, err := os.Stdin.Read(buf)
		if err != nil {
			fmt.Println("Error reading:", err)
			break
		}
		// Arrows Keys
		if bufLen == 3 && buf[0] == 0x1b && buf[1] == '[' {
			switch buf[2] {
			// Up arrow
			case 'A':
				if historyIdx > 0 {
					if historyIdx == len(commandHistory) && len(inputBuffer) > 0 {
						commandHistory = append(commandHistory, string(inputBuffer))
					} else if len(inputBuffer) > 0 {
						commandHistory[historyIdx] = string(inputBuffer)
					}
					historyIdx--
					inputBuffer = []rune(commandHistory[historyIdx])
					cursor = len(inputBuffer)
					redrawLine(inputBuffer, cursor, promptName)
				}
			// Down arrow
			case 'B':
				if historyIdx < len(commandHistory) {
					commandHistory[historyIdx] = string(inputBuffer)
					historyIdx++
					if historyIdx == len(commandHistory) {
						// reset buffer
						inputBuffer = inputBuffer[:0]
						cursor = 0
					} else {
						inputBuffer = []rune(commandHistory[historyIdx])
						cursor = len(inputBuffer)
					}
					redrawLine(inputBuffer, cursor, promptName)
				}
			// Right arrow
			case 'C':
				// Move the cursor forward
				if cursor < len(inputBuffer) {
					cursor++
					os.Stdout.Write([]byte("\033[C"))
				}
			// Left arrow
			case 'D':
				// Move cursor back
				if cursor > 0 {
					cursor--
					os.Stdout.Write([]byte("\b"))
				}
			}
			continue
		}
		// Backspace key
		if buf[0] == 127 {
			if cursor > 0 {
				copy(inputBuffer[cursor-1:], inputBuffer[cursor:]) // Shift everything left
				inputBuffer = inputBuffer[:len(inputBuffer)-1]     // decrease the buffer size by one
				cursor--
				redrawLine(inputBuffer, cursor, promptName)
			}
			continue
		}
		// Printable characters
		if buf[0] >= 32 && buf[0] <= 126 {
			r := rune(buf[0])
			inputBuffer = append(inputBuffer, 0)               // grow buffer by 1
			copy(inputBuffer[cursor+1:], inputBuffer[cursor:]) // Shift everything right
			inputBuffer[cursor] = r                            // Insert the new rune
			cursor++
			redrawLine(inputBuffer, cursor, promptName)
		}
		// Enter key
		if buf[0] == '\n' {
			switch {
			case len(inputBuffer) > 0:
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
				// Add command to history
				if historyIdx == len(commandHistory) {
					commandHistory = append(commandHistory, string(inputBuffer))
					historyIdx++
				} else {
					copy(commandHistory[historyIdx:], commandHistory[historyIdx+1:]) // Shift everything left
					commandHistory = commandHistory[:len(commandHistory)-1]          // decrease the buffer size by one
					commandHistory = append(commandHistory, string(inputBuffer))
					historyIdx = len(commandHistory)
				}
				// reset buffer
				inputBuffer = inputBuffer[:0]
				cursor = 0
				redrawLine(inputBuffer, cursor, promptName)
				continue
			default:
				historyIdx = len(commandHistory)
				fmt.Println()
				redrawLine(inputBuffer, cursor, promptName)
			}
		}
	}

}
