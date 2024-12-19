package shizuka

import (
	"bytes"
	"fmt"
	"github.com/charmbracelet/log"
	"github.com/yuin/goldmark"
	gmhtml "github.com/yuin/goldmark/renderer/html"
	"html/template"
	"io"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"
)

const dateLayout = "2006-01-02"

type Lite struct {
	Title       string
	Description string
	Author      string
	Date        string
	Tags        []string

	Path string
}

type Page struct {
	Title       string
	Description string
	Author      string
	Date        string
	Tags        []string

	MetaTitle       string
	MetaDescription string
	MetaKeywords    string

	Data map[string]any

	Location Location
	Content  template.HTML
	Template string
}

func (p Page) Lite() Lite {
	return Lite{
		Title:       p.Title,
		Description: p.Description,
		Author:      p.Author,
		Date:        p.Date,
		Tags:        p.Tags,
		Path:        p.Location.RelPath,
	}
}

type PageData struct {
	Title       string
	Description string
	Author      string
	Date        string
	Tags        []string

	MetaTitle       string
	MetaDescription string
	MetaKeywords    string

	Data map[string]any

	Content template.HTML
	PageMap map[string][]Lite
}

type BuildOpts struct {
	Dev       bool   // Whether we are using a dev build
	DevScript string // A script to bundle onto the payload when in dev mode
}

type PageBuilder struct {
	src, dst string

	dirs      []Location
	content   []Location
	static    []Location
	templates *template.Template

	pages   map[string]Page
	pageMap map[string][]Lite

	Opts BuildOpts
}

func NewPageBuilder(src, dst string) *PageBuilder {
	return &PageBuilder{
		src: src,
		dst: dst,
	}
}

func (pb *PageBuilder) Index() (err error) {
	dirs, content, static, templates, err := index(pb.src, pb.dst)
	if err != nil {
		return fmt.Errorf("NewPageBuilder: failed to index content: %w", err)
	}

	pb.dirs = dirs
	pb.static = static
	pb.templates = templates

	md := goldmark.New(
		goldmark.WithRendererOptions(
			gmhtml.WithUnsafe(),
		),
	)

	pb.pages = make(map[string]Page)
	pb.pageMap = make(map[string][]Lite)

	for _, file := range content {
		if filepath.Ext(file.SrcPath) != ".md" {
			continue
		}

		// if we are in a dev environment then ignore temp files
		if strings.HasSuffix(file.SrcPath, "~") && pb.Opts.Dev {
			continue
		}

		// /ws is used for hot-reloading in dev env, so we need to be careful if a page shares the same path
		if file.RelPath == "/ws" {
			log.Warn("/ws path is reserved for websocket communication in the dev environment. this may cause unexpected behaviour between dev and prod", "file", file.SrcPath)
		}

		fileContent, err := os.ReadFile(file.SrcPath)
		if err != nil {
			return fmt.Errorf("NewPageBuilder: failed to read %q: %w", file, err)
		}

		frontmatter, fileContent, err := extractFrontmatter(fileContent)
		if err != nil && fileContent == nil {
			log.Error("failed to parse file", "file", file, "error", err)
			continue
		} else if err != nil {
			log.Warn("failed to parse frontmatter, ignoring...", "file", file, "error", err)
		}

		htmlBuf := bytes.NewBuffer(nil)
		err = md.Convert(fileContent, htmlBuf)
		if err != nil {
			return fmt.Errorf("buildPage: failed to build html fileContent: %w", err)
		}

		fileContent = htmlBuf.Bytes()

		if pb.Opts.Dev {
			fileContent = append(fileContent, []byte(pb.Opts.DevScript)...)
		}

		pb.pages[file.RelPath] = Page{
			Title:           frontmatter.Title,
			Description:     frontmatter.Description,
			Author:          frontmatter.Author,
			Date:            frontmatter.Date,
			Tags:            frontmatter.Tags,
			MetaTitle:       frontmatter.MetaTitle,
			MetaDescription: frontmatter.MetaDescription,
			MetaKeywords:    frontmatter.MetaKeywords,
			Data:            frontmatter.Data,
			Location:        file,
			Content:         template.HTML(fileContent),
			Template:        frontmatter.Template,
		}

		pb.pageMap[file.RelPath] = make([]Lite, 0)
	}

	for s, page := range pb.pages {
		if s == "/" {
			continue // root does not have a parent
		}

		parent := filepath.Dir(page.Location.RelPath)
		pb.pageMap[parent] = append(pb.pageMap[parent], page.Lite())
	}

	for _, pages := range pb.pageMap {
		slices.SortFunc(pages, func(a, b Lite) int {
			at, _ := time.Parse(dateLayout, a.Date)
			bt, _ := time.Parse(dateLayout, b.Date)

			if at.After(bt) {
				return -1
			} else {
				return +1
			}
		})
	}

	return nil
}

