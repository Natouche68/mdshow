package main

import (
	"bufio"
	"bytes"
	"embed"
	"encoding/base64"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strings"

	"github.com/Natouche68/mdshow/themes"

	"github.com/PuerkitoBio/goquery"
	"github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"
	"github.com/charmbracelet/log"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/parser"
)

//go:embed templates/*.gohtml
var templatesFS embed.FS

type PresentationData struct {
	Content       template.HTML
	CodeBlocksCSS template.HTML
	Theme         themes.Theme
}

func main() {
	log.SetReportTimestamp(false)

	tmpl, err := template.ParseFS(templatesFS, "templates/*.gohtml")

	if err != nil {
		log.Fatal("Error parsing template", "err", err)
	}

	mdExtensions := parser.CommonExtensions | parser.OrderedListStart | parser.SuperSubscript

	themeName := ""
	if len(os.Args) < 2 {
		themeName = "catppuccin"
		log.Warn("No theme specified, using default theme")
	} else {
		themeName = os.Args[1]
	}
	theme, themeName := themes.GetTheme(themeName)
	log.Info("Using theme", "name", themeName)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		md, err := os.ReadFile("example.md")
		if err != nil {
			log.Fatal("Error reading Markdown file", "err", err)
		}

		html, css := genareteHTMLFromMarkdown(md, mdExtensions, theme.CodeStyle)

		tmpl.ExecuteTemplate(w, "index.gohtml", PresentationData{
			Content:       template.HTML(html),
			CodeBlocksCSS: template.HTML("<style>" + css + "</style>"),
			Theme:         theme,
		})
	})

	log.Info("Server started on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

func genareteHTMLFromMarkdown(md []byte, mdExtensions parser.Extensions, styleName string) (string, string) {
	mdParser := parser.NewWithExtensions(mdExtensions)
	htmlBytes := markdown.ToHTML(md, mdParser, nil)

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(htmlBytes))
	if err != nil {
		log.Fatal("Couldn't create document from Markdown", "err", err)
	}

	css := getCodeBlocksCSS(styleName)

	doc = hightlightCode(doc, styleName)
	doc = replaceImageContent(doc)

	html := getStringFromHTMLDocument(doc)

	return html, css
}

func hightlightCode(doc *goquery.Document, styleName string) *goquery.Document {
	doc.Find("code[class*=\"language-\"]").Each(func(i int, s *goquery.Selection) {
		class, _ := s.Attr("class")
		lang := strings.TrimPrefix(class, "language-")

		lexer := lexers.Get(lang)
		iterator, err := lexer.Tokenise(nil, s.Text())
		if err != nil {
			log.Error("Couldn't tokenizing code block", "err", err)
		}

		style := styles.Get(styleName)

		formatter := html.New(html.WithClasses(true))
		b := bytes.Buffer{}
		buf := bufio.NewWriter(&b)
		formatter.Format(buf, style, iterator)
		buf.Flush()

		s.SetHtml(b.String())
	})

	return doc
}

func getCodeBlocksCSS(styleName string) string {
	style := styles.Get(styleName)

	buf := bytes.Buffer{}
	writer := bufio.NewWriter(&buf)
	formatter := html.New(html.WithClasses(true))
	if err := formatter.WriteCSS(writer, style); err != nil {
		log.Error("Couldn't write CSS", "err", err)
	}
	writer.Flush()

	return buf.String()
}

func replaceImageContent(doc *goquery.Document) *goquery.Document {
	doc.Find("img").Each(func(i int, s *goquery.Selection) {
		src, _ := s.Attr("src")

		if strings.Contains(src, "http") {
			return
		}

		imgBytes, err := os.ReadFile(src)
		if err != nil {
			log.Error("Couldn't read file", "src", src, "err", err)
			return
		}

		mimeType := http.DetectContentType(imgBytes)
		encodedContent := base64.StdEncoding.EncodeToString(imgBytes)

		dataUrl := fmt.Sprintf("data:%s;base64,%s", mimeType, encodedContent)

		s.SetAttr("src", dataUrl)
	})

	return doc
}

func getStringFromHTMLDocument(doc *goquery.Document) string {
	html, err := doc.Html()
	if err != nil {
		log.Fatal("Couldn't create new HTML document", "err", err)
	}

	html = strings.TrimPrefix(html, "<html><head></head><body>")
	html = strings.TrimSuffix(html, "</body></html>")

	html = strings.ReplaceAll(html, "<pre><code class=\"language-go\"><pre class=\"chroma\"><code>", "<pre class=\"chroma\"><code>")
	html = strings.ReplaceAll(html, "</code></pre></code></pre>", "</code></pre>")

	html = strings.ReplaceAll(html, "<hr/>", "</div><div class=\"slide\">")

	return html
}
