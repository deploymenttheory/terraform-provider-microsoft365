package validators

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestDeviceNameTemplateValidator(t *testing.T) {
	t.Parallel()

	type testCase struct {
		name           string
		value          types.String
		expectError    bool
		errorSubstring string
	}

	testCases := []testCase{
		{
			name:        "empty string is valid",
			value:       types.StringValue(""),
			expectError: false,
		},
		{
			name:        "null value is valid",
			value:       types.StringNull(),
			expectError: false,
		},
		{
			name:        "unknown value is valid",
			value:       types.StringUnknown(),
			expectError: false,
		},
		{
			name:        "simple valid name",
			value:       types.StringValue("DESKTOP-01"),
			expectError: false,
		},
		{
			name:        "valid with SERIAL macro",
			value:       types.StringValue("PC-%SERIAL%"),
			expectError: false,
		},
		{
			name:        "valid with RAND macro",
			value:       types.StringValue("PC-%RAND:4%"),
			expectError: false,
		},
		{
			name:        "valid short template",
			value:       types.StringValue("A%RAND:1%"),
			expectError: false,
		},
		{
			name:        "maximum length (15 chars)",
			value:       types.StringValue("VERYLONGNAME123"),
			expectError: false,
		},
		{
			name:           "too long (16 chars)",
			value:          types.StringValue("VERYLONGNAME1234"),
			expectError:    true,
			errorSubstring: "exceeds maximum allowed length of 15 characters",
		},
		{
			name:           "contains space",
			value:          types.StringValue("DESKTOP 01"),
			expectError:    true,
			errorSubstring: "cannot contain blank spaces",
		},
		{
			name:           "only numbers",
			value:          types.StringValue("12345"),
			expectError:    true,
			errorSubstring: "cannot contain only numbers",
		},
		{
			name:           "only numbers with hyphens",
			value:          types.StringValue("123-456"),
			expectError:    true,
			errorSubstring: "cannot contain only numbers",
		},
		{
			name:           "invalid character (@)",
			value:          types.StringValue("PC@001"),
			expectError:    true,
			errorSubstring: "contains invalid character '@'",
		},
		{
			name:           "invalid macro",
			value:          types.StringValue("PC-%INVALID%"),
			expectError:    true,
			errorSubstring: "Only %SERIAL% and %RAND:x% macros are supported",
		},
		{
			name:           "invalid RAND macro format",
			value:          types.StringValue("PC-%RAND%"),
			expectError:    true,
			errorSubstring: "Only %SERIAL% and %RAND:x% macros are supported",
		},
		{
			name:           "invalid RAND macro (no digits)",
			value:          types.StringValue("PC-%RAND:abc%"),
			expectError:    true,
			errorSubstring: "Invalid %RAND:x% macro",
		},
		{
			name:        "valid with letters and numbers",
			value:       types.StringValue("PC-01A"),
			expectError: false,
		},
		{
			name:        "valid all uppercase",
			value:       types.StringValue("DESKTOP"),
			expectError: false,
		},
		{
			name:        "valid all lowercase",
			value:       types.StringValue("desktop"),
			expectError: false,
		},
		{
			name:        "valid mixed case with hyphens",
			value:       types.StringValue("DeskTop-01"),
			expectError: false,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			request := validator.StringRequest{
				ConfigValue: testCase.value,
			}
			response := validator.StringResponse{}
			DeviceNameTemplate().ValidateString(context.Background(), request, &response)

			if testCase.expectError {
				if !response.Diagnostics.HasError() {
					t.Fatalf("expected error but got none")
				}
				if testCase.errorSubstring != "" {
					found := false
					for _, diag := range response.Diagnostics.Errors() {
						if len(diag.Detail()) > 0 && contains(diag.Detail(), testCase.errorSubstring) {
							found = true
							break
						}
						if len(diag.Summary()) > 0 && contains(diag.Summary(), testCase.errorSubstring) {
							found = true
							break
						}
					}
					if !found {
						t.Fatalf("expected error containing '%s' but got: %v", testCase.errorSubstring, response.Diagnostics.Errors())
					}
				}
			} else {
				if response.Diagnostics.HasError() {
					t.Fatalf("expected no error but got: %v", response.Diagnostics.Errors())
				}
			}
		})
	}
}

// contains checks if a string contains a substring
func contains(s, substr string) bool {
	if len(substr) == 0 {
		return true
	}
	if len(s) < len(substr) {
		return false
	}
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}