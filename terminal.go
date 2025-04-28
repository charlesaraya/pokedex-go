package main

import (
	"fmt"
	"os"
	"syscall"
	"unsafe"
)

var oldState syscall.Termios

func enableRawMode() error {

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

func disableRawMode() {
	fd := int(os.Stdin.Fd())
	syscall.Syscall6(syscall.SYS_IOCTL, uintptr(fd), uintptr(syscall.TIOCSETA), uintptr(unsafe.Pointer(&oldState)), 0, 0, 0)
}

func redrawLine(buffer []rune, cursor int, prompt string) {
	fmt.Print("\r") // Move to beginning of line
	fmt.Print(prompt)
	fmt.Print(string(buffer))
	fmt.Print("\033[K") // Clear to end of line
	// Move cursor back if needed
	if cursor < len(buffer) {
		back := len(buffer) - cursor
		fmt.Printf("\033[%dD", back)
	}
}
