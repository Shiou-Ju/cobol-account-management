package subscribemessage

import (
	websocketconnection "chatroom/websocket"
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
)

// TODO: do we still need this func?
// or refactor into a long poll
func SubscribeMessages(ctx context.Context, rdb *redis.Client, channel string) {
	sub := rdb.Subscribe(ctx, channel)

	// TODO: where to close sub redis
	// maybe time like SIGTERM
	// defer sub.Close()

	go func() {
		fmt.Println("SubscribeMessages Started listening for messages")
		for {
			select {
			case msg, ok := <-sub.Channel():
				if !ok {
					fmt.Println("Channel closed")
					return
				}
				if msg == nil {
					fmt.Println("Received nil message")
					continue
				}
				// fmt.Println("Received message:", msg.Payload)
				fmt.Println("Received message in SubscribeMessages: ", msg.Payload)
				websocketconnection.BroadcastMessage(msg.Payload)
			case <-ctx.Done():
				fmt.Println("Subscription stopped")
				return
			}
		}
	}()

}
