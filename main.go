package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// 豆瓣 xml 描述结构体
type Attributes struct {
	XMLName xml.Name `xml:"rss"`
	Version string   `xml:"version,attr"`
	Channel Channel  `xml:"channel"`
}

// XML 主题结构拆分
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

// 获取 xml 文件数据
func getXMLData(url string) (data []byte, err error) {
	// 读取 xml 文件
	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	// 自定义Header
	req.Header.Set("User-Agent", "Mozilla/4.0 (compatible; MSIE 6.0; Windows NT 5.1)")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close() // 关闭文件
	// 读取所有文件内容保存至 []byte
	data, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return
}

func main() {
	url := "https://www.douban.com/feed/people/debuginn/interests"
	data, err := getXMLData(url)
	if err != nil {
		fmt.Printf("read xml file faild, err:%v\n", err)
	}

	v := Attributes{}
	unMarshalErr := xml.Unmarshal(data, &v)
	if unMarshalErr != nil {
		fmt.Printf("xml unmarshal faild, err:%v\n", err)
	}

	movieItem := v.Channel.MovieItem
	MoviesMap := make(map[int]interface{})

	for i := 0; i < len(movieItem); i++ {
		movie := make(map[string]string)
		description := strings.Split(movieItem[i].Description, "\"")
		movie["Title"] = string([]rune(movieItem[i].Title)[2:])
		movie["Link"] = movieItem[i].Link
		movie["Description"] = description[7]
		movie["Pubdate"] = movieItem[i].Pubdate

		MoviesMap[i] = movie
	}

	//循环输出
	for k, v := range MoviesMap {
		fmt.Println(k, v)
	}

}
