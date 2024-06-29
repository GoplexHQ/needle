package needle

import (
	"fmt"
	"reflect"
	"unsafe"

	"github.com/goplexhq/needle/internal"
)

// InjectStructFields injects dependencies into the fields of a struct using the global store.
// Returns an error if the injection fails.
//
// Example:
//
//	type MyStruct struct {
//	    Dep *MyDependency
//	}
//
//	var myStruct MyStruct
//	err := needle.InjectStructFields(&myStruct)
//	if err != nil {
//	    ...
//	}
func InjectStructFields[Dest any](dest *Dest) error {
	ensureGlobalStoreInitialized()

	return InjectStructFieldsFromStore[Dest](globalStore, dest)
}

// InjectStructFieldsFromStore injects dependencies into the fields of a struct using the specified store.
// Returns an error if the injection fails.
//
// Example:
//
//	store := needle.NewStore()
//
//	type MyStruct struct {
//	    Dep *MyDependency
//	}
//
//	var myStruct MyStruct
//	err := needle.InjectStructFieldsFromStore(store, &myStruct)
//	if err != nil {
//	    ...
//	}
func InjectStructFieldsFromStore[Dest any](store *Store, dest *Dest) error {
	targetType := reflect.TypeFor[Dest]()

	targetName := internal.ServiceName(targetType)

	if !internal.IsStructType(targetType) {
		return fmt.Errorf("%w: %s", ErrInvalidType, targetName)
	}

	targetValue := reflect.ValueOf(dest).Elem()
	if targetValue.Kind() == reflect.Ptr && targetValue.IsNil() {
		targetValue.Set(reflect.New(targetValue.Type().Elem()))
	}

	for i := range targetType.NumField() {
		f := targetValue.Field(i)
		ft := f.Type().Elem()

		if internal.IsStructType(ft) {
			s, err := resolveName(store, internal.ServiceName(ft))
			if err != nil {
				continue
			}
			f = reflect.NewAt(f.Type(), unsafe.Pointer(f.UnsafeAddr())).Elem()
			f.Set(reflect.ValueOf(s))
		}
	}

	return nil
}
