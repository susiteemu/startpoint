package promptui

import (
	"startpoint/tui/styles"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/rs/zerolog/log"
)

// TODO: colors from theme
var (
	promptStyle      = lipgloss.NewStyle().BorderStyle(lipgloss.RoundedBorder()).Padding(1)
	inputStyle       = lipgloss.NewStyle().BorderStyle(lipgloss.NormalBorder()).Padding(1).MarginTop(1)
	errInputStyle    = inputStyle.Copy()
	descriptionStyle = lipgloss.NewStyle()
	helpStyle        = lipgloss.NewStyle().Padding(1, 1, 0, 1)
)

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
	nameInput    textinput.Model
	initialValue string
	context      PromptContext
	label        string
	keys         keyMap
	help         help.Model
	width        int
	validator    func(s string) error
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		newWidth := min(64, msg.Width)
		m.width = newWidth
		m.nameInput.Width = m.width - 2
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEsc:
			return m, tea.Cmd(func() tea.Msg {
				return PromptCancelledMsg{}
			})
		case tea.KeyEnter:
			if m.nameInput.Err == nil {
				return m, tea.Cmd(func() tea.Msg {
					return PromptAnsweredMsg{
						Context: m.context,
						Input:   m.nameInput.Value(),
					}
				})
			}
		}
		cmds = append(cmds, tea.Cmd(func() tea.Msg {
			return promptTyped(msg.String())
		}))
	case promptTyped:
		if m.nameInput.Value() != m.initialValue {
			err := m.validator(m.nameInput.Value())
			log.Debug().Msgf("Validation result %v", err)
			m.nameInput.Err = err
		}
	}

	var cmd tea.Cmd
	m.nameInput, cmd = m.nameInput.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {

	/*
	*  WARN: contains magic numbers. Dunno exactly why, but the numbers used below make things work.
	 */

	helpView := helpStyle.Render(m.help.View(m.keys))

	inputViews := []string{}
	var descStyle = descriptionStyle
	inputViews = append(inputViews, descStyle.Render(m.label))

	var style = inputStyle.Width(m.width - 6)
	if m.nameInput.Err != nil {
		style = errInputStyle.Width(m.width - 6)
	}

	inputViews = append(inputViews, style.Render(m.nameInput.View()))
	inputViews = append(inputViews, helpView)

	return promptStyle.Width(m.width - 2).Render(lipgloss.JoinVertical(lipgloss.Left, inputViews...))

}

func New(context PromptContext, initialValue string, label string, validator func(s string) error, w int) Model {

	theme := styles.GetTheme()
	commonStyles := styles.GetCommonStyles(theme)

	promptStyle = promptStyle.Foreground(theme.TextFgColor).BorderForeground(theme.BorderFgColor)
	inputStyle = inputStyle.Foreground(theme.TextFgColor).BorderForeground(theme.BorderFgColor)
	errInputStyle = errInputStyle.Foreground(theme.TextFgColor).BorderForeground(theme.ErrorFgColor)

	nameInput := textinput.New()
	nameInput.Focus()
	nameInput.CharLimit = 64
	nameInput.Width = min(64, w-2)
	nameInput.SetValue(initialValue)
	nameInput.Prompt = ""
	nameInput.Cursor.Style = lipgloss.NewStyle().Foreground(theme.CursorFgColor).Background(theme.CursorBgColor)

	help := help.New()
	help.Styles.FullKey = commonStyles.HelpKeyStyle
	help.Styles.FullDesc = commonStyles.HelpDescStyle
	help.Styles.ShortKey = commonStyles.HelpKeyStyle
	help.Styles.ShortDesc = commonStyles.HelpDescStyle
	help.Styles.ShortSeparator = commonStyles.HelpSeparatorStyle
	help.Styles.FullSeparator = commonStyles.HelpSeparatorStyle
	help.ShortSeparator = "  "
	help.FullSeparator = "  "

	return Model{
		context:      context,
		nameInput:    nameInput,
		initialValue: initialValue,
		label:        label,
		keys:         keys,
		help:         help,
		width:        min(64, w),
		validator:    validator,
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

type promptTyped string
