package requestui

import (
	keyprompt "startpoint/tui/keyprompt"
)

type RunRequestMsg struct {
	Request Request
}

type CreateRequestMsg struct {
	Simple bool
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

type DeleteRequestMsg struct {
	Request Request
}

type PreviewRequestMsg struct {
	Request Request
}

type RunRequestFinishedMsg string

type RenameRequestMsg struct {
	Request Request
}

type CopyRequestMsg struct {
	Request Request
}

type ActivateProfile struct{}

type ShowKeyprompt struct {
	Label   string
	Entries []keyprompt.KeypromptEntry
}
