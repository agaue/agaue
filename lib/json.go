package lib

import (
	// "encoding/json"
	"os"
)

type Json struct {
	SiteName string      `json:"siteName"`
	Posts    []*PostJson `json:"posts"`
}

type PostJson struct {
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

func NewJson(siteName string) *Json {
	json := &Json{
		SiteName: siteName,
		Posts:    make([]*PostJson, 0),
	}
}

func NewPostJson(
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
	content string) *PostJson {
	return &PostJson{
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

func (json *Json) AppendPostJson(p *PostJson) {
	json.Posts = append(json.Posts, p)
}

func (json *Json) WriteToFile(path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

}
