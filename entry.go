package needle

import "reflect"

type entry struct {
	name     string
	lifetime Lifetime
	value    reflect.Value
}
