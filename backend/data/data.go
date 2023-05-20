package data

import (
	"embed"
	"io/fs"
)

//go:embed schemas
var schemas embed.FS

func Schemas() fs.FS {
	return schemas
}
