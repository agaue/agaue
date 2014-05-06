package lib

import (
	"encoding/xml"
	"os"
	"time"
)

type Rss struct {
	XMLName  xml.Name   `xml:"rss"`
	Version  string     `xml:"version,attr"`
	Channels []*Channel `xml:"channel"`
}

type Channel struct {
	Title         string   `xml:"title"`
	Description   string   `xml:"description"`
	Link          string   `xml:"link"`
	LastBuildDate string   `xml:"lastBuildDate"`
	Generator     string   `xml:"generator"`
	Image         []*Image `xml:"image"`
	Item          []*Item  `xml:"item"`
}

type Image struct {
	Url   string `xml:"url"`
	Title string `xml:"title"`
	Link  string `xml:"link"`
}

type Item struct {
	Title       string   `xml:"title"`
	Link        string   `xml:"link"`
	Description string   `xml:"description"`
	Author      string   `xml:"author"`
	Category    string   `xml:"category"`
	PublishDate string   `xml:"publishDate"`
	Image       []*Image `xml:"image"`
}

func NewRss(title string, description string, link string, author string) *Rss {
	rss := &Rss{
		Version: "2.0",
		Channels: []*Channel{
			&Channel{
				Title:       title,
				Description: description,
				Link:        link,
				Generator:   author,
				Image:       make([]*Image, 0),
				Item:        make([]*Item, 0),
			},
		},
	}

	return rss
}

func NewRssItem(title string, description string, link string, author string, category string, publishDate time.Time) *Item {
	return &Item{
		Title:       title,
		Link:        link,
		Description: description,
		Author:      author,
		Category:    category,
		PublishDate: publishDate.Format(time.RFC822),
		Image:       make([]*Image, 0),
	}
}

func (ch *Channel) AppendItem(i *Item) {
	ch.Item = append(ch.Item, i)
}

func (rss *Rss) WriteToFile(path string) error {
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
