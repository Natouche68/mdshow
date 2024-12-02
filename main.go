package main

import (
	"html/template"
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

	tmpl, err := template.ParseFiles("templates/index.gohtml")

	if err != nil {
		panic(err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl.Execute(w, template.HTML(string(html)))
	})
	http.ListenAndServe(":8080", nil)
}
