package lib

import (
	"encoding/json"
	"os"
)

type site struct {
	SiteName string      `json:"siteName"`
	Posts    []*postJSON `json:"posts"`
}

type postJSON struct {
	Slug        string `json:"slug"`
	Author      string `json:"author"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Category    string `json:"category"`
	PublishDate string `json:"publishDate"`
	ModifyDate  string `json:"modifyDate"`
	ReadingTime int    `json:"readingTime"`
	PrevSlug    string `json:"prevSlug"`
	NextSlug    string `json:"nextSlug"`
	Content     string `json:"content"`
}

func newSiteJSON(siteName string) *site {
	s := &site{
		SiteName: siteName,
		Posts:    make([]*postJSON, 0),
	}
	return s
}

func newPostJSON(
	slug string,
	author string,
	title string,
	description string,
	category string,
	publishDate string,
	modifyDate string,
	readingTime int,
	prevSlug string,
	nextSlug string,
	content string) *postJSON {
	return &postJSON{
		Slug:        slug,
		Author:      author,
		Title:       title,
		Description: description,
		Category:    category,
		PublishDate: publishDate,
		ModifyDate:  modifyDate,
		ReadingTime: readingTime,
		PrevSlug:    prevSlug,
		NextSlug:    nextSlug,
		Content:     content,
	}
}

func (s *site) appendPostJSON(p *postJSON) {
	s.Posts = append(s.Posts, p)
}

func (s *site) writeToFile(path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	encoding := json.NewEncoder(file)
	return encoding.Encode(s)
}
