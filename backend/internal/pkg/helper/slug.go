package helper

import (
	"regexp"
	"strings"
	"unicode"
)

var nonAlphanumeric = regexp.MustCompile(`[^a-z0-9]+`)

// Slugify converts a string to a URL-friendly slug.
// Example: "Hello World!" → "hello-world"
func Slugify(s string) string {
	s = strings.ToLower(s)
	s = strings.Map(func(r rune) rune {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			return r
		}
		return '-'
	}, s)
	s = regexp.MustCompile(`-+`).ReplaceAllString(s, "-")
	return strings.Trim(s, "-")
}

// GenerateProjectKey derives a Jira-style project key (2-10 uppercase letters)
// from a project name. Examples: "My Project" → "MP", "Backend Service" → "BS".
func GenerateProjectKey(name string) string {
	words := strings.Fields(strings.ToUpper(name))
	if len(words) == 0 {
		return "PROJ"
	}

	var key strings.Builder

	// Use initials if multiple words.
	if len(words) > 1 {
		for _, w := range words {
			for _, r := range w {
				if unicode.IsLetter(r) || unicode.IsDigit(r) {
					key.WriteRune(r)
					break
				}
			}
			if key.Len() >= 6 {
				break
			}
		}
	}

	// Single word: take first 4 letters.
	if key.Len() < 2 {
		key.Reset()
		for _, r := range words[0] {
			if unicode.IsLetter(r) || unicode.IsDigit(r) {
				key.WriteRune(r)
			}
			if key.Len() >= 4 {
				break
			}
		}
	}

	result := nonAlphanumeric.ReplaceAllString(key.String(), "")
	if len(result) < 2 {
		result = "PROJ"
	}
	if len(result) > 10 {
		result = result[:10]
	}
	return result
}
