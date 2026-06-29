package authUI

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"unicode"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	lipgloss "charm.land/lipgloss/v2"
)

type registerModel struct {
	width      int
	height     int
	focusIndex int
	inputs     []textinput.Model
	quitting   bool
	Err        string
}

func initializeRegisterModel(width int, height int) registerModel {
	m := registerModel{
		width:  width,
		height: height,
		inputs: make([]textinput.Model, 2),
	}
	var t textinput.Model
	for i := range m.inputs {
		t = textinput.New()
		t.CharLimit = 32
		s := t.Styles()
		s.Cursor.Color = lipgloss.Color(foreground)
		s.Focused.Prompt = focusedStyle
		s.Focused.Text = focusedStyle
		s.Blurred.Prompt = blurredStyle
		t.SetStyles(s)
		switch i {
		case 0:
			t.Placeholder = "Enter your email"
			t.SetWidth(64)
			t.CharLimit = 64
			t.Focus()
		case 1:
			t.Placeholder = "Enter your password"
			t.SetWidth(64)
			t.EchoMode = textinput.EchoPassword
			t.EchoCharacter = '•'
		}
		m.inputs[i] = t
	}
	return m
}

func (m registerModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m registerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case registerMsg:
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
				return m, SubmitReg(m.inputs[0].Value(), m.inputs[1].Value(), "http://localhost:2026/user/register")
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

func (m *registerModel) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))
	// Only text inputs with Focus() set will respond, so it's safe to simply
	// update all of them here without any further logic.
	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}
	return tea.Batch(cmds...)
}

func (m registerModel) View() tea.View {
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

func verifyPassword(s string) (length, number, upper, special bool) {
	letters := len(s)
	for _, c := range s {
		switch {
		case unicode.IsNumber(c):
			number = true
		case unicode.IsUpper(c):
			upper = true
			letters++
		case unicode.IsPunct(c) || unicode.IsSymbol(c):
			special = true
		}
	}
	length = letters > 8
	return
}

type registerMsg string

func submitCredentials(email string, password string, url string) tea.Msg {
	if len(password) == 0 {
		return registerMsg("")
	}
	length, number, upper, special := verifyPassword(password)
	if !length {
		return registerMsg("Password should have more than 8 characters")
	}
	if !number {
		return registerMsg("Password should have at least one digit")
	}
	if !upper {
		return registerMsg("Password should have at least one uppercase character")
	}
	if !special {
		return registerMsg("Password should have at least one special character")
	}
	body := map[string]string{
		"email":    email,
		"password": password,
	}
	btx, err := json.Marshal(body)
	if err != nil {
		content := fmt.Sprintf("Error occured: %v", err)
		return registerMsg(content)
	}
	res, err := http.Post(url, "application/json", bytes.NewReader(btx))
	if err != nil {
		content := fmt.Sprintf("Error occured: %v", err)
		return registerMsg(content)
	}
	if res.StatusCode != http.StatusCreated {
		httpBody, err := io.ReadAll(res.Body)
		if err != nil {
			content := fmt.Sprintf("Error occured: %v", err)
			return registerMsg(content)
		}
		content := fmt.Sprintf("Error occured: %s", httpBody)
		return registerMsg(content)
	}
	return registerMsg("Registration Successful")
}

func SubmitReg(email string, password string, url string) tea.Cmd {
	return func() tea.Msg {
		return submitCredentials(email, password, url)
	}
}
