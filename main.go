package main

import (
	"fmt"
	"log"
	"net/http"

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

		for {
			// Read message from browser
			msgType, msg, err := c.Conn.ReadMessage()
			if err != nil {
				return
			}

			// Print the message to the console
			fmt.Printf("%s sent: %s\n", c.Conn.RemoteAddr(), string(msg))

			// Write message back to browser
			if err = c.Conn.WriteMessage(msgType, msg); err != nil {
				return
			}
		}
	})
	log.Println("running...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
