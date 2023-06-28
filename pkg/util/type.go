package util

import "reflect"

func IsSameType(src, dest interface{}) bool {
	return reflect.TypeOf(src) == reflect.TypeOf(dest)
}
