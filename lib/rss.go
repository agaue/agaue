package lib

import (
	"encoding/xml"
	"os"
	"time"
)

type rss struct {
	XMLName  xml.Name   `xml:"rss"`
	Version  string     `xml:"version,attr"`
	Channels []*channel `xml:"channel"`
}

type channel struct {
	Title         string   `xml:"title"`
	Description   string   `xml:"description"`
	Link          string   `xml:"link"`
	LastBuildDate string   `xml:"lastBuildDate"`
	Generator     string   `xml:"generator"`
	Image         []*image `xml:"image"`
	Item          []*item  `xml:"item"`
}

type image struct {
	URL   string `xml:"url"`
	Title string `xml:"title"`
	Link  string `xml:"link"`
}

type item struct {
	Title       string   `xml:"title"`
	Link        string   `xml:"link"`
	Description string   `xml:"description"`
	Author      string   `xml:"author"`
	Category    string   `xml:"category"`
	PublishDate string   `xml:"publishDate"`
	Image       []*image `xml:"image"`
}

func newRss(title string, description string, link string, author string) *rss {
	r := &rss{
		Version: "2.0",
		Channels: []*channel{
			&channel{
				Title:       title,
				Description: description,
				Link:        link,
				Generator:   author,
				Image:       make([]*image, 0),
				Item:        make([]*item, 0),
			},
		},
	}

	return r
}

func newRssItem(title string, description string, link string, author string, category string, publishDate string) *item {
	return &item{
		Title:       title,
		Link:        link,
		Description: description,
		Author:      author,
		Category:    category,
		PublishDate: publishDate,
		Image:       make([]*image, 0),
	}
}

func (ch *channel) appendItem(i *item) {
	ch.Item = append(ch.Item, i)
}

func (rss *rss) writeToFile(path string) error {
	rss.Channels[0].LastBuildDate = time.Now().Format(time.RFC822)
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.WriteString(xml.Header)
	if err != nil {
		return err
	}
	encoding := xml.NewEncoder(file)
	return encoding.Encode(rss)
}
