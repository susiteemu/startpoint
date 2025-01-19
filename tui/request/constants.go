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
	CreateRequest        = "CReq"
	CreateSimpleRequest  = "CSmplReq"
	CreateComplexRequest = "CCmplxReq"
	DeleteRequest        = "DReq"
	EditRequest          = "EReq"
	PrintRequest         = "PReq"
	PrintFailedRequest   = "PEReq"
	RenameRequest        = "RnReq"
	CopyRequest          = "CpReq"
)

const (
	YamlTemplate = `# Possible request to call _before_ this one
prev_req:
# Request url, may contain template variables in a form of {var}
url:
# HTTP method
method: GET
# HTTP headers as key-val list, e.g. X-Foo-Bar: SomeValue
headers:
# Request body, e.g.
# {
#    "id": 1,
#    "name": "Jane">
# }
body: >
`

	StarlarkTemplate = `"""
prev_req: <call other request before this>
doc:url: <your url for display>
doc:method: GET
"""
# insert contents of your script here, for more see https://github.com/google/starlark-go/blob/master/doc/spec.md
# Request url
url = ""
# HTTP method
method = "GET"
# HTTP headers, e.g. { "X-Foo": "bar", "X-Foos": [ "Bar1", "Bar2" ] }
headers = {}
# Request body, e.g. { "id": 1, "people": [ {"name": "Joe"}, {"name": "Jane"}, ] }
body = {}
`

	LuaTemplate = `--[[
prev_req: <call other request before this>
doc:url: <your url for display>
doc:method: GET
--]]
return {
	-- Request url
	url = "",
	-- HTTP method
	method = "GET",
	-- HTTP headers, e.g. { ["X-Foo"]="Bar", ["X-Foos"]={ "Bar1", "Bar2" } }
	headers = {},
	-- Request body, e.g. { id=1, people={ {name="Joe"}, {name="Jane"} } }
	body = {}
}`
)
