package models

type Client struct {
	ClientID  string `json:"client_id,omitempty"`
	UserAgent string `json:"user_agent,omitempty"`
	OS        string `json:"os,omitempty"`
	Browser   string `json:"browser,omitempty"`
	Version   string `json:"version,omitempty"`
}

type Record struct {
	ID        string  `json:"id,omitempty"`
	Events    []Event `json:"events,omitempty"`
	User      User    `json:"user,omitempty"`
	Client    Client  `json:"client,omitempty"`
	UpdatedAt string  `json:"updated_at,omitempty"`
}

type Events struct {
	Events []Event `json:"events,omitempty"`
}
type Event struct {
	ID        string      `json:"id,omitempty"`
	Type      int64       `json:"type,omitempty"`
	Data      interface{} `json:"data,omitempty"`
	Timestamp int64       `json:"timestamp,omitempty"`
}

type User struct {
	ID   string `json:"user_id,omitempty"`
	Name string `json:"name,omitempty"`
}
