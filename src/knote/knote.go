package knote

import (
	"os"
	"path/filepath"
)

var Istest = false

func Read() string {
	return readfile("/knote")
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
	if Istest == true {
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
	return setfile(content, "/knote")
}
func setfile(content, filename string) bool {
	path := getpath(filename)
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	err = os.WriteFile(path, []byte(content), 755)
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
