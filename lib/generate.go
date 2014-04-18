package lib

import (
	"fmt"
	"github.com/russross/blackfriday"
	"html/template"
	"io/ioutil"
	"path/filepath"
	"strings"
)

const (
	postPath = "/Users/AriesDevil/Desktop/test/post"
)

type Page struct {
	Title    string
	Category string
	tags     string
	Date     string
	Author   string
	Content  string
}

func GenerateBlog() {
	fetchFileName(postPath)
}

func fetchFileName(postPath string) {
	fileGolb := postPath + "/*.md"
	files, _ := filepath.Glob(fileGolb)

	for _, file := range files {
		fmt.Println("File: " + file)
		readFile(file)
	}
}

func readFile(fileName string) {
	page := Page{Title: "", Category: "", tags: "", Date: "", Author: "", Content: ""}

	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		printError("Error Reading:", fileName)
		return
	}

	lines := strings.Split(string(data), "\n")
	endHead := false
	for i, line := range lines {
		line = strings.TrimSpace(line)
		if !endHead {
			//parse markdown head for params
			colonIndex := strings.Index(line, ":")

			if colonIndex > 0 {
				key := strings.ToLower(strings.TrimSpace(line[:colonIndex]))
				value := strings.TrimSpace(line[colonIndex+1:])

				switch key {
				case "title":
					page.Title = value
				case "category":
					page.Category = value
				case "tags":
					page.tags = value
				case "date":
					page.Date = value
				case "author":
					page.Author = value
				default:
					return
				}

			}

		}

		if line == "---" {
			i++
			endHead = true
			lines = lines[i:]
		}

	}
	page.Content = markdownRender([]byte(strings.Join(lines, "\n")))
	fmt.Println(page)
	fmt.Println(page.Content)
}

func markdownRender(content []byte) string {
	htmlFlags := 0
	htmlFlags |= blackfriday.HTML_USE_SMARTYPANTS
	htmlFlags |= blackfriday.HTML_SMARTYPANTS_FRACTIONS

	renderer := blackfriday.HtmlRenderer(htmlFlags, "", "")

	extensions := 0
	extensions |= blackfriday.EXTENSION_NO_INTRA_EMPHASIS
	extensions |= blackfriday.EXTENSION_TABLES
	extensions |= blackfriday.EXTENSION_FENCED_CODE
	extensions |= blackfriday.EXTENSION_AUTOLINK
	extensions |= blackfriday.EXTENSION_STRIKETHROUGH
	extensions |= blackfriday.EXTENSION_SPACE_HEADERS

	return string(blackfriday.Markdown(content, renderer, extensions))
}
