package managetui

import (
	"goful/core/client/validator"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

var selectModeKeys = []key.Binding{
	key.NewBinding(
		key.WithKeys("p"),
		key.WithHelp("p", "Preview"),
	),
	key.NewBinding(
		key.WithKeys(tea.KeyEnter.String()),
		key.WithHelp(tea.KeyEnter.String(), "Run request"),
	),
	key.NewBinding(
		key.WithKeys("i"),
		key.WithHelp("i", "Edit mode"),
	),
}

var editModeKeys = []key.Binding{
	key.NewBinding(
		key.WithKeys("a"),
		key.WithHelp("a", "Add simple"),
	),
	key.NewBinding(
		key.WithKeys("A"),
		key.WithHelp("A", "Add complex"),
	),
	key.NewBinding(
		key.WithKeys("p"),
		key.WithHelp("p", "Preview"),
	),
	key.NewBinding(
		key.WithKeys("x"),
		key.WithHelp("x", "Delete"),
	),
	key.NewBinding(
		key.WithKeys(tea.KeyEnter.String(), "e"),
		key.WithHelp(tea.KeyEnter.String()+"/e", "Edit"),
	),
	key.NewBinding(
		key.WithKeys(tea.KeyEsc.String()),
		key.WithHelp(tea.KeyEsc.String(), "View mode"),
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
			case "enter", "e":
				if !validator.IsValidUrl(request.Url) || !validator.IsValidMethod(request.Method) {
					return m.NewStatusMessage(statusMessageStyle.Render("\ue654 Invalid request."))
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
			case "x":
				deleted := request.Mold.DeleteFromFS()
				if deleted {
					index := m.Index()
					m.RemoveItem(index)
					return m.NewStatusMessage("Deleted " + request.Title())
				} else {
					return m.NewStatusMessage("Failed to delete " + request.Title())
				}
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
