package test

import "embed"

//go:embed templates
var Files embed.FS

//go:embed views
var SpecialFiles embed.FS

//go:embed without-layout
var WithoutLayout embed.FS
