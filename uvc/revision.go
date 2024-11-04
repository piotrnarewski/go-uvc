package uvc

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"os"
	"slices"
	"strings"
)

type Revision struct {
	hash    string
	message string
	objects []Object
}

func newRevision() Revision {
	var objects []Object
	currentDir := "."
	files := make([]string, 0)
	getFiles(currentDir, &files)
	for _, f := range files {
		objects = append(objects, newObject(f))
	}
	r := Revision{
		objects: objects,
	}
	r.hash = bytesSHA1(r.Bytes())
	return r
}

func getHashPath(hash string) string {
	return fmt.Sprintf(".uvc/revisions/%v", hash)
}

func getRevision(hash string) Revision {
	if _, err := os.Stat(getHashPath(hash)); errors.Is(err, os.ErrExist) {
		log.Fatal(err)
	}
	return Revision{
		hash:    hash,
		objects: loadObjects(getHashPath(hash)),
	}
}

func loadObjects(path string) []Object {
	f, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	fileContent := strings.Split(string(f), "\n")
	var result []Object
	for _, line := range fileContent {
		if line == "" {
			continue
		}
		path := strings.Split(line, ":")[0]
		hash := strings.Split(line, ":")[1]
		result = append(result, Object{hash: hash, path: path})
	}
	return result
}

func (r Revision) store() {
	if _, err := os.Stat(getHashPath(r.hash)); errors.Is(err, os.ErrExist) {
		log.Fatal(err)
	} else {
		dst, err := os.Create(getHashPath(r.hash))
		if err != nil {
			log.Fatal(err)
		}
		defer dst.Close()
		dst.Write(r.Bytes())
	}
}

func (r Revision) restore() {
	currentDir := "."
	files := make([]string, 0)
	getFiles(currentDir, &files)
	for _, obj := range r.objects {
		files = slices.DeleteFunc(files, func(s string) bool {
			return obj.path == s
		})
		obj.restore()
	}
	for _, file := range files {
		err := os.Remove(file)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func (r Revision) Bytes() []byte {
	var buffer bytes.Buffer
	for _, obj := range r.objects {
		buffer.WriteString(fmt.Sprintln(obj))
	}
	return buffer.Bytes()
}
