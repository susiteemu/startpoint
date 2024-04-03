package statusbarui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/truncate"
)

const Height = 1

var (
	statusBarStyle         = lipgloss.NewStyle()
	statusBarItemBaseStyle = lipgloss.NewStyle().PaddingLeft(1).PaddingRight(1)
)

type StatusbarItem struct {
	Text            string
	ForegroundColor lipgloss.Color
	BackgroundColor lipgloss.Color
}

type Model struct {
	items         []StatusbarItem
	itemStyles    []lipgloss.Style
	mainItemIndex int
	width         int
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	return m, nil
}

func (m Model) View() string {
	var otherItemWidths int
	var renders []string
	for i := 0; i < len(m.items); i++ {
		item := m.items[i]
		if i != m.mainItemIndex {
			renderedItem := m.itemStyles[i].Render(item.Text)
			renders = append(renders, renderedItem)
			otherItemWidths += lipgloss.Width(renderedItem)
		}
	}
	if m.mainItemIndex >= 0 && m.mainItemIndex < len(m.items) {
		renders = append(renders[:m.mainItemIndex+1], renders[m.mainItemIndex:]...)
		truncatedText := truncate.StringWithTail(m.items[m.mainItemIndex].Text, uint(m.width-otherItemWidths-3), "...")
		renders[m.mainItemIndex] = m.itemStyles[m.mainItemIndex].Copy().Width(m.width - otherItemWidths).Render(truncatedText)
	}

	bar := lipgloss.JoinHorizontal(lipgloss.Top,
		renders...,
	)

	return statusBarStyle.Width(m.width).Height(Height).Render(bar)
}

func (m *Model) ChangeText(text string, index int) {
	if index >= 0 && index < len(m.items) {
		item := m.items[index]
		item.Text = text
		m.items[index] = item
	}
}

func (m *Model) SetItem(item StatusbarItem, index int) {

	style := statusBarItemBaseStyle.Copy().Foreground(item.ForegroundColor).Background(item.BackgroundColor)
	if index >= 0 && index < len(m.items) {
		m.itemStyles[index] = style
		m.items[index] = item
	} else {
		m.items = append(m.items, item)
		m.itemStyles = append(m.itemStyles, style)
	}

}

func (m *Model) SetWidth(width int) {
	m.width = width
}

func New(items []StatusbarItem, mainItemIndex int, width int) Model {

	var styles []lipgloss.Style
	for _, v := range items {
		styles = append(styles, statusBarItemBaseStyle.Copy().Foreground(v.ForegroundColor).Background(v.BackgroundColor).PaddingLeft(1))
	}

	m := Model{
		items:         items,
		mainItemIndex: mainItemIndex,
		width:         width,
		itemStyles:    styles,
	}

	return m
}
