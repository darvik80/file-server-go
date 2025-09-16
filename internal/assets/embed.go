package assets

import (
	"embed"
	"io/fs"
)

// Embed all static files
//
//go:embed web/*
var WebFiles embed.FS

// GetWebFS returns filesystem for web resources
func GetWebFS() fs.FS {
	webFS, err := fs.Sub(WebFiles, "web")
	if err != nil {
		panic("Error creating sub filesystem: " + err.Error())
	}
	return webFS
}

// GetTemplateFS returns filesystem for templates
func GetTemplateFS() fs.FS {
	templateFS, err := fs.Sub(WebFiles, "web/templates")
	if err != nil {
		panic("Error creating template filesystem: " + err.Error())
	}
	return templateFS
}

// GetStaticFS returns filesystem for static files (CSS, JS, and root files like favicon)
func GetStaticFS() fs.FS {
	// Create filesystem that includes css, js folders and root files
	return &staticFS{webFS: WebFiles}
}

// staticFS implements fs.FS for static files
type staticFS struct {
	webFS embed.FS
}

func (s *staticFS) Open(name string) (fs.File, error) {
	// Add web/ prefix to path
	return s.webFS.Open("web/" + name)
}
