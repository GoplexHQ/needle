package needle

import (
	"fmt"
	"reflect"

	"github.com/goplexhq/needle/internal"
)

// Register registers a type with a specified lifetime to the global registry.
// Returns an error if the type is already registered or invalid.
//
// Example:
//
//	err := needle.Register[MyService](needle.Singleton)
//	if err != nil {
//	    ...
//	}
func Register[T any](lifetime Lifetime, optFuncs ...ResolutionOptionFunc) error {
	ensureGlobalRegistryInitialized()

	return RegisterToRegistry[T](globalRegistry, lifetime, optFuncs...)
}

// RegisterToRegistry registers a type with a specified lifetime to the registry.
// Returns an error if the type is already registered or invalid.
//
// Example:
//
//	registry := needle.NewRegistry()
//	err := needle.RegisterToRegistry[MyService](registry, needle.Singleton)
//	if err != nil {
//	    ...
//	}
func RegisterToRegistry[T any](registry *Registry, lifetime Lifetime, optFuncs ...ResolutionOptionFunc) error {
	opt := newResolutionOptions(optFuncs...)

	if lifetime == Scoped && opt.scope == "" {
		return ErrEmptyScope
	}

	if lifetime == ThreadLocal && opt.threadID == "" {
		opt.threadID = internal.GetGoroutineID()
	}

	typ, name, err := ensureRegistrable[T](registry, lifetime, opt)
	if err != nil {
		return err
	}

	var value reflect.Value

	if lifetime == Transient {
		value = reflect.Zero(typ)
	} else {
		value = reflect.New(typ)
	}

	registry.set(name, lifetime, value, opt)

	return nil
}

// RegisterInstance registers a pre-initialized singleton instance to the global registry.
// Returns an error if the type is already registered or invalid.
//
// Example:
//
//	err := needle.RegisterInstance(&MyService{})
//	if err != nil {
//	    ...
//	}
func RegisterInstance[T any](val *T, optFuncs ...ResolutionOptionFunc) error {
	ensureGlobalRegistryInitialized()

	return RegisterInstanceToRegistry[T](globalRegistry, val, optFuncs...)
}

// RegisterInstanceToRegistry registers a pre-initialized singleton instance to the registry.
// Returns an error if the type is already registered or invalid.
//
// Example:
//
//	registry := needle.NewRegistry()
//	err := needle.RegisterInstanceToRegistry(registry, &MyService{})
//	if err != nil {
//	    ...
//	}
func RegisterInstanceToRegistry[T any](registry *Registry, val *T, optFuncs ...ResolutionOptionFunc) error {
	opt := newResolutionOptions(optFuncs...)

	_, name, err := ensureRegistrable[T](registry, Singleton, opt)
	if err != nil {
		return err
	}

	registry.set(name, Singleton, reflect.ValueOf(val), opt)

	return nil
}

// ensureRegistrable checks if a type is registrable and not already registered in the registry.
func ensureRegistrable[T any](reg *Registry, lifetime Lifetime, opt *ResolutionOptions) (reflect.Type, string, error) {
	typ := reflect.TypeFor[T]()
	name := internal.ServiceName(typ)

	if !internal.IsStructType(typ) {
		return nil, "", fmt.Errorf("%w: %s", ErrInvalidServiceType, name)
	}

	entry, exists := reg.has(name)
	if exists && entry.lifetime == lifetime && (lifetime == Transient ||
		lifetime == Singleton ||
		(lifetime == Scoped && reg.hasScoped(opt.scope, name)) ||
		(lifetime == ThreadLocal && reg.hasThreadLocal(opt.threadID, name))) {
		return nil, "", fmt.Errorf("%w: %s", ErrRegistered, name)
	}

	return typ, name, nil
}
