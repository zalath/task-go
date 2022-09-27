
package file

import (
	"fmt"
	"os"
	"time"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

var Istest = false
var storagepath = "/../storage/dcim/taskres/"

func Del(c *gin.Context) string {
	pathp, _ := os.Executable()
	path := filepath.Dir(pathp) + storagepath
	if Istest == true {
		path = "./pic/"
	}
	oldfilename := c.PostForm("del")
	if oldfilename != "" {
		oldfilename = oldfilename[2:]
		oldfilepath := path + oldfilename
		err := os.Remove(oldfilepath)
		if err != nil {
			return "error"
		}
	}
	return "done"
}
func Upload(c *gin.Context) string {
	pathp, _ := os.Executable()
	path := filepath.Dir(pathp) + storagepath
	url := "f/"
	if Istest == true {
		path = "./pic/"
	}
	file, errLoad := c.FormFile("file")
	if errLoad != nil {
		fmt.Println(errLoad)
		return "mis"
	}
	timestamp := time.Now().Format("2006-01-02-03-04-05")
	filepath := path + timestamp + file.Filename
	url = url + timestamp + file.Filename
	err := c.SaveUploadedFile(file, filepath)
	if err != nil {
		fmt.Println(err)
		return "mis"
	}
	return url
}
