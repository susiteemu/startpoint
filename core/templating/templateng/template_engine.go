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

func DiscoverTemplateVariables(s string) []string {
	if len(s) == 0 {
		return []string{}
	}

	pattern := regexp.MustCompile(`\{([^}]*)\}`)
	matches := pattern.FindAllStringSubmatch(s, -1)
	var results []string
	for _, match := range matches {
		if len(match) > 1 {
			results = append(results, match[1])
		}
	}

	return results
}

func ProcessTemplateVariableRecursively(s string, all map[string]string) string {
	templateVars := DiscoverTemplateVariables(s)
	if len(templateVars) > 0 {
		for _, tv := range templateVars {
			matchingVariable := all[tv]
			s, _ = ProcessTemplateVariable(s, tv, matchingVariable)
			return ProcessTemplateVariableRecursively(s, all)
		}
	}
	return s
}
