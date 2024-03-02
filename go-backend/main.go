package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"chatroom/usermanagement"

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

	populateUserManager(dbpool, userManager)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		usermanagement.AllUsersHandler(w, r, dbpool, userManager)
	})

	http.HandleFunc("/user-state", func(w http.ResponseWriter, r *http.Request) {
		usermanagement.GetUserStateHandler(w, r, userManager)
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
