package html

import (
	"embed"
	"io/fs"
)

//go:embed  static
var staticFS embed.FS

var FS fs.FS

func init() {
	fs, err := fs.Sub(staticFS, "static")
	if err != nil {
		panic("can't get a subtree filesystem from staticFS")
	}
	FS = fs
}
