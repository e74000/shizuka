package shizuka

import (
	"encoding/xml"
	"os"
	"path"
)

type SitemapURL struct {
	Loc        string `xml:"loc"`
	LastMod    string `xml:"lastmod,omitempty"`
	ChangeFreq string `xml:"changefreq,omitempty"`
	Priority   string `xml:"priority,omitempty"`
}

type Sitemap struct {
	XMLName xml.Name     `xml:"urlset"`
	XMLNS   string       `xml:"xmlns,attr"`
	URLs    []SitemapURL `xml:"url"`
	Base    string       `xml:"-"`
}

func NewSitemap(baseURL string) *Sitemap {
	return &Sitemap{
		XMLNS: "http://www.sitemaps.org/schemas/sitemap/0.9",
		URLs:  make([]SitemapURL, 0),
		Base:  baseURL,
	}
}

func (s *Sitemap) AddURL(loc, lastModified, changeFreq, priority string) {
	s.URLs = append(s.URLs, SitemapURL{
		Loc:        path.Join(s.Base, loc),
		LastMod:    lastModified,
		ChangeFreq: changeFreq,
		Priority:   priority,
	})
}

func (s *Sitemap) Build(filePath string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := xml.NewEncoder(file)
	encoder.Indent("", "  ")

	return encoder.Encode(s)
}
