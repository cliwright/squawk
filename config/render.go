package config

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"
)

// Render executes the template text with the given variables and prepends mentions.
func (t *Template) Render(vars map[string]string) (string, error) {
	tmpl, err := template.New("msg").Option("missingkey=zero").Parse(t.Text)
	if err != nil {
		return "", fmt.Errorf("parsing template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, vars); err != nil {
		return "", fmt.Errorf("executing template: %w", err)
	}

	var msg strings.Builder

	if len(t.Mentions) > 0 {
		for _, m := range t.Mentions {
			msg.WriteString(fmt.Sprintf("<@%s> ", m))
		}
		msg.WriteString("\n")
	}

	msg.WriteString(buf.String())

	return msg.String(), nil
}
