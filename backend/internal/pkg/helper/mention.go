package helper

import "regexp"

// issueKeyRe — [PROJ-42] yoki [ABC-1] formatdagi mention'larni topadi.
var issueKeyRe = regexp.MustCompile(`\[([A-Z][A-Z0-9]+-\d+)\]`)

// ExtractIssueKeys — matn ichidagi barcha issue key'larni qaytaradi (masalan ["PROJ-42", "WEB-1"]).
func ExtractIssueKeys(text string) []string {
	matches := issueKeyRe.FindAllStringSubmatch(text, -1)
	seen := make(map[string]struct{})
	var keys []string
	for _, m := range matches {
		k := m[1]
		if _, ok := seen[k]; !ok {
			seen[k] = struct{}{}
			keys = append(keys, k)
		}
	}
	return keys
}
