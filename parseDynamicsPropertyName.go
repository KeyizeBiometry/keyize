package keyize

import (
	"errors"
	"regexp"
	"unicode/utf8"
)

var dynamicsPropertyNameRegex *regexp.Regexp = regexp.MustCompile(`^(DD|D|UD)\.(.|[^a])(?:\.(.|[^a]))?$`)

var kindCodeKindMap = map[string]DynamicsPropertyKind{
	"DD": DownDown,
	"UD": UpDown,
	"D":  Dwell,
}

func ParseDynamicsPropertyName(n string) (*DynamicsProperty, error) {
	components := dynamicsPropertyNameRegex.FindStringSubmatch(n)

	if components == nil {
		return nil, errors.New("failed to parse name '" + n + "'")
	}

	rune1, _ := utf8.DecodeRuneInString(components[2])

	// In a properly formatted Dwell name, there may be no components[3]
	var rune2 rune

	if len(components) >= 4 {
		rune2, _ = utf8.DecodeRuneInString(components[3])
	}

	return &DynamicsProperty{
		Kind:  kindCodeKindMap[components[1]],
		KeyA:  rune1,
		KeyB:  rune2,
		Value: 0,
	}, nil
}
