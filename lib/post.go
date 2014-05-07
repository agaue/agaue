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
	ErrEmptyPost          = fmt.Errorf("Empty Markdown File!")
	ErrInvalidFrontMatter = fmt.Errorf("Invalid Front Matter")
	ErrMissingFrontMatter = fmt.Errorf("Missing Front Matter")

	dateFormatter = map[int]string{
		//Anyone tell me what happen on 2006/01/02
		10: "2006-01-02",
		13: "2006-01-02 15h",
		14: "2006-01-02 15h",
		15: "2006-01-02 15:04",
		16: "2006-01-02 15:04",
		25: time.RFC3339,
	}
)

type PostTempalte struct {
	SiteName string
	RssURL   string
	Post     *LongPost
	Recent   []*LongPost
	Prev     *ShortPost
	Next     *ShortPost
}

type ShortPost struct {
	Slug        string
	Author      string
	Title       string
	Description string
	Category    string
	Tags        []string
	PublishDate time.Time
	ModifyDate  time.Time
}

type LongPost struct {
	*ShortPost
	ReadingTime int
	Content     template.HTML
}

func newPostTempalte(p *LongPost, i int, recent []*LongPost, all []*LongPost, config Config) *PostTempalte {
	pt := &PostTempalte{
		SiteName: config.SiteName,
		Post:     p,
		Recent:   recent,
	}

	if i > 0 {
		pt.Prev = all[i-1].ShortPost
	}

	if i < len(all)-1 {
		pt.Next = all[i+1].ShortPost
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
			} else {
				flag = true
			}
		} else if flag {
			matter := strings.SplitN(line, ":", 2)
			if len(matter) != 2 {
				//Invalid front matter line
				return nil, ErrInvalidFrontMatter
			}
			m[strings.ToLower(matter[0])] = strings.Trim(matter[1], " ") //Trim space
		} else if line != "" {
			//Empty front matter
			return nil, ErrMissingFrontMatter
		}
	}

	if err := s.Err(); err != nil {
		return nil, err
	}

	return nil, ErrEmptyPost
}

func newLongPost(file os.FileInfo) (*LongPost, error) {
	f, err := os.Open(filepath.Join(PostsDir, file.Name()))
	if err != nil {
		return nil, err
	}
	defer f.Close()

	s := bufio.NewScanner(f)
	m, err := readFrontMatter(s)
	if err != nil {
		return nil, err
	}

	slug := getSlug(strings.TrimSuffix(file.Name(), ".md"))
	fmt.Println(file.Name())
	pubDate := file.ModTime()
	if date, ok := m["Date"]; ok && len(date) > 0 {
		pubDate, err = time.Parse(dateFormatter[len(date)], date)
		if err != nil {
			return nil, err
		}
	}

	var tags []string
	if tag, ok := m["tags"]; ok {
		tags = getTags(tag)
	}

	shortPost := &ShortPost{
		slug,
		m["author"],
		m["title"],
		m["description"],
		m["category"],
		tags,
		pubDate,
		file.ModTime(),
	}

	//Read real post
	buf := bytes.NewBuffer(nil)
	for s.Scan() {
		buf.WriteString(s.Text() + "\n")
	}
	if err = s.Err(); err != nil {
		return nil, err
	}
	markdown := getMarkdownRender(buf.Bytes())

	longPost := &LongPost{
		shortPost,
		getReadingTime(string(markdown)),
		template.HTML(markdown),
	}

	return longPost, nil
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
