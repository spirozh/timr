package html

import (
	"embed"
	"html/template"
	"io"
	"log"
)

//go:embed  templates/*
var templates embed.FS

//go:embed  static/*
var Static embed.FS

// pages
var (
	index = parse("index.html")
)

func Index(w io.Writer) error {
	return index.Execute(w, nil)
}

func parse(file string) *template.Template {
	x, err := template.New(file).ParseFS(templates, "templates/layout.html", "templates/"+file)
	if err != nil {
		log.Fatalf("parsing file %s: %r", file, err)
	}
	return x
}
