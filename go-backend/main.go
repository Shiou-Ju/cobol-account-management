package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

func main() {
	var databaseURL string = "postgres://postgres:cobolexamplepw@localhost:5432/cobolexample"

	dbpool, err := pgxpool.Connect(context.Background(), databaseURL)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	defer dbpool.Close()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		allUsersHandler(w, r, dbpool)
	})

	fmt.Println("Server is running on port 3001")
	http.ListenAndServe(":3001", nil)
}

type Transaction struct {
	User        string
	Transaction float64
	Balance     float64
	Date        time.Time
}

func allUsersHandler(w http.ResponseWriter, _ *http.Request, dbpool *pgxpool.Pool) {
	const sql = `SELECT * FROM (
		SELECT 
			transactions.*,
			ROW_NUMBER() OVER(PARTITION BY "user" ORDER BY "date" DESC) as rn
		FROM 
			transactions
	) t
	WHERE t.rn = 1;`

	rows, err := dbpool.Query(context.Background(), sql)
	if err != nil {
		fmt.Fprintf(w, "Query failed: %v\n", err)
		return
	}
	defer rows.Close()

	var transactions []Transaction

	for rows.Next() {
		var t Transaction
		var rn int
		err := rows.Scan(&t.User, &t.Transaction, &t.Balance, &t.Date, &rn)
		if err != nil {
			fmt.Fprintf(w, "Failed to scan row: %v\n", err)
			return
		}
		transactions = append(transactions, t)
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(transactions); err != nil {
		fmt.Fprintf(w, "Failed to encode transactions to JSON: %v\n", err)
	}
}
