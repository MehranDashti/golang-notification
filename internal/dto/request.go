package dto

type CreateNotificationRequest struct {
	UserID   string            `json:"user_id" binding:"required"`
	Channel  string            `json:"channel" binding:"required,oneof=email sms push"`
	Title    string            `json:"title" binding:"required"`
	Body     string            `json:"body" binding:"required"`
	Metadata map[string]string `json:"metadata"`
}

type CreateProviderConfigRequest struct {
	Channel     string            `json:"channel" binding:"required,oneof=sms email push"`
	Provider    string            `json:"provider" binding:"required,oneof=kavenegar twilio smtp fcm"`
	Credentials map[string]string `json:"credentials" binding:"required"`
	IsActive    bool              `json:"is_active"`
}

type UpdateProviderConfigRequest struct {
	Credentials map[string]string `json:"credentials" binding:"required"`
	IsActive    bool              `json:"is_active"`
}
