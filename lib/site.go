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
	"time"
)

const maxGoroutines = 200

var (
	baseTemplate       *template.Template
	postTemplate       *template.Template
	indexTemplate      *template.Template
	categoryTemplate   *template.Template
	collectionTemplate *template.Template
	config             Config
	ConfigFile         string
	PublicDir          string
	PostsDir           string
	TemplatesDir       string
	RssURL             string

	specFiles = map[string]struct{}{
		"favicon.ico":          struct{}{},
		"robots.txt":           struct{}{},
		"humans.txt":           struct{}{},
		"apple-touch-icon.png": struct{}{},
	}

	funcs = template.FuncMap{
		"formattime": func(t time.Time, f string) string {
			return t.Format(f)
		},
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

type posts []LongPost

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

func getPosts(files []os.FileInfo) (allPosts []LongPost, recentPosts []LongPost) {
	fileCount := len(files)
	allPosts = make([]LongPost, 0)
	postChan := make(chan LongPost, fileCount)
	// for _, file := range files {
	// 	if runtime.NumGoroutine() > maxGoroutines {
	// 		newLongPost(file, postChan)
	// 	} else {
	// 		go newLongPost(file, postChan)
	// 	}
	// }

	for _, file := range files {
		go newLongPost(file, postChan)
	}

	for i := 0; i < fileCount; i++ {
		allPosts = append(allPosts, <-postChan)
	}

	sort.Sort(sort.Reverse(posts(allPosts)))

	for i, _ := range allPosts {
		if i > 0 {
			allPosts[i].PrevSlug = allPosts[i-1].Slug
		}
		if i < fileCount-1 {
			allPosts[i].NextSlug = allPosts[i+1].Slug
		}
	}
	recent := config.RecentPostsCount
	if fileCount < recent {
		recent = fileCount
	}
	recentPosts = allPosts[:recent]
	return
}

// func getPosts(files []os.FileInfo) (allPosts []*LongPost, recentPosts []*LongPost) {
// 	allPosts = make([]*LongPost, 0, len(files))
// 	for _, file := range files {
// 		longPost, err := newLongPost(file)
// 		if err != nil {
// 			log.Printf("Post ignored: %s; Error: %s\n", file.Name(), err)
// 		} else {
// 			allPosts = append(allPosts, longPost)
// 		}
// 	}

// 	sort.Sort(sort.Reverse(posts(allPosts)))

// 	for i, _ := range allPosts {
// 		if i > 0 {
// 			allPosts[i].PrevSlug = allPosts[i-1].Slug
// 		}
// 		if i < len(allPosts)-1 {
// 			allPosts[i].NextSlug = allPosts[i+1].Slug
// 		}
// 	}
// 	recent := config.RecentPostsCount
// 	if length := len(allPosts); length < recent {
// 		recent = length
// 	}
// 	recentPosts = allPosts[:recent]
// 	return
// }

func loadTemplates() {
	baseTemplate = template.Must(template.ParseFiles("template/base.html")).Funcs(funcs)
	postTemplate = template.Must(baseTemplate.Clone())
	postTemplate = template.Must(postTemplate.ParseFiles("template/post.html"))
	indexTemplate = template.Must(template.ParseFiles("template/index.html"))
	categoryTemplate = template.Must(template.ParseFiles("template/category.html"))
	collectionTemplate = template.Must(template.ParseFiles("template/collection.html"))
}

func GenerateSite() error {
	storeRssURL()
	loadTemplates()
	files, err := ioutil.ReadDir(PostsDir)
	if err != nil {
		return err
	}

	files = filter(files)

	allPosts, recentPosts := getPosts(files)
	collections := getCollection(allPosts)

	if err := clearPublishDir(); err != nil {
		return err
	}

	for i, p := range allPosts {
		pt := newPostTempalte(p, i, recentPosts, allPosts, config)
		if i == 0 {
			if err := generateIndexFile(pt); err != nil {
				fmt.Errorf("Generate index file error: %s", err)
				return err
			}
		}
		if err := generatePostFile(pt); err != nil {
			fmt.Errorf("Generate post file error: %s", err)
			return err
		}
		fmt.Println(i)
	}

	err = generateCategoryFile(collections)
	if err != nil {
		fmt.Errorf("Generate category file error: %s", err)
		return err
	}

	for key, _ := range collections {
		if err := generateCollectionFile(key, collections[key]); err != nil {
			fmt.Errorf("Generate collection file error: %s", err)
			return err
		}
	}

	pt := newPostTempalte(LongPost{}, 0, recentPosts, allPosts, config)

	err = generateJson(pt)
	if err != nil {
		fmt.Errorf("Generate json file error: %s", err)
		return err
	}
	return generateRss(pt)
}

func generateRss(pt *PostTempalte) error {
	rss := NewRss(config.SiteName, config.Slogan, config.BaseURL, config.Author)
	base, err := url.Parse(config.BaseURL)
	if err != nil {
		return fmt.Errorf("Error parsing base URL: %s", err)
	}

	for _, p := range pt.Recent {
		u, err := base.Parse(p.Slug)
		if err != nil {
			return fmt.Errorf("Error parsing post URL: %s", err)
		}
		rss.Channels[0].AppendItem(NewRssItem(p.Title, p.Description, u.String(), p.Author, "", p.PublishDate.Format("2006-01-02")))
	}

	return rss.WriteToFile(filepath.Join(PublicDir, "rss.xml"))
}

func generateJson(pt *PostTempalte) error {
	siteJson := NewSiteJson(config.SiteName)
	base, err := url.Parse(config.BaseURL)
	if err != nil {
		return fmt.Errorf("Error parsing base URL: %s", err)
	}

	for _, p := range pt.All {
		slug, err := base.Parse(p.Slug)
		if err != nil {
			return fmt.Errorf("Error parsing post URL: %s", err)
		}
		prevSlug, err := base.Parse(p.PrevSlug)
		if err != nil {
			return fmt.Errorf("Error parsing post URL: %s", err)
		}
		nextSlug, err := base.Parse(p.NextSlug)
		if err != nil {
			return fmt.Errorf("Error parsing post URL: %s", err)
		}
		siteJson.AppendPostJson(NewPostJson(slug.String(), p.Author, p.Title, p.Description, p.Category, p.PublishDate.Format("2006-01-02"), p.ModifyDate.Format("2006-01-02"), p.ReadingTime, prevSlug.String(), nextSlug.String(), string(p.Content)))
	}
	return siteJson.WriteToFile(filepath.Join(PublicDir, "site.json"))
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

func generateCategoryFile(c map[string][]LongPost) error {
	categoryWriter, err := os.Create(filepath.Join(PublicDir, "category.html")) //TODO: every category generate a html
	if err != nil {
		return fmt.Errorf("Error creating static file category.html: %s", err)
	}
	defer categoryWriter.Close()

	return categoryTemplate.ExecuteTemplate(categoryWriter, "category", c)
}

func generateCollectionFile(c string, posts []LongPost) error { //TODO: init reposity first
	collectionWriter, err := os.Create(filepath.Join(PublicDir, "collection", c+".html"))
	if err != nil {
		return fmt.Errorf("Error creating static file in collection %s html: %s", c, err)
	}
	defer collectionWriter.Close()

	return collectionTemplate.ExecuteTemplate(collectionWriter, "collection", posts)
}
