package git

import (
	"fmt"
	"os"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)

const (
	tmpDirFmt = "./tmp/%s_%s"
)

type Client struct {
	ID   string
	Repo string
	Ref  *git.Repository
	Send chan *Tree
}

func New(id, repo string) *Client {
	return &Client{
		ID:   id,
		Repo: repo,
		Send: make(chan *Tree, 256),
	}
}

func (c *Client) TmpDir() string {
	return fmt.Sprintf(tmpDirFmt, strings.Replace(c.Repo, "/", "_", -1), c.ID)
}

type Commit struct {
	Hash  string
	Files *object.FileIter
}

func (c *Client) Clone() error {
	r, err := git.PlainClone(c.TmpDir(), false, &git.CloneOptions{
		URL:        fmt.Sprintf("https://github.com/%s", c.Repo),
		Progress:   os.Stdout,
		NoCheckout: true,
	})
	if err != nil {
		return err
	}

	c.Ref = r
	return nil
}

func (c *Client) commits() ([]Commit, error) {
	commits := make([]Commit, 0)
	cIter, err := c.Ref.Log(&git.LogOptions{})
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

// todo: error handling (<-error)
func (c *Client) Trees() error {
	defer close(c.Send)
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
