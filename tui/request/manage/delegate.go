package managetui

import (
	"fmt"
	"goful/core/client/validator"
	"time"

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
		key.WithKeys("q", tea.KeyCtrlC.String()),
		key.WithHelp("q/ctrl+c", "quit"),
	),
}

var editModeKeys = []key.Binding{
	key.NewBinding(
		key.WithKeys("a"),
		key.WithHelp("a", "add simple"),
	),
	key.NewBinding(
		key.WithKeys("A"),
		key.WithHelp("A", "add complex"),
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
	key.NewBinding(
		key.WithKeys("q", tea.KeyCtrlC.String()),
		key.WithHelp("q/ctrl+c", "quit"),
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
						nowTime := time.Now().Format("15:04:05")
						return StatusMessage(fmt.Sprintf("%s Error! Invalid request", nowTime))
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
					return tea.Cmd(func() tea.Msg {
						nowTime := time.Now().Format("15:04:05")
						return StatusMessage(fmt.Sprintf("%s Deleted %s", nowTime, request.Title()))
					})
				} else {
					return tea.Cmd(func() tea.Msg {
						nowTime := time.Now().Format("15:04:05")
						return StatusMessage(fmt.Sprintf("%s Failed to delete %s", nowTime, request.Title()))
					})
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
