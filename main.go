package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gomarkdown/markdown"
)

func main() {
	md, err := os.ReadFile("example.md")
	if err != nil {
		panic(err)
	}

	html := markdown.ToHTML(md, nil, nil)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, string(html))
	})
	http.ListenAndServe(":8080", nil)
}
