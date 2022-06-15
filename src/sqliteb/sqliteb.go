package sqliteb

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
	"strconv"
	"github.com/jmoiron/sqlx"
	// _ "github.com/mattn/go-sqlite3" //sqlite3
)

/*Con ...*/
type Con struct {
	DB *sqlx.DB
}

//Opendb ...
func (c *Con) Opendb() {
	db, err := sqlx.Connect("sqlite3", "./buy.db")
	defer db.Close()
	if err != nil {
		log.Fatal(err)
	}
	c.DB = db
}

//NewCon ...
func NewCon() *Con {
	var c = new(Con)
	c.Opendb()
	return c
}

//buy ...
type Buy struct {
	ID        int    `db:"id" json:"id"`
	T		  string `db:"t" json:"t"`
	Money	  float64 `db:"money" json:"money"`
	Type	  string `db:"type" json:"type"`
	Ttype	  string `db:"ttype" json:"ttype"`
	Ex		  string `db:"ex" json:"ex"`
	Merchant  string `db:"merchant" json:"merchant"`
	Thing     string `db:"thing" json:"thing"` //good's name
	Trantype  int    `db:"trantype" json:"trantype"` //pay type 0 income,1 pay
	Account   string `db:"account" json:"account"` //account of bank or platform
	Order     string `db:"order" json:"order"` //in app order id
	Morder	  string `db:"morder" json:"morder"` //merchant order id
	Innout	  string `db:"innout" json:"innout"` //in or out state
}

//type ...
type Clas struct {
	ID        int    `db:"id" json:"id"`
	Name	  string `db:"name" json:"name"`
	Money	  float64 `db:"money" json:"money"`
	Ex		  string `db:"ex" json:"ex"`
}

type Sum struct {
	Key	string	`db:"key" json:"key"`
	Count	string	`db:"count" json:"count"`
	Money	string	`db:"money" json:"money"`
}

