package needle

import (
	"fmt"
	"reflect"

	"github.com/goplexhq/needle/internal"
)

// Resolve resolves an instance of the specified type from the global registry.
// Returns a pointer to the resolved instance or an error if the instance cannot be resolved.
//
// Example:
//
//	val, err := needle.Resolve[MyService]()
//	if err != nil {
//	    ...
//	}
func Resolve[T any](optFuncs ...ResolutionOptionFunc) (*T, error) {
	ensureGlobalRegistryInitialized()

	return ResolveFromRegistry[T](globalRegistry, optFuncs...)
}

// ResolveFromRegistry resolves an instance of the specified type from the given registry.
// Returns a pointer to the resolved instance or an error if the instance cannot be resolved.
//
// Example:
//
//	registry := needle.NewRegistry()
//	val, err := needle.ResolveFromRegistry[MyService](registry)
//	if err != nil {
//	    ...
//	}
func ResolveFromRegistry[T any](registry *Registry, optFuncs ...ResolutionOptionFunc) (*T, error) {
	t := reflect.TypeFor[T]()
	name := internal.ServiceName(t)

	opt := newResolutionOptions(optFuncs...)

	entry, exists := registry.has(name)
	if !exists {
		return nil, fmt.Errorf("%w: %s", ErrNotRegistered, name)
	}

	if entry.lifetime == Scoped && opt.scope == "" {
		return nil, ErrEmptyScope
	}

	if entry.lifetime == ThreadLocal && opt.threadID == "" {
		opt.threadID = internal.GetGoroutineID()
	}

	i, err := resolveName(registry, name, opt)
	if err != nil {
		return nil, err
	}

	v, valid := i.(*T)
	if !valid {
		return nil, fmt.Errorf("%w: %s", ErrServiceTypeMismatch, name)
	}

	return v, nil
}

// resolveName resolves the instance by its name from the registry.
func resolveName(registry *Registry, name string, opt *ResolutionOptions) (any, error) {
	entry, exists := registry.get(name, opt)
	if !exists {
		return nil, fmt.Errorf("%w: %s", ErrNotRegistered, name)
	}

	if entry.lifetime == Transient {
		return reflect.New(entry.value.Type()).Interface(), nil
	}

	return entry.value.Interface(), nil
}
