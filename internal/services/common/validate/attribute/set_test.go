package attribute

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
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

func TestSetRequiresStringValue(t *testing.T) {
	// Note: Full validation logic with actual config values is tested via acceptance tests.
	// Unit tests focus on early return cases (null, unknown, empty) that don't require config access.

	ctx := context.Background()
	val := SetRequiresStringValue("string_field", []string{"fido2"}, "")

	t.Run("empty_set_should_pass", func(t *testing.T) {
		elements := []attr.Value{}
		set, diags := types.SetValue(types.StringType, elements)
		assert.False(t, diags.HasError(), "unexpected error creating set")

		request := validator.SetRequest{
			ConfigValue:    set,
			Path:           path.Root("test_set"),
			PathExpression: path.MatchRoot("test_set"),
			Config:         tfsdk.Config{},
		}
		response := &validator.SetResponse{}

		val.ValidateSet(ctx, request, response)
		assert.False(t, response.Diagnostics.HasError(), "unexpected validation error for empty set")
	})
}

func TestSetRequiresStringValue_NullAndUnknown(t *testing.T) {
	ctx := context.Background()
	val := SetRequiresStringValue("string_field", []string{"fido2"}, "")

	t.Run("null_set_value", func(t *testing.T) {
		request := validator.SetRequest{
			ConfigValue:    types.SetNull(types.StringType),
			Path:           path.Root("test_set"),
			PathExpression: path.MatchRoot("test_set"),
			Config:         tfsdk.Config{},
		}
		response := &validator.SetResponse{}

		val.ValidateSet(ctx, request, response)
		assert.False(t, response.Diagnostics.HasError(), "unexpected error for null set value")
	})

	t.Run("unknown_set_value", func(t *testing.T) {
		request := validator.SetRequest{
			ConfigValue:    types.SetUnknown(types.StringType),
			Path:           path.Root("test_set"),
			PathExpression: path.MatchRoot("test_set"),
			Config:         tfsdk.Config{},
		}
		response := &validator.SetResponse{}

		val.ValidateSet(ctx, request, response)
		assert.False(t, response.Diagnostics.HasError(), "unexpected error for unknown set value")
	})

	t.Run("empty_set", func(t *testing.T) {
		elements := []attr.Value{}
		set, _ := types.SetValue(types.StringType, elements)

		request := validator.SetRequest{
			ConfigValue:    set,
			Path:           path.Root("test_set"),
			PathExpression: path.MatchRoot("test_set"),
			Config:         tfsdk.Config{},
		}
		response := &validator.SetResponse{}

		val.ValidateSet(ctx, request, response)
		assert.False(t, response.Diagnostics.HasError(), "unexpected error for empty set")
	})
}

func TestSetRequiresStringValue_Description(t *testing.T) {
	ctx := context.Background()

	t.Run("single_allowed_value", func(t *testing.T) {
		val := SetRequiresStringValue("applies_to_combinations", []string{"fido2"}, "")

		description := val.(describer).Description(ctx)
		assert.Contains(t, description, "applies_to_combinations")
		assert.Contains(t, description, "fido2")

		markdownDescription := val.(describer).MarkdownDescription(ctx)
		assert.Equal(t, description, markdownDescription)
	})

	t.Run("multiple_allowed_values", func(t *testing.T) {
		val := SetRequiresStringValue("applies_to_combinations", []string{"x509CertificateMultiFactor", "x509CertificateSingleFactor"}, "")

		description := val.(describer).Description(ctx)
		assert.Contains(t, description, "applies_to_combinations")
		assert.Contains(t, description, "x509CertificateMultiFactor")
		assert.Contains(t, description, "x509CertificateSingleFactor")

		markdownDescription := val.(describer).MarkdownDescription(ctx)
		assert.Equal(t, description, markdownDescription)
	})

	t.Run("custom_validation_message", func(t *testing.T) {
		customMessage := "This is a custom validation message"
		val := SetRequiresStringValue("field", []string{"value"}, customMessage)

		description := val.(describer).Description(ctx)
		assert.Equal(t, customMessage, description)

		markdownDescription := val.(describer).MarkdownDescription(ctx)
		assert.Equal(t, customMessage, markdownDescription)
	})
}

// Helper interface for testing Description methods
type describer interface {
	Description(context.Context) string
	MarkdownDescription(context.Context) string
}
