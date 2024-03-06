package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"chatroom/connectionmanagement"
	chatroom "chatroom/redischatroom"
	"chatroom/subscribemessage"
	"chatroom/usermanagement"

	socket "chatroom/websocket"

	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v4/pgxpool"
)

func main() {
	userManager := usermanagement.NewUserManager()

	connectionManager := connectionmanagement.NewUserConnectionManager()

	// var databaseURL string = "postgres://postgres:cobolexamplepw@localhost:5432/cobolexample"

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		databaseURL = "postgres://postgres:cobolexamplepw@localhost:5432/cobolexample"
	}

	dbpool, err := pgxpool.Connect(context.Background(), databaseURL)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	defer dbpool.Close()

	redisAddress := os.Getenv("REDIS_ADDR")
	if redisAddress == "" {
		redisAddress = "localhost:6379"
	}

	rdb := redis.NewClient(&redis.Options{
		// Addr:     "localhost:6379",
		Addr:     redisAddress,
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

	http.Handle("/go-api/users", setupCORS(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		usermanagement.AllUsersHandler(w, r, dbpool, userManager)
	})))

	http.Handle("/go-api/user-state", setupCORS(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		usermanagement.GetUserStateHandler(w, r, userManager)
	})))

	http.Handle("/go-api/set-user-state", setupCORS(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		usermanagement.SetUserStateHandler(w, r, userManager)
	})))

	http.Handle("/go-api/try-lock-user", setupCORS(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		usermanagement.TryLockUserHandler(w, r, userManager)
	})))

	http.Handle("/go-api/try-unlock-user", setupCORS(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		usermanagement.TryUnlockUserHandler(w, r, userManager)
	})))

	http.Handle("/go-api/add-connection-to-user", setupCORS((http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		connectionManager.HandleAddConnectionToUser(w, r)
	}))))

	http.Handle("/go-api/send-message", setupCORS(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		chatroom.SendChatMessage(ctx, w, r, rdb)
	})))

	http.Handle("/go-api/disconnect-all", setupCORS(http.HandlerFunc(socket.DisconnectAllHandler)))

	http.HandleFunc("/go-api/ws", func(w http.ResponseWriter, r *http.Request) {
		socket.HandleConnections(w, r, ctx, rdb, chatroom.RedisChannelName, userManager, connectionManager)
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

func setupCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		if r.Method == "OPTIONS" {
			return
		}

		next.ServeHTTP(w, r)
	})
}
