package authUI

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	lipgloss "charm.land/lipgloss/v2"
)

type loginModel struct {
	width      int
	height     int
	focusIndex int
	inputs     []textinput.Model
	quitting   bool
	Err        string
}

func initializeLoginModel(width int) loginModel {
	m := loginModel{
		inputs: make([]textinput.Model, 2),
	}
	var t textinput.Model
	for i := range m.inputs {
		t = textinput.New()
		t.CharLimit = 32
		s := t.Styles()
		s.Cursor.Color = lipgloss.Color("205")
		s.Focused.Prompt = focusedStyle
		s.Focused.Text = focusedStyle
		s.Blurred.Prompt = blurredStyle
		s.Focused.Text = focusedStyle
		t.SetStyles(s)
		switch i {
		case 0:
			t.Placeholder = "Email"
			t.CharLimit = 64
			t.SetWidth(width)
			t.Focus()
		case 1:
			t.Placeholder = "Password"
			t.SetWidth(width)
			t.EchoMode = textinput.EchoPassword
			t.EchoCharacter = '•'
		}
		m.inputs[i] = t
	}
	return m
}

func (m loginModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m loginModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case loginMsg:
		m.Err = string(msg)
	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			m.quitting = true
			return m, tea.Quit
		// Set focus to next input
		case "tab", "shift+tab", "enter", "up", "down":
			s := msg.String()
			if s == "enter" {
				return m, SubmitLogin(m.inputs[0].Value(), m.inputs[1].Value(), "http://localhost:2026/user/login")
			}
			// Cycle indexes
			if s == "up" || s == "shift+tab" {
				m.focusIndex--
			} else {
				m.focusIndex++
			}
			if m.focusIndex >= len(m.inputs) {
				m.focusIndex = 0
			}
			cmds := make([]tea.Cmd, len(m.inputs))
			for i := 0; i <= len(m.inputs)-1; i++ {
				if i == m.focusIndex {
					// Set focused state
					cmds[i] = m.inputs[i].Focus()
					continue
				}
				// Remove focused state
				m.inputs[i].Blur()
			}
			return m, tea.Batch(cmds...)
		}
	}
	// Handle character input and blinking
	cmd := m.updateInputs(msg)
	return m, cmd
}

func (m *loginModel) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))
	// Only text inputs with Focus() set will respond, so it's safe to simply
	// update all of them here without any further logic.
	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}
	return tea.Batch(cmds...)
}

func (m loginModel) View() tea.View {
	var b strings.Builder
	for i := range m.inputs {
		b.WriteString(m.inputs[i].View())
		if i < len(m.inputs)-1 {
			b.WriteRune('\n')
		}
	}
	if m.quitting {
		b.WriteRune('\n')
	}
	if m.Err != "" {
		b.WriteRune('\n')
		b.WriteRune('\n')
		b.WriteString(m.Err)
		m.Err = ""
	}
	content := lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		b.String(),
	)
	v := tea.NewView(content)
	return v
}

type loginMsg string

func submitLoginCredentials(email string, password string, url string) tea.Msg {
	body := map[string]string{
		"email":    email,
		"password": password,
	}
	btx, err := json.Marshal(body)
	if err != nil {
		content := fmt.Sprintf("Error occured: %v", err)
		return loginMsg(content)
	}
	res, err := http.Post(url, "application/json", bytes.NewReader(btx))
	if err != nil {
		content := fmt.Sprintf("Error occured: %v", err)
		return loginMsg(content)
	}
	if res.StatusCode != http.StatusAccepted {
		httpBody, err := io.ReadAll(res.Body)
		if err != nil {
			content := fmt.Sprintf("Error occured: %v", err)
			return loginMsg(content)
		}
		content := fmt.Sprintf("Error occured: %s", httpBody)
		return loginMsg(content)
	}
	return loginMsg("Login Successful")
}

func SubmitLogin(email string, password string, url string) tea.Cmd {
	return func() tea.Msg {
		return submitLoginCredentials(email, password, url)
	}
}
