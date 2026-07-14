package goroute

import (
	"errors"
	"reflect"
	"strings"
)

// ValidateStruct inspect a struct in runtime for `validate:"requared"` tags
func ValidateStruct(obj interface{}) error {
	// 1. extract the type of unknown object and value in runtime
	v := reflect.ValueOf(obj)

	// if the user passed a pointer to a struct, we need to dereference it to get the underlying struct value
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	// 3. safety check to ensure the value is indeed a struct
	if v.Kind() != reflect.Struct {
		return errors.New("BindJSON requires a struct pointers")
	}

	// 4. itarate over the fields of the struct
	for i := 0; i < v.NumField(); i++ {
		field := v.Type().Field(i)
		tag := field.Tag.Get("validate")

		if strings.Contains(tag, "required") {
			if v.Field(i).IsZero() {
				return errors.New("Missing required field: " + field.Name)
			}
		}
	}
	return nil
}
