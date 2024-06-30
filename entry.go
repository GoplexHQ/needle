package needle

import "reflect"

type serviceEntry struct {
	name     string
	lifetime Lifetime
	value    *reflect.Value
}

func (e *serviceEntry) withValue(value *reflect.Value) serviceEntry {
	e.value = value

	return *e
}
