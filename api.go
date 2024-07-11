package fcm

// easyjson:json
type sendRequest struct {
	Message *Message `json:"message"`
}

// There is not clear doc for the API sendResponse
// but it could be checked using this API playground
// https://firebase.google.com/docs/reference/fcm/rest/v1/projects.messages/send
// easyjson:json
type sendResponse struct {
	Error *responseError `json:"error,omitempty"`
}

// Mark all fields like omitempty, as we don't known API doc
// that could guarantee fields exist in the response.
// easyjson:json
type responseError struct {
	Code    int            `json:"code,omitempty"`
	Message string         `json:"message,omitempty"`
	Status  string         `json:"status,omitempty"`
	Details []errorDetails `json:"details,omitempty"`
}

// Details could be different type of structs,
// we are interested in one that has errorCode.
// easyjson:json
type errorDetails struct {
	ErrorCode       errorCode        `json:"errorCode,omitempty"`
	FieldViolations []fieldViolation `json:"field_violations"`
}

type fieldViolation struct {
	Field string `json:"field"`
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
	// Request parameters were invalid for HTTP error code = 400
	// This means that one of params was invalid, in case payload is correct this means token is invalid and a new one must be used
	errorCodeInvalidArgument errorCode = "INVALID_ARGUMENT"

	// field name that will be passed as field in fieldViolation in case problems with token on INVALID_ARGUMENT error code
	errorFieldNameToken = "message.token"
)
