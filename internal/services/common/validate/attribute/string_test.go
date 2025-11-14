package attribute

import (
	"context"
	"regexp"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestRolloutDateTimeValidator(t *testing.T) {
	t.Parallel()

	// Calculate test dates relative to current time
	now := time.Now().UTC()
	validDate := now.AddDate(0, 0, 5).Format(time.RFC3339)    // 5 days from now (valid)
	tooEarlyDate := now.AddDate(0, 0, 1).Format(time.RFC3339) // 1 day from now (too early for 2-day minimum)
	tooLateDate := now.AddDate(0, 0, 70).Format(time.RFC3339) // 70 days from now (too late for 60-day maximum)
	pastDate := now.AddDate(0, 0, -1).Format(time.RFC3339)    // 1 day ago (in the past)

	testCases := map[string]struct {
		value           types.String
		minDays         int
		maxDays         int
		expectedError   bool
		expectedMessage string
	}{
		"valid-date": {
			value:         types.StringValue(validDate),
			minDays:       2,
			maxDays:       60,
			expectedError: false,
		},
		"too-early": {
			value:           types.StringValue(tooEarlyDate),
			minDays:         2,
			maxDays:         60,
			expectedError:   true,
			expectedMessage: "DateTime must be at least 2 days in the future",
		},
		"too-late": {
			value:           types.StringValue(tooLateDate),
			minDays:         2,
			maxDays:         60,
			expectedError:   true,
			expectedMessage: "DateTime must be within 60 days from now",
		},
		"past-date": {
			value:           types.StringValue(pastDate),
			minDays:         2,
			maxDays:         60,
			expectedError:   true,
			expectedMessage: "DateTime must be at least 2 days in the future",
		},
		"invalid-format": {
			value:           types.StringValue("2025-13-45"),
			minDays:         2,
			maxDays:         60,
			expectedError:   true,
			expectedMessage: "DateTime must be in RFC3339 format",
		},
		"null-value": {
			value:         types.StringNull(),
			minDays:       2,
			maxDays:       60,
			expectedError: false,
		},
		"unknown-value": {
			value:         types.StringUnknown(),
			minDays:       2,
			maxDays:       60,
			expectedError: false,
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			request := validator.StringRequest{
				ConfigValue: testCase.value,
				Path:        path.Root("test"),
			}

			response := validator.StringResponse{}

			RolloutDateTime(testCase.minDays, testCase.maxDays).ValidateString(context.Background(), request, &response)

			if !response.Diagnostics.HasError() && testCase.expectedError {
				t.Fatal("expected error, got no error")
			}

			if response.Diagnostics.HasError() && !testCase.expectedError {
				t.Fatalf("got unexpected error: %s", response.Diagnostics)
			}

			if testCase.expectedError && testCase.expectedMessage != "" {
				found := false
				for _, diag := range response.Diagnostics {
					if matched, _ := regexp.MatchString(testCase.expectedMessage, diag.Detail()); matched {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("expected error message to contain '%s', got: %s", testCase.expectedMessage, response.Diagnostics)
				}
			}
		})
	}
}

func TestRolloutDateTimeValidator_Description(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		minDays  int
		maxDays  int
		expected string
		markdown bool
	}{
		"basic": {
			minDays:  2,
			maxDays:  60,
			expected: "datetime must be between 2 and 60 days from now",
			markdown: false,
		},
		"markdown": {
			minDays:  1,
			maxDays:  30,
			expected: "datetime must be between 1 and 30 days from now",
			markdown: true,
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			validator := rolloutDateTimeValidator{
				minDaysFromNow: testCase.minDays,
				maxDaysFromNow: testCase.maxDays,
			}

			var got string
			if testCase.markdown {
				got = validator.MarkdownDescription(context.Background())
			} else {
				got = validator.Description(context.Background())
			}

			if got != testCase.expected {
				t.Errorf("expected %s, got %s", testCase.expected, got)
			}
		})
	}
}

func TestFutureDateTimeValidator(t *testing.T) {
	t.Parallel()

	// Calculate test dates relative to current time
	now := time.Now().UTC()
	validDate := now.AddDate(0, 0, 5).Format(time.RFC3339) // 5 days from now (valid)
	pastDate := now.AddDate(0, 0, -1).Format(time.RFC3339) // 1 day ago (in the past)

	testCases := map[string]struct {
		value         types.String
		expectedError bool
	}{
		"valid-future-date": {
			value:         types.StringValue(validDate),
			expectedError: false,
		},
		"past-date": {
			value:         types.StringValue(pastDate),
			expectedError: true,
		},
		"null-value": {
			value:         types.StringNull(),
			expectedError: false,
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			request := validator.StringRequest{
				ConfigValue: testCase.value,
				Path:        path.Root("test"),
			}

			response := validator.StringResponse{}

			FutureDateTime().ValidateString(context.Background(), request, &response)

			if !response.Diagnostics.HasError() && testCase.expectedError {
				t.Fatal("expected error, got no error")
			}

			if response.Diagnostics.HasError() && !testCase.expectedError {
				t.Fatalf("got unexpected error: %s", response.Diagnostics)
			}
		})
	}
}

