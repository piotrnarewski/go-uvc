package main

import (
	"os"

	"github.com/piotrnarewski/go-uvc/uvc"
)

func main() {
	switch os.Args[1] {
	case "repo-init":
		uvc.RepoInit()
	case "status":
		uvc.Status()
	case "commit":
		uvc.Commit()
	case "restore":
		uvc.Restore()
	case "checkout":
		uvc.Checkout()
	}
}
