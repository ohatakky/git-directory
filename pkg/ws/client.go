package ws

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type Client struct {
	Conn *websocket.Conn
}

func New(w http.ResponseWriter, r *http.Request) *Client {
	protocol := r.Header.Get("Sec-Websocket-Protocol")
	upgrader := websocket.Upgrader{
		Subprotocols: []string{protocol},
		// CheckOrigin: func(r *http.Request) bool {
		// 	return r.Header.Get("Origin") == "xxxxx"
		// },
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}

	return &Client{Conn: conn}
}
