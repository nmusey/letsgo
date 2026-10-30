package names

import (
	"strings"
	"unicode"
)

// Return an idiomatic struct name
func StructName(name string) string {
	builder := strings.Builder{}

	// Split on anything that isn't a letter or digit, and before an uppercase letter following a lowercase one
	// so that existing camel case (UserProfile) is preserved
	var words []string
	current := []rune{}
	var previous rune
	for _, r := range name {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) {
			if len(current) > 0 {
				words = append(words, string(current))
				current = current[:0]
			}
			previous = 0
			continue
		}

		if unicode.IsUpper(r) && unicode.IsLower(previous) && len(current) > 0 {
			words = append(words, string(current))
			current = current[:0]
		}

		current = append(current, r)
		previous = r
	}
	if len(current) > 0 {
		words = append(words, string(current))
	}

	for _, word := range words {
		runes := []rune(word)
		builder.WriteRune(unicode.ToUpper(runes[0]))
		builder.WriteString(string(runes[1:]))
	}

	return builder.String()
}

// Return an idiomatic package name
func PackageName(name string) string {
	builder := strings.Builder{}

	for _, r := range name {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			builder.WriteRune(unicode.ToLower(r))
		}
	}

	return builder.String()
}
