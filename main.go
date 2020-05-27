package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
)

// 豆瓣 xml 描述结构体
type Attributes struct {
	XMLName xml.Name `xml:"rss"`
	Version string   `xml:"version,attr"`
	Channel Channel  `xml:"channel"`
}

type Channel struct {
	Title       string      `xml:"title"`
	Link        string      `xml:"link"`
	Description string      `xml:"description"`
	Language    string      `xml:"language"`
	Copyright   string      `xml:"copyright"`
	Pubdate     string      `xml:"pubDate"`
	MovieItem   []MovieItem `xml:"item"`
}

// 豆瓣 电影列表结构体
type MovieItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	Pubdate     string `xml:"pubDate"`
}

type MovieFmtItem struct {
	Title   string
	Link    string
	ImgURL  string
	Pubdate string
}

func main() {
	// 读取 xml 文件
	client := &http.Client{}
	req, _ := http.NewRequest("GET", "https://www.douban.com/feed/people/debuginn/interests", nil)
	// 自定义Header
	req.Header.Set("User-Agent", "Mozilla/4.0 (compatible; MSIE 6.0; Windows NT 5.1)")
	resp, err := client.Do(req)

	if err != nil {
		fmt.Printf("open xml file faild, err:%v\n", err)
		return
	}
	defer resp.Body.Close() // 关闭文件

	data, readAllErr := ioutil.ReadAll(resp.Body)

	if readAllErr != nil {
		fmt.Printf("read xml body faild, err:%v\n", readAllErr)
		return
	}

	v := Attributes{}
	unMarshalErr := xml.Unmarshal(data, &v)
	if unMarshalErr != nil {
		fmt.Printf("xml unmarshal faild, err:%v\n", err)
	}

	MovieItem := v.Channel.MovieItem

	fmt.Printf("%#v\n", MovieItem)

}
