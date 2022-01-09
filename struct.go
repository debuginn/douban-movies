package main

import "encoding/xml"

// Attributes 豆瓣 xml 描述结构体
type Attributes struct {
	XMLName xml.Name `xml:"rss"`
	Version string   `xml:"version,attr"`
	Channel Channel  `xml:"channel"`
}

// Channel XML 主题结构拆分
type Channel struct {
	Title       string      `xml:"title"`
	Link        string      `xml:"link"`
	Description string      `xml:"description"`
	Language    string      `xml:"language"`
	Copyright   string      `xml:"copyright"`
	Pubdate     string      `xml:"pubDate"`
	MovieItem   []MovieItem `xml:"item"`
}

// MovieItem 豆瓣 电影列表结构体
type MovieItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	Pubdate     string `xml:"pubDate"`
}
