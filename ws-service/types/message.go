package types

type Message struct {
	UserID    int    `json:"user_id"`
	Room      string `json:"room"`
	Content   string `json:"content"`
	Timestamp string `json:"timestamp"`
	ClientID  string `json:"client_id"`
	Type      string `json:"type"`
}
