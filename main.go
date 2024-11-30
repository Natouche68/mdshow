package main

import (
	"fmt"

	"github.com/gomarkdown/markdown"
)

func main() {
	md := []byte("# Hello world")
	html := markdown.ToHTML(md, nil, nil)
	fmt.Println(string(html))
}
