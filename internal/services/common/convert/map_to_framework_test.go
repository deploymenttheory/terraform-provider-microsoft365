package convert

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
)

func TestMapToFrameworkString(t *testing.T) {
	tests := []struct {
		name     string
		data     map[string]any
		key      string
		expected types.String
	}{
		{
			name:     "valid string value",
			data:     map[string]any{"name": "test"},
			key:      "name",
			expected: types.StringValue("test"),
		},
		{
			name:     "empty string value",
			data:     map[string]any{"name": ""},
			key:      "name",
			expected: types.StringValue(""),
		},
		{
			name:     "key does not exist",
			data:     map[string]any{"other": "value"},
			key:      "name",
			expected: types.StringNull(),
		},
		{
			name:     "value is not string",
			data:     map[string]any{"name": 123},
			key:      "name",
			expected: types.StringNull(),
		},
		{
			name:     "nil map",
			data:     nil,
			key:      "name",
			expected: types.StringNull(),
		},
		{
			name:     "empty map",
			data:     map[string]any{},
			key:      "name",
			expected: types.StringNull(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := MapToFrameworkString(tt.data, tt.key)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestMapToFrameworkBool(t *testing.T) {
	tests := []struct {
		name     string
		data     map[string]any
		key      string
		expected types.Bool
	}{
		{
			name:     "true value",
			data:     map[string]any{"enabled": true},
			key:      "enabled",
			expected: types.BoolValue(true),
		},
		{
			name:     "false value",
			data:     map[string]any{"enabled": false},
			key:      "enabled",
			expected: types.BoolValue(false),
		},
		{
			name:     "key does not exist",
			data:     map[string]any{"other": true},
			key:      "enabled",
			expected: types.BoolNull(),
		},
		{
			name:     "value is not bool",
			data:     map[string]any{"enabled": "true"},
			key:      "enabled",
			expected: types.BoolNull(),
		},
		{
			name:     "nil map",
			data:     nil,
			key:      "enabled",
			expected: types.BoolNull(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := MapToFrameworkBool(tt.data, tt.key)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestMapToFrameworkInt32(t *testing.T) {
	tests := []struct {
		name     string
		data     map[string]any
		key      string
		expected types.Int32
	}{
		{
			name:     "int32 value",
			data:     map[string]any{"count": int32(42)},
			key:      "count",
			expected: types.Int32Value(42),
		},
		{
			name:     "int value",
			data:     map[string]any{"count": 42},
			key:      "count",
			expected: types.Int32Value(42),
		},
		{
			name:     "float64 value",
			data:     map[string]any{"count": 42.0},
			key:      "count",
			expected: types.Int32Value(42),
		},
		{
			name:     "zero value",
			data:     map[string]any{"count": 0},
			key:      "count",
			expected: types.Int32Value(0),
		},
		{
			name:     "negative value",
			data:     map[string]any{"count": -5},
			key:      "count",
			expected: types.Int32Value(-5),
		},
		{
			name:     "key does not exist",
			data:     map[string]any{"other": 42},
			key:      "count",
			expected: types.Int32Null(),
		},
		{
			name:     "value is not numeric",
			data:     map[string]any{"count": "42"},
			key:      "count",
			expected: types.Int32Null(),
		},
		{
			name:     "nil map",
			data:     nil,
			key:      "count",
			expected: types.Int32Null(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := MapToFrameworkInt32(tt.data, tt.key)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestMapToFrameworkInt64(t *testing.T) {
	tests := []struct {
		name     string
		data     map[string]any
		key      string
		expected types.Int64
	}{
		{
			name:     "int64 value",
			data:     map[string]any{"count": int64(9223372036854775807)},
			key:      "count",
			expected: types.Int64Value(9223372036854775807),
		},
		{
			name:     "int value",
			data:     map[string]any{"count": 42},
			key:      "count",
			expected: types.Int64Value(42),
		},
		{
			name:     "int32 value",
			data:     map[string]any{"count": int32(42)},
			key:      "count",
			expected: types.Int64Value(42),
		},
		{
			name:     "float64 value",
			data:     map[string]any{"count": 42.0},
			key:      "count",
			expected: types.Int64Value(42),
		},
		{
			name:     "zero value",
			data:     map[string]any{"count": int64(0)},
			key:      "count",
			expected: types.Int64Value(0),
		},
		{
			name:     "negative value",
			data:     map[string]any{"count": int64(-5)},
			key:      "count",
			expected: types.Int64Value(-5),
		},
		{
			name:     "key does not exist",
			data:     map[string]any{"other": int64(42)},
			key:      "count",
			expected: types.Int64Null(),
		},
		{
			name:     "value is not numeric",
			data:     map[string]any{"count": "42"},
			key:      "count",
			expected: types.Int64Null(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := MapToFrameworkInt64(tt.data, tt.key)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestMapToFrameworkFloat64(t *testing.T) {
	tests := []struct {
		name     string
		data     map[string]any
		key      string
		expected types.Float64
	}{
		{
			name:     "float64 value",
			data:     map[string]any{"price": 99.99},
			key:      "price",
			expected: types.Float64Value(99.99),
		},
		{
			name:     "float32 value",
			data:     map[string]any{"price": float32(99.99)},
			key:      "price",
			expected: types.Float64Value(float64(float32(99.99))),
		},
		{
			name:     "int value",
			data:     map[string]any{"price": 100},
			key:      "price",
			expected: types.Float64Value(100.0),
		},
		{
			name:     "int32 value",
			data:     map[string]any{"price": int32(100)},
			key:      "price",
			expected: types.Float64Value(100.0),
		},
		{
			name:     "int64 value",
			data:     map[string]any{"price": int64(100)},
			key:      "price",
			expected: types.Float64Value(100.0),
		},
		{
			name:     "zero value",
			data:     map[string]any{"price": 0.0},
			key:      "price",
			expected: types.Float64Value(0.0),
		},
		{
			name:     "negative value",
			data:     map[string]any{"price": -5.5},
			key:      "price",
			expected: types.Float64Value(-5.5),
		},
		{
			name:     "key does not exist",
			data:     map[string]any{"other": 99.99},
			key:      "price",
			expected: types.Float64Null(),
		},
		{
			name:     "value is not numeric",
			data:     map[string]any{"price": "99.99"},
			key:      "price",
			expected: types.Float64Null(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := MapToFrameworkFloat64(tt.data, tt.key)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestMapToFrameworkStringSet(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name     string
		data     map[string]any
		key      string
		expected func() types.Set
	}{
		{
			name: "valid string slice",
			data: map[string]any{"tags": []any{"tag1", "tag2", "tag3"}},
			key:  "tags",
			expected: func() types.Set {
				set, _ := types.SetValueFrom(ctx, types.StringType, []string{"tag1", "tag2", "tag3"})
				return set
			},
		},
		{
			name: "empty string slice",
			data: map[string]any{"tags": []any{}},
			key:  "tags",
			expected: func() types.Set {
				return types.SetNull(types.StringType)
			},
		},
		{
			name: "mixed types in slice - filters out non-strings",
			data: map[string]any{"tags": []any{"tag1", 123, "tag2", true, "tag3"}},
			key:  "tags",
			expected: func() types.Set {
				set, _ := types.SetValueFrom(ctx, types.StringType, []string{"tag1", "tag2", "tag3"})
				return set
			},
		},
		{
			name: "single string in slice",
			data: map[string]any{"tags": []any{"single"}},
			key:  "tags",
			expected: func() types.Set {
				set, _ := types.SetValueFrom(ctx, types.StringType, []string{"single"})
				return set
			},
		},
		{
			name:     "key does not exist",
			data:     map[string]any{"other": []any{"tag1"}},
			key:      "tags",
			expected: func() types.Set { return types.SetNull(types.StringType) },
		},
		{
			name:     "value is not slice",
			data:     map[string]any{"tags": "not-a-slice"},
			key:      "tags",
			expected: func() types.Set { return types.SetNull(types.StringType) },
		},
		{
			name:     "nil map",
			data:     nil,
			key:      "tags",
			expected: func() types.Set { return types.SetNull(types.StringType) },
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := MapToFrameworkStringSet(ctx, tt.data, tt.key)
			expected := tt.expected()
			assert.Equal(t, expected, result)
		})
	}
}

func TestMapToFrameworkStringList(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name     string
		data     map[string]any
		key      string
		expected func() types.List
	}{
		{
			name: "valid string slice",
			data: map[string]any{"items": []any{"item1", "item2", "item3"}},
			key:  "items",
			expected: func() types.List {
				list, _ := types.ListValueFrom(ctx, types.StringType, []string{"item1", "item2", "item3"})
				return list
			},
		},
		{
			name: "empty string slice",
			data: map[string]any{"items": []any{}},
			key:  "items",
			expected: func() types.List {
				list, _ := types.ListValueFrom(ctx, types.StringType, []string{})
				return list
			},
		},
		{
			name: "mixed types in slice - filters out non-strings",
			data: map[string]any{"items": []any{"item1", 123, "item2", false, "item3"}},
			key:  "items",
			expected: func() types.List {
				list, _ := types.ListValueFrom(ctx, types.StringType, []string{"item1", "item2", "item3"})
				return list
			},
		},
		{
			name: "single string in slice",
			data: map[string]any{"items": []any{"single"}},
			key:  "items",
			expected: func() types.List {
				list, _ := types.ListValueFrom(ctx, types.StringType, []string{"single"})
				return list
			},
		},
		{
			name:     "key does not exist",
			data:     map[string]any{"other": []any{"item1"}},
			key:      "items",
			expected: func() types.List { return types.ListNull(types.StringType) },
		},
		{
			name:     "value is not slice",
			data:     map[string]any{"items": "not-a-slice"},
			key:      "items",
			expected: func() types.List { return types.ListNull(types.StringType) },
		},
		{
			name:     "nil map",
			data:     nil,
			key:      "items",
			expected: func() types.List { return types.ListNull(types.StringType) },
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := MapToFrameworkStringList(ctx, tt.data, tt.key)
			expected := tt.expected()
			assert.Equal(t, expected, result)
		})
	}
}
