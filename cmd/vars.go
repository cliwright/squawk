package cmd

import "strings"

// buildVars parses key=value strings into a map.
func buildVars(raw []string) map[string]string {
	vars := make(map[string]string)
	for _, v := range raw {
		parts := strings.SplitN(v, "=", 2)
		if len(parts) == 2 {
			vars[parts[0]] = parts[1]
		}
	}
	return vars
}
