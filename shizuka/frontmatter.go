package shizuka

import (
	"bytes"
	"errors"
	"gopkg.in/yaml.v3"
)

var (
	ErrorFrontmatterDelim   = errors.New("frontmatter delimited incorrectly")
	ErrorFrontmatterPos     = errors.New("text before frontmatter")
	ErrorFrontmatterInvalid = errors.New("invalid frontmatter")
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

	SitemapInclude    bool   `yaml:"sitemap_include"`
	SitemapChangeFreq string `yaml:"sitemap_change_freq"`
	SitemapPriority   string `yaml:"sitemap_priority"`

	RSSInclude bool `yaml:"rss_include"`

	Data     map[string]any `yaml:"data"`
	LiteData map[string]any `yaml:"lite_data"`

	Template string `yaml:"template"`
}

// extractFrontmatter parses the YAML frontmatter and returns the remaining body content.
func extractFrontmatter(content []byte) (*Frontmatter, []byte, error) {
	parts := bytes.SplitN(content, []byte("---"), 3)
	if len(parts) == 0 {
		return new(Frontmatter), content, nil
	}

	if len(parts) < 3 {
		return new(Frontmatter), content, ErrorFrontmatterDelim
	}

	if len(parts[0]) != 0 {
		return new(Frontmatter), content, ErrorFrontmatterPos
	}

	frontmatter := new(Frontmatter)
	if err := yaml.Unmarshal(parts[1], frontmatter); err != nil {
		return new(Frontmatter), content, ErrorFrontmatterInvalid
	}

	return frontmatter, bytes.Join(parts[2:], []byte{}), nil
}
