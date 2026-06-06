package helper

import (
	"regexp"
	"strings"
)

var (
	reHTMLTag    = regexp.MustCompile(`<[^>]*>`)
	reHTMLEntity = regexp.MustCompile(`&[a-zA-Z0-9#]+;`)
	reScriptTag  = regexp.MustCompile(`(?is)<script[^>]*>.*?</script>`)
	reStyleTag   = regexp.MustCompile(`(?is)<style[^>]*>.*?</style>`)
)

// StripHTML removes all HTML tags and entities from s.
func StripHTML(s string) string {
	s = reScriptTag.ReplaceAllString(s, "")
	s = reStyleTag.ReplaceAllString(s, "")
	s = reHTMLTag.ReplaceAllString(s, "")
	s = reHTMLEntity.ReplaceAllString(s, "")
	return strings.TrimSpace(s)
}

// SanitizePlainText removes leading/trailing whitespace and collapses
// internal runs of whitespace to a single space.
func SanitizePlainText(s string) string {
	s = strings.TrimSpace(s)
	return reMultiSpace.ReplaceAllString(s, " ")
}
