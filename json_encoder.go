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
	case reflect.Array:
		// arrays are never nil
		arrResult := make([]any, val.Len())
		for i := range val.Len() {
			arrResult[i] = safeToJSON(val.Index(i).Interface())
		}
		return arrResult

	case reflect.Slice:
		// slices can be nil so IsNil() is valid
		if val.IsNil() {
			return nil
		}
		sliceResult := make([]any, val.Len())
		for i := range val.Len() {
			sliceResult[i] = safeToJSON(val.Index(i).Interface())
		}
		return sliceResult

	case reflect.Map:
		if val.IsNil() {
			return nil
		}
		mapResult := make(map[string]any)
		for _, key := range val.MapKeys() {
			keyStr := fmt.Sprintf("%v", key.Interface())
			mapResult[keyStr] = safeToJSON(val.MapIndex(key).Interface())
		}
		return mapResult

	case reflect.Struct:
		if t, ok := v.(time.Time); ok {
			return t.Format(time.RFC3339Nano)
		}

		strucResult := make(map[string]any)
		t := val.Type()

		for i := range t.NumField() {
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

			// Handle omitempty logic
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
						// reflect.ValueOf(safeValue).Len() is valid if safeValue != nil
						isEmpty = reflect.ValueOf(safeValue).Len() == 0
					}
				}
				shouldOmit = isEmpty
			}

			if !shouldOmit {
				strucResult[fieldName] = safeValue
			}
		}
		return strucResult

	default:
		return unsupportedMessage
	}

}
