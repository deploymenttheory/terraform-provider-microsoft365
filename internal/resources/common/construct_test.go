package common

import (
	"errors"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
)

func TestSetStringValueFromAttributes(t *testing.T) {
	t.Run("Existing non-null string value", func(t *testing.T) {
		attrs := map[string]attr.Value{
			"test_key": types.StringValue("test_value"),
		}
		var result string
		SetStringValueFromAttributes(attrs, "test_key", func(s *string) {
			result = *s
		})
		assert.Equal(t, "test_value", result)
	})

	t.Run("Non-existing key", func(t *testing.T) {
		attrs := map[string]attr.Value{}
		var result string
		SetStringValueFromAttributes(attrs, "non_existing_key", func(s *string) {
			result = *s
		})
		assert.Equal(t, "", result)
	})

	t.Run("Null string value", func(t *testing.T) {
		attrs := map[string]attr.Value{
			"test_key": types.StringNull(),
		}
		var result string
		SetStringValueFromAttributes(attrs, "test_key", func(s *string) {
			result = *s
		})
		assert.Equal(t, "", result)
	})

	t.Run("Non-string value", func(t *testing.T) {
		attrs := map[string]attr.Value{
			"test_key": types.Int64Value(42),
		}
		var result string
		SetStringValueFromAttributes(attrs, "test_key", func(s *string) {
			result = *s
		})
		assert.Equal(t, "", result)
	})
}

func TestSetParsedValueFromAttributes(t *testing.T) {
	t.Run("Existing non-null parsable value", func(t *testing.T) {
		attrs := map[string]attr.Value{
			"test_key": types.StringValue("42"),
		}
		var result int
		err := SetParsedValueFromAttributes(attrs, "test_key", func(i *int) {
			result = *i
		}, func(s string) (interface{}, error) {
			return strconv.Atoi(s)
		})
		assert.NoError(t, err)
		assert.Equal(t, 42, result)
	})

	t.Run("Non-existing key", func(t *testing.T) {
		attrs := map[string]attr.Value{}
		var result int
		err := SetParsedValueFromAttributes(attrs, "non_existing_key", func(i *int) {
			result = *i
		}, func(s string) (interface{}, error) {
			return strconv.Atoi(s)
		})
		assert.NoError(t, err)
		assert.Equal(t, 0, result)
	})

	t.Run("Null string value", func(t *testing.T) {
		attrs := map[string]attr.Value{
			"test_key": types.StringNull(),
		}
		var result int
		err := SetParsedValueFromAttributes(attrs, "test_key", func(i *int) {
			result = *i
		}, func(s string) (interface{}, error) {
			return strconv.Atoi(s)
		})
		assert.NoError(t, err)
		assert.Equal(t, 0, result)
	})

	t.Run("Non-string value", func(t *testing.T) {
		attrs := map[string]attr.Value{
			"test_key": types.Int64Value(42),
		}
		var result int
		err := SetParsedValueFromAttributes(attrs, "test_key", func(i *int) {
			result = *i
		}, func(s string) (interface{}, error) {
			return strconv.Atoi(s)
		})
		assert.NoError(t, err)
		assert.Equal(t, 0, result)
	})

	t.Run("Parsing error", func(t *testing.T) {
		attrs := map[string]attr.Value{
			"test_key": types.StringValue("not_a_number"),
		}
		var result int
		err := SetParsedValueFromAttributes(attrs, "test_key", func(i *int) {
			result = *i
		}, func(s string) (interface{}, error) {
			return strconv.Atoi(s)
		})
		assert.Error(t, err)
		assert.Equal(t, 0, result)
	})

	t.Run("Nil parsed value", func(t *testing.T) {
		attrs := map[string]attr.Value{
			"test_key": types.StringValue("nil_value"),
		}
		var result *string
		err := SetParsedValueFromAttributes(attrs, "test_key", func(s *string) {
			result = s
		}, func(s string) (interface{}, error) {
			return nil, nil
		})
		assert.NoError(t, err)
		assert.Nil(t, result)
	})

	t.Run("Custom error", func(t *testing.T) {
		attrs := map[string]attr.Value{
			"test_key": types.StringValue("error_value"),
		}
		var result string
		customErr := errors.New("custom error")
		err := SetParsedValueFromAttributes(attrs, "test_key", func(s *string) {
			result = *s
		}, func(s string) (interface{}, error) {
			return nil, customErr
		})
		assert.Equal(t, customErr, err)
		assert.Equal(t, "", result)
	})
}
