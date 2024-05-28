package static

import (
	"embed"
)

//go:generate just static
//go:embed *
var Static embed.FS
