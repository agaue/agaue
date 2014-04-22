package lib

import (
	"fmt"
	"github.com/russross/blackfriday"
	"html/template"
	"io/ioutil"
	"path/filepath"
	"sort"
	"strings"
)

const (
	postPath     = "/Users/AriesDevil/Desktop/test/post"
	templatePath = "/Users/AriesDevil/Desktop/test/template"
	publicPath   = "/Users/AriesDevil/Desktop/test/public"
)

type Page struct {
	Title     string
	Url       string
	Category  string
	tags      string
	Date      string
	Author    string
	content   string
	NextTitle string
	NextUrl   string
	PrevTitle string
	PrevUrl   string
}

type PageSlice []Page

func GenerateBlog() {
	fetchFileName(postPath)
}

func fetchFileName(postPath string) {
	var pages PageSlice
	fileGolb := postPath + "/*.md"
	files, _ := filepath.Glob(fileGolb)

	for _, file := range files {
		fmt.Println("File: " + file)
		page := readFile(file)
		pages = append(pages, page)
	}
}

func readFile(fileName string) (page Page) {
	page = Page{Title: "", Category: "", tags: "", Date: "", Author: "", content: ""}

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
	page.content = strings.Join(lines, "\n")
	fmt.Println(page)
	fmt.Println(page.content)
	tempFileName := strings.TrimSuffix(strings.TrimPrefix(fileName, postPath+"/"), ".md")
	page.writeIndex(tempFileName)
	page.Url = slug(tempFileName)
	return page
}

func (p Page) Tags() []string {
	return strings.Split(p.tags, ",")
}

//Called by template for safety
func (p Page) Content() template.HTML {
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

	return template.HTML(blackfriday.Markdown([]byte(p.content), renderer, extensions))
}

func pageList(pages PageSlice) (list PageSlice) {
	for _, page := range pages {
		if page.Date.Format("2014") != "1970" {
			list = append(list, page)
		}
	}
	list.Sort()

	//reverse
	for i, j := 0, len(list)-1; i < j; i, j = i+1, j-1 {
		list[i], list[j] = list[j], list[i]
	}
	return list
}

func (pages PageSlice) Sort() {
	sort.Sort(pages)
}
