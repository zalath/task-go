package el

import (
	"fmt"
	sqlim "tasktask/src/sqlitem"
)

//List return direct child of selected node
func List(id, etype string) {
	fmt.Println("inlist")
	db := newdb()
	db.List(id, etype)
}

//List1 only a test
func List1(id, etype string) {
	fmt.Println("inlist")
	db := newdb()
	db.List1(id, etype)
}

//Space return all element with structure
func Space(id, dtype string) {
	// db := newdb()
}

//New create a new element into database
func New() {
	// db := newdb()
}

//El display edit page
func El(id string) {
}

//Save submit saving element
func Save() {

}

//Tik change an element's state
func Tik() {

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

func newdb() *sqlim.Con {
	db := new(sqlim.Con)
	return db
}
