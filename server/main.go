package main

import (
	"log"
	"net/http"
	"os"

	"github.com/ohatakky/git-directory/server/pkg/git"
	"github.com/ohatakky/git-directory/server/pkg/uuid"
	"github.com/ohatakky/git-directory/server/pkg/ws"
)

func main() {
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		// websocket upgrader
		c, err := ws.New(w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer c.Conn.Close()

		// git client
		org := r.FormValue("org")
		repo := r.FormValue("repo")
		g := git.New(uuid.String(), org, repo)
		err = g.Clone()
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		commits, err := g.Commits()
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// fzf
		go g.FuzzyFinder(commits)
		// tree
		// go g.Tree(commits)

		for {
			var done bool
			select {
			case send, ok := <-g.Send:
				if !ok {
					done = true
				}
				err := c.Conn.WriteJSON(send)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					done = true
				}
			case err := <-g.Error:
				http.Error(w, err.Error(), http.StatusInternalServerError)
				done = true
			}

			if done {
				break
			}
		}

		// remove tmp dir
		err = os.RemoveAll(g.TmpDir())
		if err != nil {
			log.Println(err)
		}
	})

	log.Println("running...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
