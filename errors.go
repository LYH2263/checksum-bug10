package checksum

import "errors"

var (
	ErrClosed   = errors.New("checksum: pipeline closed")
	ErrInvalid  = errors.New("checksum: invalid argument")
	ErrMismatch = errors.New("checksum: digest mismatch")
	ErrPending  = errors.New("checksum: pending chunks not flushed")
	ErrEmpty    = errors.New("checksum: empty input")
)
