package view

import (
	"embed"
	"html/template"
)

//go:embed templates/*
var templates embed.FS

func NewTemplates() (*template.Template, error) {
	return template.ParseFS(templates, "templates/*.html")
}
