package main

import (
	"log"
	"net/http"

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

		for i := 0; i < 100; i++ {
			err := c.Conn.WriteJSON(&Hoge{
				Num: 3,
			})
			if err != nil {
				return
			}
		}
	})

	log.Println("running...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
