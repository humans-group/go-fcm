package fcm

// easyjson:json
type sendRequest struct {
	Message *Message `json:"message"`
}

type Response struct {
	ErrorCode ErrorCode `json:"error_code,omitempty"`
}
