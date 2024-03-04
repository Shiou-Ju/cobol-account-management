package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	chatroom "chatroom/redischatroom"
	"chatroom/subscribemessage"
	"chatroom/usermanagement"

	socket "chatroom/websocket"

	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v4/pgxpool"
)

func main() {
	userManager := usermanagement.NewUserManager()

	var databaseURL string = "postgres://postgres:cobolexamplepw@localhost:5432/cobolexample"

	dbpool, err := pgxpool.Connect(context.Background(), databaseURL)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	defer dbpool.Close()

	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	var ctx = context.Background()

	isRedisChannelDone := make(chan bool)

	go func() {
		chatroom.PublishMessage(ctx, rdb, chatroom.RedisChannelName, "Hello, World!")
		isRedisChannelDone <- true
	}()

	<-isRedisChannelDone
	go subscribemessage.SubscribeMessages(ctx, rdb, chatroom.RedisChannelName)

	populateUserManager(dbpool, userManager)

	http.HandleFunc("/go-api/users", func(w http.ResponseWriter, r *http.Request) {
		usermanagement.AllUsersHandler(w, r, dbpool, userManager)
	})

	http.HandleFunc("/go-api/user-state", func(w http.ResponseWriter, r *http.Request) {
		usermanagement.GetUserStateHandler(w, r, userManager)
	})

	http.HandleFunc("/go-api/set-user-state", func(w http.ResponseWriter, r *http.Request) {
		usermanagement.SetUserStateHandler(w, r, userManager)
	})

	http.HandleFunc("/go-api/try-lock-user", func(w http.ResponseWriter, r *http.Request) {
		usermanagement.TryLockUserHandler(w, r, userManager)
	})

	http.HandleFunc("/go-api/send-message", func(w http.ResponseWriter, r *http.Request) {
		chatroom.SendChatMessage(ctx, w, r, rdb)
	})

	http.HandleFunc("/go-api/ws", func(w http.ResponseWriter, r *http.Request) {
		socket.HandleConnections(w, r, ctx, rdb, "chatroom")
	})

	fmt.Println("Server is running on port 3001")
	http.ListenAndServe(":3001", nil)
}

func populateUserManager(dbpool *pgxpool.Pool, userManager *usermanagement.UserManager) {
	const sql = `SELECT DISTINCT "user" FROM transactions;`

	rows, err := dbpool.Query(context.Background(), sql)
	if err != nil {
		fmt.Println("Query failed:", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var userName string
		err := rows.Scan(&userName)
		if err != nil {
			fmt.Println("Failed to scan row:", err)
			continue
		}
		userManager.AddUser(userName)
	}
}
