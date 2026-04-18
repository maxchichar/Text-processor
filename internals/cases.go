package internals

import(
	"strings"
	"strconv"
)

func ApplyCases(tokens []string) []string {
	result := make([]string, 0, len(tokens))

	for i := 0; i < len(tokens); i++ {
		token := tokens[i]

		if len(token) >= 3 && token[0] == '(' && token[len(token)-1] == ')' {
		inside := strings.ToLower(strings.Trim(token, "()"))
		
		action := inside
		count := 1

		if strings.Contains(inside, ",") {
			parts := strings.Split(inside, ",")
			
			if len(parts) == 2 {
				action = strings.TrimSpace(parts[0])
				number, err := strconv.Atoi(strings.TrimSpace(parts[1]))
				if err == nil && number > 0 {
					count = number
				}
			}
		}
		for n := 0; n < count; n++ {
			idx := len(result) - 1 - n
			if idx < 0 {
				break
			}

			switch action {
			case "up":
				result[idx] = strings.ToUpper(result[idx])
			case "low":
				result[idx] = strings.ToLower(result[idx])
			case "cap":
				if len(result[idx]) > 0 {
					result[idx] = strings.ToUpper(result[idx][:1]) + strings.ToLower(result[idx][1:])
				}
			}
		}

		continue
		}
		result = append(result, token)
	}
	return result
}