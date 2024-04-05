package print

import (
	"errors"
	"fmt"
	"goful/core/model"

	"github.com/charmbracelet/lipgloss"
)

var (
	style200   = lipgloss.NewStyle().Foreground(lipgloss.Color("#a6e3a1"))
	style300   = lipgloss.NewStyle().Foreground(lipgloss.Color("#f9e2af"))
	style400   = lipgloss.NewStyle().Foreground(lipgloss.Color("#f38ba8"))
	style500   = lipgloss.NewStyle().Foreground(lipgloss.Color("#f38ba8"))
	styleProto = lipgloss.NewStyle().Foreground(lipgloss.Color("#89b4fa"))
)

func SprintStatus(resp *model.Response, pretty bool) (string, error) {
	if resp == nil {
		return "", errors.New("Response must not be nil!")
	}

	protoStyle := lipgloss.NewStyle()
	statusStyle := lipgloss.NewStyle()
	if pretty {
		protoStyle = styleProto
		if resp.StatusCode < 300 {
			statusStyle = style200
		} else if resp.StatusCode < 400 {
			statusStyle = style300
		} else if resp.StatusCode < 500 {
			statusStyle = style400
		} else {
			statusStyle = style500
		}
	}

	return fmt.Sprintf("%v %v", protoStyle.Render(resp.Proto), statusStyle.Render(resp.Status)), nil
}
