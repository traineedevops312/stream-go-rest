package Messages

import (
	"github.com/gocql/gocql"
)

// Message struct for preparing JSON payload
type Message struct {
	ID           gocql.UUID `json:"id"`
	UserID       gocql.UUID `json:"userID"`
	UserFullName string     `json:"user_full_name"`
	Message      string     `json:"message"`
}

// GetMessageResponse struct for embedding a single message
type GetMessageResponse struct {
	Message Message `json:"message"`
}

// AllMessagesResponse struct for an array of Message structs
type AllMessagesResponse struct {
	Messages []Message `json:"messages"`
}

// NewMessageResponse struct for returning ID of message in payload
type NewMessageResponse struct {
	ID gocql.UUID `json:"id"`
}

// ErrorResponse for sending back a potential array of error strings
type ErrorResponse struct {
	Errors []string `json:"errors"`
}
