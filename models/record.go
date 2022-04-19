package models

type Client struct {
	ClientID  string `json:"client_id"`
	UserAgent string `json:"user_agent"`
	OS        string `json:"os"`
	Browser   string `json:"browser"`
	Version   string `json:"version"`
}

type Record struct {
	ID        string   `json:"id"`
	Events    []Events `json:"events"`
	User      User     `json:"user"`
	Client    Client   `json:"client"`
	UpdatedAt string   `json:"updated_at"`
}

type Events struct {
	Type      int64       `json:"type"`
	Data      interface{} `json:"data"`
	Timestamp int64       `json:"timestamp"`
}

type User struct {
	ID    string `json:"user_id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}
