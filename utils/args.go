package utils

import (
	"strings"
)

func ParseFlags(keys []string, strict bool, text string) (string, map[string]string) {
	values := make(map[string]string)
	keySet := make(map[string]struct{})
	unmatchedKeys := []string{}

	for _, k := range keys {
		keySet[k] = struct{}{}
		unmatchedKeys = append(unmatchedKeys, k)
		values[k] = ""
	}

	usedKeys := make(map[string]bool)
	args := strings.Fields(text)
	remaining := []string{}

	i := 0
	for i < len(args) {
		arg := args[i]

		// --key=val format
		if strings.HasPrefix(arg, "--") && strings.Contains(arg, "=") {
			parts := strings.SplitN(arg[2:], "=", 2)
			key := parts[0]
			val := parts[1]
			if _, ok := keySet[key]; ok {
				values[key] = val
				usedKeys[key] = true
				i++
				continue
			}
		}

		// key=val format
		if strings.Contains(arg, "=") {
			parts := strings.SplitN(arg, "=", 2)
			key := parts[0]
			val := parts[1]
			if _, ok := keySet[key]; ok {
				values[key] = val
				usedKeys[key] = true
				i++
				continue
			}
		}

		// key val format
		if _, ok := keySet[arg]; ok && i+1 < len(args) {
			values[arg] = args[i+1]
			usedKeys[arg] = true
			i += 2
			continue
		}

		// if not strict, use fallback assignment to unused keys
		if !strict {
			assigned := false
			for len(unmatchedKeys) > 0 {
				k := unmatchedKeys[0]
				unmatchedKeys = unmatchedKeys[1:]
				if !usedKeys[k] {
					values[k] = arg
					usedKeys[k] = true
					assigned = true
					break
				}
			}
			if !assigned {
				remaining = append(remaining, arg)
			}
		} else {
			remaining = append(remaining, arg)
		}

		i++
	}

	return strings.Join(remaining, " "), values
}