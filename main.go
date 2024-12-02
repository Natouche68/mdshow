package main

import (
	"html/template"
	"net/http"
	"os"

	"github.com/charmbracelet/log"
	"github.com/gomarkdown/markdown"
)

func main() {
	log.SetReportTimestamp(false)

	md, err := os.ReadFile("example.md")
	if err != nil {
		log.Fatal("Error reading Markdown file", "err", err)
	}

	html := markdown.ToHTML(md, nil, nil)

	tmpl, err := template.ParseFiles("templates/index.gohtml")

	if err != nil {
		log.Fatal("Error parsing template", "err", err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl.Execute(w, template.HTML(string(html)))
	})

	log.Info("Server started on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
