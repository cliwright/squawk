package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Template struct {
	Channel  string   `yaml:"channel"`
	Mentions []string `yaml:"mentions"`
	Text     string   `yaml:"text"`
}

type File struct {
	Templates map[string]Template `yaml:"templates"`
}

type Config struct {
	Templates map[string]Template
}

// Load reads all .yaml/.yml files from the .squawk directory,
// merges them, and returns an error if any template keys collide.
func Load(dir string) (*Config, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("reading %s: %w", dir, err)
	}

	merged := &Config{Templates: make(map[string]Template)}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		ext := filepath.Ext(entry.Name())
		if ext != ".yaml" && ext != ".yml" {
			continue
		}

		path := filepath.Join(dir, entry.Name())
		data, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("reading %s: %w", path, err)
		}

		var f File
		if err := yaml.Unmarshal(data, &f); err != nil {
			return nil, fmt.Errorf("parsing %s: %w", path, err)
		}

		for name, tmpl := range f.Templates {
			if existing, ok := merged.Templates[name]; ok {
				_ = existing
				return nil, fmt.Errorf("duplicate template %q found in %s", name, path)
			}
			merged.Templates[name] = tmpl
		}
	}

	return merged, nil
}
