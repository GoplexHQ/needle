package needle

import (
	"fmt"
	"reflect"

	"github.com/goplexhq/needle/internal"
)

// Register registers a type with a specified lifetime to the global registry.
// Returns an error if the type is already registered or invalid.
//
// The optFuncs parameter allows for optional configuration of the registration, such as setting a scope or thread ID.
//
// Available options:
//   - WithScope(scope string): Sets a scope for scoped instances. Required when lifetime is Scoped.
//   - WithThreadID(threadID string): Sets a thread ID for thread-local instances. Optional and defaults
//     to the current goroutine ID if not provided and the lifetime is ThreadLocal.
//
// Example:
//
//	err := needle.Register[MyService](needle.Singleton)
//	if err != nil {
//	    ...
//	}
//
// Example with scope:
//
//	err := needle.Register[MyService](needle.Scoped, needle.WithScope("request1"))
//	if err != nil {
//	    ...
//	}
//
// Example with thread ID:
//
//	err := needle.Register[MyService](needle.ThreadLocal, needle.WithThreadID("thread1"))
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
// The optFuncs parameter allows for optional configuration of the registration, such as setting a scope or thread ID.
//
// Available options:
//   - WithScope(scope string): Sets a scope for scoped instances. Required when lifetime is Scoped.
//   - WithThreadID(threadID string): Sets a thread ID for thread-local instances. Optional and defaults
//     to the current goroutine ID if not provided and the lifetime is ThreadLocal.
//
// Example:
//
//	registry := needle.NewRegistry()
//	err := needle.RegisterToRegistry[MyService](registry, needle.Singleton)
//	if err != nil {
//	    ...
//	}
//
// Example with scope:
//
//	registry := needle.NewRegistry()
//	err := needle.RegisterToRegistry[MyService](registry, needle.Scoped, needle.WithScope("request1"))
//	if err != nil {
//	    ...
//	}
//
// Example with thread ID:
//
//	registry := needle.NewRegistry()
//	err := needle.RegisterToRegistry[MyService](registry, needle.ThreadLocal, needle.WithThreadID("thread1"))
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

// RegisterInstance registers a pre-initialized instance with a specified lifetime to the global registry.
// Returns an error if the instance is invalid or not supported by the given lifetime.
//
// The optFuncs parameter allows for optional configuration of the registration, such as setting a scope or thread ID.
//
// Available options:
//   - WithScope(scope string): Sets a scope for scoped instances. Required when lifetime is Scoped.
//   - WithThreadID(threadID string): Sets a thread ID for thread-local instances. Optional and defaults
//     to the current goroutine ID if not provided and the lifetime is ThreadLocal.
//
// Example:
//
//	myServiceInstance := &MyService{}
//	err := needle.RegisterInstance(needle.Singleton, myServiceInstance)
//	if err != nil {
//	    ...
//	}
//
// Example with scope:
//
//	myServiceInstance := &MyService{}
//	err := needle.RegisterInstance(needle.Scoped, myServiceInstance, needle.WithScope("request1"))
//	if err != nil {
//	    ...
//	}
//
// Example with thread ID:
//
//	myServiceInstance := &MyService{}
//	err := needle.RegisterInstance(needle.ThreadLocal, myServiceInstance, needle.WithThreadID("thread1"))
//	if err != nil {
//	    ...
//	}
func RegisterInstance[T any](lifetime Lifetime, val *T, optFuncs ...ResolutionOptionFunc) error {
	ensureGlobalRegistryInitialized()

	return RegisterInstanceToRegistry[T](globalRegistry, lifetime, val, optFuncs...)
}

// RegisterInstanceToRegistry registers a pre-initialized instance with a specified lifetime to the registry.
// Returns an error if the instance is invalid or not supported by the given lifetime.
//
// The optFuncs parameter allows for optional configuration of the registration, such as setting a scope or thread ID.
//
// Available options:
//   - WithScope(scope string): Sets a scope for scoped instances. Required when lifetime is Scoped.
//   - WithThreadID(threadID string): Sets a thread ID for thread-local instances. Optional and defaults
//     to the current goroutine ID if not provided and the lifetime is ThreadLocal.
//
// Example:
//
//	registry := needle.NewRegistry()
//	myServiceInstance := &MyService{}
//	err := needle.RegisterInstanceToRegistry(registry, needle.Singleton, myServiceInstance)
//	if err != nil {
//	    ...
//	}
//
// Example with scope:
//
//	registry := needle.NewRegistry()
//	myServiceInstance := &MyService{}
//	err := needle.RegisterInstanceToRegistry(registry, needle.Scoped, myServiceInstance, needle.WithScope("request1"))
//	if err != nil {
//	    ...
//	}
//
// Example with thread ID:
//
//	registry := needle.NewRegistry()
//	myService := &MyService{}
//	err := needle.RegisterInstanceToRegistry(registry, needle.ThreadLocal, myService, needle.WithThreadID("thread1"))
//	if err != nil {
//	    ...
//	}
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

