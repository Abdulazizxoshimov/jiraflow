package helper

import (
	"regexp"
	"strings"
)

var (
	reMDHeading   = regexp.MustCompile(`(?m)^#{1,6}\s+`)
	reMDEmphasis  = regexp.MustCompile(`[*_]{1,3}([^*_]+)[*_]{1,3}`)
	reMDCode      = regexp.MustCompile("`{1,3}[^`]*`{1,3}")
	reMDLink      = regexp.MustCompile(`\[([^\]]+)\]\([^)]+\)`)
	reMDImage     = regexp.MustCompile(`!\[[^\]]*\]\([^)]+\)`)
	reMDBlockQuote = regexp.MustCompile(`(?m)^>\s+`)
	reMDHR        = regexp.MustCompile(`(?m)^[-*_]{3,}\s*$`)
	reMDListItem  = regexp.MustCompile(`(?m)^[\s]*[-*+]\s+|^[\s]*\d+\.\s+`)
	reMultiSpace  = regexp.MustCompile(`\s{2,}`)
)

// StripMarkdown converts a Markdown string to plain text.
func StripMarkdown(s string) string {
	s = reMDImage.ReplaceAllString(s, "")
	s = reMDLink.ReplaceAllString(s, "$1")
	s = reMDCode.ReplaceAllString(s, "")
	s = reMDHeading.ReplaceAllString(s, "")
	s = reMDEmphasis.ReplaceAllString(s, "$1")
	s = reMDBlockQuote.ReplaceAllString(s, "")
	s = reMDHR.ReplaceAllString(s, "")
	s = reMDListItem.ReplaceAllString(s, "")
	s = reMultiSpace.ReplaceAllString(s, " ")
	return strings.TrimSpace(s)
}

// MarkdownExcerpt strips Markdown and truncates to maxLen runes.
func MarkdownExcerpt(s string, maxLen int) string {
	return TruncateString(StripMarkdown(s), maxLen)
}
