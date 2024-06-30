package helpers

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/google/uuid"
)

var (
	ErrUnsupportedDataType = errors.New("unsupported data type")

	ErrNoSuchKey = errors.New("no such key")
)

func InterfaceValuesEqual(a, b interface{}) bool {
	if reflect.TypeOf(a) != reflect.TypeOf(b) {
		return false
	}

	switch reflect.TypeOf(a).Kind() {
	case reflect.Slice:
		return reflect.DeepEqual(a, b)
	case reflect.Map:
		return reflect.DeepEqual(a, b)
	default:
		return a == b
	}
}

// ToInterface converts an interface that may be a pointer to a non-pointer
// or a non-pointer to a pointer, based on the input.
func ToInterface(input interface{}) interface{} {
	val := reflect.ValueOf(input)

	// If the input is a pointer, dereference it.
	// Otherwise, get the address of the input.
	if val.Kind() == reflect.Ptr {
		return val.Elem().Interface()
	} else {
		return val.Addr().Interface()
	}
}

func IDField(model interface{}) interface{} {
	value := valueOf(model)

	if v := value.FieldByName("ID"); v.IsValid() {
		switch id := v.Interface().(type) {
		case uuid.UUID:
			return id
		case map[string]interface{}:
			return id["_id"]
		default:
			return id
		}
	}
	return ""
}

func valueOf(obj interface{}) reflect.Value {
	value := reflect.ValueOf(obj)

	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}
	return value
}

// GetString safely extracts a string from a map.
func GetString(m map[string]interface{}, key string) (string, error) {
	return ExtractArg[string](m, key)
}

func ExtractArg[T any](mp map[string]interface{}, key string) (T, error) {
	var zero T

	if mp == nil {
		return zero, fmt.Errorf("map is nil")
	}

	if v1, found := mp[key]; found {
		if val, ok := v1.(T); ok {
			return val, nil
		}
		return zero, fmt.Errorf("value is not the correct type: expected %T, got %T", zero, v1)
	}

	return zero, ErrNoSuchKey
}

// Define a helper function for extracting and parsing UUID
func ExtractAndParseUUID(args map[string]interface{}, key string) (uuid.UUID, error) {
	idString, err := ExtractArg[string](args, key)
	if err != nil {
		if err == ErrNoSuchKey {
			return uuid.Nil, nil
		}
		return uuid.Nil, err
	}

	return uuid.Parse(idString)
}
