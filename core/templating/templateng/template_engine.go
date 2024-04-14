package templateng

import (
	"fmt"
	"strings"
)

func ProcessTemplateVariables(s []string, variableName string, variableValue interface{}) []string {
	var processedValues []string
	for _, str := range s {
		processed, _ := ProcessTemplateVariable(str, variableName, variableValue)
		processedValues = append(processedValues, processed)
	}
	return processedValues
}

func ProcessTemplateVariable(s string, variableName string, variableValue interface{}) (string, bool) {
	templateVariable := fmt.Sprintf("{%s}", variableName)
	if strings.Contains(s, templateVariable) {
		return strings.ReplaceAll(s, templateVariable, variableValue.(string)), true
	}
	return s, false
}
