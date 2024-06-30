package needle

import "errors"

var (
	ErrRegistered          = errors.New("service already registered in the registry")
	ErrNotRegistered       = errors.New("service not registered in the registry")
	ErrInvalidServiceType  = errors.New("invalid service type: service must be a struct type")
	ErrInvalidDestType     = errors.New("invalid destination type: expected a struct type")
	ErrServiceTypeMismatch = errors.New("resolved service type does not match the expected type")
	ErrFieldPtr            = errors.New("injectable field is not a pointer")
	ErrResolveField        = errors.New("unable to resolve service for field")
	ErrEmptyScope          = errors.New("scope is required but not provided")
	ErrTransientInstance   = errors.New("transient lifetime does not support pre-initialized instances")
)
