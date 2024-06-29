package needle

import (
	"fmt"
	"reflect"

	"github.com/goplexhq/needle/internal"
)

// Resolve resolves an instance of the specified type from the global store.
// Returns a pointer to the resolved instance or an error if the instance cannot be resolved.
//
// Example:
//
//	val, err := needle.Resolve[MyService]()
//	if err != nil {
//	    ...
//	}
func Resolve[T any]() (*T, error) {
	ensureGlobalStoreInitialized()

	return ResolveFromStore[T](globalStore)
}

// ResolveFromStore resolves an instance of the specified type from the given store.
// Returns a pointer to the resolved instance or an error if the instance cannot be resolved.
//
// Example:
//
//	store := needle.NewStore()
//	val, err := needle.ResolveFromStore[MyService](store)
//	if err != nil {
//	    ...
//	}
func ResolveFromStore[T any](store *Store) (*T, error) {
	t := reflect.TypeFor[T]()
	name := internal.ServiceName(t)

	i, err := resolveName(store, name)
	if err != nil {
		return nil, err
	}

	v, valid := i.(*T)
	if !valid {
		return nil, fmt.Errorf("%w: %s", ErrServiceTypeMismatch, name)
	}

	return v, nil
}

// resolveName resolves the instance by its name from the store.
func resolveName(store *Store, name string) (any, error) {
	service, ok := store.get(name)
	if !ok {
		return nil, fmt.Errorf("%w: %s", ErrNotRegistered, name)
	}

	switch service.lifetime {
	case Transient:
		return reflect.New(service.value.Type()).Interface(), nil
	case Scoped:
		return reflect.New(service.value.Type()).Interface(), nil
	case ThreadLocal:
		return reflect.New(service.value.Type()).Interface(), nil
	case Pooled:
		return reflect.New(service.value.Type()).Interface(), nil
	case Singleton:
		return service.value.Interface(), nil
	default:
		return nil, fmt.Errorf("%w: %s", ErrInvalidLifetime, service.lifetime.String())
	}
}
