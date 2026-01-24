package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	// Initialize the board with our mock data
	board := NewBoard()

	// Start the Bubble Tea program
	// tea.WithAltScreen() allows the app to take up the full terminal window
	p := tea.NewProgram(board, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
