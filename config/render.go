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
		for i, m := range t.Mentions {
			if i > 0 {
				b.WriteString(" ")
			}
			if strings.HasPrefix(m, "S") {
				b.WriteString(fmt.Sprintf("<!subteam^%s>", m))
			} else {
				b.WriteString(fmt.Sprintf("<@%s>", m))
			}
		}
		vars["mentions"] = b.String()
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
