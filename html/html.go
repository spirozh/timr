package html

import (
	"embed"
	"html/template"
	"io"
	"log"
)

// go:embed  templates/*
var templates embed.FS

// go:embed  static/*
var static embed.FS

// pages
var (
	index = parse("index.html")
)

func Static(w io.Writer, path string) error {
	f, err := static.Open(path)
	if err != nil {
		return err
	}

	_, err = io.Copy(w, f)
	return err
}

func Index(w io.Writer) error {
	return index.Execute(w, nil)
}

func parse(file string) *template.Template {
	x, err := template.New("layout.html").ParseFS(templates, "layout.html", file)
	if err != nil {
		log.Fatalf("parsing file %s: %r", file, err)
	}
	return x
}
