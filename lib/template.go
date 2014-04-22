package lib

import (
	"fmt"
	"html/template"
	"os"
)

var (
	pageTemplate *template.Template
)

func loadTemplates() {
	pageTemplate = template.Must(template.ParseFiles("template/page.html", "template/base.html"))
}

//whether a function or a method of Page?
func (page *Page) writeIndex(fileName string) {
	loadTemplates()

	indexFile := fmt.Sprintf("%s/%s.html", publicPath, fileName)

	file, err := os.Create(indexFile)
	if err != nil {
		fmt.Errorf("Error create post html : %v", err)
	}

	if err := pageTemplate.ExecuteTemplate(file, "base", page); err != nil {
		fmt.Errorf("Error render index file for post : %v", err)
	}
}
