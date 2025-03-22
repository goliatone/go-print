package print

import (
	"fmt"
	"reflect"
	"strings"
	"time"
)

func safeToJSON(v any) any {
	if v == nil {
		return nil
	}

	val := reflect.ValueOf(v)

	// handle pointers
	if val.Kind() == reflect.Ptr {
		if val.IsNil() {
			return nil
		}
		return safeToJSON(val.Elem().Interface())
	}

	switch val.Kind() {
	case reflect.Bool, reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Int8, reflect.Uint, reflect.Uint16, reflect.Uint64, reflect.Uint32,
		reflect.Uint8, reflect.Float32, reflect.Float64, reflect.String:
		return v
	case reflect.Array, reflect.Slice:
		if val.IsNil() {
			return nil
		}
		result := make([]any, val.Len())
		for i := 0; i < val.Len(); i++ {
			result[i] = safeToJSON(val.Index(i).Interface())
		}
		return result
	case reflect.Map:
		if val.IsNil() {
			return nil
		}
		result := make(map[string]any)
		for _, key := range val.MapKeys() {
			keyStr := fmt.Sprintf("%v", key.Interface())
			result[keyStr] = safeToJSON(val.MapIndex(key).Interface())
		}
		return result
	case reflect.Struct:
		if t, ok := v.(time.Time); ok {
			return t.Format(time.RFC3339Nano)
		}

		result := make(map[string]any)
		t := val.Type()

		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)

			// skip unexported fields
			if field.PkgPath != "" {
				continue
			}
			fieldName := field.Name
			jsonTag := field.Tag.Get("json")
			if jsonTag != "" {
				parts := strings.Split(jsonTag, ",")
				if parts[0] != "-" {
					if parts[0] != "" {
						fieldName = parts[0]
					}
				} else {
					// skip json:"-" fields
					continue
				}
			}
			fieldValue := val.Field(i).Interface()
			safeValue := safeToJSON(fieldValue)

			shouldOmit := false
			if jsonTag != "" && strings.Contains(jsonTag, "omitempty") {
				isEmpty := false
				if safeValue == nil {
					isEmpty = true
				} else {
					switch reflect.ValueOf(safeValue).Kind() {
					case reflect.Bool:
						isEmpty = !reflect.ValueOf(safeValue).Bool()
					case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
						reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
						isEmpty = reflect.ValueOf(safeValue).Int() == 0
					case reflect.Float32, reflect.Float64:
						isEmpty = reflect.ValueOf(safeValue).Float() == 0
					case reflect.String:
						isEmpty = reflect.ValueOf(safeValue).String() == ""
					case reflect.Map, reflect.Slice, reflect.Array:
						isEmpty = reflect.ValueOf(safeValue).Len() == 0
					}
				}

				shouldOmit = isEmpty
			}

			if !shouldOmit {
				result[fieldName] = safeValue
			}
		}
		return result
	default:
		return unsupportedMessage
	}

}
