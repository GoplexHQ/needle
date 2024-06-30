package needle

type ResolutionOptions struct {
	scope    string
	threadID string
}

func newResolutionOptions(optFuncs ...ResolutionOptionFunc) *ResolutionOptions {
	opt := &ResolutionOptions{} //nolint:exhaustruct
	for _, optFunc := range optFuncs {
		optFunc(opt)
	}

	return opt
}

type ResolutionOptionFunc func(*ResolutionOptions)

func WithScope(scope string) ResolutionOptionFunc {
	return func(o *ResolutionOptions) {
		o.scope = scope
	}
}

func WithThreadID(id string) ResolutionOptionFunc {
	return func(o *ResolutionOptions) {
		o.threadID = id
	}
}
