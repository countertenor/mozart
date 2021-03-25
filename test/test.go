package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

func FilePathWalkDir(root string) ([]string, error) {
	var files []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		// if !info.IsDir() {
		// 	files = append(files, path)
		// }
		fmt.Println("file : ", path)
		return nil
	})
	return files, err
}
func OSReadDir(root string) ([]string, error) {
	var files []string
	f, err := os.Open(root)
	if err != nil {
		return files, err
	}
	fileInfo, err := f.Readdir(-1)
	f.Close()
	if err != nil {
		return files, err
	}

	for _, file := range fileInfo {
		files = append(files, file.Name())
	}
	return files, nil
}

func IOReadDir(root string) ([]string, error) {
	var files []string
	fileInfo, err := ioutil.ReadDir(root)
	if err != nil {
		return files, err
	}

	for _, file := range fileInfo {
		files = append(files, file.Name())
	}
	return files, nil
}

func main() {
	FilePathWalkDir("../resources/common-files")
	// f, _ := IOReadDir("../resources/common-files/bnr")
	f, _ := OSReadDir("../resources/common-files/bnr")
	fmt.Println("f : ", f)
}