func (pb *PageBuilder) Build() (err error) {
	if err := pb.replicateDirs(); err != nil {
		return fmt.Errorf("Build: failed to replicate directories: %w", err)
	}

	if err := pb.replicateStatic(); err != nil {
		return fmt.Errorf("Build: failed to replicate static content: %w", err)
	}

	for _, page := range pb.pages {
		temp := pb.templates.Lookup(page.Template)
		if temp == nil {
			log.Warn("failed to find template, skipping", "file", page.Location.SrcPath, "template", page.Template)
			continue
		}

		file, err := os.Create(page.Location.DstPath)
		if err != nil {
			return fmt.Errorf("Build: failed to create file %s: %w", page.Location.DstPath, err)
		}

		if err := temp.Execute(file, pb.pageData(page)); err != nil {
			_ = file.Close()
			return fmt.Errorf("Build: failed to render file %s: %w", page.Location.RelPath, err)
		}

		_ = file.Close()
	}

	return nil
}

func (pb *PageBuilder) pageData(page Page) PageData {
	return PageData{
		Title:           page.Title,
		Description:     page.Description,
		Author:          page.Author,
		Date:            page.Date,
		Tags:            page.Tags,
		MetaTitle:       page.MetaTitle,
		MetaDescription: page.MetaDescription,
		MetaKeywords:    page.MetaKeywords,
		Data:            page.Data,
		Content:         page.Content,
		PageMap:         pb.pageMap,
	}
}

func (pb *PageBuilder) replicateStatic() error {
	for _, location := range pb.static {

		// if we are in a dev environment then ignore temp files
		if strings.HasSuffix(location.SrcPath, "~") && pb.Opts.Dev {
			continue
		}

		// /ws is used for hot-reloading in dev env, so we need to be careful if a page shares the same path
		if location.RelPath == "/ws" {
			log.Warn("/ws path is reserved for websocket communication in the dev environment. this may cause unexpected behaviour between dev and prod", "file", location.SrcPath)
		}

		// Open the source file
		srcFile, err := os.Open(location.SrcPath)
		if err != nil {
			return fmt.Errorf("failed to open source file %s: %w", location.SrcPath, err)
		}

		// Create the destination file
		dstFile, err := os.Create(location.DstPath)
		if err != nil {
			_ = srcFile.Close()
			return fmt.Errorf("failed to create destination file %s: %w", location.DstPath, err)
		}

		// Copy the content from source to destination
		if _, err := io.Copy(dstFile, srcFile); err != nil {
			_ = srcFile.Close()
			_ = dstFile.Close()
			return fmt.Errorf("failed to copy content from %s to %s: %w", location.SrcPath, location.DstPath, err)
		}

		_ = srcFile.Close()
		_ = dstFile.Close()
	}

	return nil
}

func (pb *PageBuilder) replicateDirs() error {
	// Remove any existing destination directory
	if err := os.RemoveAll(pb.dst); err != nil {
		return fmt.Errorf("build: failed to remove %s: %w", pb.dst, err)
	}
	if err := os.MkdirAll(pb.dst, os.ModePerm); err != nil {
		return fmt.Errorf("build: failed to create directory %s: %w", pb.dst, err)
	}

	for _, dir := range pb.dirs {
		if err := os.MkdirAll(dir.DstPath, os.ModePerm); err != nil {
			return fmt.Errorf("build: failed to replicate directory %s: %w", dir.DstPath, err)
		}
	}

	return nil
}
