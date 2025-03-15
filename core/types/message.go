package types

type Message struct {
	TempID    string `json:"temp_id" validate:"required"`
	Type      string `json:"type" validate:"required"`
	UserID    int    `json:"user_id" validate:"required"`
	RoomID    string `json:"room_id" validate:"required"`
	Content   string `json:"content" validate:"required"`
	Timestamp string `json:"timestamp"`
	ClientID  string `json:"client_id"`
}
