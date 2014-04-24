package lib

import (
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

var (
	postTemplate *template.Template
)

type posts []*LongPost

func (p *posts) Len() int           { return len(p) }
func (p *posts) Less(i, j int) bool { return p[i].PublishDate.Before(p[j].PublishDate) }
func (p *posts) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

func filter(file []os.FileInfo) []os.FileInfo {
	for i := 0; i < len(file); {
		if file[i].IsDir() || filepath.Ext(file[i].Name()) != ".md" {
			file[i], file = file[len(file)-1], file[:len(file)-1]
		} else {
			i++
		}
	}

	return file
}

func loadTemplates() {
	postTemplate = template.Must(template.ParseFiles("template/post.html", "template/base.html"))
}

func writeIndex(post *LongPost, fileName string) {
	loadTemplates()

	indexFile := fmt.Sprintf("%s/%s.html", publicPath, fileName)

	file, err := os.Create(indexFile)
	if err != nil {
		fmt.Errorf("Error create post html : %v", err)
	}

	if err := pageTemplate.ExecuteTemplate(file, "base", post); err != nil {
		fmt.Errorf("Error render index file for post : %v", err)
	}
}

func clearPublishDir() error {
	files, err := ioutil.ReadDir(PublicDir)
	if err != nil {
		return fmt.Errorf("Error while getting public directory files: %s", err)
	}

	for _, file := range files {
		if !file.IsDir() && !strings.HasPrefix(file.Name(), ".") {
			if _, ok := specFiles[file.Name()]; !ok {
				err = os.Remove(filepath.Join(PublicDir, file.Name()))
				if err != nil {
					return fmt.Errorf("Error deleting file %s: %s", file.Name(), err)
				}
			}
		}
	}

	return nil
}

func getPosts(files []os.FileInfo) (allPosts []*LongPost, recentPosts []*LongPost) {
	allPosts = make([]*LongPost, 0, len(files))
	for _, file := range files {
		longPost, err := newLongPost(file)
		if err != nil {
			log.Printf("Post ignored: %s; Error: %s\n", file.Name(), err)
		} else {
			allPosts = append(allPosts, longPost)
		}
	}

	sort.Sort(sort.Reverse(posts(allPosts)))
	recent := Config.RecentPostsCount //TODO : config
	if length := len(all); length < recent {
		recent = length
	}
	recentPosts = all[:recent]
	return
}

func generateSite() error {
	loadTemplates()
	files, err := ioutil.ReadDir(PostsDir)
	if err != nil {

	}
}
