package print

import (
	"errors"
	"fmt"
	"startpoint/core/model"
	"startpoint/tui/styles"

	"github.com/charmbracelet/lipgloss"
)

func SprintStatus(resp *model.Response, pretty bool) (string, error) {
	if resp == nil {
		return "", errors.New("Response must not be nil!")
	}
	theme := styles.GetTheme()
	protoStyle := lipgloss.NewStyle()
	statusStyle := lipgloss.NewStyle()
	if pretty {
		protoStyle = lipgloss.NewStyle().Foreground(theme.ResponseProtoFgColor)
		if resp.StatusCode < 300 {
			statusStyle = lipgloss.NewStyle().Foreground(theme.ResponseStatus200FgColor)
		} else if resp.StatusCode < 400 {
			statusStyle = lipgloss.NewStyle().Foreground(theme.ResponseStatus300FgColor)
		} else if resp.StatusCode < 500 {
			statusStyle = lipgloss.NewStyle().Foreground(theme.ResponseStatus400FgColor)
		} else {
			statusStyle = lipgloss.NewStyle().Foreground(theme.ResponseStatus500FgColor)
		}
	}

	return fmt.Sprintf("%v %v", protoStyle.Render(resp.Proto), statusStyle.Render(resp.Status)), nil
}
