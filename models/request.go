package models

type Request struct {
	SessionID string  `json:"session_id"`
	ClientID  string  `json:"client_id"`
	UserID    string  `json:"user_id"`
	User      User    `json:"user"`
	Events    []Event `json:"events"`
}
