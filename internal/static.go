package internal

import (
	"embed"
)

//go:generate just tailwindcss-build
//go:embed static
var Static embed.FS
