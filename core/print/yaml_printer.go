package print

import (
	"bytes"
)

func SprintYaml(rawYaml string) (string, error) {
	buf := new(bytes.Buffer)

	lexer := resolveLexer("text/x-yaml")
	style := resolveStyle()
	formatter := resolveFormatter()
	iterator, err := lexer.Tokenise(nil, rawYaml)
	if err != nil {
		return "", err
	}
	err = formatter.Format(buf, style, iterator)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}
