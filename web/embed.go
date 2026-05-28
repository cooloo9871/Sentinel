package web

import "embed"

//go:embed dist
var StaticFiles embed.FS
