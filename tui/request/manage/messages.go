package managetui

type RunRequestMsg struct {
	Request Request
}

type EditRequestMsg struct {
	Request Request
}

type PreviewRequestMsg struct {
	Request Request
}

type RequestFinishedMsg string

type StatusMessage string

type RenameRequestMsg struct {
	Request Request
}
