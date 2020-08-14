package ws

import (
	"net/http"

	"github.com/gorilla/websocket"
)

type Client struct {
	Conn *websocket.Conn
}

func New(w http.ResponseWriter, r *http.Request) (*Client, error) {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			// return r.Header.Get("Origin") == "xxxxx"
			return true
		},
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return nil, err
	}

	return &Client{Conn: conn}, nil
}
