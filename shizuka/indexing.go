package shizuka

import (
	"fmt"
	"html/template"
	"io/fs"
	"path/filepath"
)

func walk(root string) (files []string, dirs []string, err error) {
	files = make([]string, 0)
	dirs = make([]string, 0)
	err = filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}

		if d.IsDir() {
			dirs = append(dirs, path)
		} else {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		return nil, nil, fmt.Errorf("walk: walk dir: %w", err)
	}

	return files, dirs, nil
}

func index(src, dst string) (dirs, content, static []Location, templates *template.Template, err error) {
	contentRoot := filepath.Join(src, "content")
	staticRoot := filepath.Join(src, "static")
	templatesRoot := filepath.Join(src, "templates")

	// Index content
	srcContentFiles, srcContentDirs, err := walk(contentRoot)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("index: failed to index content: %w", err)
	}

	contentDirs, err := MakeLocations(contentRoot, dst, srcContentDirs)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("index: failed to find dir paths: %w", err)
	}

	contentFiles, contentIndexDirs, err := MakeContentLocations(contentRoot, dst, srcContentFiles)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("index: failed to find file paths: %w", err)
	}

	srcStaticFiles, srcStaticDirs, err := walk(staticRoot)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("index: failed to index static files: %w", err)
	}

	staticFiles, err := MakeLocations(staticRoot, dst, srcStaticFiles)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("index: failed to find static file paths: %w", err)
	}

	staticDirs, err := MakeLocations(staticRoot, dst, srcStaticDirs)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("index: failed to find dir paths: %w", err)
	}

	dirs = locationsUnion(contentDirs, contentIndexDirs, staticDirs)
	if conflicts := locationsIntersect(contentFiles, staticFiles); len(conflicts) > 0 {
		return nil, nil, nil, nil, fmt.Errorf("index: conflicts found between static files and content: %v", conflicts)
	}

	templates, err = template.ParseGlob(filepath.Join(templatesRoot, "*.tmpl"))
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("index: failed to parse templates: %w", err)
	}

	return dirs, contentFiles, staticFiles, templates, err
}
