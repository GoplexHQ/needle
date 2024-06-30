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

// InjectStructFields injects dependencies into the fields of a struct using the global registry.
// Returns an error if the injection fails.
//
// The optFuncs parameter allows for optional configuration of the injection, such as setting a scope or thread ID.
//
// Available options:
// - WithScope(scope string): Sets a scope for resolving scoped dependencies.
// - WithThreadID(threadID string): Sets a thread ID for resolving thread-local dependencies.
//
// Fields must be annotated with `needle:"inject"` for the needle framework to inject dependencies into them.
// Additionally, the field must be a pointer to a struct.
//
// Example:
//
//	type MyDependency struct {}
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
func InjectStructFields[Dest any](dest *Dest, optFuncs ...ResolutionOptionFunc) error {
	ensureGlobalRegistryInitialized()

	return InjectStructFieldsFromRegistry[Dest](globalRegistry, dest, optFuncs...)
}

// InjectStructFieldsFromRegistry injects dependencies into the fields of a struct using the specified registry.
// Returns an error if the injection fails.
//
// The optFuncs parameter allows for optional configuration of the injection, such as setting a scope or thread ID.
//
// Available options:
// - WithScope(scope string): Sets a scope for resolving scoped dependencies.
// - WithThreadID(threadID string): Sets a thread ID for resolving thread-local dependencies.
//
// Fields must be annotated with `needle:"inject"` for the needle framework to inject dependencies into them.
// Additionally, the field must be a pointer to a struct.
//
// Example:
//
//	registry := needle.NewRegistry()
//
//	type MyDependency struct {}
//
//	type MyStruct struct {
//	    Dep *MyDependency `needle:"inject"`
//	}
//
//	var myStruct MyStruct
//	err := needle.InjectStructFieldsFromRegistry(registry, &myStruct)
//	if err != nil {
//	    ...
//	}
func InjectStructFieldsFromRegistry[Dest any](registry *Registry, dest *Dest, optFuncs ...ResolutionOptionFunc) error {
	targetType := reflect.TypeFor[Dest]()
	targetName := internal.ServiceName(targetType)

	if !internal.IsStructType(targetType) {
		return fmt.Errorf("%w: %s", ErrInvalidDestType, targetName)
	}

	targetValue := reflect.ValueOf(dest).Elem()
	initializePointerValue(&targetValue)

	opt := newResolutionOptions(optFuncs...)

	for idx := range targetType.NumField() {
		fieldType := targetType.Field(idx)
		if fieldType.Tag.Get(injectTagKey) != injectTagValue {
			continue
		}

		fieldValue := targetValue.Field(idx)
		if err := injectField(registry, fieldType, fieldValue, opt); err != nil {
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
func injectField(registry *Registry, field reflect.StructField, value reflect.Value, opt *ResolutionOptions) error {
	if !internal.IsPointerValue(value) {
		return fmt.Errorf("%w: %s", ErrFieldPtr, field.Name)
	}

	elem := value.Type().Elem()
	if internal.IsStructType(elem) {
		name := internal.ServiceName(elem)

		entry, exists := registry.has(name)
		if !exists {
			return fmt.Errorf("%w: %s", ErrNotRegistered, name)
		}

		if entry.lifetime == Scoped && opt.scope == "" {
			return ErrEmptyScope
		}

		if entry.lifetime == ThreadLocal && opt.threadID == "" {
			opt.threadID = internal.GetGoroutineID()
		}

		entryValue, err := resolveName(registry, name, opt)
		if err != nil {
			return fmt.Errorf("%w %q: %w", ErrResolveField, field.Name, err)
		}

		value = reflect.NewAt(value.Type(), unsafe.Pointer(value.UnsafeAddr())).Elem()
		value.Set(reflect.ValueOf(entryValue))
	}

	return nil
}
