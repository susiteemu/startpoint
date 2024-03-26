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
		key.WithHelp("enter", "ok"),
	),
	Quit: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "cancel"),
	),
}

type Model struct {
	nameInput textinput.Model
	context   PromptContext
	label     string
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
		case tea.KeyEsc:
			return m, tea.Cmd(func() tea.Msg {
				return PromptCancelledMsg{}
			})
		case tea.KeyEnter:
			return m, tea.Cmd(func() tea.Msg {
				return PromptAnsweredMsg{
					Context: m.context,
					Input:   m.nameInput.Value(),
				}
			})
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
	inputViews = append(inputViews, descriptionStyle.Render(m.label))

	inputViews = append(inputViews, inputStyle.Render(m.nameInput.View()))
	inputViews = append(inputViews, helpView)

	return lipgloss.JoinVertical(lipgloss.Left, inputViews...)

}

func New(context PromptContext, initialValue string, label string) Model {
	nameInput := textinput.New()
	nameInput.Focus()
	nameInput.CharLimit = 64
	nameInput.Width = 70
	nameInput.SetValue(initialValue)
	nameInput.Prompt = ""
	//inputs[ccn].Validate = ccnValidator

	return Model{
		context:   context,
		nameInput: nameInput,
		label:     label,
		keys:      keys,
		help:      help.New(),
	}
}

type PromptContext struct {
	Key        string
	Additional interface{}
}

type PromptAnsweredMsg struct {
	Context PromptContext
	Input   string
}

type PromptCancelledMsg struct{}
