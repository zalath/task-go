package el

import (
	"fmt"
	"strconv"
	"strings"
	dbt "tasktask/src/sqlitem"
	"time"

	"github.com/gin-gonic/gin"
)

//List return direct child of selected node
func List(id, etype, tik string) []dbt.El {
	db := newdb()
	data := db.List(id, etype, tik)
	idInt, _ := strconv.Atoi(id)
	if etype == "list" {
		return data
	}
	res := loopFormChild(data, idInt)
	return res
}

func loopFormChild(data []dbt.El, id int) []dbt.El {
	var res = []dbt.El{}
	for i := 0; i < len(data); i++ {
		if data[i].Pid == id {
			el := data[i]
			el.Child = loopFormChild(data, el.ID)
			res = append(res, el)
		}
	}
	return res
}

//New create a new element into database
func New(c *gin.Context) (result string) {
	db := newdb()
	el := formEl(c, db)
	db.DB.Begin()
	res, newid := db.New(el)
	result = strconv.Itoa(int(newid))
	if !res {
		db.DB.MustBegin().Rollback()
		result = "mis"
	} else {
		if !updateCt(el.Pid, "+", db) {
			db.DB.MustBegin().Rollback()
			result = "mis"
		}
	}
	db.DB.MustBegin().Commit()
	return
}

//Del delete el from db
func Del(id string) (result string) {
	db := newdb()
	db.DB.Begin()
	el := GetEl(id)
	res := db.Del(id)
	if !res {
		db.DB.MustBegin().Rollback()
		result = "mis"
	} else {
		if !updateCt(el.Pid, "-", db) {
			db.DB.MustBegin().Rollback()
			result = "mis"
		}
	}

	db.DB.MustBegin().Commit()
	result = "done"
	return
}

func formEl(c *gin.Context, db *dbt.Con) dbt.El {
	var el = dbt.El{}
	el.Title = c.PostForm("title")
	el.Cmt = c.PostForm("cmt")
	el.Pid, _ = strconv.Atoi(c.PostForm("pid"))
	if el.Pid == 0 {
		el.P = ","
	} else {
		pel := db.Get(c.PostForm("pid"))
		el.P = fmt.Sprintf("%s%s,", pel.P, c.PostForm("pid"))
	}
	el.Tik = 1
	el.Ct = 0
	el.Begintime = time.Now().Format("2006-1-2 15:04:05")
	el.Endtime = "-"
	fmt.Printf("%#v", el)
	return el
}

//GetEl ...
func GetEl(id string) dbt.El {
	db := newdb()
	res := db.Get(id)
	return res
}
//find a list of els
func Find(key string) []dbt.El {
	db := newdb()
	res := db.Find(key)
	return res
}
//Save submit saving element
func Save(id, val, col string) string {
	db := newdb()
	res := db.Update(id, val, col)
	if res {
		return "done"
	}
	return "mis"
}

//Move change an element's pid and p
func Move(id, npid string) string {
	db := newdb()
	el := db.Get(id)
	oldP := el.P
	elp := db.Get(npid)
	db.DB.Begin()

	var sb strings.Builder
	sb.WriteString(elp.P)
	sb.WriteString(strconv.Itoa(elp.ID))
	sb.WriteString(",")
	newP := sb.String()

	fmt.Println(oldP)
	fmt.Println(npid)
	fmt.Println(newP)
	res := db.Update(id, npid, "pid")
	if !res {
		fmt.Println("pid update err")
		db.DB.MustBegin().Rollback()
		return "mis"
	}
	res = updateCt(el.Pid, "-", db)
	if !res {
		fmt.Println("old p - err")
		db.DB.MustBegin().Rollback()
		return "mis"
	}
	res = updateCt(elp.ID, "+", db)
	if !res {
		fmt.Println("new p + err")
		db.DB.MustBegin().Rollback()
		return "mis"
	}

	//update new pid's ct val
	res = db.UpdateP(el.P, newP, id)
	if !res {
		fmt.Println("p update err")
		db.DB.MustBegin().Rollback()
		return "mis"
	}
	db.DB.MustBegin().Commit()
	return "done"
}
func updateCt(id int, ctype string, db *dbt.Con) bool {
	cid := strconv.Itoa(id)
	el := db.Get(cid)
	if ctype == "+" {
		res := db.Update(cid, strconv.Itoa(el.Ct+1), "ct")
		if !res {
			return false
		}
	} else if ctype == "-" {
		res := db.Update(cid, strconv.Itoa(el.Ct-1), "ct")
		if !res {
			return false
		}
	}
	return true
}
func newdb() *dbt.Con {
	db := dbt.NewCon()
	return db
}
