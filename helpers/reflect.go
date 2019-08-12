package helpers

import (
	"reflect"
)

// ReflectMethod Reflect if an interface is either a struct or a pointer to a struct
func ReflectMethod(obj interface{}, methodName string) *reflect.Value {
	value := reflect.ValueOf(obj)

	// Check if the passed interface is a pointer
	if value.Type().Kind() != reflect.Ptr {
		// Create a new type of Iface, so we have a pointer to work with
		value = reflect.New(reflect.TypeOf(obj))
	}

	// Get the method by name
	Method := value.MethodByName(methodName)
	if !Method.IsValid() {
		return nil
	}
	return &Method
}

func ExecuteFunction(f interface{}, funcName string) interface{} {
	rf := reflect.TypeOf(f)
	if rf.Kind() != reflect.Func {
		return nil
	}
	vf := reflect.ValueOf(f)
	ReflectMethod(f, funcName).Call([]reflect.Value{reflect.ValueOf(1)})
	return vf.Interface()
}
