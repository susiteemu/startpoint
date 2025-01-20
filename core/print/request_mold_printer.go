package print

import (
	"fmt"
	"startpoint/core/model"
)

func SprintRequestMold(m *model.RequestMold) (string, error) {
	if m == nil {
		return "", fmt.Errorf("RequestMold is nil")
	}

	switch m.Type {
	case model.CONTENT_TYPE_YAML:
		return SprintYaml(m.Raw())
	case model.CONTENT_TYPE_STARLARK:
		return SprintStarlark(m.Raw())
	case model.CONTENT_TYPE_LUA:
		return SprintLua(m.Raw())
	}
	return "", fmt.Errorf("Unknown RequestMold type %s", m.Type)
}
