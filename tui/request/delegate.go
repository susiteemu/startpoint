package requestui

import (
	keyprompt "startpoint/tui/keyprompt"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

var selectModeKeys = []key.Binding{
	key.NewBinding(
		key.WithKeys("p"),
		key.WithHelp("p", "preview"),
	),
	key.NewBinding(
		key.WithKeys("r"),
		key.WithHelp("r", "run"),
	),
	key.NewBinding(
		key.WithKeys("i"),
		key.WithHelp("i", "edit mode"),
	),
	key.NewBinding(
		key.WithKeys("a"),
		key.WithHelp("a", "activate profile"),
	),
	key.NewBinding(
		key.WithKeys("ctrl+n"),
		key.WithHelp("ctrl+n", "switch to Profiles"),
	),
}

var editModeKeys = []key.Binding{
	key.NewBinding(
		key.WithKeys("a"),
		key.WithHelp("a", "add"),
	),
	key.NewBinding(
		key.WithKeys("e"),
		key.WithHelp("e", "edit"),
	),
	key.NewBinding(
		key.WithKeys("d"),
		key.WithHelp("d", "delete"),
	),
	key.NewBinding(
		key.WithKeys("p"),
		key.WithHelp("p", "preview"),
	),
	key.NewBinding(
		key.WithKeys("r"),
		key.WithHelp("r", "rename"),
	),
	key.NewBinding(
		key.WithKeys("c"),
		key.WithHelp("c", "copy"),
	),
	key.NewBinding(
		key.WithKeys(tea.KeyEsc.String()),
		key.WithHelp(tea.KeyEsc.String(), "view mode"),
	),
	key.NewBinding(
		key.WithKeys("ctrl+n"),
		key.WithHelp("ctrl+n", "switch to Profiles"),
	),
}

func newBaseDelegate() list.DefaultDelegate {
	d := list.NewDefaultDelegate()

	d.Styles.SelectedTitle = d.Styles.SelectedTitle.Foreground(style.listItemTitleColor).BorderLeftForeground(style.listItemTitleColor)
	d.Styles.SelectedDesc = d.Styles.SelectedTitle.Foreground(style.listItemDescColor).BorderLeftForeground(style.listItemDescColor)

	d.Styles.NormalTitle = d.Styles.NormalTitle.Foreground(style.listItemTitleColor)
	d.Styles.NormalDesc = d.Styles.NormalTitle.Foreground(style.listItemDescColor)

	d.Styles.DimmedTitle = d.Styles.DimmedTitle.Foreground(style.listItemTitleColor)
	d.Styles.DimmedDesc = d.Styles.DimmedTitle.Foreground(style.listItemDescColor)

	return d
}

func newSelectDelegate() list.DefaultDelegate {
	d := newBaseDelegate()

	d.UpdateFunc = func(msg tea.Msg, m *list.Model) tea.Cmd {
		var request Request
		requestSelected := false
		if i, ok := m.SelectedItem().(Request); ok {
			request = i
			requestSelected = true
		}

		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch keypress := msg.String(); keypress {
			case "r":
				if requestSelected {
					return tea.Cmd(func() tea.Msg {
						return RunRequestMsg{
							Request: request,
						}
					})
				}
			case "p":
				if requestSelected {
					return tea.Cmd(func() tea.Msg {
						return PreviewRequestMsg{
							Request: request,
						}
					})
				}
			case "a":
				return tea.Cmd(func() tea.Msg {
					return ActivateProfile{}
				})
			}
		}

		return nil
	}

	d.ShortHelpFunc = func() []key.Binding {
		return selectModeKeys
	}

	d.FullHelpFunc = func() [][]key.Binding {
		return [][]key.Binding{selectModeKeys}
	}

	return d

}

func newEditModeDelegate() list.DefaultDelegate {
	d := newBaseDelegate()

	d.UpdateFunc = func(msg tea.Msg, m *list.Model) tea.Cmd {
		var request Request
		requestSelected := false
		if i, ok := m.SelectedItem().(Request); ok {
			request = i
			requestSelected = true
		}

		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch keypress := msg.String(); keypress {
			case "a":
				var keys []keyprompt.KeypromptEntry
				keys = append(keys, keyprompt.KeypromptEntry{
					Text: "yaml", Key: "y",
				})
				keys = append(keys, keyprompt.KeypromptEntry{
					Text: "starlark", Key: "s",
				})

				return tea.Cmd(func() tea.Msg {
					return ShowKeyprompt{
						Label:   "Select type of request to create",
						Entries: keys,
						Type:    CreateRequest,
					}
				})
			case "d":
				if requestSelected {
					return tea.Cmd(func() tea.Msg {
						return DeleteRequestMsg{
							Request: request,
						}
					})
				}
			case "e":
				if requestSelected {
					return tea.Cmd(func() tea.Msg {
						return EditRequestMsg{
							Request: request,
						}
					})
				}
			case "p":
				if requestSelected {
					return tea.Cmd(func() tea.Msg {
						return PreviewRequestMsg{
							Request: request,
						}
					})
				}
			case "r":
				if requestSelected {
					return tea.Cmd(func() tea.Msg {
						return RenameRequestMsg{
							Request: request,
						}
					})
				}
			case "c":
				if requestSelected {
					return tea.Cmd(func() tea.Msg {
						return CopyRequestMsg{
							Request: request,
						}
					})
				}
			}
		}

		return nil
	}

	d.ShortHelpFunc = func() []key.Binding {
		return editModeKeys
	}

	d.FullHelpFunc = func() [][]key.Binding {
		return [][]key.Binding{editModeKeys}
	}

	return d

}
