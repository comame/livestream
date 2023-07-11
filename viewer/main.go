package main

import (
	_ "embed"
	"io"
	"net/http"
	"strings"

	"github.com/comame/router-go"
)

//go:embed index.html
var indexHtml string

//go:embed hls.js
var hlsJs string

func main() {
	router.Get("/viewer/hls.js", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, hlsJs)
	})
	router.Get("/viewer/:key", func(w http.ResponseWriter, r *http.Request) {
		p := router.Params(r)
		key := p["key"]

		res := strings.Replace(indexHtml, "---stream-key---", key, 1)

		io.WriteString(w, res)
	})

	http.ListenAndServe(":8080", router.Handler())
}
