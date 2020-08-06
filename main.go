package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/ohatakky/git-directory/pkg/git"
	"github.com/ohatakky/git-directory/pkg/ws"
)

type Hoge struct {
	Num int `json:"num"`
}

func main() {
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		c := ws.New(w, r)
		defer c.Conn.Close()

		g := git.New()
		commits, err := g.Commits()
		if err != nil {
			return
		}

		// todo: pkgから溢れてるので、channelで返すようにする
		for i := len(commits) - 1; i >= 0; i-- {
			fmt.Println(commits[i].Hash)
			res := git.Tree{
				Idx: i,
				T:   make([]string, 0),
			}
			err := commits[i].Files.ForEach(func(f *object.File) error {
				res.T = append(res.T, f.Name)
				return nil
			})
			if err != nil {
				log.Fatal(err)
				return
			}
			err = c.Conn.WriteJSON(&res)
			if err != nil {
				return
			}
		}
	})

	log.Println("running...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
