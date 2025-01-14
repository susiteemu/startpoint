package ansi

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/rs/zerolog/log"
)

// ANSI escape code regex
var ansiRegex = regexp.MustCompile(`\x1b\[(\d+(;\d+)*)m`)

// ColorState holds the current foreground and background colors
type ColorState struct {
	State string
}

func ParseANSI(input string, index int) (ColorState, error) {
	var plainText strings.Builder
	colorMap := make(map[int]ColorState)

	matches := ansiRegex.FindAllStringIndex(input, -1)
	lastPos := 0
	statePositions := []int{}
	for _, match := range matches {
		start, end := match[0], match[1]
		plainText.WriteString(input[lastPos:start])

		seq := input[start:end]

		colorMap[plainText.Len()] = ColorState{
			State: seq,
		}
		statePositions = append(statePositions, plainText.Len())
		lastPos = end
	}

	// Append remaining text
	plainText.WriteString(input[lastPos:])
	finalText := plainText.String()

	if index >= len(finalText) {
		return ColorState{}, fmt.Errorf("index out of bounds")
	}

	log.Debug().Msgf("ColorMap: %v", colorMap)
	var activeState ColorState
	for _, i := range statePositions {
		if i > index {
			break
		}
		activeState = colorMap[i]
	}

	return activeState, nil
}
