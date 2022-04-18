package models

import "time"

type Client struct {
	UserAgent string `json:"userAgent"`
	OS        string `json:"os"`
	Browser   string `json:"browser"`
	Version   string `json:"version"`
}

type Record struct {
	ID        string            `json:"id" bow:"key"`
	Events    []Events          `json:"events"`
	Meta      map[string]string `json:"meta"`
	User      User              `json:"user"`
	Client    Client            `json:"client"`
	ClientID  string            `json:"client_id"`
	UpdatedAt time.Time         `json:"updatedAt"`
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
