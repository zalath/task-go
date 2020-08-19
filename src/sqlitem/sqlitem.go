package sqlitem

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3" //sqlite3
)

/*Con ...*/
type Con struct {
	DB *sqlx.DB
}

//Opendb ...
func (c *Con) Opendb() {
	db, err := sqlx.Connect("sqlite3", "./db.db")
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
	Child     interface{}
}

//List a test
func (c *Con) List(id, etype string) []El {
	db := c.DB
	var err error
	var data = []El{}
	if etype == "list" {
		// err = db.Select(&data, "select id,title,tik,p,pid,ct,cmt,begintime,endtime from e where pid = ? order by tik,id", id)
		err = db.Select(&data, "select * from e where pid = ? order by tik,id", id)
	} else {
		if id == "0" {
			id = ","
		}
		// err = db.Select(&data, "select id,title,tik,p,pid,ct,cmt,begintime,endtime from e where p like '%'||$1||'%' order by tik,id", id)
		err = db.Select(&data, "select * from e where p like '%'||$1||'%' order by tik,id", id)
	}

	if err != nil {
		c.haveErr(err)
	}
	return data
}

//Get select online from database on id
func (c *Con) Get(id string) El {
	db := c.DB
	el := El{}
	err := db.Get(&el, "select id,title,tik,p,pid,ct from e where id = ?", id)
	if err != nil {
		c.haveErr(err)
	}
	return el
}

//New create a new element
func (c *Con) New(el El) (isdone bool, newid int64) {
	isdone = true
	db := c.DB
	stmt, err := db.Prepare("insert into e (title,pid,p,tik,begintime,cmt) values(?,?,?,?,?,?)")
	if err != nil {
		c.haveErr(err)
	}
	res, er1 := stmt.Exec(el.Title, el.Pid, el.P, el.Tik, el.Begintime, el.Cmt)
	if er1 != nil {
		c.haveErr(err)
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
			PRIMARY KEY ("id" ASC)
			);
			
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
