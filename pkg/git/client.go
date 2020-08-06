package git

import (
	"fmt"
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

// todo: CommitHashも返す
type Tree struct {
	Idx int      `json:"idx"`
	T   []string `json:"t"`
}

type Client struct{}

func New() *Client {
	return &Client{}
}

func (c *Client) Commits() ([]Commit, error) {
	str := strconv.FormatInt(time.Now().UTC().Unix(), 10)
	r, err := git.PlainClone(fmt.Sprintf("./tmp/echo_%s", str), false, &git.CloneOptions{
		URL:        "https://github.com/labstack/echo",
		Progress:   os.Stdout,
		NoCheckout: true,
	})
	if err != nil {
		return nil, err
	}

	commits := make([]Commit, 0)
	cIter, err := r.Log(&git.LogOptions{})
	if err != nil {
		return nil, err
	}
	err = cIter.ForEach(func(c *object.Commit) error {
		files, err := c.Files()
		if err != nil {
			return err
		}
		commits = append(commits, Commit{
			Hash:  c.Hash.String(),
			Files: files,
		})
		return nil
	})
	if err != nil {
		return nil, err
	}

	return commits, nil
}
