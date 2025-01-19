package requestui

import (
	keyprompt "startpoint/tui/keyprompt"
)

type RunRequestMsg struct {
	Request Request
}

type CreateRequestMsg struct {
	Type string
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

type DeleteRequestConfirmedMsg struct {
	Request Request
}

type PreviewRequestMsg struct {
	Request Request
}

type RunRequestFinishedMsg string
type RunRequestFinishedWithFailureMsg string

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
	Type    string
	Payload interface{}
}
