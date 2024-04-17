package keypromptui

import (
	"fmt"
	"slices"
	"startpoint/tui/styles"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/rs/zerolog/log"
)

var promptStyle = lipgloss.NewStyle().BorderForeground(lipgloss.Color("#cdd6f44")).BorderStyle(lipgloss.RoundedBorder()).Padding(1, 2)
var descriptionStyle = lipgloss.NewStyle().PaddingBottom(1)
var entryKeyStyle = styles.HelpKeyStyle.Copy()
var entryTextStyle = styles.HelpKeyStyle.Copy().Faint(true)

type Model struct {
	width   int
	height  int
	label   string
	entries []KeypromptEntry
	keys    []string
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
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
						Key: key,
					}
				})
			}
		}
	}

	return m, nil
}

func (m Model) View() string {

	inputViews := []string{}
	inputViews = append(inputViews, descriptionStyle.Render(m.label))

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
	inputViews = append(inputViews, fmt.Sprintf("%s%s%s", escKey.Key, separator, entryTextStyle.Render(escKey.Text)))

	return promptStyle.Render(lipgloss.JoinVertical(lipgloss.Left, inputViews...))
}

func New(label string, entries []KeypromptEntry) Model {

	var keys []string
	for _, v := range entries {
		keys = append(keys, v.Key)
	}

	return Model{
		label:   label,
		keys:    keys,
		entries: entries,
	}
}

type KeypromptEntry struct {
	Text string
	Key  string
}

type KeypromptAnsweredMsg struct {
	Key string
}

type KeypromptCancelledMsg struct{}
