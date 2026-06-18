package sqlitebook

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3" //sqlite3
)

var Istest = false

type Con struct {
	DB *sqlx.DB
}

// Opendb ...
func (c *Con) Opendb() {
	pathp, _ := os.Executable()
	path := filepath.Dir(pathp)
	if Istest == true {
		path = "."
	}
	db, err := sqlx.Connect("sqlite3", path+"/db/book.db")
	if err != nil {
		log.Fatal(err)
	}
	c.DB = db
}

func NewCon() *Con {
	var c = new(Con)
	c.Opendb()
	return c
}

type Part struct {
	ID   int    `db:"id" json:"id"`
	Name string `db:"name" json:"name"`
	Age  int    `db:"age" json:"age"`
	Desc string `db:"desc" json:"desc"`
	Type int    `db:"type" json:"type"`
	Pct  int    `db:"pct" json:"pct"`
	Pic  string `db:"pic" json:"pic"`
}

type PartPs struct {
	ID   int    `db:"id" json:"id"`
	Name string `db:"name" json:"name"`
	Age  int    `db:"age" json:"age"`
	Desc string `db:"desc" json:"desc"`
	Type int    `db:"type" json:"type"`
	Pct  int    `db:"pct" json:"pct"`
	Pic  string `db:"pic" json:"pic"`

	Relationid string `db:"relationid" json:"relationid"`
}

type Relation struct {
	ID         int `db:"id" json:"id"`
	P1         int `db:"p1" json:"p1"`
	P2         int `db:"p2" json:"p2"`
	RelationID int `db:"relationid" json:"relationid"`
	Direction  int `db:"direction" json:"direction"`
	Sort       int `db:"sort" json:"sort"`
}

type RelationType struct {
	ID      int    `db:"id" json:"id"`
	Name    string `db:"name" json:"name"`
	RevName string `db:"revname" json:"revname"`
	Bookid  int    `db:"bookid" json:"bookid"`
}

// 获取书籍节点
func (c *Con) List() []Part {
	db := c.DB
	var err error
	var data = []Part{}
	err = db.Select(&data, "select * from part where type = 0 order by id desc")
	if err != nil {
		c.haveErr(err)
	}
	return data
}

// 通过书籍父节点，获取本书所有子节点
func (c *Con) Parts(id string) []Part {
	db := c.DB
	p1 := []Part{}
	err := db.Select(&p1, "select p.* from relation r join part p on r.p2=p.id  where r.p1 = ?", id)
	if err != nil {
		c.haveErr(err)
	}
	return p1
}

// 按照struct Part的所有项目更新节点内容
func (c *Con) UpdatePart(p Part) bool {
	db := c.DB
	tx, err := db.Begin()
	if err != nil {
		c.haveErr(err)
		return false
	}
	_, err = tx.Exec("update part set name=?,age=?,desc=?,type=?,pic=? where id=?", p.Name, p.Age, p.Desc, p.Type, p.Pic, p.ID)
	if err != nil {
		tx.Rollback()
		c.haveErr(err)
		return false
	}
	err = tx.Commit()
	if err != nil {
		c.haveErr(err)
		return false
	}
	return true
}

// 删除节点
func (c *Con) DeletePart(id int) bool {
	db := c.DB
	_, err := db.Exec("delete from part where id=?", id)
	if err != nil {
		c.haveErr(err)
		return false
	}
	// 删除该节点id在relation表中对应的所有p1和p2的记录
	_, err = db.Exec("delete from relation where p1=? or p2=?", id, id)
	if err != nil {
		c.haveErr(err)
		return false
	}
	return true
}

// 增加节点
// p1 父级节点id
// relationid 关系类型
// relationpos 1为正向关系（p1->id)，2为反向关系（id->p1)

