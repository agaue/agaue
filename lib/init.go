package lib

import (
	"fmt"
	"os"
	path "path/filepath"
)

// TODO: finish the init template

func CreateBlog(blogName string) {
	currentPath, _ := os.Getwd()
	fmt.Println(currentPath)
	if blogName == "" {
		// ColorLog("[ERROR] Argument [blogname] is missing\n")
		os.Exit(2)
	}

	blogpath := path.Join(currentPath, blogName)
	fmt.Println(blogpath)

	fmt.Println("[INFO] Welcome! Creating Blog...")

	os.MkdirAll(blogpath, 0775)
	fmt.Println(blogpath + string(path.Separator))
	os.Mkdir(path.Join(blogpath, "post"), 0775)
	fmt.Println(path.Join(blogpath, "post") + string(path.Separator))
	os.Mkdir(path.Join(blogpath, "theme"), 0775)
	fmt.Println(path.Join(blogpath, "theme") + string(path.Separator))
	os.Mkdir(path.Join(blogpath, "config"), 0775)
	fmt.Println(path.Join(blogpath, "config") + string(path.Separator))
	os.Mkdir(path.Join(blogpath, "test"), 0775)
	fmt.Println(path.Join(blogpath, "test") + string(path.Separator))
	os.Mkdir(path.Join(blogpath, "public"), 0775)
	fmt.Println(path.Join(blogpath, "public") + string(path.Separator))
	os.Mkdir(path.Join(blogpath+string(path.Separator)+"public", "collection"), 0775)
	fmt.Println(path.Join(blogpath+string(path.Separator)+"public", "collection") + string(path.Separator))
	os.Mkdir(path.Join(blogpath+string(path.Separator)+"public", "static"), 0775)
	fmt.Println(path.Join(blogpath+string(path.Separator)+"public", "static") + string(path.Separator))
	os.Mkdir(path.Join(blogpath+string(path.Separator)+"public"+string(path.Separator)+"static", "css"), 0775)
	fmt.Println(path.Join(blogpath+string(path.Separator)+"public"+string(path.Separator)+"static", "css") + string(path.Separator))
	os.Mkdir(path.Join(blogpath+string(path.Separator)+"public"+string(path.Separator)+"static", "img"), 0775)
	fmt.Println(path.Join(blogpath+string(path.Separator)+"public"+string(path.Separator)+"static", "img") + string(path.Separator))
	os.Mkdir(path.Join(blogpath+string(path.Separator)+"public"+string(path.Separator)+"static", "js"), 0775)
	fmt.Println(path.Join(blogpath+string(path.Separator)+"public"+string(path.Separator)+"static", "js") + string(path.Separator))
	os.Mkdir(path.Join(blogpath+string(path.Separator)+"public"+string(path.Separator)+"static", "font"), 0775)
	fmt.Println(path.Join(blogpath+string(path.Separator)+"public"+string(path.Separator)+"static", "font") + string(path.Separator))
	writeFiles(path.Join(blogpath, "post", "hello-world.md"), hello)
	fmt.Println(path.Join(blogpath, "post", "hello-world.md"))
}

var hello = `
title: Hello World

##Welcome to use Agaue! Have fun~
`

func writeFiles(filename, content string) {
	f, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	f.WriteString(content)
}
