package shizuka

import (
	"bytes"
	"fmt"
	"gopkg.in/yaml.v3"
)

type Frontmatter struct {
	Title       string   `yaml:"title"`
	Description string   `yaml:"description"`
	Author      string   `yaml:"author"`
	Date        string   `yaml:"date"`
	Tags        []string `yaml:"tags"`

	MetaTitle       string `yaml:"meta_title"`
	MetaDescription string `yaml:"meta_description"`
	MetaKeywords    string `yaml:"meta_keywords"`

	Data map[string]any `yaml:"data"`

	Template string `yaml:"template"`
}

// extractFrontmatter parses the YAML frontmatter and returns the remaining body content.
func extractFrontmatter(content []byte) (*Frontmatter, []byte, error) {
	parts := bytes.SplitN(content, []byte("---"), 3)
	if len(parts) == 0 {
		return new(Frontmatter), content, nil
	}

	if len(parts) < 3 {
		return new(Frontmatter), content, fmt.Errorf("frontmatter not properly delimited by '---'")
	}

	if len(parts[0]) != 0 {
		return new(Frontmatter), content, fmt.Errorf("text before frontmatter")
	}

	frontmatter := new(Frontmatter)
	if err := yaml.Unmarshal(parts[1], frontmatter); err != nil {
		return new(Frontmatter), content, fmt.Errorf("failed to parse YAML frontmatter: %w", err)
	}

	return frontmatter, bytes.Join(parts[2:], []byte{}), nil
}
