package utils

import (
	"reflect"
	"strings"
)

// StructToMap converts a struct to map[string]any for use with GORM .Updates().
//
// Key resolution (in order):
//  1. json tag name  → used as map key
//  2. json:"-"       → field is skipped entirely
//  3. no json tag    → field name is used as map key
//
// Value rules:
//   - Pointer field (nil)     → excluded  (field not sent by client = don't touch)
//   - Pointer field (non-nil) → included with dereferenced value
//   - Non-pointer zero value  → excluded  (empty string, 0, false, etc.)
//   - Non-pointer non-zero    → included as-is
//
// Use *T for optional fields: nil = skip, &false = update to false, &true = update to true.
//
// Example:
//
//	type UpdateReq struct {
//	    Name     string  `json:"name"`
//	    IsActive *bool   `json:"is_active"`
//	    LogoURL  *string `json:"logo_url"`
//	}
//
//	req := UpdateReq{Name: "Acme", IsActive: nil, LogoURL: ptr("https://...")}
//	result := StructToMap(&req)
//	// map[string]any{
//	//     "name":     "Acme",
//	//     "logo_url": "https://...",
//	//     // is_active skipped — nil pointer
//	// }
func StructToMap(v any) map[string]any {
	result := make(map[string]any)

	val := reflect.ValueOf(v)
	typ := reflect.TypeOf(v)

	// unwrap pointer-to-struct
	if val.Kind() == reflect.Ptr {
		if val.IsNil() {
			return result
		}
		val = val.Elem()
		typ = typ.Elem()
	}

	if val.Kind() != reflect.Struct {
		return result
	}

	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		fieldVal := val.Field(i)

		if !field.IsExported() {
			continue
		}

		// resolve key from json tag
		key := field.Name
		if tag, ok := field.Tag.Lookup("json"); ok {
			name, _, _ := strings.Cut(tag, ",")
			if name == "-" {
				continue // explicitly excluded
			}
			if name != "" {
				key = name
			}
		}

		// Check if it's a Nullable type (has IsSet and IsNull methods)
		if isNullable := fieldVal.MethodByName("IsSet"); isNullable.IsValid() {
			if isSetResult := isNullable.Call(nil); len(isSetResult) > 0 && isSetResult[0].Bool() {
				// Field was in JSON
				if isNull := fieldVal.MethodByName("IsNull"); isNull.IsValid() {
					if isNullResult := isNull.Call(nil); len(isNullResult) > 0 && isNullResult[0].Bool() {
						// Explicitly set to NULL
						result[key] = nil
					} else {
						// Has a value, get it from the Value field
						if valueField := fieldVal.FieldByName("Value"); valueField.IsValid() {
							result[key] = valueField.Interface()
						}
					}
				}
			}
			continue
		}

		if fieldVal.Kind() == reflect.Ptr {
			// pointer nil → skip (client didn't send this field)
			// pointer non-nil → include with dereferenced value
			if !fieldVal.IsNil() {
				result[key] = fieldVal.Elem().Interface()
			}
		} else {
			// non-pointer: skip zero values
			if !fieldVal.IsZero() {
				result[key] = fieldVal.Interface()
			}
		}
	}

	return result
}
