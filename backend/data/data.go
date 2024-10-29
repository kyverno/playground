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

func Schemas() (fs.FS, error) {
	return fs.Sub(schemas, "schemas")
}

func BuiltInCrds(crd string) (fs.FS, error) {
	return fs.Sub(builtInCrds, filepath.Join(crdsFolder, crd))
}
