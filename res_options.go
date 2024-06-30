package needle

// ResolutionOptions holds configuration options for resolving services.
type ResolutionOptions struct {
	scope    string
	threadID string
}

// ResolutionOptionFunc is a function that modifies a ResolutionOptions struct.
type ResolutionOptionFunc func(*ResolutionOptions)

// WithScope sets the scope for a ResolutionOptions struct.
//
// Example:
//
//	opt := needle.WithScope("scope1")
func WithScope(scope string) ResolutionOptionFunc {
	return func(o *ResolutionOptions) {
		o.scope = scope
	}
}

// WithThreadID sets the thread ID for a ResolutionOptions struct.
//
// Example:
//
//	opt := needle.WithThreadID("thread1")
func WithThreadID(id string) ResolutionOptionFunc {
	return func(o *ResolutionOptions) {
		o.threadID = id
	}
}

// newResolutionOptions creates a new ResolutionOptions struct from the provided option functions.
//
// Example:
//
//	opt := newResolutionOptions(needle.WithScope("scope1"), needle.WithThreadID("thread1"))
func newResolutionOptions(optFuncs ...ResolutionOptionFunc) *ResolutionOptions {
	opt := &ResolutionOptions{} //nolint:exhaustruct
	for _, optFunc := range optFuncs {
		optFunc(opt)
	}

	return opt
}
