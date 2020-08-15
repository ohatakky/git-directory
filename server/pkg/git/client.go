package git

import (
	"fmt"
	"io"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

const (
	tmpDirFmt = "./tmp/%s_%s_%s"
)

type Client struct {
	ID    string
	Org   string
	Repo  string
	Ref   *git.Repository
	Send  chan *Commit
	Error chan error
}

func New(id, org, repo string) *Client {
	return &Client{
		ID:   id,
		Org:  org,
		Repo: repo,
		Send: make(chan *Commit, 256),
	}
}

func (c *Client) TmpDir() string {
	return fmt.Sprintf(tmpDirFmt, c.Org, c.Repo, c.ID)
}

func (c *Client) Clone() error {
	r, err := git.PlainClone(c.TmpDir(), false, &git.CloneOptions{
		URL:          fmt.Sprintf("https://github.com/%s/%s", c.Org, c.Repo),
		SingleBranch: true,
		NoCheckout:   true,
		Tags:         git.NoTags,
	})
	if err != nil {
		return err
	}
	c.Ref = r
	return nil
}

type Object struct {
	IsFile bool   `json:"is_file"`
	Name   string `json:"name"`
}

type Commit struct {
	Hash    string    `json:"hash"`
	Message string    `json:"message"`
	Author  string    `json:"author"`
	Objects []*Object `json:"objects"`
}

func (c *Client) Commits() {
	cIter, err := c.Ref.Log(&git.LogOptions{
		Order: git.LogOrderCommitterTime,
	})
	if err != nil {
		c.Error <- err
	}
	defer cIter.Close()

	rev := make([]*object.Commit, 0)
	err = cIter.ForEach(func(co *object.Commit) error {
		rev = append(rev, co)
		return nil
	})
	if err != nil {
		c.Error <- err
	}

	for i := len(rev) - 1; i >= 0; i-- {
		tree, err := rev[i].Tree()
		if err != nil {
			c.Error <- err
		}
		send := &Commit{
			Hash:    rev[i].Hash.String(),
			Message: rev[i].Message,
			Author:  rev[i].Author.String(),
			Objects: make([]*Object, 0),
		}
		seen := make(map[plumbing.Hash]bool)
		walker := object.NewTreeWalker(tree, true, seen)
		defer walker.Close()
		for {
			name, entry, err := walker.Next()
			if err == io.EOF {
				break
			}
			send.Objects = append(send.Objects, &Object{
				IsFile: entry.Mode.IsFile(),
				Name:   name,
			})
		}
		c.Send <- send
	}
}
