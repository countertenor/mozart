package statik

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"

	tmplStatik "github.com/countertenor/mozart/statik/tmpl/statik"
	webStatik "github.com/countertenor/mozart/statik/web/statik"
	"github.com/rakyll/statik/fs"
)

//constant namespaces declared
const (
	Template = tmplStatik.Template
	Webapp   = webStatik.Webapp
)

//GetStaticFS gets the static FS according to namespace
func GetStaticFS(namespace string) (http.FileSystem, error) {
	statikFS, err := fs.NewWithNamespace(namespace)
	if err != nil {
		return nil, fmt.Errorf("could not open statikFS err: %v ", err)
	}
	return statikFS, nil
}

//OpenFileFromStaticFS gets the FS used for static files
func OpenFileFromStaticFS(namespace string, filename string) (http.File, error) {
	statikFS, err := GetStaticFS(namespace)
	if err != nil {
		return nil, err
	}
	file, err := statikFS.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("error opening file %v err: %v", filename, err)
	}
	return file, nil
}

//GetActualDirName gets actual dir name from inside a dir
func GetActualDirName(namespace string, directory, dirToLookIn string) (string, error) {
	if directory == "" {
		return "", nil
	}
	var fullDirPath string
	statikFS, err := GetStaticFS(namespace)
	if err != nil {
		return "", err
	}

	var fullReg strings.Builder
	regExp := "[/a-z0-9-]*"
	b := "\\b"
	sp := strings.Split(directory, "/")
	for _, str := range sp {
		//fmt.Println(s)
		fullReg.WriteString(b)
		fullReg.WriteString(str)
		fullReg.WriteString(b)
		fullReg.WriteString(regExp)
	}
	// fmt.Println("fullReg : ", fullReg.String())
	r, err := regexp.Compile(fullReg.String())

	err = fs.Walk(statikFS, dirToLookIn, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			// fmt.Println("info : ", path)
			if r.MatchString(path) {
				fullDirPath = path
				fullDirPath = strings.TrimPrefix(fullDirPath, dirToLookIn)
				// fmt.Println("match : ", fullDirPath)
				return io.EOF //return a known error, to exit out of walk
			}
		}
		return nil
	})
	if err != io.EOF && err != nil {
		return "", fmt.Errorf("error getting static files : %v", err)
	}
	return fullDirPath, nil
}

//GetAllDirsInDir gets all dirs inside a directory
func GetAllDirsInDir(namespace string, dirToLookIn string) ([]string, error) {
	var dirs []string
	if dirToLookIn == "" {
		return dirs, nil
	}

	statikFS, err := GetStaticFS(namespace)
	if err != nil {
		return nil, err
	}

	err = fs.Walk(statikFS, dirToLookIn, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			dirToReturn := path
			dirToReturn = strings.TrimPrefix(dirToReturn, dirToLookIn)
			dirToReturn = strings.TrimSpace(strings.Join(strings.Split(dirToReturn, "/"), " "))
			dirToReturn = regexp.MustCompile("([0-9]+-)").ReplaceAllString(dirToReturn, "")
			if dirToReturn != "" {
				dirs = append(dirs, dirToReturn)
			}
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("error getting files : %v", err)
	}
	// sort.Strings(dirs) //done through the completion script itself
	return dirs, nil
}
