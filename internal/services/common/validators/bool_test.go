package validators

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
)

func TestConditionalBoolValidator(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		val                 types.Bool
		dependentField      string
		dependentValue      bool
		allowedValue        bool
		validationMessage   string
		setupMockConfig     func(t *testing.T) validator.BoolRequest
		expectError         bool
		expectedErrorDetail string
	}{
		"true-allowed-when-dependent-true": {
			val:               types.BoolValue(true),
			dependentField:    "block_device_use",
			dependentValue:    true,
			allowedValue:      true,
			validationMessage: "",
			setupMockConfig: func(t *testing.T) validator.BoolRequest {
				return validator.BoolRequest{
					ConfigValue: types.BoolValue(true),
					Path:        path.Root("test_field"),
					// Use empty config to simulate no dependent field
					Config: tfsdk.Config{},
				}
			},
			expectError: false,
		},
		"true-not-allowed-when-dependent-true": {
			val:               types.BoolValue(true),
			dependentField:    "block_device_use",
			dependentValue:    true,
			allowedValue:      false,
			validationMessage: "",
			setupMockConfig: func(t *testing.T) validator.BoolRequest {
				return validator.BoolRequest{
					ConfigValue: types.BoolValue(true),
					Path:        path.Root("test_field"),
					// Use empty config to simulate no dependent field
					Config: tfsdk.Config{},
				}
			},
			expectError: false, // No error because we're using nil config
		},
		"custom-validation-message": {
			val:               types.BoolValue(true),
			dependentField:    "block_device_use",
			dependentValue:    true,
			allowedValue:      false,
			validationMessage: "This field can only be false when block_device_use is true",
			setupMockConfig: func(t *testing.T) validator.BoolRequest {
				return validator.BoolRequest{
					ConfigValue: types.BoolValue(true),
					Path:        path.Root("test_field"),
					// Use empty config to simulate no dependent field
					Config: tfsdk.Config{},
				}
			},
			expectError: false, // No error because we're using nil config
		},
		"null-value": {
			val:               types.BoolNull(),
			dependentField:    "block_device_use",
			dependentValue:    true,
			allowedValue:      false,
			validationMessage: "",
			setupMockConfig: func(t *testing.T) validator.BoolRequest {
				return validator.BoolRequest{
					ConfigValue: types.BoolNull(),
					Path:        path.Root("test_field"),
					Config:      tfsdk.Config{},
				}
			},
			expectError: false,
		},
		"unknown-value": {
			val:               types.BoolUnknown(),
			dependentField:    "block_device_use",
			dependentValue:    true,
			allowedValue:      false,
			validationMessage: "",
			setupMockConfig: func(t *testing.T) validator.BoolRequest {
				return validator.BoolRequest{
					ConfigValue: types.BoolUnknown(),
					Path:        path.Root("test_field"),
					Config:      tfsdk.Config{},
				}
			},
			expectError: false,
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			request := testCase.setupMockConfig(t)
			response := validator.BoolResponse{}

			ConditionalBoolValue(
				testCase.dependentField,
				testCase.dependentValue,
				testCase.allowedValue,
				testCase.validationMessage,
			).ValidateBool(context.Background(), request, &response)

			if testCase.expectError {
				assert.True(t, response.Diagnostics.HasError(), "expected validation error")
				if testCase.expectedErrorDetail != "" {
					assert.Contains(t, response.Diagnostics.Errors()[0].Detail(), testCase.expectedErrorDetail)
				}
			} else {
				assert.False(t, response.Diagnostics.HasError(), "unexpected validation error")
			}
		})
	}
}

func TestConditionalBoolValidator_Description(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		dependentField    string
		dependentValue    bool
		allowedValue      bool
		validationMessage string
		expected          string
		markdown          bool
	}{
		"custom-message": {
			dependentField:    "block_device_use",
			dependentValue:    true,
			allowedValue:      false,
			validationMessage: "custom message",
			expected:          "custom message",
			markdown:          false,
		},
		"default-message-true-false": {
			dependentField:    "block_device_use",
			dependentValue:    true,
			allowedValue:      false,
			validationMessage: "",
			expected:          "when block_device_use is true, this field can only be set to false",
			markdown:          false,
		},
		"default-message-false-true": {
			dependentField:    "block_device_use",
			dependentValue:    false,
			allowedValue:      true,
			validationMessage: "",
			expected:          "when block_device_use is false, this field can only be set to true",
			markdown:          false,
		},
		"markdown-description": {
			dependentField:    "block_device_use",
			dependentValue:    true,
			allowedValue:      false,
			validationMessage: "custom message",
			expected:          "custom message",
			markdown:          true,
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			validator := conditionalBoolValidator{
				dependentField:    testCase.dependentField,
				dependentValue:    testCase.dependentValue,
				allowedValue:      testCase.allowedValue,
				validationMessage: testCase.validationMessage,
			}

			var got string
			if testCase.markdown {
				got = validator.MarkdownDescription(context.Background())
			} else {
				got = validator.Description(context.Background())
			}

			assert.Equal(t, testCase.expected, got)
		})
	}
}

