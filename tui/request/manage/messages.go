package managetui

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type RunRequestMsg struct {
	Request Request
}

type CreateRequestFinishedMsg struct {
	root     string
	filename string
	err      error
}

type EditRequestMsg struct {
	Request Request
}

type EditRequestFinishedMsg struct {
	Request Request
	err     error
}

type PreviewRequestMsg struct {
	Request Request
}

type RequestFinishedMsg string

type StatusMessage string

type RenameRequestMsg struct {
	Request Request
}

type CopyRequestMsg struct {
	Request Request
}

func createStatusMsg(msg string) tea.Cmd {
	return tea.Cmd(func() tea.Msg {
		nowTime := time.Now().Format("15:04:05")
		return StatusMessage(fmt.Sprintf("%s %s", nowTime, msg))
	})
}
