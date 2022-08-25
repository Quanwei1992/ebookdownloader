package ebookdownloader

import (
	"embed"
)

//go:embed tpls
var templateFS embed.FS

//go:embed fonts
var fontFS embed.FS
