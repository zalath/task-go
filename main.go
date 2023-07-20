package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"tasktask/src/buy"
	"tasktask/src/el"
	"tasktask/src/file"
	"tasktask/src/middleware"
	"tasktask/src/note"
	dbb "tasktask/src/sqliteb"
	dbm "tasktask/src/sqlitem"
	dbn "tasktask/src/sqliten"
	"time"

	"github.com/gin-gonic/gin"
)

var istest = false

func main() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.LoggerWithWriter(gin.DefaultWriter, "/h"), gin.Recovery())
	r.Use(middleware.Cors())
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	//--test start--测试环境打开本行，调整数据库位置
	setTest("test")
	//--test end--------
	if istest {
		r.StaticFS("/f", http.Dir("./pic"))
	} else {
		pathp, _ := os.Executable()
		path := filepath.Dir(pathp) + "/../storage/dcim/taskres/"
		r.StaticFS("/f", http.Dir(path))
	}
	r.POST("/bList", bList)
	r.POST("/bBuy", bBuy)
	r.POST("/bClas", bClas)
	r.POST("/bDel", bDel)
	r.POST("/bNewClas", bNewClas)
	r.POST("/bDelClas", bDelClas)
	r.POST("/bSumMonth", bSumMonth)
	r.POST("/bSumType", bSumType)
	r.POST("/bCsv", bCsv)

	r.POST("/nlist", nlist)
	r.POST("/nel", ngetel)
	r.POST("/nnew", nnew)
	r.POST("/ntik", ntik)
	r.POST("/nsave", nsave)
	r.POST("/nmove", nmove)
	r.POST("/nspace", nspace)
	r.POST("/ndel", ndel)
	r.POST("/nfind", nfind)

	r.POST("/list", list)
	r.POST("/el", getel)
	r.POST("/new", new)
	r.POST("/tik", tik)
	r.POST("/save", save)
	r.POST("/move", move)
	r.POST("/space", space)
	r.POST("/del", del)
	r.POST("/find", find)

	r.POST("/keylist", keylist)
	r.POST("/keynew", keynew)
	r.POST("/keydel", keydel)
	r.POST("/keyupdate", keyupdate)
	r.POST("/keyget", keyget)
	r.GET("/keytest", keytest)

	// get file upload
	r.POST("/fupload", fupload)
	r.POST("/fdel", fdel)
	fmt.Println("running at 10488")

	r.Run(":10488") // listen and serve on 0.0.0.0:8888
}
func bList(c *gin.Context) {
	res := buy.List(c.PostForm("t"), c.PostForm("page"))
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
	res := buy.Sum(c.PostForm("month"), c.PostForm("typeid"), "month")
	c.JSON(http.StatusOK, res)
}
func bSumType(c *gin.Context) {
	res := buy.Sum(c.PostForm("month"), c.PostForm("typeid"), "type")
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
	fmt.Println(c)
	res := note.Save(c.PostForm("id"), c.PostForm("title"), "title")
	if res == "done" {
		res = note.Save(c.PostForm("id"), c.PostForm("cmt"), "cmt")
		if res == "done" {
			res = note.Save(c.PostForm("id"), c.PostForm("content"), "content")
			if res == "done" {
				res = note.Save(c.PostForm("id"), c.PostForm("file"), "file")
			}
		}
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
		if res == "done" {
			res = el.Save(c.PostForm("id"), c.PostForm("file"), "file")
		}
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

func keylist(c *gin.Context) {
	id := c.PostForm("id")
	res := el.Keylist(id)
	c.JSON(http.StatusOK, res)
}

func keynew(c *gin.Context) {
	res := el.Keynew(c.PostForm("name"), c.PostForm("val"), c.PostForm("pid"))
	c.JSON(http.StatusOK, res)
}

func keydel(c *gin.Context) {
	id := c.PostForm("id")
	res := el.KeyDel(id)
	c.JSON(http.StatusOK, res)
}

func keyupdate(c *gin.Context) {
	id := c.PostForm("id")
	name := c.PostForm("name")
	val := c.PostForm("val")
	res := el.Keyupdate(id, name, "name")
	if res == "done" {
		res = el.Keyupdate(id, val, "val")
	}
	c.JSON(http.StatusOK, res)
}

func keyget(c *gin.Context) {
	id := c.PostForm("id")
	res := el.KeyGet(id)
	c.JSON(http.StatusOK, res)
}

func keytest(c *gin.Context) {
	id := el.Keynew("tnode", "tval", "0")
	fmt.Println(id)
	fmt.Println("new ready")
	res1 := el.Keylist("0")
	fmt.Println(res1)
	fmt.Println("list ready")
	res2 := el.Keyupdate(id, "tn1", "name")
	fmt.Println("update" + id)
	fmt.Println(res2)
	res := el.KeyGet(id)
	fmt.Println(res)
	fmt.Println("get ready")
	res4 := el.KeyDel(id)
	fmt.Println("del" + id)
	fmt.Println(res4)
}

func setTest(v string) {
	istest = true
	dbm.Istest = true
	dbn.Istest = true
	dbb.Istest = true
	file.Istest = true
}

func fupload(c *gin.Context) {
	res := file.Del(c)
	res = file.Upload(c)
	c.JSON(http.StatusOK, res)
}
func fdel(c *gin.Context) {
	res := file.Del(c)
	c.JSON(http.StatusOK, res)
}
