package lib

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

var (
	postTemplate  *template.Template
	indexTemplate *template.Template
	config        Config
	ConfigFile    string
	PublicDir     string
	PostsDir      string
	TemplatesDir  string
	RssURL        string

	specFiles = map[string]struct{}{
		"favicon.ico":          struct{}{},
		"robots.txt":           struct{}{},
		"humans.txt":           struct{}{},
		"apple-touch-icon.png": struct{}{},
	}
)

func init() {
	pwd, err := os.Getwd()
	if err != nil {
		log.Fatal("FATAL", err)
	}
	PublicDir = filepath.Join(pwd, "public")
	PostsDir = filepath.Join(pwd, "post")
	TemplatesDir = filepath.Join(pwd, "template")
	ConfigFile = filepath.Join(pwd, "config.json")
	config = GetConfig(ConfigFile)
}

func storeRssURL() {
	base, err := url.Parse(config.BaseURL)
	if err != nil {
		fmt.Errorf("Error parsing the baseurl: %s", err)
	}
	rss, err := base.Parse("/rss")
	if err != nil {
		fmt.Errorf("Error parsing the rss url: %s", err)
	}

	RssURL = rss.String()
}

type posts []*LongPost

func (p posts) Len() int           { return len(p) }
func (p posts) Less(i, j int) bool { return p[i].PublishDate.Before(p[j].PublishDate) }
func (p posts) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

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

func clearPublishDir() error {
	files, err := ioutil.ReadDir(PublicDir)
	if err != nil {
		return fmt.Errorf("Error getting public directory files: %s", err)
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
	recent := config.RecentPostsCount
	if length := len(allPosts); length < recent {
		recent = length
	}
	recentPosts = allPosts[:recent]
	return
}

func loadTemplates() {
	postTemplate = template.Must(template.ParseFiles("template/post.html", "template/base.html"))
	indexTemplate = template.Must(template.ParseFiles("template/index.html"))
}

func GenerateSite() error {
	//TODO: format date in the template
	storeRssURL()
	loadTemplates()
	files, err := ioutil.ReadDir(PostsDir)
	if err != nil {
		return err
	}

	files = filter(files)

	allPosts, recentPosts := getPosts(files)

	if err := clearPublishDir(); err != nil {
		return err
	}

	for i, p := range allPosts {
		pt := newPostTempalte(p, i, recentPosts, allPosts)
		if i == 0 {
			if err := generateIndexFile(pt); err != nil {
				return err
			}
		}
		if err := generatePostFile(pt); err != nil {
			return err
		}
	}

	pt := newPostTempalte(nil, 0, recentPosts, allPosts)
	return generateRss(pt)
}

func generateRss(pt *PostTempalte) error {
	rss := NewRss(config.SiteName, config.Slogan, config.BaseURL)
	base, err := url.Parse(config.BaseURL)
	if err != nil {
		return fmt.Errorf("Error parsing base URL: %s", err)
	}

	for _, p := range pt.Recent {
		u, err := base.Parse(p.Slug)
		if err != nil {
			return fmt.Errorf("Error parsing post URL: %s", err)
		}
		rss.Channels[0].AppendItem(NewRssItem(p.Title, p.Description, u.String(), p.Author, "", p.PublishDate))
	}

	return rss.WriteToFile(filepath.Join(PublicDir, "rss.xml"))
}

func generatePostFile(pt *PostTempalte) error {

	fileWriter, err := os.Create(filepath.Join(PublicDir, pt.Post.Slug))
	if err != nil {
		return fmt.Errorf("Error creating static file %s: %s", pt.Post.Slug, err)
	}
	defer fileWriter.Close()

	return postTemplate.ExecuteTemplate(fileWriter, "base", pt)
}

func generateIndexFile(pt *PostTempalte) error {

	indexWriter, err := os.Create(filepath.Join(PublicDir, "index.html"))
	if err != nil {
		return fmt.Errorf("Error creating static file index.html: %s", err)
	}
	defer indexWriter.Close()

	return indexTemplate.ExecuteTemplate(indexWriter, "index", pt)
}
