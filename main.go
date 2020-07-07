package main

import (
	"net/http"
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
	r.POST("/els", list)
	r.POST("/el", getel)
	r.POST("/new", new)
	r.POST("/tik", tik)
	r.POST("/save", save)
	r.POST("/move", move)
	r.POST("/space", space)
	r.Run(":8888") // listen and serve on 0.0.0.0:8080
}

func list(c *gin.Context) {
	id := c.PostForm("id")
	res := el.List(id, "list")
	c.JSON(http.StatusOK, res)
}
func space(c *gin.Context) {
	id := c.PostForm("id")
	res := el.List(id, "")
	c.JSON(http.StatusOK, res)
}
func new(c *gin.Context) {
	res := el.New(c)
	c.JSON(http.StatusOK, res)
}
func getel(c *gin.Context) {
	res := el.GetEl(c.PostForm("id"))
	c.JSON(http.StatusOK, res)
}
func tik(c *gin.Context) {
	res := el.Save(c.PostForm("id"), c.PostForm("tik"), "tik")
	c.JSON(http.StatusOK, res)
}
func save(c *gin.Context) {
	res := el.Save(c.PostForm("id"), c.PostForm("title"), "title")
	c.JSON(http.StatusOK, res)
}
func move(c *gin.Context) {
	res := el.Move(c.PostForm("id"), c.PostForm("npid"))
	c.JSON(http.StatusOK, res)
}
