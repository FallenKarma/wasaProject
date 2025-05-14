package models

import (
	"time"
)

// User represents a WASAText user
type User struct {
	ID       string `json:"identifier"`
	Name     string `json:"name"`
	PhotoURL string `json:"photo,omitempty"`
}

// MessageType defines the type of message
type MessageType string

const (
	TextMessage  MessageType = "text"
	PhotoMessage MessageType = "photo"
)

// MessageStatus defines the status of a message
type MessageStatus string

const (
	Sent     MessageStatus = "sent"
	Received MessageStatus = "received"
	Read     MessageStatus = "read"
)

// Message represents a message in WASAText
type Message struct {
	ID        string        `json:"id"`
	Sender    string        `json:"sender"`
	Timestamp time.Time     `json:"timestamp"`
	Content   string        `json:"content"`
	Type      MessageType   `json:"type"`
	Status    MessageStatus `json:"status"`
	ReplyTo   string        `json:"replyTo,omitempty"` // ID of message being replied to
}

// Reaction represents a user's reaction to a message
type Reaction struct {
	MessageID string `json:"messageId"`
	UserID    string `json:"user"`
	Emoji     string `json:"emoji"`
}

// ConversationType defines the type of conversation
type ConversationType string

const (
	DirectConversation ConversationType = "direct"
	GroupConversation  ConversationType = "group"
)

// Conversation represents a chat conversation
type Conversation struct {
	ID           string          `json:"id"`
	Name         string          `json:"name"`
	Type         ConversationType `json:"type"`
	PhotoURL     string          `json:"photo,omitempty"`
	Participants []Participant        `json:"participants"`
	LastMessage  *Message        `json:"lastMessage,omitempty"`
	Messages     []Message       `json:"messages,omitempty"`
}

type Participant struct {
    ID   string `json:"id"`
    Name string `json:"name"`
}

// CreateConversationRequest represents the request to create a new conversation
type CreateConversationRequest struct {
	Participants []string `json:"participants"` // UserIDs of participants (excluding the creator)
	Type         ConversationType `json:"type"`
	Name         string          `json:"name"`
}

// LoginRequest represents the login request body
type LoginRequest struct {
	Name string `json:"name"`
}

// LoginResponse represents the login response
type LoginResponse struct {
	Identifier string `json:"identifier"`
}

// UpdateUsernameRequest represents the update username request body
type UpdateUsernameRequest struct {
	Name string `json:"name"`
}

// AddToGroupRequest represents the request to add a user to a group
type AddToGroupRequest struct {
	UserID string `json:"userId"`
}

// SetGroupNameRequest represents the request to set a group name
type SetGroupNameRequest struct {
	Name string `json:"name"`
}

// ForwardMessageRequest represents the request to forward a message
type ForwardMessageRequest struct {
	MessageID            string `json:"messageId"`
	TargetConversationID string `json:"targetConversationId"`
}