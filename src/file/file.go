package file

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

var Istest = false
var storagepath = "/../storage/dcim/taskres/"

func Del(c *gin.Context) string {
	pathp, _ := os.Executable()
	path := filepath.Dir(pathp) + storagepath
	if Istest == true {
		path = "."
	}
	oldfilename := c.PostForm("del")
	if oldfilename != "" {
		oldfilepath := path + oldfilename
		err := os.Remove(oldfilepath)
		if err != nil {
			fmt.Println(err)
			return "error"
		}
	}
	return "done"
}
func Upload(c *gin.Context) string {
	pathp, _ := os.Executable()
	path := filepath.Dir(pathp) + storagepath
	url := storagepath
	if Istest == true {
		path = "."
	}
	file, errLoad := c.FormFile("pic")
	if errLoad != nil {
		fmt.Println(errLoad)
		return "error"
	}
	filepath := path + file.Filename
	url = url + file.Filename
	err := c.SaveUploadedFile(file, filepath)
	if err != nil {
		fmt.Println(err)
		return "error"
	}
	return url
}
