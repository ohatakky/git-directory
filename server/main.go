package main

import (
	"fmt"
	"log"
	"net/http"

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
			fmt.Println(err)
			return
		}

		c := ws.New(w, r)
		defer c.Conn.Close()

		go g.Trees()

		for tree := range g.Send {
			err := c.Conn.WriteJSON(tree)
			if err != nil {
				return
			}
		}
	})

	log.Println("running...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