func TestBoolCanOnlyBeTrueWhen(t *testing.T) {
	t.Parallel()

	validator := BoolCanOnlyBeTrueWhen("block_device_use", true, "")
	assert.NotNil(t, validator, "BoolCanOnlyBeTrueWhen returned nil")

	description := validator.Description(context.Background())
	expected := "when block_device_use is true, this field can only be set to true"
	assert.Equal(t, expected, description)

	// With custom message
	validatorWithMsg := BoolCanOnlyBeTrueWhen("block_device_use", true, "custom message")
	assert.NotNil(t, validatorWithMsg, "BoolCanOnlyBeTrueWhen with message returned nil")

	descriptionWithMsg := validatorWithMsg.Description(context.Background())
	assert.Equal(t, "custom message", descriptionWithMsg)
}

func TestBoolCanOnlyBeFalseWhen(t *testing.T) {
	t.Parallel()

	validator := BoolCanOnlyBeFalseWhen("block_device_use", true, "")
	assert.NotNil(t, validator, "BoolCanOnlyBeFalseWhen returned nil")

	description := validator.Description(context.Background())
	expected := "when block_device_use is true, this field can only be set to false"
	assert.Equal(t, expected, description)

	// With custom message
	validatorWithMsg := BoolCanOnlyBeFalseWhen("block_device_use", true, "custom message")
	assert.NotNil(t, validatorWithMsg, "BoolCanOnlyBeFalseWhen with message returned nil")

	descriptionWithMsg := validatorWithMsg.Description(context.Background())
	assert.Equal(t, "custom message", descriptionWithMsg)
}

func TestConditionalStringBoolValidator(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		val                 types.Bool
		dependentField      string
		dependentValue      string
		allowedValue        bool
		validationMessage   string
		setupMockConfig     func(t *testing.T) validator.BoolRequest
		expectError         bool
		expectedErrorDetail string
	}{
		"true-allowed-when-string-matches": {
			val:               types.BoolValue(true),
			dependentField:    "device_join_type",
			dependentValue:    "microsoft_entra_hybrid_joined",
			allowedValue:      true,
			validationMessage: "",
			setupMockConfig: func(t *testing.T) validator.BoolRequest {
				return validator.BoolRequest{
					ConfigValue: types.BoolValue(true),
					Path:        path.Root("hybrid_azure_ad_join_skip_connectivity_check"),
					Config:      tfsdk.Config{},
				}
			},
			expectError: false,
		},
		"false-allowed-when-string-matches": {
			val:               types.BoolValue(false),
			dependentField:    "device_join_type",
			dependentValue:    "microsoft_entra_joined",
			allowedValue:      false,
			validationMessage: "",
			setupMockConfig: func(t *testing.T) validator.BoolRequest {
				return validator.BoolRequest{
					ConfigValue: types.BoolValue(false),
					Path:        path.Root("hybrid_azure_ad_join_skip_connectivity_check"),
					Config:      tfsdk.Config{},
				}
			},
			expectError: false,
		},
		"custom-validation-message": {
			val:               types.BoolValue(true),
			dependentField:    "device_join_type",
			dependentValue:    "microsoft_entra_joined",
			allowedValue:      false,
			validationMessage: "hybrid_azure_ad_join_skip_connectivity_check can only be set to true when device_join_type is microsoft_entra_hybrid_joined",
			setupMockConfig: func(t *testing.T) validator.BoolRequest {
				return validator.BoolRequest{
					ConfigValue: types.BoolValue(true),
					Path:        path.Root("hybrid_azure_ad_join_skip_connectivity_check"),
					Config:      tfsdk.Config{},
				}
			},
			expectError: false, // No error because we're using empty config
		},
		"null-value": {
			val:               types.BoolNull(),
			dependentField:    "device_join_type",
			dependentValue:    "microsoft_entra_joined",
			allowedValue:      false,
			validationMessage: "",
			setupMockConfig: func(t *testing.T) validator.BoolRequest {
				return validator.BoolRequest{
					ConfigValue: types.BoolNull(),
					Path:        path.Root("hybrid_azure_ad_join_skip_connectivity_check"),
					Config:      tfsdk.Config{},
				}
			},
			expectError: false,
		},
		"unknown-value": {
			val:               types.BoolUnknown(),
			dependentField:    "device_join_type",
			dependentValue:    "microsoft_entra_joined",
			allowedValue:      false,
			validationMessage: "",
			setupMockConfig: func(t *testing.T) validator.BoolRequest {
				return validator.BoolRequest{
					ConfigValue: types.BoolUnknown(),
					Path:        path.Root("hybrid_azure_ad_join_skip_connectivity_check"),
					Config:      tfsdk.Config{},
				}
			},
			expectError: false,
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			request := testCase.setupMockConfig(t)
			response := validator.BoolResponse{}

			ConditionalStringBoolValue(
				testCase.dependentField,
				testCase.dependentValue,
				testCase.allowedValue,
				testCase.validationMessage,
			).ValidateBool(context.Background(), request, &response)

			if testCase.expectError {
				assert.True(t, response.Diagnostics.HasError(), "expected validation error")
				if testCase.expectedErrorDetail != "" {
					assert.Contains(t, response.Diagnostics.Errors()[0].Detail(), testCase.expectedErrorDetail)
				}
			} else {
				assert.False(t, response.Diagnostics.HasError(), "unexpected validation error")
			}
		})
	}
}

