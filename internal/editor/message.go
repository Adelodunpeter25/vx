package editor

import "time"

// MessageType defines the type of status message
type MessageType int

const (
	// MessageTransient shows briefly then disappears (e.g., "Copied")
	MessageTransient MessageType = iota
	// MessagePersistent stays until replaced (e.g., file info)
	MessagePersistent
	// MessageError stays until user action (e.g., save errors)
	MessageError
)

// Message represents a status bar message
type Message struct {
	text      string
	msgType   MessageType
	timestamp time.Time
}

// MessageManager handles status bar messages with auto-expiry
type MessageManager struct {
	current *Message
}

// NewMessageManager creates a new message manager
func NewMessageManager() *MessageManager {
	return &MessageManager{}
}

// Set sets a message with the given type
func (m *MessageManager) Set(text string, msgType MessageType) {
	m.current = &Message{
		text:      text,
		msgType:   msgType,
		timestamp: time.Now(),
	}
}

// SetTransient sets a transient message (auto-clears after 2 seconds)
func (m *MessageManager) SetTransient(text string) {
	m.Set(text, MessageTransient)
}

// SetPersistent sets a persistent message (stays until replaced)
func (m *MessageManager) SetPersistent(text string) {
	m.Set(text, MessagePersistent)
}

// SetError sets an error message (stays until user action)
func (m *MessageManager) SetError(text string) {
	m.Set(text, MessageError)
}

// Get returns the current message text, or empty if expired
func (m *MessageManager) Get() string {
	if m.current == nil {
		return ""
	}
	
	// Check if transient message has expired (2 seconds)
	if m.current.msgType == MessageTransient {
		if time.Since(m.current.timestamp) > 2*time.Second {
			m.current = nil
			return ""
		}
	}
	
	return m.current.text
}

// Clear clears the current message
func (m *MessageManager) Clear() {
	m.current = nil
}

// ClearIfTransient clears only if current message is transient
func (m *MessageManager) ClearIfTransient() {
	if m.current != nil && m.current.msgType == MessageTransient {
		m.current = nil
	}
}
