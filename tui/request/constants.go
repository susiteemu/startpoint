package requestui

const (
	List ActiveView = iota
	Prompt
	Keyprompt
	Duplicate
	Preview
	Stopwatch
	Profiles
)

const (
	Select Mode = iota
	Edit
)

const (
	CreateRequestLabel = "Choose a name for your request. Make it filename compatible and unique within this workspace. After choosing \"ok\" your $EDITOR will open and you will be able to write the contents of the request. Remember to quit your editor window to return back."
	RenameRequestLabel = "Rename your request."
	CopyRequestLabel   = "Choose name for your request."
)

const (
	CreateSimpleRequest  = "CSmplReq"
	CreateComplexRequest = "CCmplxReq"
	EditRequest          = "EReq"
	PrintRequest         = "PReq"
	RenameRequest        = "RnReq"
	CopyRequest          = "CpReq"
)
