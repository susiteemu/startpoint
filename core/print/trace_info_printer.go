package print

import (
	"fmt"
	"startpoint/core/model"

	"github.com/charmbracelet/lipgloss"
)

var traceInfoStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#cdd6f4")).Faint(true)

func SprintTraceInfo(traceInfo model.TraceInfo, pretty bool) (string, error) {

	ti := fmt.Sprintf(`
DNSLookup: %s
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
		return traceInfoStyle.Render(ti), nil
	}

	return ti, nil

}
