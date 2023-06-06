package htmx

import (
	"net/http"
	"strings"
)

type (
	HTMX struct{}
)

func New() *HTMX {
	return &HTMX{}
}

func (s *HTMX) NewHandler(w http.ResponseWriter, r *http.Request) *Handler {
	return &Handler{
		w:       w,
		r:       r,
		request: s.HxHeader(r.Context()),
		response: &HxResponseHeader{
			Headers: make(map[HxResponseKey]string),
		},
		statusCode: http.StatusOK,
	}
}

func HxStrToBool(str string) bool {
	return strings.EqualFold(str, "true")
}

func HxBoolToStr(b bool) string {
	if b {
		return "true"
	}

	return "false"
}
