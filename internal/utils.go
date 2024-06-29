package internal

import "reflect"

func ServiceName(rt reflect.Type) string {
	if rt == nil || (rt.PkgPath() == "" && rt.Name() == "") {
		return ""
	}

	return rt.PkgPath() + "." + rt.Name()
}

func IsStructType(rt reflect.Type) bool {
	return rt != nil && rt.Kind() == reflect.Struct
}
