package git

import (
	"github.com/go-git/go-git/v5/plumbing/object"
)

type Fzf struct {
	Hash   string   `json:"hash"`
	Files  []string `json:"files"`
	Commit struct {
		Message string `json:"message"`
		Author  string `json:"author"`
	} `json:"commit"`
}

func (c *Client) FuzzyFinder(commits []Commit) {
	defer close(c.Send)
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
			c.Error <- err
		}
		c.Send <- &fzf
	}
}
