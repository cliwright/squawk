package config

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"
)

// Render executes the template text with the given variables.
// Mentions are available as {{ .mentions }} within the template.
func (t *Template) Render(vars map[string]string) (string, error) {
	if _, ok := vars["mentions"]; !ok {
		var b strings.Builder
		for _, m := range t.Mentions {
			b.WriteString(fmt.Sprintf("• <@%s>\n", m))
		}
		vars["mentions"] = strings.TrimRight(b.String(), "\n")
	}

	tmpl, err := template.New("msg").Option("missingkey=zero").Parse(t.Text)
	if err != nil {
		return "", fmt.Errorf("parsing template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, vars); err != nil {
		return "", fmt.Errorf("executing template: %w", err)
	}

	return buf.String(), nil
}
