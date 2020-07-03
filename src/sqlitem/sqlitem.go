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

//New ...
func (c *Con) New() {
	c.Opendb()
}

//RR ...
type RR struct {
	ID    int    `db:"id"`
	Title string `db:"title"`
	Tik   int    `db:"tik"`
	P     string `db:"p"`
}

// //List ...
// func (c *Con) List(id, etype string) []RR {
// 	c.Opendb()
// 	var data = []RR{}
// 	fmt.Println(c.DB)
// 	fmt.Println("indb")
// 	if etype == "list" {
// 		// c.DB.Select(&data, "select id,title,tik,p from e where pid = '?' order by tik", id)
// 		fmt.Println("list")
// 		c.DB.Select(&data, "select id,title,tik,p from e")
// 		fmt.Println(len(data))
// 		for i := 0; i < len(data); i++ {
// 			fmt.Println(i)
// 			// fmt.Println(data[i].title)
// 		}
// 	} else {
// 		fmt.Println("space")
// 		c.DB.Select(&data, "select id,title,tik,p from e where p like '&?&' order by tik", id)
// 		fmt.Println(len(data))
// 		for i := 0; i < len(data); i++ {
// 			fmt.Println(fmt.Sprintf("s+%b", i))
// 			// fmt.Println(data[i].title)
// 		}
// 	}
// 	return data
// }

//List a test
func (c *Con) List(id, etype string) {
	db, err := sqlx.Connect("sqlite3", "./db.db")
	if err != nil {
		log.Fatal(err)
		fmt.Println(err)
	}
	var data = []RR{}
	err = db.Select(&data, "select id,title,tik,p from e")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%#v\n", data)
	fmt.Println(len(data))
}

// List1 ...
func (c *Con) List1(id, etype string) {
	db, err := sql.Open("sqlite3", "./db.db")
	if err != nil {
		log.Fatal(err)
	}
	stmt, _ := db.Prepare("select id,title,tik,p from e where pid = ?")
	defer stmt.Close()
	res := stmt.QueryRow("66")
	var cc = new(RR)
	res.Scan(&cc.ID, &cc.Title, &cc.Tik, &cc.P)
	fmt.Println("res is:")
	fmt.Println(cc.Title)
}

// c.Opendb()
// res, _ := c.DB.Query("select id,title,tik,p from e")
// for res.Next() {
// 	var title string
// 	var tik int
// 	var p string
// 	res.Scan(&id, &title, &tik, &p)
// 	fmt.Println("res is:")
// 	fmt.Println(id)
// 	fmt.Println(title)
// 	fmt.Println(tik)
// 	fmt.Println(p)
// }
// }

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
