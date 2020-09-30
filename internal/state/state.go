package state

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io"
	"os"
	"sync"
)

var lock sync.Mutex
var once sync.Once

//singleton reference
var fileTypeInstance *fileType

//fileType is the main struct for file database
type fileType struct {
	path         string
	isRegistered bool
}

//getFileTypeInstance returns the singleton instance of the filedb object
func getFileTypeInstance() (*fileType, error) {

	if fileTypeInstance == nil {
		return nil, fmt.Errorf("state DB nil")
	}
	return fileTypeInstance, nil
}

//InitiateFileTypeInstance initiates the instance and sets the path for the DB
func InitiateFileTypeInstance(filePath string) {
	if fileTypeInstance == nil {
		once.Do(func() {
			fileTypeInstance = &fileType{path: filePath}
		})
	}
}

// Save saves a representation of v to the file at path.
func Save(v interface{}) error {
	lock.Lock()
	defer lock.Unlock()
	fileType, err := getFileTypeInstance()
	if err != nil {
		return err
	}

	if !fileType.isRegistered {
		gob.Register(v)
		fileType.isRegistered = true
	}

	f, err := os.Create(fileType.path)
	if err != nil {
		return err
	}
	defer f.Close()
	r, err := marshal(v)
	if err != nil {
		return err
	}
	_, err = io.Copy(f, r)
	return err
}

// Load loads the file at path into v.
func Load(v interface{}) error {
	fileType, err := getFileTypeInstance()
	if err != nil {
		return err
	}
	if fileExists(fileType.path) {
		lock.Lock()
		defer lock.Unlock()
		f, err := os.Open(fileType.path)
		if err != nil {
			return err
		}
		defer f.Close()
		return unmarshal(f, v)
	}
	fmt.Printf("Db file %v not found", fileType.path)
	return nil
}

// marshal is a function that marshals the object into an
// io.Reader.
var marshal = func(v interface{}) (io.Reader, error) {
	var buf bytes.Buffer
	e := gob.NewEncoder(&buf)
	err := e.Encode(v)
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(buf.Bytes()), nil
}

// unmarshal is a function that unmarshals the data from the
// reader into the specified value.
var unmarshal = func(r io.Reader, v interface{}) error {
	d := gob.NewDecoder(r)
	err := d.Decode(v)
	if err != nil {
		return err
	}
	return nil
}

// fileExists checks if a file exists and is not a directory before we
// try using it to prevent further errors.
func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
