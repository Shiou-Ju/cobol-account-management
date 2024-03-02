package usermanagement

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

type Transaction struct {
	User        string
	Transaction float64
	Balance     float64
	Date        time.Time
}

type UserState bool

const (
	Picked    UserState = true
	Available UserState = false
)

type User struct {
	Name  string
	State UserState
	Mutex sync.Mutex
}

type UserManager struct {
	Users map[string]*User
	Mutex sync.RWMutex
}

func NewUserManager() *UserManager {
	return &UserManager{
		Users: make(map[string]*User),
	}
}

func (m *UserManager) AddUser(name string) {
	m.Mutex.Lock()
	defer m.Mutex.Unlock()
	m.Users[name] = &User{
		Name:  name,
		State: Available,
	}
}

func (m *UserManager) SetUserState(name string, state UserState) {
	m.Mutex.RLock()
	user, exists := m.Users[name]
	m.Mutex.RUnlock()

	if exists {
		user.Mutex.Lock()
		user.State = state
		user.Mutex.Unlock()
	}
}

func (m *UserManager) Exists(name string) bool {
	m.Mutex.RLock()
	defer m.Mutex.RUnlock()
	_, exists := m.Users[name]
	return exists
}

func (m *UserManager) GetUserState(name string) UserState {
	m.Mutex.RLock()
	user, exists := m.Users[name]
	m.Mutex.RUnlock()

	// TODO: maybe not should only limited to 3 users
	if exists {
		return user.State
	}

	return Available
}

func (m *UserManager) TryLockUser(name string) bool {
	m.Mutex.Lock()
	defer m.Mutex.Unlock()

	user, exists := m.Users[name]
	if !exists || user.State == Picked {
		return false
	}

	user.State = Picked
	return true
}
func AllUsersHandler(w http.ResponseWriter, _ *http.Request, dbpool *pgxpool.Pool, userManager *UserManager) {

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

		userManager.AddUser(t.User)
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(transactions); err != nil {
		fmt.Fprintf(w, "Failed to encode transactions to JSON: %v\n", err)
	}
}

func GetUserStateHandler(w http.ResponseWriter, r *http.Request, userManager *UserManager) {
	userName := r.URL.Query().Get("username")
	if userName == "" {
		http.Error(w, "Username is required", http.StatusBadRequest)
		return
	}

	state := userManager.GetUserState(userName)

	response := struct {
		Username string `json:"username"`
		Status   string `json:"status"`
	}{
		Username: userName,
		Status:   "available",
	}

	if state == Picked {
		response.Status = "picked"
	}

	if state == Available && !userManager.Exists(userName) {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		fmt.Fprintf(w, "Failed to encode response to JSON: %v\n", err)
	}
}

func SetUserStateHandler(w http.ResponseWriter, r *http.Request, userManager *UserManager) {
	var requestData struct {
		Username string `json:"username"`
		Status   string `json:"status"`
	}
	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	var isInvalidBody = requestData.Username == "" || (requestData.Status != "picked" && requestData.Status != "available")

	if isInvalidBody {
		http.Error(w, "Invalid username or status", http.StatusBadRequest)
		return
	}

	var newState UserState
	if requestData.Status == "picked" {
		newState = Picked
	} else {
		newState = Available
	}
	userManager.SetUserState(requestData.Username, newState)

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "User state updated successfully")
}

func TryLockUserHandler(w http.ResponseWriter, r *http.Request, userManager *UserManager) {
	var requestData struct {
		Username string `json:"username"`
	}
	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if requestData.Username == "" {
		http.Error(w, "Username is required", http.StatusBadRequest)
		return
	}

	if !userManager.Exists(requestData.Username) {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	if !userManager.TryLockUser(requestData.Username) {
		http.Error(w, "User already locked", http.StatusConflict)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "User locked successfully")
}
