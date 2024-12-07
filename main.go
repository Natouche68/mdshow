package main

import (
	"bufio"
	"bytes"
	"html/template"
	"net/http"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"
	"github.com/charmbracelet/log"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/parser"
)

type PresentationData struct {
	Content       template.HTML
	CodeBlocksCSS template.HTML
}

func main() {
	log.SetReportTimestamp(false)

	tmpl, err := template.ParseFiles("templates/index.gohtml")

	if err != nil {
		log.Fatal("Error parsing template", "err", err)
	}

	mdExtensions := parser.CommonExtensions | parser.OrderedListStart | parser.SuperSubscript

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		md, err := os.ReadFile("example.md")
		if err != nil {
			log.Fatal("Error reading Markdown file", "err", err)
		}

		mdParser := parser.NewWithExtensions(mdExtensions)
		htmlBytes := markdown.ToHTML(md, mdParser, nil)

		css := getCodeBlocksCSS()

		html := hightlightCodeToHTML(htmlBytes)
		html = strings.ReplaceAll(html, "<hr/>", "</div><div class=\"slide\">")

		tmpl.Execute(w, PresentationData{
			Content:       template.HTML(html),
			CodeBlocksCSS: template.HTML("<style>" + css + "</style>"),
		})
	})

	log.Info("Server started on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

func hightlightCodeToHTML(htmlBytes []byte) string {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(htmlBytes))
	if err != nil {
		log.Fatal("Couldn't create document from Markdown", "err", err)
	}

	doc.Find("code[class*=\"language-\"]").Each(func(i int, s *goquery.Selection) {
		class, _ := s.Attr("class")
		lang := strings.TrimPrefix(class, "language-")

		lexer := lexers.Get(lang)
		iterator, err := lexer.Tokenise(nil, s.Text())
		if err != nil {
			log.Error("Couldn't tokenizing code block", "err", err)
		}

		style := styles.Get("catppuccin-macchiato")

		formatter := html.New(html.WithClasses(true))
		b := bytes.Buffer{}
		buf := bufio.NewWriter(&b)
		formatter.Format(buf, style, iterator)
		buf.Flush()

		s.SetHtml(b.String())
	})

	html, err := doc.Html()
	if err != nil {
		log.Error("Couldn't create new HTML document", "err", err)
	}

	html = strings.ReplaceAll(html, "<pre><code class=\"language-go\"><pre class=\"chroma\"><code>", "<pre class=\"chroma\"><code>")
	html = strings.ReplaceAll(html, "</code></pre></code></pre>", "</code></pre>")

	return html
}

func getCodeBlocksCSS() string {
	style := styles.Get("catppuccin-macchiato")

	buf := bytes.Buffer{}
	writer := bufio.NewWriter(&buf)
	formatter := html.New(html.WithClasses(true))
	if err := formatter.WriteCSS(writer, style); err != nil {
		log.Error("Couldn't write CSS", "err", err)
	}
	writer.Flush()

	return buf.String()
}
