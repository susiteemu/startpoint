package editui

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var inputStyle = lipgloss.NewStyle().BorderForeground(lipgloss.Color("36")).BorderStyle(lipgloss.NormalBorder()).Padding(1).Width(70)
var descriptionStyle = lipgloss.NewStyle().Padding(1).Width(70)

type keyMap struct {
	Save key.Binding
	Quit key.Binding
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Save, k.Quit}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Save, k.Quit}, // first column
	}
}

var keys = keyMap{
	Save: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "create"),
	),
	Quit: key.NewBinding(
		key.WithKeys("esc", "ctrl+c"),
		key.WithHelp("esc", "quit"),
	),
}

type Model struct {
	nameInput textinput.Model
	Name      string
	Complex   bool
	keys      keyMap
	help      help.Model
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyEnter:
			m.Name = m.nameInput.Value()
			return m, tea.Cmd(func() tea.Msg { return CreateMsg{} })
		}
	}

	var cmd tea.Cmd
	m.nameInput, cmd = m.nameInput.Update(msg)

	return m, cmd
}

func (m Model) View() string {

	helpView := m.help.View(m.keys)

	inputViews := []string{}
	inputViews = append(inputViews, "Name")

	if m.Complex {
		inputViews = append(inputViews, descriptionStyle.Render("Choose a name for your complex request. Make it filename compatible and unique within this workspace. After pressing <enter> program will open your $EDITOR and quit. You will then be able to write the contents of the request."))

	} else {
		inputViews = append(inputViews, descriptionStyle.Render("Choose a name for your request. Make it filename compatible and unique within this workspace. After pressing <enter> program will open your $EDITOR and quit. You will then be able to write the contents of the request."))

	}

	inputViews = append(inputViews, inputStyle.Render(m.nameInput.View()))
	inputViews = append(inputViews, helpView)

	return lipgloss.JoinVertical(lipgloss.Left, inputViews...)

}

func New(complex bool) Model {
	nameInput := textinput.New()
	nameInput.Placeholder = "Name your request"
	nameInput.Focus()
	nameInput.CharLimit = 64
	nameInput.Width = 70
	nameInput.Prompt = ""
	//inputs[ccn].Validate = ccnValidator

	return Model{
		nameInput: nameInput,
		Complex:   complex,
		keys:      keys,
		help:      help.New(),
	}
}

type CreateMsg struct{}
