package buy

import (
	"fmt"
	"encoding/csv"
	"strconv"
	"os"
	"io"
	dbb "tasktask/src/sqliteb"
	"github.com/axgle/mahonia"
)

func Csv(filetype string) string {
	// f, err := os.Open("./alipay_record_20220414_094715.csv")
	f, err := os.Open("./微信支付账单(20220112-20220412).csv")
	if err != nil {
		fmt.Println(err)
		return ""
	}
	reader := csv.NewReader(f)
	reader.FieldsPerRecord = -1
	if filetype == "ali" {
		de_ali(reader)
	}
	if filetype == "wx" {
		de_wx(reader)
	}
	return ""
}
func de_ali(reader *csv.Reader){
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
		dt.Type = ConvertToString(row[7],"gbk","utf-8")
		dt.Ex = ConvertToString(row[6],"gbk","utf-8")
		dt.Merchant = ConvertToString(row[1],"gbk","utf-8")
		dt.Thing = ConvertToString(row[3],"gbk","utf-8")
		dt.Trantype = 0
		dt.Account = ConvertToString(row[4],"gbk","utf-8")
		dt.Order = ConvertToString(row[8],"gbk","utf-8")
		dt.Morder = ConvertToString(row[9],"gbk","utf-8")
		d=append(d,dt)
    }
	fmt.Println(d)
}

func de_wx(reader *csv.Reader){
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
		dt.T = row[10]
		dt.Money, _ = strconv.ParseFloat(row[5],64)
		dt.Type = row[7]
		dt.Ex = row[6]
		dt.Merchant = row[1]
		dt.Thing = row[3]
		dt.Trantype = 0
		dt.Account = row[4]
		dt.Order = row[8]
		dt.Morder = row[9]


		// dt.T = ConvertToString(row[10],"gbk","utf-8")
		// dt.Money, _ = strconv.ParseFloat(row[5],64)
		// dt.Type = ConvertToString(row[7],"gbk","utf-8")
		// dt.Ex = ConvertToString(row[6],"gbk","utf-8")
		// dt.Merchant = ConvertToString(row[1],"gbk","utf-8")
		// dt.Thing = ConvertToString(row[3],"gbk","utf-8")
		// dt.Trantype = 0
		// dt.Account = ConvertToString(row[4],"gbk","utf-8")
		// dt.Order = ConvertToString(row[8],"gbk","utf-8")
		// dt.Morder = ConvertToString(row[9],"gbk","utf-8")
		d=append(d,dt)
    }
	fmt.Println(d)
}
func ConvertToString(src, srcCode, targetCode string) string {
	srcCoder := mahonia.NewDecoder(srcCode)
	srcResult := srcCoder.ConvertString(src)
	targetCoder := mahonia.NewDecoder(targetCode)
	_, cdata, _ := targetCoder.Translate([]byte(srcResult), true)
	result := string(cdata)
	return result
}