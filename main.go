package main

import (
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/ohatakky/git-directory/pkg/git"
	"github.com/ohatakky/git-directory/pkg/ws"
)

func main() {
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		c := ws.New(w, r)
		defer c.Conn.Close()

		uuid, err := uuid.NewUUID()
		if err != nil {
			return
		}
		repo := r.FormValue("repo")
		g := git.New(uuid.String(), repo)
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
