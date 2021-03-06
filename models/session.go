package models

type Client struct {
	ClientID  string `json:"client_id,omitempty"`
	UserAgent string `json:"user_agent,omitempty"`
	OS        string `json:"os,omitempty"`
	Browser   string `json:"browser,omitempty"`
	Version   string `json:"version,omitempty"`
}

type Session struct {
	SessionID string  `json:"session_id,omitempty"`
	Events    []Event `json:"events,omitempty"`
	User      User    `json:"user,omitempty"`
	Client    Client  `json:"client,omitempty"`
	UpdatedAt string  `json:"updated_at,omitempty"`
}

type Events struct {
	Events []Event `json:"events,omitempty"`
}
type Event struct {
	Type      int64       `json:"type,omitempty"`
	Data      interface{} `json:"data,omitempty"`
	Timestamp int64       `json:"timestamp,omitempty"`
}

type User struct {
	UserID   string `json:"user_id,omitempty"`
	UserName string `json:"user_name,omitempty"`
}
