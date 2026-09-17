package generate

import (
	"strings"
	"unicode"
)

type Name struct {
	Struct string
	Package string
}

func NewName(raw string) Name {
	words := splitWords(raw)

	var structName strings.Builder
	var packageName strings.Builder
	for _, word := range words {
		structName.WriteString(capitalize(word))
		packageName.WriteString(strings.ToLower(word))
	}

	return Name{
		Struct:  structName.String(),
		Package: packageName.String(),
	}
}

func splitWords(raw string) []string {
	var words []string
	var current []rune

	flush := func() {
		if len(current) > 0 {
			words = append(words, string(current))
			current = nil
		}
	}

	runes := []rune(raw)
	for i, r := range runes {
		switch {
		case r == ' ' || r == '_' || r == '-':
			flush()
		case unicode.IsUpper(r) && i > 0 && !isWordBreak(runes[i-1]):
			flush()
			current = append(current, r)
		default:
			current = append(current, r)
		}
	}
	flush()

	return words
}

func isWordBreak(r rune) bool {
	return r == ' ' || r == '_' || r == '-' || unicode.IsUpper(r)
}

func capitalize(word string) string {
	if word == "" {
		return word
	}

	runes := []rune(word)
	return string(unicode.ToUpper(runes[0])) + strings.ToLower(string(runes[1:]))
}
