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
	branch := getCurrentBranch()
	revision := newRevision()
	revision.store()
	branch.updateRevisions(revision)
}

func Restore() {
	branch := getCurrentBranch()
	revision := getRevision(branch.headRevision.hash)
	revision.restore()
}
