package book

import (
	"strconv"
	dbbook "tasktask/src/sqlitebook"

	"github.com/gin-gonic/gin"
)

// 列举根节点
func List() []dbbook.Part {
	db := newdb()
	p := db.List()
	return p
}

// 根据根节点列举一级子节点
func Parts(id string) []dbbook.Part {
	db := newdb()
	p := db.Parts(id)
	return p
}

func PartRelationList(p1, relationid, p2 string) []dbbook.Relation {
	db := newdb()
	p := db.PartRelationList(p1, relationid, p2)
	return p
}
func PartRelationMap(bookid string) []dbbook.Relation {
	db := newdb()
	p := db.PartRelationMap(bookid)
	return p
}

// 修改关系排序
func PartRelationSetOrder(id, order string) string {
	db := newdb()
	if !db.PartRelationSetOrder(id, order) {
		return "mis"
	}
	return "done"
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
func AddPart(name, desc, pic, age, parttype, pct, bookid string) dbbook.PartPs {
	db := newdb()
	ageval, _ := strconv.Atoi(age)
	typeval, _ := strconv.Atoi(parttype)
	pctval, _ := strconv.Atoi(pct)
	p := dbbook.PartPs{
		Name:       name,
		Age:        ageval,
		Desc:       desc,
		Pic:        pic,
		Type:       typeval,
		Pct:        pctval,
		Relationid: "0",
	}
	preturn := db.AddPart(p, bookid)
	return preturn
}

// 获取所有关系表记录
func GetRelationType(bookid string) []dbbook.RelationType {
	db := newdb()
	return db.GetRelationType(bookid)
}

// 创建关系类型
func CreateRelationType(name, revname, bookid string) (bool, dbbook.RelationType) {
	db := newdb()
	res, id := db.CreateRelationType(name, revname, bookid)
	if res {
		bookidval, _ := strconv.Atoi(bookid)
		relationType := dbbook.RelationType{ID: id, Name: name, RevName: revname, Bookid: bookidval}
		return true, relationType
	} else {
		return false, dbbook.RelationType{}
	}
}
func UpdateRelationType(id, name, revname string) string {
	db := newdb()
	if !db.UpdateRelationType(id, name, revname) {
		return "mis"
	}
	return "done"
}
func DeleteRelationType(id string) string {
	db := newdb()
	if !db.DeleteRelationType(id) {
		return "mis"
	}
	return "done"
}

// 删除关系表记录
func DeleteRelation(id string) string {
	db := newdb()
	if !db.DeleteRelation(id) {
		return "mis"
	}
	return "done"
}

// 创建关系表记录
func Makerelation(p1, p2, relationid, direction string) string {
	db := newdb()
	if !db.Makerelation(p1, p2, relationid, direction) {
		return "mis"
	}
	return "done"
}

func newdb() *dbbook.Con {
	db := dbbook.NewCon()
	return db
}
