package print

import (
	"fmt"
	"github.com/susiteemu/startpoint/core/model"
	"github.com/susiteemu/startpoint/tui/styles"

	"github.com/charmbracelet/lipgloss"
)

func SprintTraceInfo(traceInfo model.TraceInfo, pretty bool) (string, error) {
	ti := fmt.Sprintf(`DNSLookup: %s
ConnTime: %s
TCPConnTime: %s
TLSHandshake: %s
ServerTime: %s
ResponseTime: %s
TotalTime: %s
IsConnReused: %v
IsConnWasIdle: %v
ConnIdleTime: %s
RequestAttempt: %d
RemoteAddr: %s
	`,
		traceInfo.DNSLookup, traceInfo.ConnTime, traceInfo.TCPConnTime, traceInfo.TLSHandshake, traceInfo.ServerTime, traceInfo.ResponseTime, traceInfo.TotalTime, traceInfo.IsConnReused, traceInfo.IsConnWasIdle, traceInfo.ConnIdleTime, traceInfo.RequestAttempt, traceInfo.RemoteAddr,
	)

	if pretty {
		theme := styles.GetTheme()
		traceInfoStyle := lipgloss.NewStyle().Foreground(theme.TextFgColor).Faint(true)
		return traceInfoStyle.Render(ti), nil
	}

	return ti, nil

}
