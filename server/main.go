package main

import (
	"log"
	"net/http"
	"os"

	"github.com/ohatakky/git-directory/server/pkg/git"
	"github.com/ohatakky/git-directory/server/pkg/uuid"
	"github.com/ohatakky/git-directory/server/pkg/ws"
)

// websocat ws://localhost:8080/ws?repo=gorilla/websocket | jq .

func main() {
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		repo := r.FormValue("repo")
		g := git.New(uuid.String(), repo)
		err := g.Clone()
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			log.Fatal(err)
			return
		}

		c, err := ws.New(w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Fatal(err)
			return
		}
		defer c.Conn.Close()

		go g.FuzzyFinder()

		for tree := range g.Send {
			err := c.Conn.WriteJSON(tree)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				log.Fatal(err)
				err = os.RemoveAll(g.TmpDir())
				if err != nil {
					log.Fatal(err)
				}
				return
			}
		}

		err = os.RemoveAll(g.TmpDir())
		if err != nil {
			log.Fatal(err)
			return
		}
	})

	log.Println("running...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
