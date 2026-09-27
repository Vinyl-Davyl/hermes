package web

import (
	"embed"
	"io/fs"
	"net/http"
)

// Built by `npm --prefix web run build` into out/.
//
//go:embed all:out
var files embed.FS

func Serve(addr string) error {
	sub, err := fs.Sub(files, "out")
	if err != nil {
		return err
	}
	fileServer := http.FileServer(http.FS(sub))
	return http.ListenAndServe(addr, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/docs" || r.URL.Path == "/docs/" {
			r.URL.Path = "/docs.html"
		}
		fileServer.ServeHTTP(w, r)
	}))
}
