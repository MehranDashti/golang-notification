package kafka

type NotificationMessage struct {
	NotificationID string            `json:"notification_id"`
	UserID         string            `json:"user_id"`
	Channel        string            `json:"channel"`
	Title          string            `json:"title"`
	Body           string            `json:"body"`
	Metadata       map[string]string `json:"metadata"`
	Attempt        int               `json:"attempt"`
}
