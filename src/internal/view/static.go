package view

import (
	"embed"
	"net/http"
)

//go:embed static/*
var static embed.FS

func StaticHandler() http.Handler {
	return http.FileServer(http.FS(static))
}
