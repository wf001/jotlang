package util

import "reflect"

// NOTE: there might be the better way
func EqualType(v interface{}, ty interface{}) bool {
	return reflect.TypeOf(v) == reflect.TypeOf(ty)
}
