package needle

import (
	"fmt"
	"reflect"
	"unsafe"

	"github.com/goplexhq/needle/internal"
)

const (
	injectTagKey   = "needle"
	injectTagValue = "inject"
)

// InjectStructFields injects dependencies into the fields of a struct using the global store.
// Returns an error if the injection fails.
//
// Example:
//
//	type MyStruct struct {
//	    Dep *MyDependency `needle:"inject"`
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
//	    Dep *MyDependency `needle:"inject"`
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
		return fmt.Errorf("type %s is not a struct: %w", targetName, ErrInvalidType)
	}

	targetValue := reflect.ValueOf(dest).Elem()
	initializePointerValue(&targetValue)

	for i := range targetType.NumField() {
		fieldType := targetType.Field(i)
		if fieldType.Tag.Get(injectTagKey) != injectTagValue {
			continue
		}

		fieldValue := targetValue.Field(i)
		if err := injectField(store, fieldType, fieldValue); err != nil {
			return err
		}
	}

	return nil
}

// initializePointerValue ensures the pointer value is not nil by initializing it.
func initializePointerValue(value *reflect.Value) {
	if internal.IsPointerValue(*value) && value.IsNil() {
		value.Set(reflect.New(value.Type().Elem()))
	}
}

// injectField injects a dependency into a single struct field.
func injectField(store *Store, fieldType reflect.StructField, fieldValue reflect.Value) error {
	if !internal.IsPointerValue(fieldValue) {
		return fmt.Errorf("%w: %s", ErrFieldPtr, fieldType.Name)
	}

	elemType := fieldValue.Type().Elem()
	if internal.IsStructType(elemType) {
		service, err := resolveName(store, internal.ServiceName(elemType))
		if err != nil {
			return fmt.Errorf("unable to resolve service for field %q: %w", fieldType.Name, err)
		}

		fieldValue = reflect.NewAt(fieldValue.Type(), unsafe.Pointer(fieldValue.UnsafeAddr())).Elem()
		fieldValue.Set(reflect.ValueOf(service))
	}

	return nil
}
