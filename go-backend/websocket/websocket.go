package websocketconnection

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	redisChatroom "chatroom/redischatroom"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/websocket"
)

var clients = make(map[*websocket.Conn]bool)

var lock = sync.Mutex{}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		var canCrossOrigin = true
		return canCrossOrigin
	},
}

func HandleConnections(w http.ResponseWriter, r *http.Request, ctx context.Context, rdb *redis.Client, channel string) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
	}
	defer ws.Close()

	// TODO: why do we need lock here
	lock.Lock()
	var newClientStatus = true
	clients[ws] = newClientStatus
	fmt.Printf("clients\n")
	fmt.Print(clients)
	fmt.Printf("\n")
	lock.Unlock()

	for {
		var msg redisChatroom.ChatMessage
		err := ws.ReadJSON(&msg)
		if err != nil {
			lock.Lock()

			delete(clients, ws)
			lock.Unlock()
			fmt.Printf("error: %v", err)
			break
		}
		fmt.Printf("Received message in HandleConnections: %+v\n", msg)

		jsonMsg, err := json.Marshal(msg)
		if err != nil {
			fmt.Printf("Failed to serialize message: %v", err)
			continue
		}

		redisChatroom.PublishMessage(ctx, rdb, channel, string(jsonMsg))
	}

}

// TODO: make sure this works
func BroadcastMessage(message string) {
	lock.Lock()
	defer lock.Unlock()

	fmt.Printf("inside BroadcastMessage")

	for client := range clients {
		err := client.WriteMessage(websocket.TextMessage, []byte(message))
		if err != nil {
			fmt.Printf("error: %v", err)
			client.Close()
			delete(clients, client)
		}
	}
}
