package dto

type CreateNotificationRequest struct {
	UserID   string            `json:"user_id" binding:"required"`
	Channel  string            `json:"channel" binding:"required"`
	Title    string            `json:"title" binding:"required"`
	Body     string            `json:"body" binding:"required"`
	Metadata map[string]string `json:"meta_title"`
}
