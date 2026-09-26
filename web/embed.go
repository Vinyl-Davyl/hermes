package web

import (
	"embed"
	"io/fs"
	"net/http"
)

//go:embed index.html styles.css
var Files embed.FS

func Serve(addr string) error {
	sub, err := fs.Sub(Files, ".")
	if err != nil {
		return err
	}
	return http.ListenAndServe(addr, http.FileServer(http.FS(sub)))
}
