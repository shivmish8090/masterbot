package utils

import (
	"strings"
)

func ParseFlags(keys []string, text string) (string, map[string]string) {
	values := make(map[string]string)
	keySet := make(map[string]struct{})

	for _, k := range keys {
		keySet[k] = struct{}{}
		values[k] = ""
	}

	lines := strings.Split(text, "\n")
	var remainingLines []string

	for _, line := range lines {
		args := strings.Fields(line)
		var remaining []string
		i := 0

		for i < len(args) {
			arg := args[i]

			// --key=val format
			if strings.HasPrefix(arg, "--") && strings.Contains(arg, "=") {
				parts := strings.SplitN(arg[2:], "=", 2)
				key, val := parts[0], parts[1]
				if _, ok := keySet[key]; ok {
					values[key] = val
					i++
					continue
				}
			}

			// key=val format
			if strings.Contains(arg, "=") {
				parts := strings.SplitN(arg, "=", 2)
				key, val := parts[0], parts[1]
				if _, ok := keySet[key]; ok {
					values[key] = val
					i++
					continue
				}
			}

			// key val format
			if _, ok := keySet[arg]; ok && i+1 < len(args) {
				values[arg] = args[i+1]
				i += 2
				continue
			}

			remaining = append(remaining, arg)
			i++
		}

		remainingLines = append(remainingLines, strings.Join(remaining, " "))
	}

	return strings.Join(remainingLines, "\n"), values
}