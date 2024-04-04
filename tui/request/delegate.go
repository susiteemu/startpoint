package requestui

import (
	"goful/core/client/validator"
	keyprompt "goful/tui/keyprompt"

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
		key.WithKeys(tea.KeyEnter.String()),
		key.WithHelp(tea.KeyEnter.String(), "run request"),
	),
	key.NewBinding(
		key.WithKeys("i"),
		key.WithHelp("i", "edit mode"),
	),
	key.NewBinding(
		key.WithKeys("a"),
		key.WithHelp("a", "activate profile"),
	),
}

var editModeKeys = []key.Binding{
	key.NewBinding(
		key.WithKeys("a"),
		key.WithHelp("a", "add"),
	),
	key.NewBinding(
		key.WithKeys("p"),
		key.WithHelp("p", "preview"),
	),
	key.NewBinding(
		key.WithKeys("x"),
		key.WithHelp("x", "delete"),
	),
	key.NewBinding(
		key.WithKeys(tea.KeyEnter.String(), "e"),
		key.WithHelp(tea.KeyEnter.String()+"/e", "edit"),
	),
	key.NewBinding(
		key.WithKeys("r"),
		key.WithHelp("r", "rename"),
	),
	key.NewBinding(
		key.WithKeys("c"),
		key.WithHelp("c", "copy request"),
	),
	key.NewBinding(
		key.WithKeys(tea.KeyEsc.String()),
		key.WithHelp(tea.KeyEsc.String(), "view mode"),
	),
}

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
		var request Request
		if i, ok := m.SelectedItem().(Request); ok {
			request = i
		} else {
			return nil
		}

		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch keypress := msg.String(); keypress {
			case "enter":
				if !validator.IsValidUrl(request.Url) || !validator.IsValidMethod(request.Method) {
					return tea.Cmd(func() tea.Msg {
						return createStatusMsg("Invalid request.")
					})
				}
				return tea.Cmd(func() tea.Msg {
					return RunRequestMsg{
						Request: request,
					}
				})
			case "p":
				return tea.Cmd(func() tea.Msg {
					return PreviewRequestMsg{
						Request: request,
					}
				})
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
		if i, ok := m.SelectedItem().(Request); ok {
			request = i
		} else {
			return nil
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
					}
				})
			case "x":
				return tea.Cmd(func() tea.Msg {
					return DeleteRequestMsg{
						Request: request,
					}
				})
			case "enter", "e":
				return tea.Cmd(func() tea.Msg {
					return EditRequestMsg{
						Request: request,
					}
				})
			case "p":
				return tea.Cmd(func() tea.Msg {
					return PreviewRequestMsg{
						Request: request,
					}
				})
			case "r":
				return tea.Cmd(func() tea.Msg {
					return RenameRequestMsg{
						Request: request,
					}
				})
			case "c":
				return tea.Cmd(func() tea.Msg {
					return CopyRequestMsg{
						Request: request,
					}
				})
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
