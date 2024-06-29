package needle

import "errors"

var (
	ErrRegistered          = errors.New("service already registered in the store")
	ErrNotRegistered       = errors.New("service not registered in the store")
	ErrInvalidServiceType  = errors.New("invalid service type: service must be a struct type")
	ErrInvalidDestType     = errors.New("invalid destination type: expected a struct type")
	ErrServiceTypeMismatch = errors.New("resolved service type does not match the expected type")
	ErrInvalidLifetime     = errors.New("invalid lifetime value")
	ErrFieldPtr            = errors.New("injectable field is not a pointer")
	ErrResolveField        = errors.New("unable to resolve service for field")
)
