package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/gin-gonic/gin"
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
		movie["Img"] = description[7]
		movie["Pubdate"] = movieItem[i].Pubdate

		MoviesMap[i] = movie
	}
	//
	//循环输出
	//for k, v := range MoviesMap {
	//	fmt.Println(k, v)
	//}

	data, _ = json.Marshal(MoviesMap)

	r := gin.Default()
	r.Use(Cors())
	r.GET("/doubanmovies", func(context *gin.Context) {
		context.JSON(http.StatusOK, MoviesMap)
	})

	_ = r.Run(":8080")
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

// 跨域
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method               //请求方法
		origin := c.Request.Header.Get("Origin") //请求头部
		var headerKeys []string                  // 声明请求头keys
		for k, _ := range c.Request.Header {
			headerKeys = append(headerKeys, k)
		}
		headerStr := strings.Join(headerKeys, ", ")
		if headerStr != "" {
			headerStr = fmt.Sprintf("access-control-allow-origin, access-control-allow-headers, %s", headerStr)
		} else {
			headerStr = "access-control-allow-origin, access-control-allow-headers"
		}
		if origin != "" {
			c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
			c.Header("Access-Control-Allow-Origin", "*")                                       // 这是允许访问所有域
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE,UPDATE") //服务器支持的所有跨域请求的方法,为了避免浏览次请求的多次'预检'请求
			//  header的类型
			c.Header("Access-Control-Allow-Headers", "Authorization, Content-Length, X-CSRF-Token, Token,session,X_Requested_With,Accept, Origin, Host, Connection, Accept-Encoding, Accept-Language,DNT, X-CustomHeader, Keep-Alive, User-Agent, X-Requested-With, If-Modified-Since, Cache-Control, Content-Type, Pragma")
			//              允许跨域设置                                                                                                      可以返回其他子段
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers,Cache-Control,Content-Language,Content-Type,Expires,Last-Modified,Pragma,FooBar") // 跨域关键设置 让浏览器可以解析
			c.Header("Access-Control-Max-Age", "172800")                                                                                                                                                           // 缓存请求信息 单位为秒
			c.Header("Access-Control-Allow-Credentials", "false")                                                                                                                                                  //  跨域请求是否需要带cookie信息 默认设置为true
			c.Set("content-type", "application/json")                                                                                                                                                              // 设置返回格式是json
		}

		//放行所有OPTIONS方法
		if method == "OPTIONS" {
			c.JSON(http.StatusOK, "Options Request!")
		}
		// 处理请求
		c.Next() //  处理请求
	}
}
