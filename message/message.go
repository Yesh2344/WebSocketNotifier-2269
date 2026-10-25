package message
// minor polish

import "encoding/json"

// Payload represents the generic data part of a notification.
type Payload map[string]interface{}

// Message defines the wire format for notifications.
type Message struct {
	Type string  `json:"type"` // e.g., "notification"
	Data Payload `json:"data"`
}

// Encode marshals a Message into JSON bytes.
func (m *Message) Encode() ([]byte, error) {
	return json.Marshal(m)
}

// Decode unmarshals JSON bytes into a Message.
func Decode(data []byte) (*Message, error) {
	var m Message
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, err
	}
	return &m, nil
}