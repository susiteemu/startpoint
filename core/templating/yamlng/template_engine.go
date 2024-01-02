package yamlng

import (
	"fmt"
	"strings"
)

func ProcessTemplateVariables(s []string, variableName string, variableValue interface{}) []string {
	var processed []string
	for _, str := range s {
		processed = append(processed, ProcessTemplateVariable(str, variableName, variableValue))
	}
	return processed
}

func ProcessTemplateVariable(s string, variableName string, variableValue interface{}) string {
	templateVariable := fmt.Sprintf("{%s}", variableName)
	if strings.Contains(s, templateVariable) {
		return strings.ReplaceAll(s, templateVariable, variableValue.(string))
	}
	return s
}
