package uvc

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
)

type Object struct {
	hash string
	path string
}

func getObjectHashPath(hash string) string {
	return fmt.Sprintf(".uvc/objects/%v", hash)
}

func newObject(path string) Object {
	hash := fileSHA1(path)
	obj := Object{
		hash: hash,
		path: path,
	}
	obj.store()
	return obj
}

func getObject(hash, path string) Object {
	if _, err := os.Stat(getObjectHashPath(hash)); errors.Is(err, os.ErrNotExist) {
		log.Fatal(err)
	}
	return Object{
		hash: hash,
		path: path,
	}
}

func (obj Object) store() {
	if _, err := os.Stat(getObjectHashPath(obj.hash)); errors.Is(err, os.ErrNotExist) {
		dst, err := os.Create(getObjectHashPath(obj.hash))
		if err != nil {
			log.Fatal(err)
		}
		defer dst.Close()
		src, err := os.Open(obj.path)
		if err != nil {
			log.Fatal(err)
		}
		defer src.Close()
		_, err = io.Copy(dst, src)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func (obj Object) restore() {
	if _, err := os.Stat(getObjectHashPath(obj.hash)); errors.Is(err, os.ErrNotExist) {
		log.Fatal(err)
	}
	src, err := os.Open(getObjectHashPath(obj.hash))
	if err != nil {
		log.Fatal(err)
	}
	dst, err := os.Create(obj.path)
	if err != nil {
		log.Fatal(err)
	}
	defer dst.Close()
	io.Copy(dst, src)
}

func (obj Object) String() string {
	return fmt.Sprintf("%v:%v", obj.path, obj.hash)
}
