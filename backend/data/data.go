package data

import "embed"

//go:embed dist
var StaticFiles embed.FS

//go:embed schemas
var Schemas embed.FS
