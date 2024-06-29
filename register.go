package needle

import (
	"fmt"
	"reflect"

	"github.com/goplexhq/needle/internal"
)

// ensureRegistrable checks if a type is registrable and not already registered in the store.
func ensureRegistrable[T any](store *Store) (reflect.Type, string, error) {
	typ := reflect.TypeFor[T]()
	name := internal.ServiceName(typ)

	if !internal.IsStructType(typ) {
		return nil, "", fmt.Errorf("%w: %s", ErrInvalidType, name)
	}

	if store.has(name) {
		return nil, "", fmt.Errorf("%w: %s", ErrRegistered, name)
	}

	return typ, name, nil
}

// RegisterToStore registers a type with a specified lifetime to the store.
// Returns an error if the type is already registered or invalid.
//
// Example:
//
//	store := needle.NewStore()
//	err := needle.RegisterToStore[MyService](store, needle.Singleton)
//	if err != nil {
//	    ...
//	}
func RegisterToStore[T any](store *Store, lifetime Lifetime) error {
	typ, name, err := ensureRegistrable[T](store)
	if err != nil {
		return err
	}

	newEntry := entry{ //nolint:exhaustruct
		name:     name,
		lifetime: lifetime,
	}

	if lifetime == Singleton {
		newEntry.value = reflect.New(typ)
	} else {
		newEntry.value = reflect.Zero(typ)
	}

	store.set(name, newEntry)

	return nil
}

// Register registers a type with a specified lifetime to the global store.
// Returns an error if the type is already registered or invalid.
//
// Example:
//
//	err := needle.Register[MyService](needle.Singleton)
//	if err != nil {
//	    ...
//	}
func Register[T any](lifetime Lifetime) error {
	ensureGlobalStoreInitialized()

	return RegisterToStore[T](globalStore, lifetime)
}

// RegisterInstanceToStore registers a pre-initialized singleton instance to the store.
// Returns an error if the type is already registered or invalid.
//
// Example:
//
//	store := needle.NewStore()
//	err := needle.RegisterInstanceToStore(store, &MyService{})
//	if err != nil {
//	    ...
//	}
func RegisterInstanceToStore[T any](store *Store, val *T) error {
	_, name, err := ensureRegistrable[T](store)
	if err != nil {
		return err
	}

	store.set(name, entry{
		name:     name,
		lifetime: Singleton,
		value:    reflect.ValueOf(val),
	})

	return nil
}

// RegisterInstance registers a pre-initialized singleton instance to the global store.
// Returns an error if the type is already registered or invalid.
//
// Example:
//
//	err := needle.RegisterInstance(&MyService{})
//	if err != nil {
//	    ...
//	}
func RegisterInstance[T any](val *T) error {
	ensureGlobalStoreInitialized()

	return RegisterInstanceToStore[T](globalStore, val)
}
