package managetui

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

var keys = []key.Binding{
	key.NewBinding(
		key.WithKeys("a"),
		key.WithHelp("a", "Add new profile"),
	)}

func newBaseDelegate() list.DefaultDelegate {
	d := list.NewDefaultDelegate()
	d.SetHeight(3)
	d.Styles.SelectedTitle = d.Styles.SelectedTitle.Foreground(requestTitleColor).BorderLeftForeground(requestTitleColor)
	d.Styles.SelectedDesc = d.Styles.SelectedTitle.Foreground(requestDescColor).BorderLeftForeground(requestDescColor)
	return d
}

func newSelectDelegate() list.DefaultDelegate {
	d := newBaseDelegate()
	d.UpdateFunc = func(msg tea.Msg, m *list.Model) tea.Cmd {
		var profile Profile
		if i, ok := m.SelectedItem().(Profile); ok {
			profile = i
		} else {
			return nil
		}

		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch keypress := msg.String(); keypress {
			case "enter":
				return tea.Cmd(func() tea.Msg {
					return ProfileSelectedMsg{
						Profile: profile,
					}
				})
			}
		}

		return nil
	}

	d.ShortHelpFunc = func() []key.Binding {
		return keys
	}

	d.FullHelpFunc = func() [][]key.Binding {
		return [][]key.Binding{keys}
	}

	return d

}
