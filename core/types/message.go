package types

type Message struct {
	TempID     string `json:"temp_id" validate:"required"`
	SenderID   int    `json:"sender_id" validate:"required"`
	ReceiverID int    `json:"receiver_id" validate:"required"`
	RoomID     string `json:"room_id" validate:"required"`
	Content    string `json:"content" validate:"required"`
	Timestamp  string `json:"timestamp"`
	ClientID   string `json:"client_id"`
}
