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

func RegisterInstance[T any](lifetime Lifetime, val *T, optFuncs ...ResolutionOptionFunc) error {
	ensureGlobalRegistryInitialized()

	return RegisterInstanceToRegistry[T](globalRegistry, lifetime, val, optFuncs...)
}

func RegisterInstanceToRegistry[T any](reg *Registry, lifetime Lifetime, val *T, optFns ...ResolutionOptionFunc) error {
	if lifetime == Transient {
		return ErrTransientInstance
	}

	opt := newResolutionOptions(optFns...)

	if lifetime == Scoped && opt.scope == "" {
		return ErrEmptyScope
	}

	if lifetime == ThreadLocal && opt.threadID == "" {
		opt.threadID = internal.GetGoroutineID()
	}

	_, name, err := ensureRegistrable[T](reg, lifetime, opt)
	if err != nil {
		return err
	}

	reg.set(name, lifetime, reflect.ValueOf(val), opt)

	return nil
}

func RegisterSingletonInstance[T any](val *T) error {
	return RegisterInstance(Singleton, val)
}

func RegisterSingletonInstanceToRegistry[T any](registry *Registry, val *T) error {
	return RegisterInstanceToRegistry(registry, Singleton, val)
}

func RegisterScopedInstance[T any](val *T, optFuncs ...ResolutionOptionFunc) error {
	return RegisterInstance(Scoped, val, optFuncs...)
}

func RegisterScopedInstanceToRegistry[T any](registry *Registry, val *T, optFuncs ...ResolutionOptionFunc) error {
	return RegisterInstanceToRegistry(registry, Scoped, val, optFuncs...)
}

func RegisterThreadLocalInstance[T any](val *T, optFuncs ...ResolutionOptionFunc) error {
	return RegisterInstance(ThreadLocal, val, optFuncs...)
}

func RegisterThreadLocalInstanceToRegistry[T any](registry *Registry, val *T, optFuncs ...ResolutionOptionFunc) error {
	return RegisterInstanceToRegistry(registry, ThreadLocal, val, optFuncs...)
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
