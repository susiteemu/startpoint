package keypromptui

import (
	"fmt"
	"github.com/susiteemu/startpoint/tui/styles"
	"slices"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/rs/zerolog/log"
)

var (
	promptStyle      lipgloss.Style
	descriptionStyle lipgloss.Style
	entryKeyStyle    lipgloss.Style
	entryTextStyle   lipgloss.Style
)

type Model struct {
	width      int
	label      string
	entries    []KeypromptEntry
	keys       []string
	promptType string
	payload    interface{}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		newWidth := min(40, msg.Width-2)
		m.width = newWidth
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case tea.KeyEsc.String():
			return m, tea.Cmd(func() tea.Msg {
				return KeypromptCancelledMsg{}
			})
		default:
			index := slices.Index(m.keys, keypress)
			if index != -1 {
				key := m.keys[index]
				log.Debug().Msgf("Key %s", key)
				return m, tea.Cmd(func() tea.Msg {
					return KeypromptAnsweredMsg{
						Key:     key,
						Type:    m.promptType,
						Payload: m.payload,
					}
				})
			}
		}
	}

	return m, nil
}

func (m Model) View() string {

	inputViews := []string{}
	inputViews = append(inputViews, descriptionStyle.Width(m.width).Render(m.label))

	renderItems := []KeypromptEntry{}
	maxKeyLen := 0
	for _, v := range m.entries {
		if len(v.Key) > maxKeyLen {
			maxKeyLen = len(v.Key)
		}
		renderItems = append(renderItems, v)
	}
	escKey := KeypromptEntry{
		Text: "cancel", Key: tea.KeyEsc.String(),
	}
	if len(escKey.Key) > maxKeyLen {
		maxKeyLen = len(escKey.Key)
	}

	for _, v := range renderItems {
		key := fmt.Sprintf("%s", entryKeyStyle.Render(v.Key))
		separator := strings.Repeat(" ", 1+maxKeyLen-len(v.Key))
		inputViews = append(inputViews, fmt.Sprintf("%s%s%s", key, separator, entryTextStyle.Render(v.Text)))
	}

	inputViews = append(inputViews, "")

	separator := strings.Repeat(" ", 1+maxKeyLen-len(escKey.Key))
	inputViews = append(inputViews, fmt.Sprintf("%s%s%s", entryKeyStyle.Render(escKey.Key), separator, entryTextStyle.Render(escKey.Text)))

	return promptStyle.Width(m.width).Render(lipgloss.JoinVertical(lipgloss.Left, inputViews...))
}

func New(label string, entries []KeypromptEntry, promptType string, payload interface{}, width int) Model {

	theme := styles.GetTheme()
	commonStyles := styles.GetCommonStyles(theme)

	promptStyle = lipgloss.NewStyle().BorderForeground(theme.BorderFgColor).BorderStyle(lipgloss.RoundedBorder()).Padding(1, 2)
	descriptionStyle = lipgloss.NewStyle().Foreground(theme.TextFgColor).PaddingBottom(1)
	entryKeyStyle = commonStyles.HelpKeyStyle
	entryTextStyle = commonStyles.HelpKeyStyle.Faint(true)

	var keys []string
	for _, v := range entries {
		keys = append(keys, v.Key)
	}

	return Model{
		label:      label,
		keys:       keys,
		entries:    entries,
		promptType: promptType,
		payload:    payload,
		width:      min(40, width-2),
	}
}

type KeypromptEntry struct {
	Text string
	Key  string
}

type KeypromptAnsweredMsg struct {
	Key     string
	Type    string
	Payload interface{}
}

type KeypromptCancelledMsg struct{}
