package events

type DocumentEvent struct {
	EventType string `json:"eventType"`
	ID        uint   `json:"id"`
	Title     string `json:"title"`
	Author    string `json:"author"`
	Content   string `json:"content,omitempty"`
}
