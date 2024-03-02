package chatroom

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
)

type ChatMessage struct {
	Username string `json:"username"`
	Message  string `json:"message"`
	Time     int64  `json:"time"`
}

var ctx = context.Background()

func PublishMessage(rdb *redis.Client, channel, message string) error {
	err := rdb.Publish(ctx, channel, message).Err()

	if err != nil {
		fmt.Printf("Error publishing message to channel %s: %v\n", channel, err)
		return err
	}

	fmt.Println("Message published:", message)
	return nil
}

func SubscribeMessages(rdb *redis.Client, channel string) {
	sub := rdb.Subscribe(ctx, channel)
	defer sub.Close()

	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()
		for {
			select {
			case msg := <-sub.Channel():
				fmt.Println("Received message:", msg.Payload)
			case <-ctx.Done():
				fmt.Println("Subscription stopped")
				return
			}
		}
	}()
	wg.Wait()
}

func SendChatMessage(w http.ResponseWriter, r *http.Request, rdb *redis.Client) {
	var msg ChatMessage
	if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	msg.Time = time.Now().Unix()

	messageJSON, err := json.Marshal(msg)
	if err != nil {
		http.Error(w, "Failed to encode message", http.StatusInternalServerError)
		return
	}

	if err := PublishMessage(rdb, "chatroom", string(messageJSON)); err != nil {
		fmt.Printf("Error publishing message: %v\n", err)
		http.Error(w, "Failed to publish message", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Message sent successfully")
}
