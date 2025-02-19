package print

import (
	"errors"
	"fmt"

	"github.com/susiteemu/startpoint/core/model"
	"github.com/susiteemu/startpoint/tui/styles"

	"github.com/charmbracelet/lipgloss"
)

func SprintStatus(resp *model.Response, pretty bool) (string, string, error) {
	if resp == nil {
		return "", "", errors.New("Response must not be nil!")
	}
	theme := styles.LoadTheme()
	protoStyle := lipgloss.NewStyle()
	statusStyle := lipgloss.NewStyle()
	status := fmt.Sprintf("%v %v", resp.Proto, resp.Status)
	prettyStatus := ""
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
		prettyStatus = fmt.Sprintf("%v %v", protoStyle.Render(resp.Proto), statusStyle.Render(resp.Status))
	}

	return status, prettyStatus, nil
}
