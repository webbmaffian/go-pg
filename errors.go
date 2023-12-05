package pg

const (
	ErrNotFound     Error = "row not found"
	ErrNoColumns    Error = "no columns to select"
	ErrNoConnection Error = "missing db connection"
	ErrResultClosed Error = "result is closed"
	ErrNotSupported Error = "not supported yet"
	ErrNoPointer    Error = "dst must be a pointer"
	ErrInvalidDst   Error = "invalid dst"
)
