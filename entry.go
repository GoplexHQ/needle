package needle

import "reflect"

// serviceEntry holds metadata about a registered service.
type serviceEntry struct {
	name     string
	lifetime Lifetime
	value    *reflect.Value
}

// withValue sets the value of a serviceEntry and returns the updated entry.
//
// Example:
//
//	entry := serviceEntry{name: "MyService", lifetime: needle.Singleton}
//	entry = entry.withValue(reflect.ValueOf(&MyService{}))
func (e *serviceEntry) withValue(value *reflect.Value) serviceEntry {
	e.value = value

	return *e
}
