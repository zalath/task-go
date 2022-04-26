package main

import (
	"fmt"
	"net/http"
	"tasktask/src/el"
	"tasktask/src/middleware"
	"tasktask/src/buy"
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
	r.POST("/bList", bList)
	r.POST("/bBuy", bBuy)
	r.POST("/bClas", bClas)
	r.POST("/bDel", bDel)
	r.POST("/bNewClas", bNewClas)
	r.POST("/bDelClas", bDelClas)
	r.POST("/bSumMonth", bSumMonth)
	r.POST("/bSumType", bSumType)
	r.POST("/bCsv",bCsv)

	r.POST("/list", list) //get a list of els
	r.POST("/el", getel)
	r.POST("/new", new)
	r.POST("/tik", tik)
	r.POST("/save", save)
	r.POST("/move", move)
	r.POST("/space", space) // get a formed tree of els
	r.POST("/del", del)
	fmt.Printf("running");
	r.Run(":10488") // listen and serve on 0.0.0.0:8888
}
func bList(c *gin.Context) {
	res := buy.List(c.PostForm("t"),c.PostForm("page"))
	c.JSON(http.StatusOK, res)
}
func bBuy(c *gin.Context) {
	res := buy.Buy(c)
	c.JSON(http.StatusOK, res)
}
func bClas(c *gin.Context) {
	res := buy.Clas()
	c.JSON(http.StatusOK, res)
}
func bDel(c *gin.Context) {
	res := buy.Del(c.PostForm("id"))
	c.JSON(http.StatusOK, res)
}
func bNewClas(c *gin.Context) {
	res := buy.NewClas(c)
	c.JSON(http.StatusOK, res)
}
func bDelClas(c *gin.Context) {
	res := buy.DelClas(c.PostForm("id"))
	c.JSON(http.StatusOK, res)
}
func bSumMonth(c *gin.Context) {
	res := buy.Sum(c.PostForm("month"),c.PostForm("typeid"),"month")
	c.JSON(http.StatusOK, res)
}
func bSumType(c *gin.Context) {
	res := buy.Sum(c.PostForm("month"), c.PostForm("typeid"),"type")
	c.JSON(http.StatusOK, res)
}
func bCsv(c *gin.Context) {
	res := buy.Csv("wx")
	c.JSON(http.StatusOK, res)
}

func list(c *gin.Context) {
	id := c.PostForm("id")
	tik := c.PostForm("tik")
	res := el.List(id, "list", tik)
	c.JSON(http.StatusOK, res)
}
func space(c *gin.Context) {
	id := c.PostForm("id")
	tik := c.PostForm("tik")
	res := el.List(id, "", tik)
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
	if res == "done" {
		res = el.Save(c.PostForm("id"), c.PostForm("cmt"), "cmt")
	}
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
