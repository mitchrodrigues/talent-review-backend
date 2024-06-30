package helpers

import (
	"reflect"
	"testing"
)

func TestInterfaceValuesEqual(t *testing.T) {
	tests := []struct {
		name string
		a, b interface{}
		want bool
	}{
		{"Equal strings", "hello", "hello", true},
		{"Different strings", "hello", "world", false},
		{"Equal int slices", []int{1, 2, 3}, []int{1, 2, 3}, true},
		{"Different int slices", []int{1, 2, 3}, []int{1, 2, 4}, false},
		{"Equal maps", map[string]int{"a": 1, "b": 2}, map[string]int{"a": 1, "b": 2}, true},
		{"Different maps", map[string]int{"a": 1, "b": 2}, map[string]int{"a": 1, "b": 3}, false},
		{"Different types (string vs int)", "hello", 123, false},
		{"Equal floats", 1.23, 1.23, true},
		{"Different floats", 1.23, 4.56, false},
		{"Nil vs non-nil slice", nil, []int{1, 2, 3}, false},
		{"Nil slices", []int(nil), []int(nil), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := InterfaceValuesEqual(tt.a, tt.b)
			if got != tt.want {
				t.Errorf("for %v and %v, got %v, want %v", tt.a, tt.b, got, tt.want)
			}
		})
	}
}

func TestExtractArg(t *testing.T) {
	tests := []struct {
		name      string
		mp        map[string]interface{}
		key       string
		expected  interface{}
		expectErr bool
	}{
		{
			name: "String type - success",
			mp: map[string]interface{}{
				"cursor": "some-cursor",
			},
			key:       "cursor",
			expected:  "some-cursor",
			expectErr: false,
		},
		{
			name: "Int type - success",
			mp: map[string]interface{}{
				"limit": 20,
			},
			key:       "limit",
			expected:  20,
			expectErr: false,
		},
		{
			name: "Wrong type - error",
			mp: map[string]interface{}{
				"limit": 20,
			},
			key:       "limit",
			expected:  "20", // Expecting a string here to induce a type error
			expectErr: true,
		},
		{
			name: "Non-existent key - error",
			mp: map[string]interface{}{
				"limit": 20,
			},
			key:       "cursor",
			expected:  nil,
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var val interface{}
			var err error

			switch tt.expected.(type) {
			case string:
				val, err = ExtractArg[string](tt.mp, tt.key)
			case int:
				val, err = ExtractArg[int](tt.mp, tt.key)
			default:
				val, err = ExtractArg[interface{}](tt.mp, tt.key)
			}

			if (err != nil) != tt.expectErr {
				t.Errorf("expected error: %v, got: %v", tt.expectErr, err)
			}
			if !tt.expectErr && val != tt.expected {
				t.Errorf("expected value: %v, got: %v", tt.expected, val)
			}
			if tt.expectErr && val != nil {
				var zeroVal interface{}
				switch tt.expected.(type) {
				case string:
					zeroVal = ""
				case int:
					zeroVal = 0
				}
				if !reflect.DeepEqual(val, zeroVal) {
					t.Errorf("expected zero value: %v, got: %v", zeroVal, val)
				}
			}
		})
	}
}
