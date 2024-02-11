package static

import (
	"embed"
)

//go:generate just _tailwindcss-build
//go:generate just _htmx-download
//go:generate just _feather-icons-download
//go:embed *
var Static embed.FS
