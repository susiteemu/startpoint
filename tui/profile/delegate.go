package profileui

import (
	"fmt"
	"github.com/susiteemu/startpoint/tui/messages"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
)

type embeddedKeyMap struct {
	Select  key.Binding
	Cancel  key.Binding
	Preview key.Binding
}

var embeddedKeys = embeddedKeyMap{
	Select: key.NewBinding(
		key.WithKeys(tea.KeyEnter.String()),
		key.WithHelp(tea.KeyEnter.String(), "select"),
	),
	Cancel: key.NewBinding(
		key.WithKeys(tea.KeyEsc.String()),
		key.WithHelp(tea.KeyEsc.String(), "cancel"),
	),
}

var editKeys = []key.Binding{
	key.NewBinding(
		key.WithKeys("a"),
		key.WithHelp("a", "add"),
	),
	key.NewBinding(
		key.WithKeys("d"),
		key.WithHelp("d", "delete"),
	),
	key.NewBinding(
		key.WithKeys("c"),
		key.WithHelp("c", "copy profile"),
	),
	key.NewBinding(
		key.WithKeys("r"),
		key.WithHelp("r", "rename"),
	),
	key.NewBinding(
		key.WithKeys("e"),
		key.WithHelp("e", "edit"),
	),
	key.NewBinding(
		key.WithKeys("p"),
		key.WithHelp("p", "preview"),
	),
	key.NewBinding(
		key.WithKeys("ctrl+n"),
		key.WithHelp("ctrl+n", "switch to Requests"),
	),
}

type profileItemDelegate struct {
	normalTitle   lipgloss.Style
	normalDesc    lipgloss.Style
	selectedTitle lipgloss.Style
	selectedDesc  lipgloss.Style
	list.DefaultDelegate
}

