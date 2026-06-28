package authUI

import (
	tea "charm.land/bubbletea/v2"
	lipgloss "charm.land/lipgloss/v2"
	"fmt"
)

type sessionState int

const (
	home sessionState = iota
	login
	register
)

type authIntro struct {
	width        int
	height       int
	login        loginModel
	register     registerModel
	sessionState sessionState
}

func InitializeAuthIntro() authIntro {
	return authIntro{}
}

func (a authIntro) Init() tea.Cmd {
	// this print makes the thing go full screen,
	// somehow making view take the entire width and height wasn't that good of an idea it seems.
	fmt.Print("\x1bc\x1b[H")
	return nil
}

func (a authIntro) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		a.width = msg.Width
		a.height = msg.Height
		a.register = initializeRegisterModel(a.width, a.height)
		a.login = initializeLoginModel(a.width)
	}

	switch a.sessionState {
	case home:
		switch msg := msg.(type) {
		case tea.KeyPressMsg:
			switch msg.String() {
			case "q":
				return a, tea.Quit
			case "l", "L":
				a.sessionState = login
			case "r", "R":
				a.sessionState = register
			}
		}
	case login:
		if key, ok := msg.(tea.KeyPressMsg); ok && key.String() == "esc" {
			a.sessionState = home
			return a, nil
		}
		model, cmd := a.login.Update(msg)
		a.login = model.(loginModel)
		return a, cmd

	case register:
		if key, ok := msg.(tea.KeyPressMsg); ok && key.String() == "esc" {
			a.sessionState = home
			return a, nil
		}
		model, cmd := a.register.Update(msg)
		a.register = model.(registerModel)
		return a, cmd

	}
	return a, nil
}

func (a authIntro) View() tea.View {
	switch a.sessionState {
	case login:
		return a.login.View()
	case register:
		return a.register.View()
	}
	login := fmt.Sprintf("%-9s [L]", "Login")
	register := fmt.Sprintf("%-9s [R]", "Register")
	content := lipgloss.JoinVertical(
		lipgloss.Center,
		asciiStyle.Render(logo),
		selectStyle.Render(login),
		selectStyle.Render(register),
	)
	return tea.NewView(
		lipgloss.Place(
			a.width,
			a.height,
			lipgloss.Center,
			lipgloss.Center,
			content,
		),
	)
}
