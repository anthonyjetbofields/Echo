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
	_, err := bindRecursive(val, params)
	return err
}

func bindRecursive(val reflect.Value, params url.Values) (bool, error) {
	typ := val.Type()
	bound := false

	for i := 0; i < val.NumField(); i++ {
		fieldVal := val.Field(i)
		fieldType := typ.Field(i)

		if fieldType.Anonymous {
			isPtr := fieldVal.Kind() == reflect.Ptr
			var structVal reflect.Value

			if isPtr {
				if fieldVal.IsNil() {
					structVal = reflect.New(fieldVal.Type().Elem()).Elem()
				} else {
					structVal = fieldVal.Elem()
				}
			} else {
				structVal = fieldVal
			}

			if structVal.Kind() == reflect.Struct {
				bnd, err := bindRecursive(structVal, params)
				if err != nil {
					return false, err
				}
				if bnd {
					bound = true
					if isPtr && fieldVal.IsNil() {
						fieldVal.Set(structVal.Addr())
					}
				}
				continue
			}
		}

		tag := fieldType.Tag.Get("query")
		if tag == "" || tag == "-" {
			continue
		}

		if _, ok := params[tag]; !ok {
			continue
		}

		valStr := params.Get(tag)

		bound = true
		switch fieldVal.Kind() {
		case reflect.String:
			fieldVal.SetString(valStr)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			if valStr != "" {
				intVal, err := strconv.ParseInt(valStr, 10, 64)
				if err == nil {
					fieldVal.SetInt(intVal)
				}
			}
		// Add other types as needed
		}
	}
	return bound, nil
}