func (c *Con) AddPart(p PartPs, bookid string) PartPs {
	// 插入节点，并根据新id和P1插入relation表
	db := c.DB
	tx, err := db.Begin()
	if err != nil {
		c.haveErr(err)
	}
	res, err := tx.Exec("insert into part (`name`,`age`,`desc`,`pic`,`type`,`pct`) values (?,?,?,?,?,?)", p.Name, p.Age, p.Desc, p.Pic, p.Type, p.Pct)
	if err != nil {
		tx.Rollback()
		c.haveErr(err)
	}
	id, _ := res.LastInsertId()
	p.ID = int(id)
	idStr := strconv.FormatInt(id, 10)
	// 将节点归到书目下
	rid := c.Makerelationsql(tx, bookid, idStr, strconv.Itoa(p.Type), "1")
	if rid == 0 {
		tx.Rollback()
		c.haveErr(err)
	}

	err = tx.Commit()
	if err != nil {
		c.haveErr(err)
	}

	p.Relationid = strconv.Itoa(rid)
	return p
}

// 查询关系类型表
func (c *Con) GetRelationType(bookid string) []RelationType {
	db := c.DB
	var rt []RelationType
	err := db.Select(&rt, "select id,name,revname,bookid from relationtype where bookid=? order by id asc", bookid)
	if err != nil {
		c.haveErr(err)
	}
	return rt
}

// 创建关系类型
func (c *Con) CreateRelationType(name, revname, bookid string) (bool, int) {
	db := c.DB
	tx, err := db.Begin()
	if err != nil {
		c.haveErr(err)
		return false, 0
	}
	sql, err := tx.Exec("insert into relationtype (name,revname,bookid) values (?,?,?)", name, revname, bookid)
	if err != nil {
		tx.Rollback()
		c.haveErr(err)
		return false, 0
	}
	id, _ := sql.LastInsertId()
	err = tx.Commit()
	if err != nil {
		c.haveErr(err)
		return false, 0
	}
	return true, int(id)
}

// 删除关系类型，并将已创建的关系改为不确定，方便后续查找更新。
func (c *Con) DeleteRelationType(id string) bool {
	db := c.DB
	tx, err := db.Begin()
	if err != nil {
		c.haveErr(err)
		return false
	}
	_, err = tx.Exec("delete from relationtype where id=?", id)
	if err != nil {
		c.haveErr(err)
		return false
	}
	_, err = tx.Exec("update relation set relationid=0 where relationid=?", id)
	if err != nil {
		c.haveErr(err)
		return false
	}
	err = tx.Commit()
	if err != nil {
		c.haveErr(err)
		return false
	}
	return true
}

func (c *Con) UpdateRelationType(id, name, revname string) bool {
	db := c.DB
	tx, err := db.Begin()
	if err != nil {
		c.haveErr(err)
		return false
	}
	_, err = tx.Exec("update relationtype set name=?,revname=? where id=?", name, revname, id)
	if err != nil {
		tx.Rollback()
		c.haveErr(err)
		return false
	}
	err = tx.Commit()
	if err != nil {
		c.haveErr(err)
		return false
	}
	return true
}

// 通过父级节点，获得关系列表
func (c *Con) PartRelationList(p1, relationid, p2 string) []Relation {
	db := c.DB
	p := []Relation{}
	query := "select * from relation where 1=1"
	args := []interface{}{}
	var orParts []string
	if p1 != "0" {
		orParts = append(orParts, "p1 = ?")
		args = append(args, p1)
	}
	if p2 != "0" {
		orParts = append(orParts, "p2 = ?")
		args = append(args, p2)
	}
	if len(orParts) > 0 {
		query += " AND (" + strings.Join(orParts, " OR ") + ")"
	}
	if relationid == "0" {
		query += " AND relationid <> ?"
		args = append(args, 2)
	} else {
		query += " AND relationid = ?"
		args = append(args, relationid)
	}
	query += " ORDER BY sort ASC"
	err := db.Select(&p, query, args...)
	if err != nil {
		c.haveErr(err)
	}
	return p
}

func (c *Con) PartRelationMap(bookid string) []Relation {
	db := c.DB
	p1list := []Relation{}
	var err error
	err = db.Select(&p1list, "select r.* from relation r left outer join relationtype rt on r.relationid = rt.id where r.relationid > 2 and rt.bookid=?", bookid)
	if err != nil {
		c.haveErr(err)
	}
	return p1list
}

