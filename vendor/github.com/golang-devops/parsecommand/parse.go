package parsecommand

import (
	"strings"
)

func Parse(line string) ([]string, error) {
	args := [][]rune{}

	var quoteChar rune
	var isInQuote bool = false

	trimmedLine := strings.TrimSpace(line)

	currentArg := []rune{}
	for i, c := range trimmedLine {
		isLastChar := i == len(trimmedLine)-1

		if !isInQuote && (c == '"' || c == '\'') {
			isInQuote = true
			quoteChar = c
			if isLastChar {
				args = append(args, currentArg)
			}
			continue
		}

		if isInQuote && c == quoteChar {
			//Ensure it is not escaped with a slash beforehand
			if i == 0 || trimmedLine[i-1] != '\\' {
				isInQuote = false
				if isLastChar {
					args = append(args, currentArg)
				}
				continue
			}
		}

		if !isInQuote && c == ' ' {
			//Ignore multiple spaces
			if i > 0 && trimmedLine[i-1] != ' ' {
				args = append(args, currentArg)
				currentArg = []rune{}
				continue
			}
		}

		currentArg = append(currentArg, c)

		if isLastChar {
			args = append(args, currentArg)
		}
	}

	strArgs := []string{}
	for _, a := range args {
		strArgs = append(strArgs, strings.TrimSpace(string(a)))
	}
	return strArgs, nil
}