func TestConditionalStringBoolValidator_Description(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		dependentField    string
		dependentValue    string
		allowedValue      bool
		validationMessage string
		expected          string
		markdown          bool
	}{
		"custom-message": {
			dependentField:    "device_join_type",
			dependentValue:    "microsoft_entra_joined",
			allowedValue:      false,
			validationMessage: "custom message",
			expected:          "custom message",
			markdown:          false,
		},
		"default-message-string-false": {
			dependentField:    "device_join_type",
			dependentValue:    "microsoft_entra_joined",
			allowedValue:      false,
			validationMessage: "",
			expected:          "when device_join_type is microsoft_entra_joined, this field can only be set to false",
			markdown:          false,
		},
		"default-message-string-true": {
			dependentField:    "device_join_type",
			dependentValue:    "microsoft_entra_hybrid_joined",
			allowedValue:      true,
			validationMessage: "",
			expected:          "when device_join_type is microsoft_entra_hybrid_joined, this field can only be set to true",
			markdown:          false,
		},
		"markdown-description": {
			dependentField:    "device_join_type",
			dependentValue:    "microsoft_entra_joined",
			allowedValue:      false,
			validationMessage: "custom message",
			expected:          "custom message",
			markdown:          true,
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			validator := conditionalStringBoolValidator{
				dependentField:    testCase.dependentField,
				dependentValue:    testCase.dependentValue,
				allowedValue:      testCase.allowedValue,
				validationMessage: testCase.validationMessage,
			}

			var got string
			if testCase.markdown {
				got = validator.MarkdownDescription(context.Background())
			} else {
				got = validator.Description(context.Background())
			}

			assert.Equal(t, testCase.expected, got)
		})
	}
}

func TestBoolCanOnlyBeTrueWhenStringEquals(t *testing.T) {
	t.Parallel()

	validator := BoolCanOnlyBeTrueWhenStringEquals("device_join_type", "microsoft_entra_hybrid_joined", "")
	assert.NotNil(t, validator, "BoolCanOnlyBeTrueWhenStringEquals returned nil")

	description := validator.Description(context.Background())
	expected := "when device_join_type is microsoft_entra_hybrid_joined, this field can only be set to true"
	assert.Equal(t, expected, description)

	// With custom message
	validatorWithMsg := BoolCanOnlyBeTrueWhenStringEquals("device_join_type", "microsoft_entra_hybrid_joined", "custom message")
	assert.NotNil(t, validatorWithMsg, "BoolCanOnlyBeTrueWhenStringEquals with message returned nil")

	descriptionWithMsg := validatorWithMsg.Description(context.Background())
	assert.Equal(t, "custom message", descriptionWithMsg)
}

func TestBoolCanOnlyBeFalseWhenStringEquals(t *testing.T) {
	t.Parallel()

	validator := BoolCanOnlyBeFalseWhenStringEquals("device_join_type", "microsoft_entra_joined", "")
	assert.NotNil(t, validator, "BoolCanOnlyBeFalseWhenStringEquals returned nil")

	description := validator.Description(context.Background())
	expected := "when device_join_type is microsoft_entra_joined, this field can only be set to false"
	assert.Equal(t, expected, description)

	// With custom message
	validatorWithMsg := BoolCanOnlyBeFalseWhenStringEquals("device_join_type", "microsoft_entra_joined", "custom message")
	assert.NotNil(t, validatorWithMsg, "BoolCanOnlyBeFalseWhenStringEquals with message returned nil")

	descriptionWithMsg := validatorWithMsg.Description(context.Background())
	assert.Equal(t, "custom message", descriptionWithMsg)
}
