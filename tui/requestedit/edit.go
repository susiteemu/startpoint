package editui

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var editViewStyle = lipgloss.NewStyle().Margin(1, 2).BorderBackground(lipgloss.Color("#cdd6f4")).Align(lipgloss.Center)
var contentAreaStyle = lipgloss.NewStyle().Align(lipgloss.Center)

// keyMap defines a set of keybindings. To work for help it must satisfy
// key.Map. It could also very easily be a map[string]key.Binding.
type keyMap struct {
	Save key.Binding
	Quit key.Binding
}

// ShortHelp returns keybindings to be shown in the mini help view. It's part
// of the key.Map interface.
func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Save, k.Quit}
}

// FullHelp returns keybindings for the expanded help view. It's part of the
// key.Map interface.
func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Save, k.Quit}, // first column
	}
}

var keys = keyMap{
	Save: key.NewBinding(
		key.WithKeys("ctrl+s"),
		key.WithHelp("ctrl+s", "save"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "esc", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
}

type Model struct {
	nameInput   textinput.Model
	contentArea textarea.Model
	focused     int
	err         error
	keys        keyMap
	help        help.Model
}

type (
	errMsg error
)

var (
	inputStyle = lipgloss.NewStyle()
)

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmds []tea.Cmd = make([]tea.Cmd, 2)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyShiftTab, tea.KeyCtrlP:
			m.prevInput()
		case tea.KeyTab, tea.KeyCtrlN:
			m.nextInput()
		}
		m.nameInput.Blur()
		m.contentArea.Blur()
		if m.focused == 0 {
			m.nameInput.Focus()
		} else {
			m.contentArea.Focus()
		}

	// We handle errors just like any other message
	case errMsg:
		m.err = msg
		return m, nil
	}

	m.nameInput, cmds[0] = m.nameInput.Update(msg)
	m.contentArea, cmds[1] = m.contentArea.Update(msg)

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {

	helpView := m.help.View(m.keys)

	inputViews := []string{}
	inputViews = append(inputViews, inputStyle.Width(10).Render("Name"))
	inputViews = append(inputViews, m.nameInput.View())
	inputViews = append(inputViews, inputStyle.Width(10).Render("Request"))
	inputViews = append(inputViews, contentAreaStyle.Render(m.contentArea.View()))
	inputViews = append(inputViews, helpView)

	return lipgloss.JoinVertical(lipgloss.Left, inputViews...)

}

/* 	return fmt.Sprintf(
		`
%s
 %s

%s
%s

%s
`,
		inputStyle.Width(30).Render("Name"),
		m.nameInput.View(),
		inputStyle.Width(30).Render("Content"),
		m.contentArea.View(),
		continueStyle.Render("Continue ->"),
		) + "\n" */

// nextInput focuses the next input field
func (m *Model) nextInput() {
	m.focused = (m.focused + 1) % 2
}

// prevInput focuses the previous input field
func (m *Model) prevInput() {
	m.focused--
	// Wrap around
	if m.focused < 0 {
		m.focused = 2 - 1
	}
}

func New() Model {
	nameInput := textinput.New()
	nameInput.Placeholder = "Name your request"
	nameInput.Focus()
	nameInput.CharLimit = 64
	nameInput.Width = 30
	nameInput.Prompt = ""
	//inputs[ccn].Validate = ccnValidator

	contentArea := textarea.New()
	contentArea.SetHeight(15)
	contentArea.SetWidth(80)
	contentArea.Placeholder = ""

	return Model{
		nameInput:   nameInput,
		contentArea: contentArea,
		focused:     0,
		err:         nil,
		keys:        keys,
		help:        help.New(),
	}
}
