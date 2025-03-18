package types

import (
	"time"

	coreTypes "github.com/hoyci/ms-chat/core/types"
)

type BroadcastMessage struct {
	UserID    int                 `json:"user_id"`
	Messages  []coreTypes.Message `json:"messages"`
	Timestamp time.Time           `json:"timestamp"`
}
