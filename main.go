package main

import (
	"html/template"
	"net/http"
	"os"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/gomarkdown/markdown"
)

type PresentationData struct {
	Content template.HTML
}

func main() {
	log.SetReportTimestamp(false)

	tmpl, err := template.ParseFiles("templates/index.gohtml")

	if err != nil {
		log.Fatal("Error parsing template", "err", err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		md, err := os.ReadFile("example.md")
		if err != nil {
			log.Fatal("Error reading Markdown file", "err", err)
		}

		html := string(markdown.ToHTML(md, nil, nil))

		html = strings.ReplaceAll(html, "<hr>", "</div><div class=\"slide\">")

		tmpl.Execute(w, PresentationData{
			Content: template.HTML(html),
		})
	})

	log.Info("Server started on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
