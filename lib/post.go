package lib

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/russross/blackfriday"
	"html/template"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

var (
	errEmptyPost          = fmt.Errorf("Empty Markdown File!")
	errInvalidFrontMatter = fmt.Errorf("Invalid Front Matter")
	errMissingFrontMatter = fmt.Errorf("Missing Front Matter")

	dateFormatter = map[int]string{
		10: "2006-01-02",
		13: "2006-01-02 15h",
		14: "2006-01-02 15h",
		15: "2006-01-02 15:04",
		16: "2006-01-02 15:04",
		25: time.RFC3339,
	}
)

type postTempalte struct {
	SiteName string
	RssURL   string
	Post     longPost
	Recent   []longPost
	All      []longPost
	Prev     shortPost
	Next     shortPost
}

type shortPost struct {
	Slug        string
	Author      string
	Title       string
	Description string
	Category    string
	Tags        []string
	PublishDate time.Time
	ModifyDate  time.Time
	PrevSlug    string
	NextSlug    string
}

type longPost struct {
	shortPost
	ReadingTime int
	Content     template.HTML
}

func newPostTemplate(p longPost, i int, recent []longPost, all []longPost, conf *config) *postTempalte {
	pt := &postTempalte{
		SiteName: conf.siteName,
		Post:     p,
		Recent:   recent,
		All:      all,
	}

	if i > 0 {
		pt.Prev = all[i-1].shortPost
	}

	if i < len(all)-1 {
		pt.Next = all[i+1].shortPost
	}

	return pt
}

func readFrontMatter(s *bufio.Scanner) (map[string]string, error) {
	m := make(map[string]string)
	flag := false
	for s.Scan() {
		line := strings.Trim(s.Text(), " ") //Trim space
		if line == "---" {
			if flag {
				return m, nil
			}
			flag = true
		} else if flag {
			matter := strings.SplitN(line, ":", 2)
			if len(matter) != 2 {
				//Invalid front matter line
				return nil, errInvalidFrontMatter
			}
			m[strings.ToLower(matter[0])] = strings.Trim(matter[1], " ") //Trim space
		} else if line != "" {
			//Empty front matter
			return nil, errMissingFrontMatter
		}
	}

	if err := s.Err(); err != nil {
		return nil, err
	}

	return nil, errEmptyPost
}

func newLongPost(file os.FileInfo, postChan chan<- longPost) {
	f, err := os.Open(filepath.Join(postsDir, file.Name()))
	if err != nil {
		return
	}
	defer f.Close()

	s := bufio.NewScanner(f)
	m, err := readFrontMatter(s)
	if err != nil {
		return
	}

	slug := getSlug(strings.TrimSuffix(file.Name(), ".md"))
	fmt.Println(file.Name())
	// fmt.Println(m["Date"])
	pubDate := file.ModTime()
	if date, ok := m["date"]; ok && len(date) > 0 {
		pubDate, err = time.Parse(dateFormatter[len(date)], date)
		if err != nil {
			return
		}
	}

	var tags []string
	if tag, ok := m["tags"]; ok {
		tags = getTags(tag)
	}

	shortpost := shortPost{
		slug,
		m["author"],
		m["title"],
		m["description"],
		m["category"],
		tags,
		pubDate,
		file.ModTime(),
		"",
		"",
	}
	fmt.Println(m["category"])

	//Read real post
	buf := bytes.NewBuffer(nil)
	for s.Scan() {
		buf.WriteString(s.Text() + "\n")
	}
	if err = s.Err(); err != nil {
		return
	}
	markdown := getMarkdownRender(buf.Bytes())

	longpost := longPost{
		shortpost,
		getReadingTime(string(markdown)),
		template.HTML(markdown),
	}

	postChan <- longpost
}

func getSlug(filename string) (slug string) {
	//TODO: remove date from final slug ?
	re, _ := regexp.Compile(`[^\w\s-]`)
	slug = re.ReplaceAllLiteralString(filename, "")

	re, _ = regexp.Compile(`[-\s]+`)
	slug = re.ReplaceAllLiteralString(slug, "-")

	slug = strings.ToLower(slug) + ".html"
	return
}

func getTags(tag string) []string {
	t := strings.Replace(tag, " ", "", -1) //Why Trim function doesn't work here ?
	return strings.Split(t, ",")
}

func getReadingTime(content string) int {
	readingTime := strings.Count(content, "") / 400
	if readingTime < 1 {
		readingTime = 1
	}

	return readingTime
}

func getMarkdownRender(content []byte) []byte {
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

	return blackfriday.Markdown(content, renderer, extensions)
}
