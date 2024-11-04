package uvc

import (
	"fmt"
	"log"
	"os"
)

const (
	RepoDir       string = ".uvc"
	DefaultBranch string = "main"
)

func RepoInit() {
	if repositoryInitialized() {
		fmt.Fprintf(os.Stderr, "Repository already initialized.")
	}
	if err := os.Mkdir(RepoDir, os.FileMode(0700)); err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
	}
	repoPaths := []string{"branches", "objects", "revisions"}
	for _, path := range repoPaths {
		os.Mkdir(fmt.Sprintf("%v/%v", RepoDir, path), os.FileMode(0700))
	}
	f, err := os.Create(fmt.Sprintf("%v/%v", RepoDir, "current"))
	if err != nil {
		log.Fatal(err)
	}
	f.Write([]byte(DefaultBranch))

}

func Status() {
	currentDir := "."
	files := make([]string, 0)
	getFiles(currentDir, &files)
	for _, file := range files {
		fmt.Printf("%v: %v\n", fileSHA1(file), file)
	}
}

func Commit() {
	revision := newRevision()
	revision.store()
}

func Restore() {
	revision := getRevision("81146d1604796f0432b2b59ab9ed6ff3e5bf75f1")
	revision.restore()
}
