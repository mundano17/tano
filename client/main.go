package main

import (
	"fmt"
	"os"
	authUI "tanoclient/command/auth"

	tea "charm.land/bubbletea/v2"
)

func main() {
	// client := auth.InitializeNewAPIClient()
	p := tea.NewProgram(authUI.InitializeAuthIntro())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
