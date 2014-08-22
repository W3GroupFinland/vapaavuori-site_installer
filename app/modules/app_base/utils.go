package app_base

import (
	"reflect"
	"runtime"
)

func GetFunctionName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}

func ReflectValueFunctionName(value reflect.Value) string {
	return runtime.FuncForPC(value.Pointer()).Name()
}