// 删除关系表记录
func (c *Con) DeleteRelation(id string) bool {
	db := c.DB
	_, err := db.Exec("delete from relation where id=?", id)
	if err != nil {
		c.haveErr(err)
		return false
	}
	return true
}

// 创建关系表记录
func (c *Con) Makerelation(p1, p2, relationid, direction string) bool {
	db := c.DB
	tx, err := db.Begin()
	if err != nil {
		c.haveErr(err)
		return false
	}
	id := c.Makerelationsql(tx, p1, p2, relationid, direction)
	if id == 0 {
		tx.Rollback()
		c.haveErr(err)
		return false
	}
	err = tx.Commit()
	if err != nil {
		c.haveErr(err)
		return false
	}
	return true
}
func (c *Con) Makerelationsql(tx *sql.Tx, p1, p2, relationid, direction string) int {
	sql, err := tx.Exec("insert into relation (p1,p2,relationid,direction) values (?,?,?,?)", p1, p2, relationid, direction)
	if err != nil {
		c.haveErr(err)
		return 0
	}
	id, _ := sql.LastInsertId()
	_, err = tx.Exec("update relation set sort=? where id=?", id, id)
	if err != nil {
		c.haveErr(err)
		return 0
	}
	return int(id)
}

// 更新节点关系的顺序
func (c *Con) PartRelationSetOrder(id, order string) bool {
	// 更新表中的order，使对应的id的order变为新的值, 同时将order等于或者大于该值的order值加1
	db := c.DB
	tx, err := db.Begin()
	if err != nil {
		c.haveErr(err)
		return false
	}
	var currentSort int
	err = db.Get(&currentSort, "SELECT sort FROM relation WHERE id = ?", id)
	if err != nil {
		c.haveErr(err)
		return false
	}
	// 获取id对应的那行数据，将里面的sort作为最大值，下面这条更新仅更新从参数order到最大值之间的数据
	_, err = tx.Exec("update relation set `sort`=`sort`+1 where `sort`>=? and `sort`<?", order, currentSort)
	if err != nil {
		tx.Rollback()
		c.haveErr(err)
		return false
	}
	_, err = tx.Exec("update relation set `sort`=? where id=?", order, id)
	if err != nil {
		tx.Rollback()
		c.haveErr(err)
		return false
	}
	err = tx.Commit()
	if err != nil {
		c.haveErr(err)
		return false
	}
	return true
}
func (c *Con) haveErr(err error) {
	if err.Error() == "no such table: part" {
		db := c.DB
		sql := `CREATE TABLE "part" (
			"id"  INTEGER NOT NULL,
			"name"  TEXT NOT NULL,
			"age"  INTEGER NOT NULL,
      		"desc"  TEXT NOT NULL,
      		"type"  INTEGER NOT NULL,
			"pct"  INTEGER NOT NULL DEFAULT 0,
			"pic" TEXT,`
		_, err := db.Exec(sql)
		if err != nil {
			log.Fatal("database error e")
			return
		}
	} else if err.Error() == "no such table: relation" {
		db := c.DB
		sql := `CREATE TABLE "relation" (
			"id"  INTEGER NOT NULL,
			"p1"  INTEGER NOT NULL,
			"relationid"  INTEGER NOT NULL,
			"p2"  INTEGER NOT NULL,
			"sort" INTEGER NOT NULL`
		_, err := db.Exec(sql)
		if err != nil {
			log.Fatal("database error e")
			return
		}
	} else if err.Error() == "no such table: relationtype" {
		db := c.DB
		sql := `CREATE TABLE "relaitontype" (
			"id"  INTEGER NOT NULL,
			"name"  TEXT NOT NULL,
			"revname"  TEXT NOT NULL,
			"bookid"  INTEGER NOT NULL`
		_, err := db.Exec(sql)
		if err != nil {
			log.Fatal("database error e")
			return
		}
	} else if err != nil {
		log.Println("have error !")
		log.Println(err)
	}
}
