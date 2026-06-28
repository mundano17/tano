package authUI

import (
	"fmt"

	lipgloss "charm.land/lipgloss/v2"
)

const (
	background string = "#131313" // Onyx
	foreground string = "#f4f1de" // Eggshell
	accent     string = "#d72638" // Classic Crimson
)

var logo string = `
████████╗ █████╗ ███╗   ██╗ ██████╗
╚══██╔══╝██╔══██╗████╗  ██║██╔═══██╗
   ██║   ███████║██╔██╗ ██║██║   ██║
   ██║   ██╔══██║██║╚██╗██║██║   ██║
   ██║   ██║  ██║██║ ╚████║╚██████╔╝
   ╚═╝   ╚═╝  ╚═╝╚═╝  ╚═══╝ ╚═════╝
`

var asciiStyle lipgloss.Style = lipgloss.NewStyle().
	Foreground(lipgloss.Color(accent)).
	//	Background(lipgloss.Color(background)).
	Align(lipgloss.Center)

var selectStyle lipgloss.Style = lipgloss.NewStyle().
	Foreground(lipgloss.Color(foreground)).
	//Background(lipgloss.Color(background)).
	Align(lipgloss.Center)

var label lipgloss.Style = lipgloss.NewStyle().
	Foreground(lipgloss.Color(foreground)).
	Bold(true)

var input lipgloss.Style = lipgloss.NewStyle().
	Foreground(lipgloss.Color(foreground))

// Input colors for Login and Register Model
// TODO: colora mathanum
var (
	focusedStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color(accent))
	blurredStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color(foreground))
	cursorStyle         = focusedStyle
	noStyle             = lipgloss.NewStyle()
	helpStyle           = blurredStyle
	cursorModeHelpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(accent))
	focusedButton       = focusedStyle.Render("[ Submit ]")
	blurredButton       = fmt.Sprintf("[ %s ]", blurredStyle.Render("Submit"))
)
