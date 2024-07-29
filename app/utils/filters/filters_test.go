package filters

import (
	"reflect"
	"testing"

	"github.com/golly-go/golly"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// DummyModel is a placeholder model for executing the query
type DummyModel struct {
	ID int
}

func TestConvertValue(t *testing.T) {
	tests := []struct {
		kind    reflect.Kind
		input   interface{}
		want    interface{}
		wantErr bool
	}{
		{reflect.String, "test", "test", false},
		{reflect.Int, "123", int64(123), false},
		{reflect.Float64, "123.45", 123.45, false},
		{reflect.Bool, "true", true, false},
		{reflect.Invalid, "invalid", nil, true},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			got, err := convertValue(tt.kind, tt.input)
			if (err != nil) != tt.wantErr {
				assert.Equal(t, tt.wantErr, err != nil, "convertValue() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got, "convertValue() = %v, want %v", got, tt.want)
		})
	}
}

func TestExtractOperatorAndValue(t *testing.T) {
	tests := []struct {
		operator string
		value    interface{}
		wantOp   string
		wantVal  interface{}
	}{
		{"contains", "test", "LIKE", "%test%"},
		{"starts_with", "test", "LIKE", "test%"},
		{"ends_with", "test", "LIKE", "%test"},
		{"has", nil, "IS NOT NULL", nil},
		{"is", "123", "=", "123"},
		{"invalid", "test", "=", "test"},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			gotOp, gotVal := extractOperatorAndValue(tt.operator, tt.value)
			assert.Equal(t, tt.wantOp, gotOp, "extractOperatorAndValue() gotOp = %v, want %v", gotOp, tt.wantOp)
			assert.Equal(t, tt.wantVal, gotVal, "extractOperatorAndValue() gotVal = %v, want %v", gotVal, tt.wantVal)
		})
	}
}

func TestScopes(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err, "failed to open sqlite database")

	fields := map[string]FieldType{
		"isManager":  {Kind: reflect.Bool, DBFieldName: "is_manager"},
		"age":        {Kind: reflect.Int, DBFieldName: "age"},
		"name":       {Kind: reflect.String, DBFieldName: "name"},
		"createdAt":  {Kind: reflect.String, DBFieldName: "created_at"},
		"department": {Kind: reflect.String, DBFieldName: "department"},
	}

	filter := NewFilter("User", fields)

	tests := []struct {
		name      string
		filterArg []interface{}
		wantQuery string
		wantArgs  []interface{}
		wantErr   bool
	}{
		{
			name: "test isManager true",
			filterArg: []interface{}{
				map[string]interface{}{
					"field":    "isManager",
					"operator": "is",
					"value":    "true",
				},
			},
			wantQuery: "SELECT * FROM `dummy_models` WHERE is_manager = ?",
			wantArgs:  []interface{}{true},
		},
		{
			name: "test age greater than",
			filterArg: []interface{}{
				map[string]interface{}{
					"field":    "age",
					"operator": "greater than",
					"value":    "30",
				},
			},
			wantQuery: "SELECT * FROM `dummy_models` WHERE age > ?",
			wantArgs:  []interface{}{int64(30)},
		},
		{
			name: "test name contains",
			filterArg: []interface{}{
				map[string]interface{}{
					"field":    "name",
					"operator": "contains",
					"value":    "John",
				},
			},
			wantQuery: "SELECT * FROM `dummy_models` WHERE name LIKE ?",
			wantArgs:  []interface{}{"%John%"},
		},
		{
			name: "test createdAt before",
			filterArg: []interface{}{
				map[string]interface{}{
					"field":    "createdAt",
					"operator": "less than",
					"value":    "2023-01-01",
				},
			},
			wantQuery: "SELECT * FROM `dummy_models` WHERE created_at < ?",
			wantArgs:  []interface{}{"2023-01-01"},
		},
		{
			name: "test department starts_with",
			filterArg: []interface{}{
				map[string]interface{}{
					"field":    "department",
					"operator": "starts_with",
					"value":    "Sales",
				},
			},
			wantQuery: "SELECT * FROM `dummy_models` WHERE department LIKE ?",
			wantArgs:  []interface{}{"Sales%"},
		},
		{
			name: "test complex filters",
			filterArg: []interface{}{
				map[string]interface{}{
					"field":    "isManager",
					"operator": "is",
					"value":    "true",
				},
				map[string]interface{}{
					"field":    "age",
					"operator": "greater than",
					"value":    "30",
				},
				map[string]interface{}{
					"field":    "department",
					"operator": "starts_with",
					"value":    "Sales",
				},
			},
			wantQuery: "SELECT * FROM `dummy_models` WHERE is_manager = ? AND age > ? AND department LIKE ?",
			wantArgs:  []interface{}{true, int64(30), "Sales%"},
		},
		{
			name: "test invalid field",
			filterArg: []interface{}{
				map[string]interface{}{
					"field":    "invalidField",
					"operator": "is",
					"value":    "true",
				},
			},
			wantErr: true,
		},
		{
			name: "test invalid value",
			filterArg: []interface{}{
				map[string]interface{}{
					"field":    "isManager",
					"operator": "is",
					"value":    nil,
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scopes, err := filter.Scopes(golly.Context{}, tt.filterArg)

			if tt.wantErr {
				assert.Error(t, err, "expected an error but got none")
				return
			}

			assert.NoError(t, err, "unexpected error: %v", err)

			tx := db.Session(&gorm.Session{DryRun: true}).Scopes(scopes...)

			// Execute a dummy query to generate the SQL
			var results []DummyModel
			tx.Find(&results)

			stmt := tx.Statement

			// Debugging output
			t.Logf("Generated SQL: %v", stmt.SQL.String())
			t.Logf("Generated Args: %v", stmt.Vars)

			assert.Equal(t, tt.wantQuery, stmt.SQL.String(), "expected query %v, got %v", tt.wantQuery, stmt.SQL.String())
			assert.Equal(t, tt.wantArgs, stmt.Vars, "expected args %v, got %v", tt.wantArgs, stmt.Vars)
		})
	}
}
