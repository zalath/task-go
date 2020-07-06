package el

import (
	"fmt"
	"strconv"
	"tasktask/src/sqlitem"
	dbt "tasktask/src/sqlitem"

	"github.com/gin-gonic/gin"
)

//List return direct child of selected node
func List(id, etype string) []dbt.El {
	db := newdb()
	data := db.List(id, etype)
	idInt, _ := strconv.Atoi(id)
	var res = loopFormChild(data, idInt)
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
func New(c *gin.Context) string {
	db := newdb()
	el := formEl(c, db)
	res, _ := db.New(el)
	result := "done"
	if !res {
		result = "mis"
	}
	return result
}
func formEl(c *gin.Context, db *dbt.Con) sqlitem.El {
	var el = sqlitem.El{}
	el.Title = c.PostForm("name")
	el.Pid, _ = strconv.Atoi(c.PostForm("pid"))
	if el.Pid == 0 {
		el.P = ","
	} else {
		pel := db.Get(c.PostForm("pid"))
		el.P = fmt.Sprintf("%s%s,", pel.P, c.PostForm("pid"))
	}
	el.Tik = 0
	return el
}

//GetEl ...
func GetEl(id string) dbt.El {
	db := newdb()
	res := db.Get(id)
	return res
}

//EltoEdit display edit page
func EltoEdit(id string) {
}

//Save submit saving element
func Save() string {
	return ""
}

//Tik change an element's state
func Tik(id string) string {
	db := newdb()
	data := db.Get(id)
	var tik int
	switch data.Tik {
	case 1:
		tik = 2
	case 2:
		tik = 3
	case 3:
		tik = 1
	}
	return "done"
}

//Move change an element's pid and p
func Move() {

}

//Ct get child node's count
func Ct() {

}

//RefreshCt refresh an element's ct count
func RefreshCt() {

}

func newdb() *dbt.Con {
	db := dbt.NewCon()
	return db
}
