package data

import (
	"embed"
	"io/fs"
)

//go:embed dist
var staticFiles embed.FS

//go:embed schemas
var schemas embed.FS

func StaticFiles() fs.FS {
	return staticFiles
}

func Schemas() fs.FS {
	return schemas
}
