package main

import (
	"fmt"
	"net/http"
	"tasktask/src/el"
	"tasktask/src/middleware"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.Use(middleware.Cors())
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.POST("/list", list) //get a list of els
	r.POST("/el", getel)
	r.POST("/new", new)
	r.POST("/tik", tik)
	r.POST("/save", save)
	r.POST("/move", move)
	r.POST("/space", space) // get a formed tree of els
	r.POST("/del", del)
	r.Run(":8888") // listen and serve on 0.0.0.0:8888
}

func list(c *gin.Context) {
	id := c.PostForm("id")
	res := el.List(id, "list")
	c.JSON(http.StatusOK, res)
}
func space(c *gin.Context) {
	id := c.PostForm("id")
	fmt.Println(id)
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
	if c.PostForm("tik") == "2" {
		el.Save(c.PostForm("id"), time.Now().Format("2006-1-2 15:04:05"), "endtime")
	} else {
		el.Save(c.PostForm("id"), "", "endtime")
	}
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
func del(c *gin.Context) {
	res := el.Del(c.PostForm("id"))
	c.JSON(http.StatusOK, res)
}