func (c *Con) List(t, page string) []Buy {
	db := c.DB
	var err error
	var data = []Buy{}
	var sb strings.Builder
	// sb.WriteString("select * from list")
	pagenow, _ := strconv.Atoi(page)
	limit := strconv.Itoa((pagenow-1) * 10)
	sb.WriteString("select * from list where t like '%")
	sb.WriteString(t)
	sb.WriteString("%' order by t desc ")
	sb.WriteString("limit ")
	sb.WriteString(limit)
	sb.WriteString(",10")
	fmt.Println(sb.String())
	err = db.Select(&data, sb.String())
	if err != nil {
		c.haveErr(err)
		return []Buy{}
	}
	return data
}
func (c *Con) New(b Buy) (isdone bool, newid int64){
	isdone = true
	db := c.DB
	st, err := db.Begin()
	if err != nil {
		fmt.Println(err)
		isdone = false
		return 
	}
	isdone, newid = c.doNew(b, st)
	if isdone {
		st.Commit()
	}else{
		st.Rollback()
	}
	return 
}
func (c *Con) BatchNew(bs []Buy) (isdone bool){
	db := c.DB
	st, err := db.Begin()
	if err != nil {
		fmt.Println(err)
		return false
	}
	for i := 0; i < len(bs); i++ {
		b := bs[i]
		res, _ := c.doNew(b, st)
		if res != true{
			st.Rollback()
			return false
		}
	}
	fmt.Println("fined")
	st.Commit()
	return true
}
func (c *Con) doNew(b Buy, st *sql.Tx) (isdone bool, newid int64){
	isdone = true
	res, er1 := st.Exec("insert into list (t,money,type,ex,merchant,thing,trantype,account,`order`,morder,ttype,innout) values (?,?,?,?,?,?,?,?,?,?,?,?)",
	b.T, b.Money, b.Type, b.Ex, b.Merchant, b.Thing, b.Trantype, b.Account, b.Order, b.Morder, b.Ttype, b.Innout)
	if er1 != nil {
		fmt.Println(er1)
		isdone = false
		return 
	}
	newid, _ = res.LastInsertId()
	return
}
// update any table any col any val in db
func (c *Con) Update(table, id, val, col string) (isdone bool) {
	fmt.Println("update tal :", id, col, val)
	isdone = true
	db := c.DB

	var sb strings.Builder
	sb.WriteString("update ")
	sb.WriteString(table)
	sb.WriteString(" set ")
	sb.WriteString(col)
	sb.WriteString("=? where id=?")

	stmt, err := db.Prepare(sb.String())
	defer stmt.Close()
	if err != nil {
		c.haveErr(err)
	}
	_, er1 := stmt.Exec(val, id)
	if er1 != nil {
		c.haveErr(er1)
		isdone = false
		return
	}
	return
}
// delete any table's any line
func (c *Con) Del(table, id string) (isdone bool) {
	isdone = true
	db := c.DB
	var sb strings.Builder
	sb.WriteString("delete from ")
	sb.WriteString(table)
	sb.WriteString(" where id=$1")
	_, er1 := db.Exec(sb.String(), id)
	if er1 != nil {
		c.haveErr(er1)
		isdone = false
		return
	}
	return
}
func (c *Con) ClasList() (types []Clas) {
	db := c.DB
	var err error
	types = []Clas{}
	err = db.Select(&types, "select * from type")
	if err != nil {
		c.haveErr(err)
		return
	}
	return
}
func (c *Con) ClasNew(cs Clas) (isdone bool, newid int64) {
	isdone = true
	db := c.DB
	stmt, err := db.Prepare("insert into type (name,money,ex) values (?,?,?)")
	defer stmt.Close()
	if err != nil {
		c.haveErr(err)
		isdone = false
		return
	}
	res, er1 := stmt.Exec(cs.Name, cs.Money, cs.Ex)
	if er1 != nil {
		c.haveErr(er1)
		isdone = false
		return
	}
	newid, _ = res.LastInsertId()
	return
}
func (c *Con) Sum(month,typeid,ver string) (data []Sum) {
	db := c.DB
	var sb strings.Builder
	para := ""
	group := ""
	if ver == "type" {
		para = "type as key"
		group = "type"
	} else if ver == "month" {
		para = "substr(t,0,11) as key"
		group = "substr(t,0,11)"
	}
	sb.WriteString("select ")
	sb.WriteString(para)
	sb.WriteString(",sum(1) as count,sum(money) as money from list ")
	in := ""
	if month != "" {
		sb.WriteString("where t like '%")
		sb.WriteString(month)
		sb.WriteString("%' ")
		in = "and"
	}
	if typeid != "" {
		if in == "" {
			in = "where"
		}
		sb.WriteString(in)
		sb.WriteString(" type like '%")
		sb.WriteString(typeid)
		sb.WriteString("%' ")
	}
	sb.WriteString("group by ")
	sb.WriteString(group)
	sb.WriteString(" order by t desc")
	var sql = sb.String()
	fmt.Println(sql)
	err := db.Select(&data, sql)
	if err != nil {
		c.haveErr(err)
		return
	}
	return
}
func (c *Con) haveErr(err error) {
	if err.Error() == "no such table: buy" {
		db := c.DB
		sql := `CREATE TABLE "list" (
				"id"  INTEGER NOT NULL,
				"t"  TEXT NOT NULL,
				"money"  REAL NOT NULL,
				"type"  TEXT NOT NULL,
				"ex"  TEXT,
				"merchant" TEXT,
				"thing" TEXT,
				"trantype" TEXT,
				"account" TEXT,
				"order" TEXT,
				"morder" TEXT,
				"ttype" TEXT,
				"innout" TEXT,
				PRIMARY KEY ("id" ASC)
			);
			CREATE TABLE "type" (
				"id"  INTEGER NOT NULL,
				"name"  TEXT NOT NULL,
				"money"  REAL NOT NULL,
				"ex"  TEXT,
				PRIMARY KEY ("id" ASC)
			)
			
			`
		_, err := db.Exec(sql)
		if err != nil {
			log.Fatal("database error")
			return
		}
	} else {
		fmt.Println(err)
	}
}
