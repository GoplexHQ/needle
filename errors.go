package needle

import "errors"

var (
	ErrRegistered          = errors.New("service already registered")
	ErrNotRegistered       = errors.New("service not registered")
	ErrInvalidType         = errors.New("service type is invalid")
	ErrServiceTypeMismatch = errors.New("resolved service type does not match the given type")
	ErrInvalidLifetime     = errors.New("invalid lifetime")
)
