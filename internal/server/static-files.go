package server

import (
	"embed"
	"io/fs"
	"os"
	"path/filepath"
)

// Embed static files at compile time
// This embeds files from the static-source directory at build time
//
//go:embed static-source/*
var staticFS embed.FS

// embeddedAssetNames returns the names of all embedded static files
func embeddedAssetNames() []string {
	var names []string

	// Walk the embedded filesystem
	fs.WalkDir(staticFS, "static-source", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			// Remove the "static-source/" prefix as path is determined by config
			relPath := "/" + filepath.Base(path)
			names = append(names, relPath)
		}
		return nil
	})

	return names
}

// embedAssets loads and returns the asset for the given name
func embedAssets(name string) ([]byte, error) {
	// Remove leading slash if present to match embed.FS expectations
	if len(name) > 0 && name[0] == '/' {
		name = name[1:]
	}

	// Read from embedded filesystem
	return staticFS.ReadFile("static-source/" + name)
}

// restoreEmbeddedAsset writes an embedded asset to the filesystem
func restoreEmbeddedAsset(dir, name string) error {
	data, err := embedAssets(name)
	if err != nil {
		return err
	}

	// Remove leading slash for file path
	if len(name) > 0 && name[0] == '/' {
		name = name[1:]
	}

	targetPath := filepath.Join(dir, name)

	// Ensure parent directory exists
	if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
		return err
	}

	// Write file
	return os.WriteFile(targetPath, data, 0644)
}
