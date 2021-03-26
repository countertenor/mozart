package static

import (
	"embed"
	"io/fs"
	"strings"
)

type staticType string

const (
	ResourceType staticType = "resources"
	WebappBuild  staticType = "webapp/build"
)

//go:embed resources
var Resources embed.FS

//go:embed webapp/build
var Webapp embed.FS

func OpenFileFromStaticFS(staticType staticType, filename string) (fs.File, error) {
	filePath := filename
	if !strings.Contains(filename, string(ResourceType)) {
		filePath = string(staticType) + "/" + filename
	}
	file, err := Resources.Open(filePath)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func Walk(staticType staticType, dir string, fn fs.WalkDirFunc) error {
	return fs.WalkDir(Resources, string(staticType)+"/"+dir, fn)
}
