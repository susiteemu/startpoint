package templateng

import (
	"fmt"
	"regexp"
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
	templateVariable := regexp.MustCompile(fmt.Sprintf(`{\s*%s\s*}`, variableName))
	return templateVariable.ReplaceAllString(s, variableValue.(string)), true
}
