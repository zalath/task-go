package buy

import (
	"fmt"
	"strconv"
	dbb "tasktask/src/sqliteb"
	"time"

	"github.com/gin-gonic/gin"
)
// import order list
func Imp(){
	// import file
	// for list
	// save data to db
}
// bill list
func List(t, page string) []dbb.Buy {
	db := newdb()
	fmt.Println(time.Now())
	if t == "" {
		t = time.Now().Format("2006-01")
	}
	if page == "" {
		page = "1"
	}
	data := db.List(t, page)//fetch out this month's record
	// getlist and return ,with paginated
	return data
}
// new bill
func Buy(c *gin.Context) (result string) {
	// make new bill insert
	db := newdb()
	bu := formBuy(c, db)
	db.DB.Begin()
	fmt.Println(bu)
	res, newid := db.New(bu)
	if !res {
		db.DB.MustBegin().Rollback()
		result = "mis"
	} else {
		result = strconv.Itoa(int(newid))
	}
	db.DB.MustBegin().Commit()
	return
}
//del bill
func Del(id string) bool{
	db := newdb()
	res := db.Del("list", id)
	return res
}
// type list
func Clas() []dbb.Clas{
	// get class list and return, with paginated
	db := newdb()
	data := db.ClasList()
	return data
}
// new type
func NewClas(c *gin.Context) (result string) {
	// make new class insert
	db := newdb()
	cls := formClas(c, db)
	db.DB.Begin()
	res, newid := db.ClasNew(cls)
	if !res {
		db.DB.MustBegin().Rollback()
		result = "mis"
	} else {
		result = strconv.Itoa(int(newid))
	}
	db.DB.MustBegin().Commit()
	return 
}
func DelClas(id string) bool{
	db := newdb()
	res := db.Del("type", id)
	return res
}
// form summary data
// month "2022-02-02"
// typeid "12"
// ver "type","date"按哪个字段group by
func Sum(month,typeid,ver string) []dbb.Sum {
	db := newdb()
	res := db.Sum(month, typeid, ver)
	return res
}
func formBuy(c *gin.Context, db *dbb.Con) dbb.Buy {
	var bu = dbb.Buy{}
	bu.T = c.PostForm("t")
	v, _ := strconv.ParseFloat(c.PostForm("money"),64)
	bu.Money = v
	bu.Type = c.PostForm("type")
	bu.Ex = c.PostForm("ex")
	bu.Merchant = c.PostForm("merchant")
	bu.Thing = c.PostForm("thing")
	bu.Trantype, _ = strconv.Atoi(c.PostForm("trantype"))
	bu.Account = c.PostForm("account")
	bu.Order = c.PostForm("order")
	bu.Morder = c.PostForm("morder")
	fmt.Println("buy structure formed")
	return bu
}
func formClas(c *gin.Context, db *dbb.Con) dbb.Clas {
	var cl = dbb.Clas{}
	cl.Name = c.PostForm("name")
	v, _ := strconv.ParseFloat(c.PostForm("money"),64)
	cl.Money = v
	cl.Ex = c.PostForm("ex")
	return cl
}
func newdb() *dbb.Con {
	db := dbb.NewCon()
	return db
}
