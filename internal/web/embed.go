package web

import (
	"embed"
	"io/fs"
)

//go:embed all:dist
var distFS embed.FS

// GetDistFS returns the embedded frontend filesystem
// Returns error if dist directory is empty (dev mode)
func GetDistFS() (fs.FS, error) {
	return fs.Sub(distFS, "dist")
}
