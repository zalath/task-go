package book

import (
	"strconv"
	dbbook "tasktask/src/sqlitebook"

	"github.com/gin-gonic/gin"
)

// 列举根节点
func GetRoot() []dbbook.Part {
	db := newdb()
	p := db.Partlist()
	return p
}

// 根据根节点列举一级子节点
func GetParts(id string) []dbbook.PartR {
	idval, _ := strconv.Atoi(id)
	db := newdb()
	p := db.Parts(idval)
	return p
}

// 更新part表中的数据
func UpdatePart(c *gin.Context) string {
	idval, _ := strconv.Atoi(c.PostForm("id"))
	ageval, _ := strconv.Atoi(c.PostForm("age"))
	typeval, _ := strconv.Atoi(c.PostForm("type"))

	p := dbbook.Part{
		ID:   idval,
		Name: c.PostForm("name"),
		Age:  ageval,
		Desc: c.PostForm("desc"),
		Type: typeval,
		Pic:  c.PostForm("pic"),
	}
	db := newdb()
	if !db.UpdatePart(p) {
		return "mis"
	}
	return "done"
}

// 删除节点
func DeletePart(id string) string {
	idval, _ := strconv.Atoi(id)
	db := newdb()
	if !db.DeletePart(idval) {
		return "mis"
	}
	return "done"
}

// 增加节点
// p1 父级节点id
// relationid 关系类型
// relationpos 1为正向关系（p1->id)，2为反向关系（id->p1)
func AddPart(name, desc, pic, age, parttype, p1, relationid, relationpos string) dbbook.Part {
	db := newdb()
	ageval, _ := strconv.Atoi(age)
	typeval, _ := strconv.Atoi(parttype)
	p := dbbook.Part{
		Name: name,
		Age:  ageval,
		Desc: desc,
		Pic:  pic,
		Type: typeval,
	}

	preturn := db.AddPart(p, p1, relationid, relationpos)
	return preturn
}

// 修改排序
func UpdateOrder(id, order string) string {
	db := newdb()
	if !db.UpdateOrder(id, order) {
		return "mis"
	}
	return "done"
}

// 获取所有关系表记录
func GetRelationType(bookid string) []dbbook.Relation {
	db := newdb()
	return db.GetRelationType(bookid)
}

// 创建关系类型
func CreateRelationType(name, revname, bookid string) (bool, dbbook.Relation) {
	db := newdb()
	res, id := db.CreateRelationType(name, revname, bookid)
	if res {
		bookidval, _ := strconv.Atoi(bookid)
		relation := dbbook.Relation{ID: id, Name: name, RevName: revname, Bookid: bookidval}
		return true, relation
	} else {
		return false, dbbook.Relation{}
	}
}

func DeleteRelationType(id string) string {
	db := newdb()
	if !db.DeleteRelationType(id) {
		return "mis"
	}
	return "done"
}

// 删除关系表记录
func DeleteRelation(p1, p2 string) string {
	db := newdb()
	if !db.DeleteRelation(p1, p2) {
		return "mis"
	}
	return "done"
}

// 创建关系表记录
func Makerelation(p1, p2, relationid string) string {
	db := newdb()
	if !db.Makerelation(p1, p2, relationid) {
		return "mis"
	}
	return "done"
}

func newdb() *dbbook.Con {
	db := dbbook.NewCon()
	return db
}
