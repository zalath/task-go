package sqlitebook

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"
	"strconv"

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
	db, err := sqlx.Connect("sqlite3", path+"/book.db")
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
	ID    int    `db:"id" json:"id"`
	Name  string `db:"name" json:"name"`
	Age   int    `db:"age" json:"age"`
	Desc  string `db:"desc" json:"desc"`
	Order int    `db:"order" json:"order"`
	Type  int    `db:"type" json:"type"`
	Pct   int    `db:"pct" json:"pct"`
	Pic   string `db:"pic" json:"pic"`
}

// 继承Part，并添加额外的属性
type PartR struct {
	Part
	Relationname string `db:"relationname" json:"relationname"`
}

// 获取书籍节点
func (c *Con) Partlist() []Part {
	db := c.DB
	var err error
	var data = []Part{}
	err = db.Select(&data, "select * from part where type = 0 order by `order` asc")
	if err != nil {
		c.haveErr(err)
	}
	return data
}

// 通过父级节点，获取子节点列表
func (c *Con) Parts(id int) []PartR {
	db := c.DB
	p1 := []PartR{}
	err := db.Get(&p1, "select p.*,rt.name as relationname from relation r join part p on r.p2=p.id join relationtype rt on r.relationid = rt.id where r.p1 = ?", id)
	if err != nil {
		c.haveErr(err)
	}
	p2 := []PartR{}
	err = db.Get(&p2, "select p.*,rt.revname as relationname from relation r join part p on r.p1=p.id join relationtype rt on r.relationid = rt.id where r.p2 = ?", id)
	if err != nil {
		c.haveErr(err)
	}
	// 合并两个切片，并返回一个新的切片
	p := append(p1, p2...)
	return p
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

func (c *Con) AddPart(p Part, p1, relationid, relationpos string) Part {
	// 插入节点，并根据新id和P1插入relation表
	db := c.DB
	tx, err := db.Begin()
	if err != nil {
		c.haveErr(err)
	}
	res, err := tx.Exec("insert into part (`name`,`age`,`desc`,`pic`,`type`,`pct`) values (?,?,?,?,?,0)", p.Name, p.Age, p.Desc, p.Pic, p.Type)
	if err != nil {
		tx.Rollback()
		c.haveErr(err)
	}
	id, _ := res.LastInsertId()

	if relationpos == "1" {
		err = c.Makerelationsql(tx, p1, strconv.FormatInt(id, 10), relationid)
	} else {
		err = c.Makerelationsql(tx, strconv.FormatInt(id, 10), p1, relationid)
	}
	if err != nil {
		tx.Rollback()
		c.haveErr(err)
	}
	err = tx.Commit()
	if err != nil {
		c.haveErr(err)
	}
	p.ID = int(id)
	return p
}

// 更新书籍节点顺序
func (c *Con) UpdateOrder(id, order string) bool {
	// 更新表中的order，使对应的id的order变为新的值, 同时将order等于或者大于该值的order值加1
	db := c.DB
	tx, err := db.Begin()
	if err != nil {
		c.haveErr(err)
		return false
	}
	_, err = tx.Exec("update part set `order`=`order`+1 where `order`>=?", order)
	if err != nil {
		tx.Rollback()
		c.haveErr(err)
		return false
	}
	_, err = tx.Exec("update part set `order`=? where id=?", order, id)
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

type Relation struct {
	ID      int    `db:"id",json:"id"`
	Name    string `db:"name",json:"name"`
	RevName string `db:"revname",json:"revname"`
	Bookid  int    `db:"bookid",json:"bookid"`
}

// 查询关系类型表
func (c *Con) GetRelationType(bookid string) []Relation {
	db := c.DB
	var rt []Relation
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

// 删除关系表记录
func (c *Con) DeleteRelation(p1, p2 string) bool {
	db := c.DB
	_, err := db.Exec("delete from relation where (p1=? and p2=?) or (p1=? and p2=?)", p1, p2, p2, p1)
	if err != nil {
		c.haveErr(err)
		return false
	}
	return true
}

// 创建关系表记录
func (c *Con) Makerelation(p1, p2, relationid string) bool {
	db := c.DB
	tx, err := db.Begin()
	if err != nil {
		c.haveErr(err)
		return false
	}
	err = c.Makerelationsql(tx, p1, p2, relationid)
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
func (c *Con) Makerelationsql(tx *sql.Tx, p1, p2, relationid string) error {
	_, err := tx.Exec("insert into relation (p1,p2,relationid) values (?,?,?)", p1, p2, relationid)
	return err
}

func (c *Con) haveErr(err error) {
	if err.Error() == "no such table: part" {
		db := c.DB
		sql := `CREATE TABLE "part" (
			"id"  INTEGER NOT NULL,
			"name"  TEXT NOT NULL,
			"age"  INTEGER NOT NULL,
      "desc"  TEXT NOT NULL,
			"order"  INTEGER NOT NULL,
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
		sql := `CREATE TABLE "part" (
			"id"  INTEGER NOT NULL,
			"p2"  INTEGER NOT NULL,
      "relationid"  INTEGER NOT NULL,
      "p2"  INTEGER NOT NULL`
		_, err := db.Exec(sql)
		if err != nil {
			log.Fatal("database error e")
			return
		}
	} else if err.Error() == "no such table: relationtype" {
		db := c.DB
		sql := `CREATE TABLE "part" (
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
		log.Println(err)
	}
}
