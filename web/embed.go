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
	return http.ListenAndServe(addr, http.FileServer(http.FS(sub)))
}
