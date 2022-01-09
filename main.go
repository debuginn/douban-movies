package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
)

func main() {

	m := new(MoviesData)

	// first
	m.SetMoviesData()

	// crond
	c := cron.New()
	c.AddFunc("60 * * * *", func() {
		m.SetMoviesData()
	})

	r := gin.Default()
	r.Use(Cors())
	r.GET("/doubanmovies", func(context *gin.Context) {
		context.JSON(http.StatusOK, m.GetMoviesData())
	})

	r.Run(":8080")
}
