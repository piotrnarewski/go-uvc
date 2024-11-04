package uvc

import (
	"crypto/sha1"
	"fmt"
	"io"
	"log"
	"os"
)

func fileSHA1(filePath string) string {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	h := sha1.New()
	if _, err := io.Copy(h, file); err != nil {
		log.Fatal(err)
	}
	return fmt.Sprintf("%x", h.Sum(nil))
}

func bytesSHA1(bytes []byte) string {
	h := sha1.New()
	h.Write(bytes)
	return fmt.Sprintf("%x", h.Sum(nil))
}

func getFiles(path string, result *[]string) {
	dirEntry, err := os.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}
	for _, entry := range dirEntry {
		if entry.Name()[0] == '.' {
			continue
		}
		if entry.IsDir() {
			getFiles(fmt.Sprintf("%v/%v", path, entry.Name()), result)
		} else {
			*result = append(*result, fmt.Sprintf("%v/%v", path, entry.Name()))
		}
	}
}

func repositoryInitialized() (result bool) {
	_, err := os.Stat(RepoDir)
	if err == nil {
		result = true
	}
	return result
}
