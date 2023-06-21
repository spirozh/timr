package http

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/spirozh/timr/html"
)

func Selma(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "hello selma!!\n")
}

func Long(w http.ResponseWriter, r *http.Request) {
	// Make sure that the writer supports flushing.
	//
	flusher, ok := w.(http.Flusher)

	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	for i := 0; ; i++ {

		fmt.Fprintf(w, "data: %d\n\n", i)
		flusher.Flush()

		time.Sleep(time.Second)
	}
}

func TemplateRoutes() http.Handler {
	m := http.NewServeMux()
	return m
}

func staticHandler(w http.ResponseWriter, r *http.Request) {
	// assume GET
	io.WriteString(w, "static handler\n")

}

func APIRoutes() http.Handler {
	m := http.NewServeMux()

	return m
}

func SSERoutes() http.Handler {
	m := http.NewServeMux()

	m.HandleFunc("/api/sse/long/", Long)

	return m
}

func static() http.HandlerFunc {
	// http.FileServer(http.FS(html.Static))

	efs := html.Static
	hfs := http.FS(efs)

	fs := http.FileServer(hfs)

	return func(w http.ResponseWriter, r *http.Request) {
		files := []string{"/static/hello.txt", "static/hello.txt", "/hello.txt", "hello.txt"}
		for _, file := range files {
			_, e := efs.Open(file)
			if e != nil {
				fmt.Println("error!", e)
				continue
			}
			fmt.Println("ok: ", file)
		}

		fmt.Println(r.URL.Path)
		fs.ServeHTTP(w, r)
	}

}

func routes() http.Handler {
	m := http.NewServeMux()

	// templates
	m.Handle("/", TemplateRoutes())

	// static
	//	m.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.FS(html.Static))))
	m.Handle("/static/", static())

	// api
	m.Handle("/api/", APIRoutes())

	// sse
	m.Handle("/api/sse/", SSERoutes())

	m.HandleFunc("/selma", Selma)

	return m
}
