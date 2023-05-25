package data

import (
	"embed"
	"io/fs"
	"path/filepath"
)

const crdsFolder = "crds"

//go:embed schemas
var schemas embed.FS

//go:embed crds
var builtInCrds embed.FS

func Schemas() fs.FS {
	return schemas
}

func BuiltInCrds(crd string) (fs.FS, string) {
	return builtInCrds, filepath.Join(crdsFolder, crd)
}
