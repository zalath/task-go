package main

import (
	"fmt"
	"net/http"
	"tasktask/src/note"
	"tasktask/src/el"
	"tasktask/src/buy"
	"tasktask/src/middleware"
	"time"
	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.LoggerWithWriter(gin.DefaultWriter,"/h"),gin.Recovery())
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

	
	r.POST("/nlist", nlist) //get a list of els
	r.POST("/nel", ngetel)
	r.POST("/nnew", nnew)
	r.POST("/ntik", ntik)
	r.POST("/nsave", nsave)
	r.POST("/nmove", nmove)
	r.POST("/nspace", nspace) // get a formed tree of els
	r.POST("/ndel", ndel)
	r.POST("/nfind", nfind)

	r.POST("/list", list) //get a list of els
	r.POST("/el", getel)
	r.POST("/new", new)
	r.POST("/tik", tik)
	r.POST("/save", save)
	r.POST("/move", move)
	r.POST("/space", space) // get a formed tree of els
	r.POST("/del", del)
	r.POST("/find", find)
	fmt.Println("running at 10488");
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
	res := buy.Csv(c.PostForm("type"), c)
	c.JSON(http.StatusOK, res)
}

func nlist(c *gin.Context) {
	id := c.PostForm("id")
	tik := c.PostForm("tik")
	res := note.List(id, "list", tik)
	c.JSON(http.StatusOK, res)
}
func nspace(c *gin.Context) {
	id := c.PostForm("id")
	tik := c.PostForm("tik")
	res := note.List(id, "", tik)
	c.JSON(http.StatusOK, res)
}
func nnew(c *gin.Context) {
	res := note.New(c)
	c.JSON(http.StatusOK, res)
}
func ngetel(c *gin.Context) {
	res := note.GetEl(c.PostForm("id"))
	c.JSON(http.StatusOK, res)
}
func ntik(c *gin.Context) {
	res := note.Save(c.PostForm("id"), c.PostForm("tik"), "tik")
	if c.PostForm("tik") == "2" {
		note.Save(c.PostForm("id"), time.Now().Format("2006-1-2 15:04:05"), "endtime")
	} else {
		note.Save(c.PostForm("id"), "", "endtime")
	}
	c.JSON(http.StatusOK, res)
}
func nsave(c *gin.Context) {
	res := note.Save(c.PostForm("id"), c.PostForm("title"), "title")
	if res == "done" {
		res = note.Save(c.PostForm("id"), c.PostForm("cmt"), "cmt")
	}
	c.JSON(http.StatusOK, res)
}
func nmove(c *gin.Context) {
	res := note.Move(c.PostForm("id"), c.PostForm("npid"))
	c.JSON(http.StatusOK, res)
}
func ndel(c *gin.Context) {
	res := note.Del(c.PostForm("id"))
	c.JSON(http.StatusOK, res)
}
func nfind(c *gin.Context) {
	res := note.Find(c.PostForm("key"))
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
func find(c *gin.Context) {
	res := el.Find(c.PostForm("key"))
	c.JSON(http.StatusOK, res)
}
