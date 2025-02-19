package tui

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/help"
	"github.com/susiteemu/startpoint/core/loader"
	"github.com/susiteemu/startpoint/core/tools/paths"
	"github.com/susiteemu/startpoint/tui/overlay"
	profileUI "github.com/susiteemu/startpoint/tui/profile"
	requestUI "github.com/susiteemu/startpoint/tui/request"
	statusbar "github.com/susiteemu/startpoint/tui/statusbar"
	"github.com/susiteemu/startpoint/tui/styles"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
	"github.com/rs/zerolog/log"
)

type ActiveView int

const (
	Requests ActiveView = iota
	Profiles
)

type Model struct {
	active         ActiveView
	requests       requestUI.Model
	profiles       profileUI.Model
	topbar         statusbar.Model
	help           help.Model
	width          int
	height         int
	runningRequest bool
	reloadProfiles bool
	workspace      string
}

type topbarColors struct {
	requestsBg lipgloss.Color
	requestsFg lipgloss.Color
	profilesBg lipgloss.Color
	profilesFg lipgloss.Color
}

func getTopbarColors(activeView ActiveView) topbarColors {

	theme := styles.LoadTheme()

	var requestsBg, requestsFg, profilesBg, profilesFg lipgloss.Color
	switch activeView {
	case Requests:
		requestsBg = theme.TitleBgColor
		requestsFg = theme.TitleFgColor
		profilesBg = theme.StatusbarPrimaryBgColor
		profilesFg = theme.StatusbarPrimaryFgColor
	case Profiles:
		profilesBg = theme.TitleBgColor
		profilesFg = theme.TitleFgColor
		requestsBg = theme.StatusbarPrimaryBgColor
		requestsFg = theme.StatusbarPrimaryFgColor
	}
	return topbarColors{
		requestsBg: requestsBg,
		requestsFg: requestsFg,
		profilesBg: profilesBg,
		profilesFg: profilesFg,
	}
}

func updateTopbar(m *Model) {
	colors := getTopbarColors(m.active)
	requestsItem := statusbar.StatusbarItem{
		Text: "Requests", BackgroundColor: colors.requestsBg, ForegroundColor: colors.requestsFg,
	}

	profilesItem := statusbar.StatusbarItem{
		Text: "Profiles", BackgroundColor: colors.profilesBg, ForegroundColor: colors.profilesFg,
	}

	m.topbar.SetItem(requestsItem, 0)
	m.topbar.SetItem(profilesItem, 1)
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.topbar.SetWidth(msg.Width)
		// request gets -1 for height because mainview has topbar
		m.requests.SetSize(msg.Width, msg.Height-1)
		m.profiles.SetSize(msg.Width, msg.Height-1)
		m.width = msg.Width
		m.height = msg.Height
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case tea.KeyCtrlC.String():
			return m, tea.Quit
		case "ctrl+n":
			if !m.runningRequest {
				switch m.active {
				case Requests:
					m.active = Profiles
				case Profiles:
					if m.reloadProfiles {
						loadedProfiles, err := loader.ReadProfiles(m.workspace)
						if err != nil {
							log.Error().Err(err).Msgf("Failed to read profiles")
						} else {
							requestUI.RefreshProfiles(loadedProfiles)
						}
						m.reloadProfiles = false
					}
					m.active = Requests
				}
				updateTopbar(&m)
				return m, nil
			}
		case "?":
			m.help.ShowAll = !m.help.ShowAll
			return m, nil
		}
	case requestUI.RunRequestMsg:
		m.runningRequest = true
	case requestUI.RunRequestFinishedMsg:
		m.runningRequest = false
	case profileUI.ProfilesChangedMsg:
		m.reloadProfiles = true
	}

	var cmds []tea.Cmd
	var cmd tea.Cmd
	switch m.active {
	case Requests:
		m.requests, cmd = m.requests.Update(msg)
		cmds = append(cmds, cmd)
	case Profiles:
		m.profiles, cmd = m.profiles.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	switch m.active {
	case Requests:
		return renderRequests(m)
	case Profiles:
		return renderProfiles(m)
	default:
		return renderRequests(m)
	}
}

