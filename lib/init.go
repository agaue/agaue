package lib

import (
	"fmt"
	"os"
	path "path/filepath"
	"strings"
)

func CreateBlog(blogName string) {
	currentPath, _ := os.Getwd()
	if blogName == "" {
		// ColorLog("[ERROR] Argument [blogname] is missing\n")
		os.Exit(2)
	}

	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		os.Exit(2)
	}

	blogpath := path.Join(currentPath, blogName)

	fmt.Println("[INFO] Creating Blog...")

	os.MkdirAll(blogpath, 0775)
	fmt.Println(blogpath + string(path.Separator))
	os.Mkdir(path.Join(blogpath, "public"), 0775)
	fmt.Println(path.Join(blogpath, "public") + string(path.Separator))
	writeFiles(path.Join(blogpath, "public", "hello-world.md"), strings.Replace(hello, "{{.Appname}}", strings.Join(strings.Split(blogpath[len(blogsrcpath)+1:], string(path.Separator)), "/"), -1))
	fmt.Println(path.Join(blogpath, "public", "hello-world.md"))
}

var hello = `
title: Hello World

##Welcome to Agaue blog engine! Hello World!
`

func writeFiles(filename, content string) {
	f, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	f.WriteString(content)
}
