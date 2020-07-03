package main

import (
	"tasktask/src/el"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.POST("/list", list)
	r.POST("/list1", list1)
	r.Run(":8888") // listen and serve on 0.0.0.0:8080
}

func list(c *gin.Context) {
	id := c.PostForm("id")
	etype := c.PostForm("type")
	el.List(id, etype)
}
func list1(c *gin.Context) {
	id := c.PostForm("id")
	etype := c.PostForm("type")
	el.List1(id, etype)
}
