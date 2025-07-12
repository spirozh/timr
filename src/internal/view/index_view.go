package view

import (
	"html/template"
	"net/http"
)

type IndexView struct {
	templ *template.Template
}

func NewIndexView(templ *template.Template) *IndexView {
	return &IndexView{templ: templ}
}

func (t *IndexView) Index(w http.ResponseWriter, _ *http.Request) {
	var viewModel any = struct{ Foo string }{"FOO"}
	if err := t.templ.ExecuteTemplate(w, "index", viewModel); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
