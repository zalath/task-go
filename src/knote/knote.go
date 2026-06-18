package knote

import (
	// "fmt"
	"os"
	"path/filepath"
)

var Istest = false

var path = "/db/knote"

var confpath = "/conf/"

func Read() string {
	return readfile(path)
}

func readfile(filename string) string {
	path := getpath(filename)
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	content, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	return string(content)
}

func getpath(filename string) string {
	pathp, _ := os.Executable()
	path := filepath.Dir(pathp)
	if Istest {
		path = "."
	}
	path += filename
	file, err := os.Open(path)
	if err != nil {
		os.Create(path)
	}
	file.Close()
	return path
}
func Set(content string) bool {
	return setfile(content, path)
}
func setfile(content, filename string) bool {
	path := getpath(filename)
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	err = os.WriteFile(path, []byte(content), 0755)
	if err != nil {
		panic(err)
	}
	return true
}
func Readkeyfile() string {
	return readfile("/keyfile")
}
func Setkeyfile(content string) bool {
	return setfile(content, "/keyfile")
}

// 服务器配置文件部分-------------------------------------
//
// 获取配置文件的列表
func Conffilelist() []string {
	// 正则表达式匹配指定目录中所有以'conf-'开头的文件
	// fmt.Println(confpath + "conf-*")
	pathp, _ := os.Executable()
	path := filepath.Dir(pathp)
	if Istest {
		path = "."
	}
	path += confpath
	pattern := filepath.Join(path, "conf-*")
	files, err := filepath.Glob(pattern)

	// fmt.Println(files)
	if err != nil {
		panic(err)
	}
	
	re := make([]string, 0)
	for _, file := range files {
		filename := filepath.Base(file)
		if len(filename) > 5 {
			re = append(re, filename[5:])
		}
	}

	return re
}

// 读取指定的配置文件的内容
func Getconffile(filename string) string {
	return readfile(confpath + "conf-" + filename)
}

// 将配置内容写入配置文件
func Setconffile(content, filename string) bool {
	return setfile(content, confpath+"conf-"+filename)
}
