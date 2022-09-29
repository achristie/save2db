package configure

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	focusedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	blurredStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	cursorStyle  = focusedStyle.Copy()
	noStyle      = lipgloss.NewStyle()

	focusedButton = focusedStyle.Copy().Render("[ Submit ]")
	blurredButton = fmt.Sprintf("[ %s ]", blurredStyle.Render("Submit"))
)

type model struct {
	focusIndex int
	inputs     []textinput.Model
}

func NewSQLiteModel() model {
	m := model{
		inputs: make([]textinput.Model, 1),
	}

	t := textinput.New()
	t.Placeholder = "Path"
	t.Focus()
	t.TextStyle = focusedStyle
	t.PromptStyle = focusedStyle

	m.inputs[0] = t
	return m
}

func NewConfigureDBModel() model {
	m := model{
		inputs: make([]textinput.Model, 5),
	}

	var t textinput.Model
	for i := range m.inputs {
		t = textinput.New()
		t.CursorStyle = cursorStyle
		t.CharLimit = 32
		switch i {
		case 0:
			t.Placeholder = "DBHost"
			t.Focus()
			t.PromptStyle = focusedStyle
			t.TextStyle = focusedStyle
		case 1:
			t.Placeholder = "DBPort"
			t.PromptStyle = focusedStyle
			t.TextStyle = focusedStyle
		case 2:
			t.Placeholder = "DBUsername"
			t.PromptStyle = focusedStyle
			t.TextStyle = focusedStyle
		case 3:
			t.Placeholder = "DBPassword"
			t.EchoMode = textinput.EchoPassword
			t.EchoCharacter = '•'
		case 4:
			t.Placeholder = "DBName"
			t.PromptStyle = focusedStyle
			t.TextStyle = focusedStyle
		}

		m.inputs[i] = t
	}

	return m
}

type BackMsg struct{}

func backCmd() tea.Cmd {
	return func() tea.Msg {
		return BackMsg{}
	}
}

type SubmitMsg struct {
	Values map[string]string
}

func submitCmd(values map[string]string) tea.Cmd {
	return func() tea.Msg {
		return SubmitMsg{Values: values}
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, backCmd()

		// Set focus to next input
		case "tab", "shift+tab", "enter", "up", "down":
			s := msg.String()

			// Did the user press enter while the submit button was focused?
			if s == "enter" && m.focusIndex == len(m.inputs) {
				d := make(map[string]string)
				for _, v := range m.inputs {
					d[v.Placeholder] = v.Value()
				}
				return m, submitCmd(d)
			}

			// Cycle indexes
			if s == "up" || s == "shift+tab" {
				m.focusIndex--
			} else {
				m.focusIndex++
			}

			if m.focusIndex > len(m.inputs) {
				m.focusIndex = 0
			} else if m.focusIndex < 0 {
				m.focusIndex = len(m.inputs)
			}

			cmds := make([]tea.Cmd, len(m.inputs))
			for i := 0; i <= len(m.inputs)-1; i++ {
				if i == m.focusIndex {
					// Set focused state
					cmds[i] = m.inputs[i].Focus()
					m.inputs[i].PromptStyle = focusedStyle
					m.inputs[i].TextStyle = focusedStyle
					continue
				}
				// Remove focused state
				m.inputs[i].Blur()
				m.inputs[i].PromptStyle = noStyle
				m.inputs[i].TextStyle = noStyle
			}

			return m, tea.Batch(cmds...)
		}
	}

	// Handle character input and blinking
	cmd := m.updateInputs(msg)

	return m, cmd
}

func (m *model) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))

	// Only text inputs with Focus() set will respond, so it's safe to simply
	// update all of them here without any further logic.
	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

func (m model) View() string {
	var b strings.Builder

	for i := range m.inputs {
		b.WriteString(m.inputs[i].View())
		if i < len(m.inputs)-1 {
			b.WriteRune('\n')
		}
	}

	button := &blurredButton
	if m.focusIndex == len(m.inputs) {
		button = &focusedButton
	}
	fmt.Fprintf(&b, "\n\n%s\n\n", *button)

	return b.String()
}

func NewProgram(model model) {
	if err := tea.NewProgram(model).Start(); err != nil {
		fmt.Printf("could not start program: %s\n", err)
		os.Exit(1)
	}
}
