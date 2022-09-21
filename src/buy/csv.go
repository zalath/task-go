package buy
//个人中心-导出历史账单的解析
import (
	"fmt"
	"encoding/csv"
	"strconv"
	"io"
	dbb "tasktask/src/sqliteb"
	"github.com/axgle/mahonia"
	"github.com/gin-gonic/gin"
)
/*
filetype: wx ali
*/
func Csv(filetype string,c *gin.Context) string {
	uf, err := c.FormFile("file")
	if err != nil {
		fmt.Println(err)
		return "file mis"
	}
	f, er1 := uf.Open()
	if er1 != nil {
		fmt.Println(er1)
		return "open mis"
	}
	reader := csv.NewReader(f)
	reader.FieldsPerRecord = -1
	var dat = []dbb.Buy{}
	if filetype == "ali" {
		dat = deAli(reader)
	}
	if filetype == "wx" {
		dat = deWx(reader)
	}
	db := dbb.NewCon()
	res := db.BatchNew(dat)
	if res {
		return "fin"
	}
	return "mis"
}
//解析支付宝账单
func deAli(reader *csv.Reader) []dbb.Buy{
	var d = []dbb.Buy{}
	for {
        row, err := reader.Read()
        if err != nil && err != io.EOF {
            fmt.Println(err)
        }
        if err == io.EOF {
            break
        }
		if  len(row) < 8 {
			continue
		}
		var dt = dbb.Buy{}
		dt.T = ConvertToString(row[10],"gbk","utf-8")
		dt.Money, _ = strconv.ParseFloat(row[5],64)
		dt.Ttype = ConvertToString(row[7],"gbk","utf-8")
		dt.Ex = ConvertToString(row[6],"gbk","utf-8")
		dt.Merchant = ConvertToString(row[1],"gbk","utf-8")
		dt.Thing = ConvertToString(row[3],"gbk","utf-8")
		dt.Trantype = 0
		dt.Account = ConvertToString(row[4],"gbk","utf-8")
		dt.Order = ConvertToString(row[8],"gbk","utf-8")
		dt.Morder = ConvertToString(row[9],"gbk","utf-8")
		dt.Innout = ConvertToString(row[0],"gbk","utf-8")
		d=append(d,dt)
  }
	return d
}
//解析微信账单
func deWx(reader *csv.Reader) []dbb.Buy{
	var d = []dbb.Buy{}
	st := false
	for {
        row, err := reader.Read()
        if err != nil && err != io.EOF {
            fmt.Println(err)
        }
        if err == io.EOF {
            break
        }
		if !st {
			if row[0] == "交易时间" {
				st = true
			}
			continue
		}
		var dt = dbb.Buy{}
		dt.T = row[0]
		money := string([]rune(row[5])[1:])
		dt.Money, _ = strconv.ParseFloat(money,64)
		dt.Ttype = row[1]
		dt.Ex = row[7]
		dt.Merchant = row[2]
		dt.Thing = row[3]
		dt.Trantype = 0
		dt.Account = row[6]
		dt.Order = row[8]
		dt.Morder = row[9]
		dt.Innout = row[4]
		d=append(d,dt)
  }
	return d
}
func ConvertToString(src, srcCode, targetCode string) string {
	srcCoder := mahonia.NewDecoder(srcCode)
	srcResult := srcCoder.ConvertString(src)
	targetCoder := mahonia.NewDecoder(targetCode)
	_, cdata, _ := targetCoder.Translate([]byte(srcResult), true)
	result := string(cdata)
	return result
}