package internals

import(
	"regexp"
)

func ApplyTokenizer(token string) []string {
	if len(token) == 0 {
		return []string{}
	}

	token = regexp.MustCompile(`(\S+)\s*(\([^)]+\))`).ReplaceAllString(token, "$1 $2")
	t := regexp.MustCompile(`\([^)]+\)|\S+`)
	return t.FindAllString(token, -1)
}