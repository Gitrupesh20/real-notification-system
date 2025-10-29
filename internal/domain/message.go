package domain

import "fmt"

type Message struct {
	Id        string `json:"id"`
	Type      string `json:"type"`
	Body      string `json:"body"`
	Timestamp int64  `json:"timestamp"`
	UserId    string `json:"userId,omitempty"` // if empty means fan out
}

func NewMessage() {
	//wil complete this fns
}

func (m *Message) Validate() (bool, error) {
	if len(m.Id) == 0 {
		return false, fmt.Errorf("id is empty")
	} else if len(m.Type) == 0 {
		return false, fmt.Errorf("type is empty")
	} else if m.Timestamp == 0 {
		return false, fmt.Errorf("timestamp is empty")
	}

	return true, nil
}