func (d profileItemDelegate) Height() int  { return 2 }
func (d profileItemDelegate) Spacing() int { return 1 }
func (d profileItemDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd {
	if d.UpdateFunc == nil {
		return nil
	}
	return d.UpdateFunc(msg, m)
}
func (d profileItemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(Profile)
	if !ok {
		return
	}
	title := i.Name
	desc := fmt.Sprintf("Vars: %d", i.Variables)
	textwidth := m.Width() - d.normalTitle.GetPaddingLeft() - d.normalTitle.GetPaddingRight()
	title = ansi.Truncate(title, textwidth, "...")

	titleFn := d.normalTitle.Render
	descFn := d.normalDesc.Render
	if index == m.Index() {
		titleFn = d.selectedTitle.Render
		descFn = d.selectedDesc.Render
	}

	content := []string{}
	content = append(content, titleFn(title))
	content = append(content, descFn(desc))
	fmt.Fprint(w, strings.Join(content, "\n"))
}

type embeddedProfileItemDelegate struct {
	normalTitle   lipgloss.Style
	normalDesc    lipgloss.Style
	selectedTitle lipgloss.Style
	selectedDesc  lipgloss.Style
	list.DefaultDelegate
}

func (d embeddedProfileItemDelegate) Height() int  { return 1 }
func (d embeddedProfileItemDelegate) Spacing() int { return 1 }
func (d embeddedProfileItemDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd {
	if d.UpdateFunc == nil {
		return nil
	}
	return d.UpdateFunc(msg, m)
}
func (d embeddedProfileItemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(Profile)
	if !ok {
		return
	}
	title := i.Name
	textwidth := m.Width() - d.normalTitle.GetPaddingLeft() - d.normalTitle.GetPaddingRight()
	title = ansi.Truncate(title, textwidth, "...")

	titleFn := d.normalTitle.Render
	if index == m.Index() {
		titleFn = d.selectedTitle.Render
	}

	content := []string{}
	content = append(content, titleFn(title))
	fmt.Fprint(w, strings.Join(content, "\n"))
}

func newEmbeddedDelegate() list.ItemDelegate {
	d := embeddedProfileItemDelegate{
		normalTitle: lipgloss.NewStyle().Foreground(style.listItemTitleColor).Padding(0, 0, 0, 2),
		normalDesc:  lipgloss.NewStyle().Foreground(style.listItemTitleColor).Padding(0, 0, 0, 2),
		selectedTitle: lipgloss.NewStyle().Border(lipgloss.NormalBorder(), false, false, false, true).
			BorderForeground(style.listItemTitleColor).
			Foreground(style.listItemTitleColor).
			Padding(0, 0, 0, 1),
		selectedDesc: lipgloss.NewStyle().Border(lipgloss.NormalBorder(), false, false, false, true).
			BorderForeground(style.listItemTitleColor).
			Foreground(style.listItemTitleColor).
			Padding(0, 0, 0, 1),
	}
	d.UpdateFunc = func(msg tea.Msg, m *list.Model) tea.Cmd {
		var profile Profile
		var profileSelected bool
		if i, ok := m.SelectedItem().(Profile); ok {
			profile = i
			profileSelected = true
		} else {
			profileSelected = false
		}

		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch keypress := msg.String(); keypress {
			case tea.KeyEnter.String():
				if profileSelected {
					return tea.Cmd(func() tea.Msg {
						return ProfileSelectedMsg{
							Profile: profile,
						}
					})
				}
			case tea.KeyEsc.String():
				return tea.Cmd(func() tea.Msg {
					return ProfileSelectCancelledMsg{}
				})
			}
		}
		return nil
	}
	return d
}
func newNormalDelegate() list.ItemDelegate {
	d := profileItemDelegate{
		normalTitle: lipgloss.NewStyle().Foreground(style.listItemTitleColor).Padding(0, 0, 0, 2),
		normalDesc:  lipgloss.NewStyle().Foreground(style.listItemTitleColor).Padding(0, 0, 0, 2),
		selectedTitle: lipgloss.NewStyle().Border(lipgloss.NormalBorder(), false, false, false, true).
			BorderForeground(style.listItemTitleColor).
			Foreground(style.listItemTitleColor).
			Padding(0, 0, 0, 1),
		selectedDesc: lipgloss.NewStyle().Border(lipgloss.NormalBorder(), false, false, false, true).
			BorderForeground(style.listItemTitleColor).
			Foreground(style.listItemTitleColor).
			Padding(0, 0, 0, 1),
	}
	d.UpdateFunc = func(msg tea.Msg, m *list.Model) tea.Cmd {
		var profile Profile
		var profileSelected bool
		if i, ok := m.SelectedItem().(Profile); ok {
			profile = i
			profileSelected = true
		} else {
			profileSelected = false
		}

		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch keypress := msg.String(); keypress {
			case "a":
				return tea.Cmd(func() tea.Msg {
					return CreateProfileMsg{}
				})
			case "d":
				if profileSelected {
					if profile.Name == "default" {
						return messages.CreateStatusMsg("You can't delete default profile")
					} else {
						return tea.Cmd(func() tea.Msg {
							return DeleteProfileMsg{
								Profile: profile,
							}
						})
					}
				}
			case "c":
				if profileSelected {
					return tea.Cmd(func() tea.Msg {
						return CopyProfileMsg{
							Profile: profile,
						}
					})
				}
			case "r":
				if profileSelected {
					if profile.Name == "default" {
						return messages.CreateStatusMsg("You can't rename default profile")
					} else {
						return tea.Cmd(func() tea.Msg {
							return RenameProfileMsg{
								Profile: profile,
							}
						})
					}
				}
			case "e":
				if profileSelected {
					return tea.Cmd(func() tea.Msg {
						return EditProfileMsg{
							Profile: profile,
						}
					})
				}
			case "p":
				if profileSelected {
					return tea.Cmd(func() tea.Msg {
						return PreviewProfileMsg{
							Profile: profile,
						}
					})
				}
			}
		}
		return nil
	}

	d.ShortHelpFunc = func() []key.Binding {
		return editKeys
	}

	d.FullHelpFunc = func() [][]key.Binding {
		return [][]key.Binding{editKeys}
	}

	return d
}
