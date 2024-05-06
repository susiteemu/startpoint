package profileui

import (
	"startpoint/tui/messages"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type embeddedKeyMap struct {
	Select key.Binding
	Cancel key.Binding
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

func newEmbeddedDelegate() list.DefaultDelegate {
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
			case tea.KeyEnter.String():
				return tea.Cmd(func() tea.Msg {
					return ProfileSelectedMsg{
						Profile: profile,
					}
				})
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
func newNormalDelegate() list.DefaultDelegate {
	d := newBaseDelegate()
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
