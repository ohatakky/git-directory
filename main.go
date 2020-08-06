package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)

type Commit struct {
	Hash  string
	Files *object.FileIter
}

type Response struct {
	Idx  int
	Tree []string
}

func main() {
	str := strconv.FormatInt(time.Now().UTC().Unix(), 10)
	r, err := git.PlainClone(fmt.Sprintf("./tmp/echo_%s", str), false, &git.CloneOptions{
		URL:        "https://github.com/labstack/echo",
		Progress:   os.Stdout,
		NoCheckout: true,
	})
	if err != nil {
		log.Fatal(err)
	}

	commits := make([]Commit, 0)
	cIter, err := r.Log(&git.LogOptions{})
	if err != nil {
		log.Fatal(err)
	}
	err = cIter.ForEach(func(c *object.Commit) error {
		files, err := c.Files()
		if err != nil {
			log.Fatal(err)
		}
		commits = append(commits, Commit{
			Hash:  c.Hash.String(),
			Files: files,
		})
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	for i := len(commits) - 1; i >= 0; i-- {
		fmt.Println("-----------", commits[i].Hash, "---------------")
		// todo: websocketでResponse返す
		err := commits[i].Files.ForEach(func(f *object.File) error {
			fmt.Println(f.Name)
			return nil
		})
		if err != nil {
			log.Fatal(err)
		}
	}
}
