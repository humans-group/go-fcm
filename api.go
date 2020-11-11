package fcm

// easyjson:json
type sendRequest struct {
	Message *Message `json:"message"`
}

// easyjson:json
type response struct {
	Error *responseError `json:"error,omitempty"`
}

// easyjson:json
type responseError struct {
	Code    int       `json:"code"`
	Message string    `json:"message,omitempty"`
	Status  string    `json:"status"`
	Details []details `json:"details,omitempty"`
}

// details could be different type of structs,
// we are interested in one that has errorCode.
// easyjson:json
type details struct {
	ErrorCode errorCode `json:"errorCode,omitempty"`
}

// Possible error codes are listed here
// https://firebase.google.com/docs/reference/fcm/rest/v1/ErrorCode
type errorCode string

const (
	// Unknown error
	errorCodeUnspecified errorCode = "UNSPECIFIED_ERROR"

	// App instance was unregistered from FCM for HTTP error code = 404
	// This usually means that the token used is no longer valid (i.e. expired) and a new one must be used.
	errorCodeUnregistered errorCode = "UNREGISTERED"
)