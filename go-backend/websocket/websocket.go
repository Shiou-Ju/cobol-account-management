package websocketconnection

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"chatroom/connectionmanagement"
	redisChatroom "chatroom/redischatroom"
	"chatroom/usermanagement"

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

func HandleConnections(w http.ResponseWriter, r *http.Request, ctx context.Context, rdb *redis.Client, channel string, userManager *usermanagement.UserManager, connectionManager *connectionmanagement.UserConnectionManager) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
	}
	defer ws.Close()

	// TODO: why do we need lock here
	lock.Lock()
	var newClientStatus = true
	clients[ws] = newClientStatus

	connectionAddress := fmt.Sprintf("%p", ws)
	hash := sha256.Sum256([]byte(connectionAddress))
	// TODO: security issue
	hashedAddress := hex.EncodeToString(hash[:])

	initialMessage := map[string]string{"connection": hashedAddress, "isMessage": "false"}

	err = ws.WriteJSON(initialMessage)
	if err != nil {
		fmt.Printf("Failed to send initial message: %v\n", err)
	}

	fmt.Printf("clients map:\n")
	fmt.Print(clients)
	fmt.Printf("clients map finished\n")
	lock.Unlock()

	go func(conn *websocket.Conn, hashedAddress string) {
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()

		// TODO: for select warning
		for {
			select {
			case <-ticker.C:
				if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
					fmt.Printf("Failed to send ping to %s: %v\n", hashedAddress, err)

					handleFailedPing(hashedAddress, userManager, connectionManager)

					return
				}
			}
		}
	}(ws, hashedAddress)

	for {
		var msg redisChatroom.ChatMessage
		err := ws.ReadJSON(&msg)
		if err != nil {
			lock.Lock()

			delete(clients, ws)
			lock.Unlock()
			fmt.Printf("error in HandleConnections: %v", err)
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

// TODO: add more error logs
func handleFailedPing(hashedAddress string, userManager *usermanagement.UserManager, connectionManager *connectionmanagement.UserConnectionManager) {
	fmt.Printf("Handling failed ping for connection: %s\n", hashedAddress)

	username, isUserFindingSuccess := connectionManager.GetUserByConnection(hashedAddress)

	if !isUserFindingSuccess {
		fmt.Println("faliled to find user name using hash")
	}

	isUnlockSuccess := userManager.TryUnlockUser(username)

	fmt.Println("TryUnlockUser success?: ")
	fmt.Println(isUnlockSuccess)
	fmt.Println("")

	connectionManager.RemoveConnectionByUser(username)

	fmt.Println("removed user from connection map")
	fmt.Println("")

}

func DisconnectAll() bool {
	lock.Lock()
	defer lock.Unlock()

	for client := range clients {
		err := client.Close()
		if err != nil {

			return false
		}
		delete(clients, client)
	}

	return true
}

func DisconnectAllHandler(w http.ResponseWriter, r *http.Request) {
	isSuccess := DisconnectAll()

	if !isSuccess {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Error disconnecting all connections.")
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "All connections have been successfully disconnected.")
}
