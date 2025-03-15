package websocket

import (
	"sync"

	"github.com/hoyci/ms-chat/ws-service/types"
)

var (
	connections = make(map[string]types.Connection)
	mu          sync.RWMutex
)

func AddUserDeviceConnection(clientID string, conn types.Connection) {
	mu.Lock()
	defer mu.Unlock()
	connections[clientID] = conn
}

func RemoveConnection(clientID string) {
	mu.Lock()
	defer mu.Unlock()
	delete(connections, clientID)
}

func GetUserDevicesConnections(userID int) []types.Connection {
	mu.RLock()
	defer mu.RUnlock()

	result := make([]types.Connection, 0)

	for _, conn := range connections {
		if ok := conn.UserID == userID; ok {
			result = append(result, types.Connection{
				ClientID: conn.ClientID,
				UserID:   conn.UserID,
				Channel:  conn.Channel,
			})
		}
	}

	return result
}
