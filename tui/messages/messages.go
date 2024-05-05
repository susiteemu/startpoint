package messages

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type StatusMessage string

func CreateStatusMsg(msg string) tea.Cmd {
	return tea.Cmd(func() tea.Msg {
		nowTime := time.Now().Format("15:04:05")
		return StatusMessage(fmt.Sprintf("%s %s", nowTime, msg))
	})
}
