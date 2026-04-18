package internals

import(
	"strconv"
)

func ApplyHex(tokens []string) []string {
	for i := 0; i < len(tokens); i++ {
		if tokens[i] == "(hex)" && i > 0 {
			n, err := strconv.ParseInt(tokens[i-1], 16, 64)
			if err == nil {
				tokens[i-1] = strconv.FormatInt(n, 10)
			}
			tokens[i] = ""
		}
	}
	return tokens
}

func ApplyBin(tokens []string) []string {
	for i := 0; i < len(tokens); i++ {
		if tokens[i] == "(bin)" && i > 0 {
			n, err := strconv.ParseInt(tokens[i-1], 2, 64)
			if err == nil {
				tokens[i-1] = strconv.FormatInt(n, 10)
			}
			tokens[i] = ""
		}
	}
	return tokens
}
