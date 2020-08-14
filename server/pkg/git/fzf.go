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
	Send chan *Fzf
}

func New(id, repo string) *Client {
	return &Client{
		ID:   id,
		Repo: repo,
		Send: make(chan *Fzf, 256),
	}
}

func (c *Client) TmpDir() string {
	return fmt.Sprintf(tmpDirFmt, strings.Replace(c.Repo, "/", "_", -1), c.ID)
}

type Commit struct {
	Hash    string
	Message string
	Author  string
	Files   *object.FileIter
}

func (c *Client) Clone() error {
	r, err := git.PlainClone(c.TmpDir(), false, &git.CloneOptions{
		URL:          fmt.Sprintf("https://github.com/%s", c.Repo),
		Progress:     os.Stdout,
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

func (c *Client) commits() ([]Commit, error) {
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

type Fzf struct {
	Hash   string   `json:"hash"`
	Files  []string `json:"files"`
	Commit struct {
		Message string `json:"message"`
		Author  string `json:"author"`
	} `json:"commit"`
}

// todo: error handling (<-error)
func (c *Client) FuzzyFinder() error {
	defer close(c.Send)
	commits, err := c.commits()
	if err != nil {
		return err
	}

	for i := len(commits) - 1; i >= 0; i-- {
		fzf := Fzf{
			Hash:  commits[i].Hash,
			Files: make([]string, 0),
			Commit: struct {
				Message string `json:"message"`
				Author  string `json:"author"`
			}{
				Message: commits[i].Message,
				Author:  commits[i].Author,
			},
		}
		err := commits[i].Files.ForEach(func(f *object.File) error {
			fzf.Files = append(fzf.Files, f.Name)
			return nil
		})
		if err != nil {
			return err
		}
		c.Send <- &fzf
	}

	return nil
}
