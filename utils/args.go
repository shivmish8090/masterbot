package utils

import (
	"strings"
)

func ParseValues(keys []string, args ...string) map[string]string {
	values := make(map[string]string)
	keySet := make(map[string]struct{})
	unmatchedKeys := []string{}

	for _, k := range keys {
		keySet[k] = struct{}{}
		unmatchedKeys = append(unmatchedKeys, k)
	}

	usedKeys := make(map[string]bool)
	i := 0
	for i < len(args) {
		arg := args[i]

		if strings.HasPrefix(arg, "--") && strings.Contains(arg, "=") {
			parts := strings.SplitN(arg[2:], "=", 2)
			key := parts[0]
			val := parts[1]
			if _, ok := keySet[key]; ok {
				values[key] = val
				usedKeys[key] = true
			}
			i++
			continue
		}

		if strings.Contains(arg, "=") {
			parts := strings.SplitN(arg, "=", 2)
			key := parts[0]
			val := parts[1]
			if _, ok := keySet[key]; ok {
				values[key] = val
				usedKeys[key] = true
			}
			i++
			continue
		}

		if _, ok := keySet[arg]; ok && i+1 < len(args) {
			values[arg] = args[i+1]
			usedKeys[arg] = true
			i += 2
			continue
		}

		for len(unmatchedKeys) > 0 {
			k := unmatchedKeys[0]
			if usedKeys[k] {
				unmatchedKeys = unmatchedKeys[1:]
				continue
			}
			values[k] = arg
			usedKeys[k] = true
			unmatchedKeys = unmatchedKeys[1:]
			break
		}

		i++
	}

	return values
}