// RegisterSingletonInstance registers pre-initialized singleton instance to the global registry.
//
// Example:
//
//	myServiceInstance := &MyService{}
//	err := needle.RegisterSingletonInstance(myServiceInstance)
//	if err != nil {
//	    ...
//	}
func RegisterSingletonInstance[T any](val *T) error {
	return RegisterInstance(Singleton, val)
}

// RegisterSingletonInstanceToRegistry registers pre-initialized singleton instance to the specified registry.
//
// Example:
//
//	registry := needle.NewRegistry()
//	myServiceInstance := &MyService{}
//	err := needle.RegisterSingletonInstanceToRegistry(registry, myServiceInstance)
//	if err != nil {
//	    ...
//	}
func RegisterSingletonInstanceToRegistry[T any](registry *Registry, val *T) error {
	return RegisterInstanceToRegistry(registry, Singleton, val)
}

// RegisterScopedInstance registers a pre-initialized scoped instance to the global registry.
//
// The optFuncs parameter allows for optional configuration of the registration, such as setting a scope.
//
// Available options:
// - WithScope(scope string): Sets a scope for scoped instances. Required.
//
// Example:
//
//	myServiceInstance := &MyService{}
//	err := needle.RegisterScopedInstance(myServiceInstance, needle.WithScope("request1"))
//	if err != nil {
//	    ...
//	}
func RegisterScopedInstance[T any](val *T, optFuncs ...ResolutionOptionFunc) error {
	return RegisterInstance(Scoped, val, optFuncs...)
}

// RegisterScopedInstanceToRegistry registers a pre-initialized scoped instance to the specified registry.
//
// The optFuncs parameter allows for optional configuration of the registration, such as setting a scope.
//
// Available options:
// - WithScope(scope string): Sets a scope for scoped instances. Required.
//
// Example:
//
//	registry := needle.NewRegistry()
//	myServiceInstance := &MyService{}
//	err := needle.RegisterScopedInstanceToRegistry(registry, myServiceInstance, needle.WithScope("request1"))
//	if err != nil {
//	    ...
//	}
func RegisterScopedInstanceToRegistry[T any](registry *Registry, val *T, optFuncs ...ResolutionOptionFunc) error {
	return RegisterInstanceToRegistry(registry, Scoped, val, optFuncs...)
}

// RegisterThreadLocalInstance registers a pre-initialized thread-local instance to the global registry.
//
// The optFuncs parameter allows for optional configuration of the registration, such as setting a thread ID.
//
// Available options:
//   - WithThreadID(threadID string): Sets a thread ID for thread-local instances. Optional and defaults
//     to the current goroutine ID if not provided and the lifetime is ThreadLocal.
//
// Example:
//
//	myServiceInstance := &MyService{}
//	err := needle.RegisterThreadLocalInstance(myServiceInstance, needle.WithThreadID("thread1"))
//	if err != nil {
//	    ...
//	}
func RegisterThreadLocalInstance[T any](val *T, optFuncs ...ResolutionOptionFunc) error {
	return RegisterInstance(ThreadLocal, val, optFuncs...)
}

// RegisterThreadLocalInstanceToRegistry registers a pre-initialized thread-local instance to the specified registry.
//
// The optFuncs parameter allows for optional configuration of the registration, such as setting a thread ID.
//
// Available options:
//   - WithThreadID(threadID string): Sets a thread ID for thread-local instances. Optional and defaults
//     to the current goroutine ID if not provided and the lifetime is ThreadLocal.
//
// Example:
//
//	registry := needle.NewRegistry()
//	myServiceInstance := &MyService{}
//	err := needle.RegisterThreadLocalInstanceToRegistry(registry, myServiceInstance, needle.WithThreadID("thread1"))
//	if err != nil {
//	    ...
//	}
func RegisterThreadLocalInstanceToRegistry[T any](registry *Registry, val *T, optFuncs ...ResolutionOptionFunc) error {
	return RegisterInstanceToRegistry(registry, ThreadLocal, val, optFuncs...)
}

// ensureRegistrable checks if a type is registrable and not already registered in the registry.
// Returns the type, name, and an error if the type is not registrable or already registered.
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
