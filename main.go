package main

import (
	"fmt"
	"os"

	"github.com/gomarkdown/markdown"
)

func main() {
	md, err := os.ReadFile("example.md")
	if err != nil {
		panic(err)
	}

	html := markdown.ToHTML(md, nil, nil)
	fmt.Println(string(html))
}
