package shizuka

import (
	"fmt"
	"path/filepath"
	"strings"
)

type Location struct {
	SrcPath string
	DstPath string
	RelPath string
}

func NewLocation(srcRoot, dstRoot, srcPath string) (*Location, error) {
	relPath, err := filepath.Rel(srcRoot, srcPath)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate relative path: %w", err)
	}

	relPath = "/" + filepath.ToSlash(relPath)

	dstPath := filepath.Join(dstRoot, relPath)

	return &Location{
		SrcPath: srcPath,
		DstPath: dstPath,
		RelPath: relPath,
	}, nil
}

func ContentLocation(srcRoot, dstRoot, srcPath string) (*Location, string, error) {
	relPath, err := filepath.Rel(srcRoot, srcPath)
	if err != nil {
		return nil, "", fmt.Errorf("failed to calculate relative path: %w", err)
	}

	relPath = "/" + filepath.ToSlash(relPath)

	ext := filepath.Ext(srcPath)
	if ext == ".md" {
		if filepath.Base(srcPath) == "index.md" {
			relPath = strings.TrimSuffix(relPath, ".md")
			dstPath := filepath.Join(dstRoot, relPath) + ".html"
			relPath = filepath.Dir(relPath)

			return &Location{
				SrcPath: srcPath,
				DstPath: dstPath,
				RelPath: relPath,
			}, "", nil
		} else {
			relPath = strings.TrimSuffix(relPath, ext)
			dirPath := filepath.Join(dstRoot, relPath)
			dstPath := filepath.Join(dirPath, "index.html") // Use .html for output

			return &Location{
				SrcPath: srcPath,
				DstPath: dstPath,
				RelPath: relPath,
			}, dirPath, nil
		}

	}

	// Non-Markdown files: keep original extension
	dstPath := filepath.Join(dstRoot, relPath)

	return &Location{
		SrcPath: srcPath,
		DstPath: dstPath,
		RelPath: relPath,
	}, "", nil
}

func MakeLocations(srcRoot, dstRoot string, srcPaths []string) (locations []Location, err error) {
	locations = make([]Location, len(srcPaths))

	for i, srcPath := range srcPaths {
		loc, err := NewLocation(srcRoot, dstRoot, srcPath)
		if err != nil {
			return nil, err
		}

		locations[i] = *loc
	}

	return locations, nil
}

func MakeContentLocations(srcRoot, dstRoot string, srcPaths []string) (locations []Location, dirs []Location, err error) {
	locations = make([]Location, len(srcPaths))
	dirs = make([]Location, 0, len(srcPaths))

	for i, path := range srcPaths {
		loc, dir, err := ContentLocation(srcRoot, dstRoot, path)
		if err != nil {
			return nil, nil, err
		}

		if dir != "" {
			dirs = append(dirs, Location{
				DstPath: dir, // i think this is the only bit we care about
			})
		}

		locations[i] = *loc
	}

	return locations, dirs, nil
}

func locationsUnion(l ...[]Location) []Location {
	locations := make(map[string]Location)

	for _, slice := range l {
		for _, location := range slice {
			locations[location.DstPath] = location
		}
	}

	result := make([]Location, 0, len(locations))
	for _, location := range locations {
		result = append(result, location)
	}

	return result
}

func locationsIntersect(a, b []Location) []Location {
	locations := make(map[string]Location)
	counts := make(map[string]int)

	for _, location := range a {
		locations[location.DstPath] = location
		counts[location.DstPath]++
	}

	for _, location := range b {
		locations[location.DstPath] = location
		counts[location.DstPath]++
	}

	result := make([]Location, 0, len(locations))
	for _, location := range locations {
		if counts[location.DstPath] > 1 {
			result = append(result, location)
		}
	}

	return result
}
