package git

import (
	"fmt"

	"github.com/go-git/go-git/v5"
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
	Send  chan *Fzf
	Error chan error
}

func New(id, org, repo string) *Client {
	return &Client{
		ID:   id,
		Org:  org,
		Repo: repo,
		Send: make(chan *Fzf, 256),
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

type Commit struct {
	Hash    string
	Message string
	Author  string
	Files   *object.FileIter
}

func (c *Client) Commits() ([]Commit, error) {
	commits := make([]Commit, 0)
	cIter, err := c.Ref.Log(&git.LogOptions{
		Order: git.LogOrderCommitterTime,
	})
	if err != nil {
		return nil, err
	}
	err = cIter.ForEach(func(c *object.Commit) error {
		files, err := c.Files()
		if err != nil {
			return err
		}

		commits = append(commits, Commit{
			Hash:    c.Hash.String(),
			Author:  c.Author.Name,
			Message: c.Message,
			Files:   files,
		})
		return nil
	})
	if err != nil {
		return nil, err
	}

	return commits, nil
}
