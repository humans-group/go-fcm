package fcm

// This error can be caused by missing registration tokens, unregistered or expired tokens.
const (
	ErrUnregistered Error = "Unregistered"
	ErrInvalidToken Error = "InvalidToken"
)

type Error string

func (e Error) Error() string {
	return string(e)
}
