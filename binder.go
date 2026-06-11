package echo

import (
	"net/url"
	"reflect"
	"strconv"
)

// DefaultBinder is the default implementation of the Binder interface.
type DefaultBinder struct{}

// BindQueryParams binds query parameters to the destination struct.
func (b *DefaultBinder) BindQueryParams(c Context, i interface{}) error {
	params := c.QueryParams()
	return bindStruct(i, params)
}

func bindStruct(dst interface{}, params url.Values) error {
	val := reflect.ValueOf(dst).Elem()
	return bindRecursive(val, params)
}

func bindRecursive(val reflect.Value, params url.Values) error {
	typ := val.Type()
	for i := 0; i < val.NumField(); i++ {
		fieldVal := val.Field(i)
		fieldType := typ.Field(i)

		if fieldType.Anonymous {
			if fieldVal.Kind() == reflect.Ptr {
				if fieldVal.IsNil() {
					fieldVal.Set(reflect.New(fieldVal.Type().Elem()))
				}
				fieldVal = fieldVal.Elem()
			}
			if fieldVal.Kind() == reflect.Struct {
				if err := bindRecursive(fieldVal, params); err != nil {
					return err
				}
				continue
			}
		}

		tag := fieldType.Tag.Get("query")
		if tag == "" || tag == "-" {
			continue
		}

		valStr := params.Get(tag)
		if valStr == "" {
			continue
		}

		switch fieldVal.Kind() {
		case reflect.String:
			fieldVal.SetString(valStr)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			intVal, err := strconv.ParseInt(valStr, 10, 64)
			if err == nil {
				fieldVal.SetInt(intVal)
			}
		// Add other types as needed
		}
	}
	return nil
}