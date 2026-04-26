package utils

import (
	"encoding/json"
	"fmt"
	"reflect"
)

// Nullable represents a value that can be NULL in the database.
// It distinguishes between three states:
//   - Not provided (omitted from JSON) → Zero value, will be skipped in StructToMap
//   - Provided as null (null in JSON) → null=true, value=nil
//   - Provided with value → null=false, value=actual value
type Nullable[T any] struct {
	Value T
	null  bool // true if explicitly set to null
	set   bool // true if field was in JSON
}

// NewNullable creates a Nullable with a value
func NewNullable[T any](value T) Nullable[T] {
	return Nullable[T]{Value: value, set: true}
}

// NewNull creates a Nullable that represents NULL
func NewNull[T any]() Nullable[T] {
	return Nullable[T]{null: true, set: true}
}

// IsNull returns true if this represents a NULL value
func (n Nullable[T]) IsNull() bool {
	return n.null
}

// IsSet returns true if this field was in the JSON
func (n Nullable[T]) IsSet() bool {
	return n.set
}

// UnmarshalJSON implements json.Unmarshaler
func (n *Nullable[T]) UnmarshalJSON(data []byte) error {
	n.set = true

	if string(data) == "null" {
		n.null = true
		var zero T
		n.Value = zero
		return nil
	}

	n.null = false
	return json.Unmarshal(data, &n.Value)
}

// MarshalJSON implements json.Marshaler
func (n Nullable[T]) MarshalJSON() ([]byte, error) {
	if n.null {
		return []byte("null"), nil
	}
	return json.Marshal(n.Value)
}

// Scan implements sql.Scanner for database/sql
func (n *Nullable[T]) Scan(src interface{}) error {
	if src == nil {
		n.null = true
		n.set = true
		var zero T
		n.Value = zero
		return nil
	}

	n.null = false
	n.set = true

	// Try to unmarshal from JSON (for []byte, string)
	switch v := src.(type) {
	case []byte:
		return json.Unmarshal(v, &n.Value)
	case string:
		return json.Unmarshal([]byte(v), &n.Value)
	}

	// Try direct assignment for compatible types
	srcVal := reflect.ValueOf(src)
	dstType := reflect.TypeOf((*T)(nil)).Elem()

	if srcVal.Type().ConvertibleTo(dstType) {
		reflect.ValueOf(&n.Value).Elem().Set(srcVal.Convert(dstType))
		return nil
	}

	return fmt.Errorf("cannot scan %T into Nullable[%v]", src, dstType)
}
