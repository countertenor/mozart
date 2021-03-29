package static

import (
	"embed"
	"fmt"
	"io"
	"io/fs"
	"path/filepath"
	"regexp"
	"strings"
)

type staticType string

//embedding from web_static using conditional build tags
var WebappBuildType staticType
var Webapp embed.FS

//go:embed resources
var Resources embed.FS
var ResourceType staticType = "resources"

func GetEmbedFS(staticType staticType) embed.FS {
	switch staticType {
	case ResourceType:
		return Resources
	case WebappBuildType:
		return Webapp
	//future use cases
	default:
		return Resources
	}
}

func OpenFileFromStaticFS(staticType staticType, filename string) (fs.File, error) {
	filePath := filename
	if !strings.Contains(filename, string(ResourceType)) {
		filePath = filepath.Join(string(staticType), filename)
	}
	file, err := Resources.Open(filePath)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func Walk(staticType staticType, dir string, fn fs.WalkDirFunc) error {
	return fs.WalkDir(GetEmbedFS(staticType), filepath.Join(string(staticType), dir), fn)
}

//GetActualDirName gets actual dir name from inside a dir
func GetActualDirName(staticType staticType, dirToGenerateFrom, dirToLookIn string) (string, error) {
	if dirToGenerateFrom == "" {
		return "", nil
	}
	var fullDirPath string
	var fullReg strings.Builder
	regExp := "[/a-z0-9-]*"
	b := "\\b"
	sp := strings.Split(dirToGenerateFrom, "/")
	for _, str := range sp {
		//fmt.Println(s)
		fullReg.WriteString(b)
		fullReg.WriteString(str)
		fullReg.WriteString(b)
		fullReg.WriteString(regExp)
	}
	// fmt.Println("fullReg : ", fullReg.String())
	r, err := regexp.Compile(fullReg.String())
	if err != nil {
		return "", fmt.Errorf("error compiling regex %v: %v", fullReg.String(), err)
	}
	// fmt.Println("dirToGenerateFrom : ", dirToGenerateFrom)
	// fmt.Println("dirToLookIn : ", dirToLookIn)

	err = Walk(staticType, dirToLookIn, func(path string, info fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			// fmt.Println("info : ", path)
			if r.MatchString(path) {
				fullDirPath, err = GetRelativePath(path, dirToLookIn)
				if err != nil {
					return err
				}
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
func GetAllDirsInDir(staticType staticType, dirToLookIn string) ([]string, error) {
	var dirs []string
	if dirToLookIn == "" {
		return dirs, nil
	}

	err := Walk(staticType, dirToLookIn, func(path string, info fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			dirToReturn, err := GetRelativePath(path, dirToLookIn)
			if err != nil {
				return err
			}
			dirToReturn = strings.TrimSpace(strings.Join(strings.Split(dirToReturn, "/"), " "))
			dirToReturn = regexp.MustCompile("([0-9]+-)").ReplaceAllString(dirToReturn, "")
			if dirToReturn != "" && dirToReturn != "." {
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

func GetRelativePath(path, dir string) (string, error) {
	base := filepath.Join(string(ResourceType), dir)
	relativePath, err := filepath.Rel(base, path)
	if err != nil {
		return "", fmt.Errorf("unable to remove prefix from %v, err :%v", path, err)
	}
	return relativePath, nil
}
