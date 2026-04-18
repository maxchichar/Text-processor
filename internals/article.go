package internals

import(
	"regexp"
	"strings"
)


func ApplyAnRule(tokens []string) []string {
	token := strings.Join(tokens, " ")
	token = regexp.MustCompile(`(?i)\b(a)\s+([aeiouhAEIOUH])`).ReplaceAllString(token, "an $2")
	token = regexp.MustCompile(`(?i)\ban (university|useful|useless|unicorn|uniform|union|unique|user|European)\b`).ReplaceAllString(token, "a $1")
	return strings.Fields(token)
}