func renderRequests(m Model) string {
	var views []string
	views = append(views, m.topbar.View())
	views = append(views, m.requests.View())

	joined := lipgloss.JoinVertical(
		lipgloss.Top,
		views...,
	)
	if m.help.ShowAll {
		helpModal := style.helpPaneStyle.Render(m.help.View(m.requests.GetHelpKeys()))
		// position at the bottom
		x := (m.width / 2) - (lipgloss.Width(helpModal) / 2)
		y := m.height - lipgloss.Height(helpModal) - 1
		joined = overlay.PlaceOverlay(x, y, helpModal, joined)
	}
	return joined
}

func renderProfiles(m Model) string {
	var views []string
	views = append(views, m.topbar.View())
	views = append(views, m.profiles.View())

	joined := lipgloss.JoinVertical(
		lipgloss.Top,
		views...,
	)
	if m.help.ShowAll {
		helpModal := style.helpPaneStyle.Render(m.help.View(m.profiles.GetHelpKeys()))
		// position at the bottom
		x := (m.width / 2) - (lipgloss.Width(helpModal) / 2)
		y := m.height - lipgloss.Height(helpModal) - 1
		joined = overlay.PlaceOverlay(x, y, helpModal, joined)
	}
	return joined
}

func Start(workspace string, activeView ActiveView) {

	theme := styles.LoadTheme()

	commonStyles := styles.GetCommonStyles(theme)
	InitStyle(theme, commonStyles)

	loadedRequests, err := loader.ReadRequests(workspace)
	if err != nil {
		log.Error().Err(err).Msgf("Failed to read requests")
		fmt.Printf("Failed to read requests %v", err)
		return
	}
	loadedProfiles, err := loader.ReadProfiles(workspace)
	if err != nil {
		log.Error().Err(err).Msgf("Failed to read profiles")
		fmt.Printf("Failed to read profiles %v", err)
		return
	}
	log.Info().Msgf("Loaded %d requests and %d profiles", len(loadedRequests), len(loadedProfiles))

	topbarColors := getTopbarColors(activeView)

	topbarItems := []statusbar.StatusbarItem{
		{Text: "Requests", BackgroundColor: topbarColors.requestsBg, ForegroundColor: topbarColors.requestsFg},
		{Text: "Profiles", BackgroundColor: topbarColors.profilesBg, ForegroundColor: topbarColors.profilesFg},
		{Text: "", BackgroundColor: theme.StatusbarPrimaryBgColor, ForegroundColor: theme.StatusbarPrimaryFgColor},
		{Text: fmt.Sprintf("Workspace: %s", paths.ShortenPath(workspace)), BackgroundColor: theme.StatusbarFourthColBgColor, ForegroundColor: theme.StatusbarSecondaryFgColor},
	}

	tb := statusbar.New(topbarItems, 2, 0)

	help := help.New()
	help.Styles.ShortKey = style.helpKeyStyle
	help.Styles.ShortDesc = style.helpDescStyle
	help.Styles.FullKey = style.helpKeyStyle
	help.Styles.FullDesc = style.helpDescStyle
	help.ShortSeparator = "  "
	m := Model{
		active:    activeView,
		requests:  requestUI.New(loadedRequests, loadedProfiles),
		profiles:  profileUI.New(loadedProfiles),
		topbar:    tb,
		workspace: workspace,
		help:      help,
	}

	output := termenv.NewOutput(os.Stdout)
	originalBackground := output.BackgroundColor()
	originalForeground := output.ForegroundColor()

	output.SetBackgroundColor(output.Color(theme.BgColorStr))
	output.SetForegroundColor(output.Color(theme.TextFgColorStr))

	p := tea.NewProgram(m, tea.WithAltScreen())
	r, err := p.Run()
	output.SetBackgroundColor(originalBackground)
	output.SetForegroundColor(originalForeground)
	output.Reset()
	if err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}

	if m, ok := r.(Model); ok {
		m.requests.HandlePostAction()
	}

}