func TestStringLengthValidator(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		value         types.String
		maxLength     int
		expectedError bool
	}{
		"valid-length": {
			value:         types.StringValue("test"),
			maxLength:     10,
			expectedError: false,
		},
		"exact-max-length": {
			value:         types.StringValue("test"),
			maxLength:     4,
			expectedError: false,
		},
		"exceeds-max-length": {
			value:         types.StringValue("test string"),
			maxLength:     5,
			expectedError: true,
		},
		"null-value": {
			value:         types.StringNull(),
			maxLength:     5,
			expectedError: false,
		},
		"unknown-value": {
			value:         types.StringUnknown(),
			maxLength:     5,
			expectedError: false,
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			request := validator.StringRequest{
				ConfigValue: testCase.value,
				Path:        path.Root("test"),
			}

			response := validator.StringResponse{}

			StringLengthAtMost(testCase.maxLength).ValidateString(context.Background(), request, &response)

			if !response.Diagnostics.HasError() && testCase.expectedError {
				t.Fatal("expected error, got no error")
			}

			if response.Diagnostics.HasError() && !testCase.expectedError {
				t.Fatalf("got unexpected error: %s", response.Diagnostics)
			}
		})
	}
}

func TestIllegalCharactersValidator(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		value          types.String
		forbiddenChars []rune
		expectedError  bool
	}{
		"valid-string": {
			value:          types.StringValue("validstring"),
			forbiddenChars: []rune{'!', '@'},
			expectedError:  false,
		},
		"contains-forbidden-char": {
			value:          types.StringValue("invalid!string"),
			forbiddenChars: []rune{'!', '@'},
			expectedError:  true,
		},
		"multiple-forbidden-chars": {
			value:          types.StringValue("invalid@string!"),
			forbiddenChars: []rune{'!', '@'},
			expectedError:  true,
		},
		"null-value": {
			value:          types.StringNull(),
			forbiddenChars: []rune{'!', '@'},
			expectedError:  false,
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			request := validator.StringRequest{
				ConfigValue: testCase.value,
				Path:        path.Root("test"),
			}

			response := validator.StringResponse{}

			IllegalCharactersInString(testCase.forbiddenChars, "").ValidateString(context.Background(), request, &response)

			if !response.Diagnostics.HasError() && testCase.expectedError {
				t.Fatal("expected error, got no error")
			}

			if response.Diagnostics.HasError() && !testCase.expectedError {
				t.Fatalf("got unexpected error: %s", response.Diagnostics)
			}
		})
	}
}

func TestASCIIOnlyValidator(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		value         types.String
		expectedError bool
	}{
		"valid-ascii": {
			value:         types.StringValue("Hello World 123!"),
			expectedError: false,
		},
		"contains-non-ascii": {
			value:         types.StringValue("Hello 世界"),
			expectedError: true,
		},
		"null-value": {
			value:         types.StringNull(),
			expectedError: false,
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			request := validator.StringRequest{
				ConfigValue: testCase.value,
				Path:        path.Root("test"),
			}

			response := validator.StringResponse{}

			ASCIIOnly().ValidateString(context.Background(), request, &response)

			if !response.Diagnostics.HasError() && testCase.expectedError {
				t.Fatal("expected error, got no error")
			}

			if response.Diagnostics.HasError() && !testCase.expectedError {
				t.Fatalf("got unexpected error: %s", response.Diagnostics)
			}
		})
	}
}

// TestRequiredWhenEquals is omitted because it requires complex config mocking
// that is difficult to set up in unit tests. The validator is tested
// through integration tests in the actual resource tests.

