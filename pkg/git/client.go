package git

import (
	"fmt"
	"os"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)

type Client struct {
	ID   string
	Repo string
	Send chan *Tree
}

func New(id, repo string) *Client {
	return &Client{
		ID:   id,
		Repo: repo,
		Send: make(chan *Tree, 256),
	}
}

type Commit struct {
	Hash  string
	Files *object.FileIter
}

func (c *Client) commits() ([]Commit, error) {
	r, err := git.PlainClone(fmt.Sprintf("./tmp/%s_%s", c.Repo, c.ID), false, &git.CloneOptions{
		URL:        fmt.Sprintf("https://github.com/%s", c.Repo),
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

type Tree struct {
	I    int      `json:"i"`
	Hash string   `json:"hash"`
	T    []string `json:"t"`
}

func (c *Client) Trees() error {
	commits, err := c.commits()
	if err != nil {
		return err
	}

	for i := len(commits) - 1; i >= 0; i-- {
		tree := Tree{
			I:    i,
			Hash: commits[i].Hash,
			T:    make([]string, 0),
		}
		err := commits[i].Files.ForEach(func(f *object.File) error {
			tree.T = append(tree.T, f.Name)
			return nil
		})
		if err != nil {
			return err
		}
		c.Send <- &tree
	}

	return nil
}
