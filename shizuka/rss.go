package shizuka

import (
    "encoding/xml"
    "os"
    "path/filepath"
    "time"
)

type RSS struct {
    XMLName xml.Name   `xml:"rss"`
    Version string     `xml:"version,attr"`
    XMLNS   string     `xml:"xmlns:atom,attr"` // Add Atom namespace
    Channel RssChannel `xml:"channel"`
}

type RssChannel struct {
    Title         string    `xml:"title"`
    Link          string    `xml:"link"`
    AtomLink      AtomLink  `xml:"atom:link"` // Add Atom link
    Description   string    `xml:"description"`
    Language      string    `xml:"language"`
    LastBuildDate string    `xml:"lastBuildDate"`
    Items         []RSSItem `xml:"item"`
}

type AtomLink struct {
    Href string `xml:"href,attr"`
    Rel  string `xml:"rel,attr"`
    Type string `xml:"type,attr"`
}

type RSSItem struct {
    Title       string `xml:"title"`
    Link        string `xml:"link"`
    Description string `xml:"description"`
    PubDate     string `xml:"pubDate"`
    GUID        string `xml:"guid"`
}

func NewRSS(baseURL, title, description, language string) *RSS {
    return &RSS{
        Version: "2.0",
        XMLNS:   "http://www.w3.org/2005/Atom",
        Channel: RssChannel{
            Title:       title,
            Link:        baseURL,
            Description: description,
            Language:    language,
            AtomLink: AtomLink{
                Href: baseURL + "/rss.xml",
                Rel:  "self",
                Type: "application/rss+xml",
            },
            LastBuildDate: time.Now().Format(time.RFC1123Z),
            Items:         make([]RSSItem, 0),
        },
    }
}

func (r *RSS) AddItem(link, publishDate, title, description string) {
    date, err := time.Parse("2006-01-02", publishDate)
    if err != nil {
        date = time.Now() // fallback date if improperly formatted
    }

    fullLink := filepath.Join(r.Channel.Link, link)
    r.Channel.Items = append(r.Channel.Items, RSSItem{
        Title:       title,
        Link:        fullLink,
        Description: description,
        PubDate:     date.Format(time.RFC1123Z),
        GUID:        fullLink,
    })
}

func (r *RSS) Build(filePath string) error {
    file, err := os.Create(filePath)
    if err != nil {
        return err
    }
    defer file.Close()

    encoder := xml.NewEncoder(file)
    encoder.Indent("", "  ")

    return encoder.Encode(r)
}