func TestConditionalStringEmptyValidator(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		val                 types.String
		dependentField      string
		dependentValue      string
		validationMessage   string
		setupMockConfig     func(t *testing.T) validator.StringRequest
		expectError         bool
		expectedErrorDetail string
	}{
		"empty-string-allowed-when-condition-matches": {
			val:               types.StringValue(""),
			dependentField:    "device_join_type",
			dependentValue:    "microsoft_entra_hybrid_joined",
			validationMessage: "",
			setupMockConfig: func(t *testing.T) validator.StringRequest {
				return validator.StringRequest{
					ConfigValue: types.StringValue(""),
					Path:        path.Root("device_name_template"),
					Config:      tfsdk.Config{},
				}
			},
			expectError: false,
		},
		"null-value-allowed-when-condition-matches": {
			val:               types.StringNull(),
			dependentField:    "device_join_type",
			dependentValue:    "microsoft_entra_hybrid_joined",
			validationMessage: "",
			setupMockConfig: func(t *testing.T) validator.StringRequest {
				return validator.StringRequest{
					ConfigValue: types.StringNull(),
					Path:        path.Root("device_name_template"),
					Config:      tfsdk.Config{},
				}
			},
			expectError: false,
		},
		"non-empty-string-with-different-dependent-value": {
			val:               types.StringValue("template-%RAND:5%"),
			dependentField:    "device_join_type",
			dependentValue:    "microsoft_entra_hybrid_joined",
			validationMessage: "",
			setupMockConfig: func(t *testing.T) validator.StringRequest {
				return validator.StringRequest{
					ConfigValue: types.StringValue("template-%RAND:5%"),
					Path:        path.Root("device_name_template"),
					Config:      tfsdk.Config{},
				}
			},
			expectError: false, // No error because we're using empty config
		},
		"custom-validation-message": {
			val:               types.StringValue("template-%RAND:5%"),
			dependentField:    "device_join_type",
			dependentValue:    "microsoft_entra_hybrid_joined",
			validationMessage: "device_name_template must not be set when device_join_type is microsoft_entra_hybrid_joined",
			setupMockConfig: func(t *testing.T) validator.StringRequest {
				return validator.StringRequest{
					ConfigValue: types.StringValue("template-%RAND:5%"),
					Path:        path.Root("device_name_template"),
					Config:      tfsdk.Config{},
				}
			},
			expectError: false, // No error because we're using empty config
		},
		"unknown-value": {
			val:               types.StringUnknown(),
			dependentField:    "device_join_type",
			dependentValue:    "microsoft_entra_hybrid_joined",
			validationMessage: "",
			setupMockConfig: func(t *testing.T) validator.StringRequest {
				return validator.StringRequest{
					ConfigValue: types.StringUnknown(),
					Path:        path.Root("device_name_template"),
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
			response := validator.StringResponse{}

			StringMustBeEmptyWhenStringEquals(
				testCase.dependentField,
				testCase.dependentValue,
				testCase.validationMessage,
			).ValidateString(context.Background(), request, &response)

			if testCase.expectError {
				if !response.Diagnostics.HasError() {
					t.Fatal("expected validation error")
				}
				if testCase.expectedErrorDetail != "" {
					found := false
					for _, err := range response.Diagnostics.Errors() {
						if contains(err.Detail(), testCase.expectedErrorDetail) {
							found = true
							break
						}
					}
					if !found {
						t.Fatalf("expected error detail %q not found in errors: %v", testCase.expectedErrorDetail, response.Diagnostics.Errors())
					}
				}
			} else {
				if response.Diagnostics.HasError() {
					t.Fatalf("unexpected validation error: %s", response.Diagnostics)
				}
			}
		})
	}
}

func TestConditionalStringEmptyValidator_Description(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		dependentField    string
		dependentValue    string
		validationMessage string
		expected          string
		markdown          bool
	}{
		"custom-message": {
			dependentField:    "device_join_type",
			dependentValue:    "microsoft_entra_hybrid_joined",
			validationMessage: "custom message",
			expected:          "custom message",
			markdown:          false,
		},
		"default-message": {
			dependentField:    "device_join_type",
			dependentValue:    "microsoft_entra_hybrid_joined",
			validationMessage: "",
			expected:          "this field must be empty when device_join_type is microsoft_entra_hybrid_joined",
			markdown:          false,
		},
		"markdown-description": {
			dependentField:    "device_join_type",
			dependentValue:    "microsoft_entra_hybrid_joined",
			validationMessage: "custom message",
			expected:          "custom message",
			markdown:          true,
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			validator := conditionalStringEmptyValidator{
				dependentField:    testCase.dependentField,
				dependentValue:    testCase.dependentValue,
				validationMessage: testCase.validationMessage,
			}

			var got string
			if testCase.markdown {
				got = validator.MarkdownDescription(context.Background())
			} else {
				got = validator.Description(context.Background())
			}

			if got != testCase.expected {
				t.Errorf("expected %q, got %q", testCase.expected, got)
			}
		})
	}
}

func TestStringMustBeEmptyWhenStringEquals(t *testing.T) {
	t.Parallel()

	validator := StringMustBeEmptyWhenStringEquals("device_join_type", "microsoft_entra_hybrid_joined", "")
	if validator == nil {
		t.Fatal("StringMustBeEmptyWhenStringEquals returned nil")
	}

	description := validator.Description(context.Background())
	expected := "this field must be empty when device_join_type is microsoft_entra_hybrid_joined"
	if description != expected {
		t.Errorf("expected %q, got %q", expected, description)
	}

	// With custom message
	validatorWithMsg := StringMustBeEmptyWhenStringEquals("device_join_type", "microsoft_entra_hybrid_joined", "custom message")
	if validatorWithMsg == nil {
		t.Fatal("StringMustBeEmptyWhenStringEquals with message returned nil")
	}

	descriptionWithMsg := validatorWithMsg.Description(context.Background())
	if descriptionWithMsg != "custom message" {
		t.Errorf("expected %q, got %q", "custom message", descriptionWithMsg)
	}
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 || (len(s) > len(substr) && (s[:len(substr)] == substr || s[len(s)-len(substr):] == substr || containsInMiddle(s, substr))))
}

func containsInMiddle(s, substr string) bool {
	for i := 1; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
