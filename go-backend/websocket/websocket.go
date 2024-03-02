package websocketconnection

import (
	"fmt"
	"net/http"

	chatroom "chatroom/redischatroom"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		var canCrossOrigin = true
		return canCrossOrigin
	},
}

func HandleConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
	}
	defer ws.Close()

	for {
		var msg chatroom.ChatMessage
		err := ws.ReadJSON(&msg)
		if err != nil {
			fmt.Printf("error: %v", err)
			break
		}
		fmt.Printf("Received message: %+v\n", msg)
	}
}
