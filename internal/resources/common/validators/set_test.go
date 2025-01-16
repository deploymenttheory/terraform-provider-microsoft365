package validators

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
)

func TestStringSetValidator(t *testing.T) {
	allowedValues := []string{
		"none", "mdm", "windows10XManagement", "configManager",
		"intuneManagementExtension", "thirdParty", "documentGateway",
		"appleRemoteManagement", "microsoftSense", "exchangeOnline",
		"mobileApplicationManagement", "linuxMdm", "enrollment",
		"endpointPrivilegeManagement", "unknownFutureValue",
		"windowsOsRecovery", "android",
	}

	testCases := []struct {
		name          string
		values        []string
		expectError   bool
		errorContains string
	}{
		{
			name:          "valid_single_value",
			values:        []string{"mdm"},
			expectError:   false,
			errorContains: "",
		},
		{
			name:          "valid_multiple_values",
			values:        []string{"mdm", "microsoftSense", "configManager"},
			expectError:   false,
			errorContains: "",
		},
		{
			name:          "invalid_value",
			values:        []string{"mdm", "invalid_value"},
			expectError:   true,
			errorContains: "Set element value must be one of:",
		},
		{
			name:          "empty_set",
			values:        []string{},
			expectError:   false,
			errorContains: "",
		},
		{
			name:          "duplicate_values",
			values:        []string{"mdm", "mdm", "microsoftSense"},
			expectError:   false,
			errorContains: "",
		},
		{
			name:          "all_valid_values",
			values:        allowedValues,
			expectError:   false,
			errorContains: "",
		},
	}

	ctx := context.Background()
	val := StringSetAllowedValues(allowedValues...)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Convert string values to types.String
			var elements []attr.Value
			for _, v := range tc.values {
				elements = append(elements, types.StringValue(v))
			}

			// Create set value
			set, diags := types.SetValue(types.StringType, elements)
			assert.False(t, diags.HasError(), "unexpected error creating set")

			// Create request and response
			request := validator.SetRequest{
				ConfigValue:    set,
				Path:           path.Root("test"),
				PathExpression: path.MatchRoot("test"),
			}
			response := &validator.SetResponse{}

			// Validate
			val.ValidateSet(ctx, request, response)

			if tc.expectError {
				assert.True(t, response.Diagnostics.HasError(), "expected validation error")
				assert.Contains(t, response.Diagnostics.Errors()[0].Detail(), tc.errorContains)
			} else {
				assert.False(t, response.Diagnostics.HasError(), "unexpected validation error")
			}
		})
	}
}

func TestStringSetValidator_NullAndUnknown(t *testing.T) {
	allowedValues := []string{"mdm", "microsoftSense"}
	ctx := context.Background()
	val := StringSetAllowedValues(allowedValues...)

	t.Run("null_value", func(t *testing.T) {
		request := validator.SetRequest{
			ConfigValue:    types.SetNull(types.StringType),
			Path:           path.Root("test"),
			PathExpression: path.MatchRoot("test"),
		}
		response := &validator.SetResponse{}

		val.ValidateSet(ctx, request, response)
		assert.False(t, response.Diagnostics.HasError(), "unexpected error for null value")
	})

	t.Run("unknown_value", func(t *testing.T) {
		request := validator.SetRequest{
			ConfigValue:    types.SetUnknown(types.StringType),
			Path:           path.Root("test"),
			PathExpression: path.MatchRoot("test"),
		}
		response := &validator.SetResponse{}

		val.ValidateSet(ctx, request, response)
		assert.False(t, response.Diagnostics.HasError(), "unexpected error for unknown value")
	})
}

func TestStringSetValidator_Description(t *testing.T) {
	allowedValues := []string{"mdm", "microsoftSense"}
	val := StringSetAllowedValues(allowedValues...)
	ctx := context.Background()

	description := val.(describer).Description(ctx)
	assert.Contains(t, description, "mdm")
	assert.Contains(t, description, "microsoftSense")

	markdownDescription := val.(describer).MarkdownDescription(ctx)
	assert.Equal(t, description, markdownDescription)
}

// Helper interface for testing Description methods
type describer interface {
	Description(context.Context) string
	MarkdownDescription(context.Context) string
}
