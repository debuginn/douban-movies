package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
)

const (
	feedUrl = "https://www.douban.com/feed/people/debuginn/interests"
)

type MoviesData struct {
	Movies map[int]interface{}
	sync.RWMutex
}

// SetMoviesData .
func (m *MoviesData) SetMoviesData() {

	data, err := m.GetXMLData(feedUrl)
	if err != nil {
		fmt.Printf("read xml file faild, err:%v\n", err)
	}

	v := Attributes{}
	unMarshalErr := xml.Unmarshal(data, &v)
	if unMarshalErr != nil {
		fmt.Printf("xml unmarshal faild, err:%v\n", err)
	}

	movieItem := v.Channel.MovieItem
	MoviesMap := make(map[int]interface{}, len(movieItem))
	for i := 0; i < len(movieItem); i++ {
		movie := make(map[string]string)
		description := strings.Split(movieItem[i].Description, "\"")

		movie["Title"] = string([]rune(movieItem[i].Title)[2:])
		movie["Link"] = movieItem[i].Link
		movie["Img"] = description[7]
		movie["Pubdate"] = movieItem[i].Pubdate

		MoviesMap[i] = movie
	}

	m.Lock()
	m.Movies = MoviesMap
	m.Unlock()
}

// GetMoviesData .
func (m *MoviesData) GetMoviesData() map[int]interface{} {
	m.RLock()
	defer m.RUnlock()
	return m.Movies
}

// GetXMLData 获取 xml 文件数据
func (m *MoviesData) GetXMLData(url string) (data []byte, err error) {
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
