package fcm

type ErrorCode string

const (
	// Unknown error
	errCodeUnspecified ErrorCode = "UNSPECIFIED_ERROR"

	// App instance was unregistered from FCM for HTTP error code = 404
	// This usually means that the token used is no longer valid (i.e. expired) and a new one must be used.
	errCodeUnregistered ErrorCode = "UNREGISTERED"
)
