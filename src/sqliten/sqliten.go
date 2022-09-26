package sqliten

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
	"os"
	"path/filepath"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3" //sqlite3
)
var Istest = false
/*Con ...*/
type Con struct {
	DB *sqlx.DB
}

//Opendb ...
func (c *Con) Opendb() {
	pathp, _ := os.Executable()
	path := filepath.Dir(pathp)
	if Istest == true {
		path = "."
	}
	db, err := sqlx.Connect("sqlite3", path+"/note.db")
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

//El ...
type El struct {
	ID        int    `db:"id" json:"id"`
	Title     string `db:"title" json:"title"`
	Tik       int    `db:"tik" json:"tik"`
	P         string `db:"p" json:"p"`
	Pid       int    `db:"pid" json:"pid"`
	Ct        int    `db:"ct" json:"ct"`
	Begintime string `db:"begintime" json:"begintime"`
	Endtime   string `db:"endtime" json:"endtime"`
	Cmt       string `db:"cmt" json:"cmt"`
	Content		string `db:"content" json:"content"`
	File			string `db:"file" json:"file"`
	Tikc	  []Tikc	`json:"tikc"`
	Child     interface{}
}
type Tikc struct {
	Tik int `db:"tik" json:"tik"`
	C int `db:"c" json:"c"`
}

//List a test
func (c *Con) List(id, etype, tik string) []El {
	db := c.DB
	var err error
	var data = []El{}
	var where = ""
	if tik != "" {
		where = " and tik = " + tik
	}	
	if etype == "list" {
		err = db.Select(&data, "select * from e where pid = ? " + where + " order by tik asc,id desc", id)
	} else {
		if id == "0" {
			id = ","
		}
		err = db.Select(&data, "select * from e where p like '%'||$1||'%' " + where + "  order by tik asc,id desc", id)
	}
	for i := 0; i < len(data); i++ {
		data[i].Tikc  = c.Count(data[i].ID)
	}
	if err != nil {
		c.haveErr(err)
	}
	return data
}

//Get countings for one line
func (c *Con) Count(id int) []Tikc {
	db := c.DB
	var d = []Tikc{}
	err := db.Select(&d, "select tik,sum(1) as c from e where pid = ? group by tik order by tik", id)
	if err != nil {
		c.haveErr(err)
	}
	return d
}

//Get select online from database on id
func (c *Con) Get(id string) El {
	db := c.DB
	el := El{}
	err := db.Get(&el, "select id,title,tik,p,pid,ct,cmt,content,file from e where id = ?", id)
	if err != nil {
		c.haveErr(err)
	}
	el.Tikc = c.Count(el.ID)
	return el
}
func (c *Con) Find(key string) []El {
	db := c.DB
	els := []El{}
	err := db.Select(&els, "select * from e where title like '%'||$1||'%' order by pid desc", key)
	if err != nil {
		c.haveErr(err)
	}
	for i := 0; i < len(els); i++ {
		els[i].Tikc  = c.Count(els[i].ID)
	}
	return els
}
//New create a new element
func (c *Con) New(el El) (isdone bool, newid int64) {
	isdone = true
	db := c.DB
	stmt, err := db.Prepare("insert into e (title,pid,p,tik,begintime,endtime,cmt,content,file) values(?,?,?,?,?,?,?,?,?)")
	defer stmt.Close()
	if err != nil {
		c.haveErr(err)
		isdone = false
		return
	}
	res, er1 := stmt.Exec(el.Title, el.Pid, el.P, el.Tik, el.Begintime, el.Endtime, el.Cmt, el.Content, el.File)
	if er1 != nil {
		c.haveErr(er1)
		isdone = false
		return
	}
	newid, _ = res.LastInsertId()
	return
}

//Del delete an element
func (c *Con) Del(id string) (isdone bool) {
	isdone = true
	db := c.DB
	_, err := db.Exec("delete from e where p like '%'||$1||'%'", id)
	if err != nil {
		c.haveErr(err)
		isdone = false
		return
	}
	_, er1 := db.Exec("delete from e where id=$1", id)
	if er1 != nil {
		c.haveErr(err)
		isdone = false
		return
	}
	return
}

//Update ...
func (c *Con) Update(id, val, col string) (isdone bool) {
	fmt.Println("update el :", id, col, val)
	isdone = true
	db := c.DB

	var sb strings.Builder
	sb.WriteString("update e set ")
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

//UpdateP ...
func (c *Con) UpdateP(p, np, id string) (isdone bool) {
	fmt.Println("update p:", p, np)
	isdone = true
	db := c.DB
	var stmt *sql.Stmt
	var err error
	if p == "," {
		stmt, err = db.Prepare("update e set `p` = replace(`p`,?,?) where `id`=?")
	} else {
		stmt, err = db.Prepare("update e set `p` = replace(`p`,?,?) where `p` like '%'||?||'%'")
	}
	defer stmt.Close()
	if err != nil {
		c.haveErr(err)
	}
	var er1 error
	if p == "," {
		_, er1 = stmt.Exec(p, np, id)
	} else {
		_, er1 = stmt.Exec(p, np, p)
	}

	if er1 != nil {
		c.haveErr(er1)
		isdone = false
		return
	}
	return
}

func (c *Con) haveErr(err error) {
	if err.Error() == "no such table: e" {
		db := c.DB
		sql := `CREATE TABLE "e" (
			"id"  INTEGER NOT NULL,
			"title"  TEXT NOT NULL,
			"tik"  INTEGER NOT NULL,
			"pid"  INTEGER NOT NULL,
			"p"  TEXT,
			"ct"  INTEGER NOT NULL DEFAULT 0,
			"begintime" TEXT,
			"endtime" TEXT,
			"cmt" TEXT,
			"content" TEXT,
			"file" TEXT,
			PRIMARY KEY ("id" ASC)
			);
			insert into e(id,title,tik,pid,p,ct,begintime,endtime,cmt,content,file) values('0','======','0','-1','','0','','','','','');
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
