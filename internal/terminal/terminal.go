package terminal

import (
	"fmt"
	"os"
	"strings"
	"syscall"
	"text/tabwriter"
	"unsafe"
)

var oldState syscall.Termios

const (
	KEY_UP        string = "key_up"
	KEY_DOWN      string = "key_down"
	KEY_LEFT      string = "key_left"
	KEY_RIGHT     string = "key_right"
	KEY_ENTER     string = "key_enter"
	KEY_BACKSPACE string = "key_backspace"
	KEY_PRINTABLE string = "key_printable"
	KEY_UNKNOWN   string = "key_unknown"
)

const (
	PROMPT      string = "Pokedex > "
	MAX_COLS    int    = 3
	MAX_ROWS    int    = 5
	MAX_COL_PAD int    = 4
)

// Preprocesses a string and returns lowercased words
func CleanInput(text string) []string {
	if len(text) == 0 {
		return []string{}
	}
	lowercased := strings.ToLower(text)
	fields := strings.Fields(lowercased)
	return fields
}

func EnableRawMode() error {

	fd := int(os.Stdin.Fd())

	// Get current terminal settings
	_, _, err := syscall.Syscall6(syscall.SYS_IOCTL, uintptr(fd), uintptr(syscall.TIOCGETA), uintptr(unsafe.Pointer(&oldState)), 0, 0, 0)
	if err != 0 {
		return fmt.Errorf("failed to get terminal settings: %v", err)
	}

	// Create a copy to modify
	newState := oldState

	// Turn off canonical mode and echo
	newState.Lflag &^= syscall.ICANON | syscall.ECHO

	// Set new terminal attributes
	_, _, err = syscall.Syscall6(syscall.SYS_IOCTL, uintptr(fd), uintptr(syscall.TIOCSETA), uintptr(unsafe.Pointer(&newState)), 0, 0, 0)
	if err != 0 {
		return fmt.Errorf("failed to set terminal to raw mode: %v", err)
	}

	return nil
}

func DisableRawMode() {
	fd := int(os.Stdin.Fd())
	syscall.Syscall6(syscall.SYS_IOCTL, uintptr(fd), uintptr(syscall.TIOCSETA), uintptr(unsafe.Pointer(&oldState)), 0, 0, 0)
}

func GetKey(buffer []byte) string {
	// Arrow keys
	if len(buffer) == 3 && buffer[0] == 0x1b && buffer[1] == '[' {
		switch buffer[2] {
		case 'A':
			return KEY_UP
		case 'B':
			return KEY_DOWN
		case 'C':
			return KEY_RIGHT
		case 'D':
			return KEY_LEFT
		}
	}
	// Backspace key
	if buffer[0] == 127 {
		return KEY_BACKSPACE
	}
	// Printable characters
	if buffer[0] >= 32 && buffer[0] <= 126 {
		return KEY_PRINTABLE
	}
	// Enter key
	if buffer[0] == '\n' {
		return KEY_ENTER
	}
	return KEY_UNKNOWN
}

func InitBuffer() ([]rune, int) {
	return []rune{}, 0
}

func AddToBuffer(r rune, buffer *[]rune, cursor *int) {
	*buffer = append(*buffer, 0)                     // grow buffer by 1
	copy((*buffer)[*cursor+1:], (*buffer)[*cursor:]) // Shift everything right
	(*buffer)[*cursor] = r                           // Insert the new rune
	*cursor++
}

func DeleteFromBuffer(buffer *[]rune, cursor *int) bool {
	if *cursor > 0 {
		copy((*buffer)[*cursor-1:], (*buffer)[*cursor:]) // Shift everything left
		*buffer = (*buffer)[:len(*buffer)-1]             // decrease the buffer size by one
		*cursor--
		return true
	}
	return false
}

func UpdateBuffer(command string, buffer *[]rune, cursor *int) {
	*buffer = []rune(command)
	*cursor = len(*buffer)
}

func ResetBuffer(buffer *[]rune, cursor *int) {
	*buffer = (*buffer)[:0]
	*cursor = 0
}

func MoveCursorLeft(cursor *int) bool {
	if *cursor > 0 {
		*cursor--
		os.Stdout.Write([]byte("\b")) // prints the cursor one slot back
		return true
	}
	return false
}

func MoveCursorRight(cursor *int, max int) bool {
	if *cursor < max {
		*cursor++
		os.Stdout.Write([]byte("\033[C")) // prints the cursor one slot fwd
		return true
	}
	return false
}

func InitCommandHistory() ([]string, int) {
	return []string{}, 0
}

func AddCommand(command string, history *[]string, idx *int) {
	if *idx == len(*history) {
		*history = append(*history, command)
		*idx++
	} else {
		copy((*history)[*idx:], (*history)[*idx+1:]) // Shift everything left
		*history = (*history)[:len(*history)-1]      // decrease the buffer size by one
		*history = append(*history, string(command))
		*idx = len(*history)
	}
}

func RedrawLine(buffer []rune, cursor int) {
	fmt.Print("\r") // Move to beginning of line
	fmt.Print(PROMPT)
	fmt.Print(string(buffer))
	fmt.Print("\033[K") // Clear to end of line
	// Move cursor back if needed
	if cursor < len(buffer) {
		back := len(buffer) - cursor
		fmt.Printf("\033[%dD", back)
	}
}

func PrettyPrint(list []string) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, MAX_COL_PAD, ' ', 0)

	for i, item := range list {
		fmt.Fprintf(w, "- %s\t", item)
		if len(list) > MAX_ROWS && (i+1)%MAX_COLS == 0 {
			fmt.Fprintln(w)
		}
	}
	fmt.Fprintln(w)
	w.Flush()
}
