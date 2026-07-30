package fs

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Resolver struct {
	searchPaths []string
}

func NewResolver(searchPaths []string) *Resolver {
	return &Resolver{searchPaths: searchPaths}
}

func (r *Resolver) Resolve(importPath, fromDir string) (string, error) {
	if strings.HasPrefix(importPath, "./") || strings.HasPrefix(importPath, "../") {
		return r.resolveRelative(importPath, fromDir)
	}
	return r.resolveSearch(importPath)
}

func (r *Resolver) resolveRelative(importPath, fromDir string) (string, error) {
	candidate := filepath.Clean(filepath.Join(fromDir, importPath))
	if !strings.HasSuffix(candidate, SourceExtension) {
		candidate += SourceExtension
	}
	if _, err := os.Stat(candidate); err != nil {
		return "", fmt.Errorf("import %q not found", importPath)
	}
	return candidate, nil
}

func (r *Resolver) resolveSearch(importPath string) (string, error) {
	for _, sp := range r.searchPaths {
		candidate := filepath.Join(sp, importPath+SourceExtension)
		if _, err := os.Stat(candidate); err == nil {
			return candidate, nil
		}
	}
	return "", fmt.Errorf("import %q not found in search paths", importPath)
}
