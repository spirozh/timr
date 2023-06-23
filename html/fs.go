package html

import (
	"embed"
	"io/fs"
)

//go:embed  static
var staticFS embed.FS

var FS = static{staticFS}

type static struct {
	fs embed.FS
}

func (f static) Open(name string) (fs.File, error) {
	return f.fs.Open("static/" + name)
}
