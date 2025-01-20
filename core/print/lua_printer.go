package print

import (
	"bytes"
)

func SprintLua(rawYaml string) (string, error) {
	buf := new(bytes.Buffer)

	lexer := resolveLexer("text/x-lua")
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
