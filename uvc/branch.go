package uvc

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
)

type Branch struct {
	name         string
	headRevision Revision
	revisions    []Revision
}

func getCurrentBranch() Branch {
	branch, err := os.ReadFile(".uvc/current")
	if err != nil {
		log.Fatal(err)
	}
	return getBranch(string(branch))
}

func setCurrentBranch(branchName string) {
	f, err := os.Create(".uvc/current")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	_, err = f.WriteString(branchName)
	if err != nil {
		log.Fatal(err)
	}
}

func getBranchPath(hash string) string {
	return fmt.Sprintf(".uvc/branches/%v", hash)
}

func getBranch(name string) Branch {
	var b Branch
	revisions, err := loadRevisions(getBranchPath(name))
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			b = Branch{
				name: name,
				headRevision: Revision{},
				revisions: []Revision{},
			}
		} else {
			log.Fatal(err)
		}
	} else {
		b = Branch{
			name: name,
			headRevision: revisions[len(revisions)-1],
			revisions: revisions,
		}
	}
	return b
}

func loadRevisions(path string) ([]Revision, error) {
	f, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	fileContent := strings.Split(string(f), "\n")
	var result []Revision
	for _, line := range fileContent {
		if line == "" {
			continue
		}
		result = append(result, getRevision(line))
	}
	return result, nil
}

func (b Branch) store() {
	f, err := os.Create(getBranchPath(b.name))
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	for _, r := range b.revisions {
		f.WriteString(fmt.Sprintf("%v\n", r))
	}
}

func (b *Branch) updateRevisions(r Revision) {
	b.revisions = append(b.revisions, r)
	b.headRevision = r
	b.store()
}