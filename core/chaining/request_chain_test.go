package requestchain

import (
	"goful/core/model"
	"testing"
)

var testData = []*model.RequestMold{
	{Yaml: &model.YamlRequest{Name: "Req1", PrevReq: ""}},
	{Yaml: &model.YamlRequest{Name: "Req2", PrevReq: ""}},
	{Yaml: &model.YamlRequest{Name: "Req3", PrevReq: "Req5"}},
	{Yaml: &model.YamlRequest{Name: "Req4", PrevReq: "Req3"}},
	{Yaml: &model.YamlRequest{Name: "Req5", PrevReq: "Req2"}},
}

func TestResolveRequestChain(t *testing.T) {

	gotMolds := ResolveRequestChain(testData[3], testData)
	wantedNames := []string{
		"Req2", "Req5", "Req3", "Req4",
	}
	assertResult(gotMolds, wantedNames, t)

	gotMolds = ResolveRequestChain(testData[0], testData)
	wantedNames = []string{
		"Req1",
	}
	assertResult(gotMolds, wantedNames, t)

	gotMolds = ResolveRequestChain(testData[4], testData)
	wantedNames = []string{
		"Req2", "Req5",
	}
	assertResult(gotMolds, wantedNames, t)
}

func assertResult(gotMolds []*model.RequestMold, wantedNames []string, t *testing.T) {
	if len(gotMolds) != len(wantedNames) {
		t.Errorf("got %d, wanted %d", len(gotMolds), len(wantedNames))
		return
	}
	for i := 0; i < len(gotMolds); i++ {
		got := gotMolds[i].Name()
		wanted := wantedNames[i]
		if got != wanted {
			t.Errorf("got %s, wanted %s", got, wanted)
			return
		}
	}
}
