package fcm

// easyjson:json
type sendRequest struct {
	Message *Message `json:"message"`
}

// easyjson:json
type Response struct {
	ErrorCode ErrorCode `json:"error_code,omitempty"`
}
