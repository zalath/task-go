package sqlitem

import (
	"database/sql"
	"fmt"
	"log"
	"os"

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
	ID    int    `db:"id" json:"id"`
	Title string `db:"title" json:"title"`
	Tik   int    `db:"tik" json:"tik"`
	P     string `db:"p" json:"p"`
	Pid   int    `db:"pid" json:"pid"`
	Child interface{}
}

//List a test
func (c *Con) List(id, etype string) []El {
	db := c.DB
	var err error
	var data = []El{}
	if etype == "list" {
		err = db.Select(&data, "select id,title,tik,p,pid from e where pid = ?", id)
	} else {
		err = db.Select(&data, "select id,title,tik,p,pid from e where p like '%'||$1||'%'", id)
	}

	if err != nil {
		fmt.Println(err)
	}
	return data
}

//Get select online from database on id
func (c *Con) Get(id string) El {
	db := c.DB
	el := El{}
	err := db.Get(&el, "select id,title,tik,p,pid from e where id = ?", id)
	if err != nil {
		fmt.Println(err)
	}
	return el
}

//New create a new element
func (c *Con) New(el El) (isdone bool, newid int64) {
	isdone = true
	db := c.DB
	stmt, err := db.Prepare("insert into e (title,pid,p,tik) values(?,?,?,?)")
	if err != nil {
		fmt.Println(err)
	}
	res, er1 := stmt.Exec(el.Title, el.Pid, el.P, el.Tik)
	if er1 != nil {
		fmt.Println(err)
		isdone = false
		return
	}
	newid, _ = res.LastInsertId()
	return
}

func lite() {
	os.Remove("./foo.db")

	db, err := sql.Open("sqlite3", "./foo.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	sqlStmt := `
	create table foo (id integer not null primary key, name text);
	delete from foo;
	`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return
	}

	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	stmt, err := tx.Prepare("insert into foo(id, name) values(?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	for i := 0; i < 100; i++ {
		_, err = stmt.Exec(i, fmt.Sprintf("こんにちわ世界%03d", i))
		if err != nil {
			log.Fatal(err)
		}
	}
	tx.Commit()

	rows, err := db.Query("select id, name from foo")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		var name string
		err = rows.Scan(&id, &name)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(id, name)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	stmt, err = db.Prepare("select name from foo where id = ?")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	var name string
	err = stmt.QueryRow("3").Scan(&name)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(name)

	// _, err = db.Exec("delete from foo")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	_, err = db.Exec("insert into foo(id, name) values(1, 'foo'), (2, 'bar'), (3, 'baz')")
	if err != nil {
		log.Fatal(err)
	}

	rows, err = db.Query("select id, name from foo")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		var name string
		err = rows.Scan(&id, &name)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(id, name)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
}
