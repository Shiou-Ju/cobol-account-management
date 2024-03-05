package connectionmanagement

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

type UserConnectionManager struct {
	connections map[string]string
	mutex       sync.RWMutex
}

func NewUserConnectionManager() *UserConnectionManager {
	return &UserConnectionManager{
		connections: make(map[string]string),
	}
}

func (m *UserConnectionManager) AddConnectionToUser(username, connectionHash string) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.connections[username] = connectionHash
}

func (m *UserConnectionManager) RemoveConnectionByUser(username string) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	delete(m.connections, username)
}

func (m *UserConnectionManager) GetConnectionByUser(username string) (string, bool) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	hash, exists := m.connections[username]
	return hash, exists
}

func (m *UserConnectionManager) GetUserByConnection(hashedAddress string) (string, bool) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	for username, hash := range m.connections {
		if hash == hashedAddress {
			return username, true
		}
	}
	return "", false
}

func (m *UserConnectionManager) HandleAddConnectionToUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var data struct {
		Username   string `json:"username"`
		Connection string `json:"connection"`
	}

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	fmt.Println("data in HandleAddConnectionToUser")
	fmt.Println(data)

	m.AddConnectionToUser(data.Username, data.Connection)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Connection added successfully"))
}
