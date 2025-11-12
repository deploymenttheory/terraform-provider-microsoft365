package utilityEntraIdSidConverter

import (
	"context"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestSidRidRangeValidator(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		value         types.String
		expectError   bool
		errorContains string
	}{
		"valid-all-rids-within-range": {
			value:       types.StringValue("S-1-12-1-1943430372-1249052806-2496021943-3034400218"),
			expectError: false,
		},
		"valid-maximum-uint32-values": {
			value:       types.StringValue("S-1-12-1-4294967295-4294967295-4294967295-4294967295"),
			expectError: false,
		},
		"valid-zero-values": {
			value:       types.StringValue("S-1-12-1-0-0-0-0"),
			expectError: false,
		},
		"invalid-rid-exceeds-uint32-max": {
			value:         types.StringValue("S-1-12-1-1234567890-9876543210-1111111111-2222222222"),
			expectError:   true,
			errorContains: "exceeds the maximum uint32 value",
		},
		"invalid-first-rid-exceeds-uint32-max": {
			value:         types.StringValue("S-1-12-1-4294967296-1000000000-2000000000-3000000000"),
			expectError:   true,
			errorContains: "exceeds the maximum uint32 value",
		},
		"invalid-last-rid-exceeds-uint32-max": {
			value:         types.StringValue("S-1-12-1-1000000000-2000000000-3000000000-5000000000"),
			expectError:   true,
			errorContains: "exceeds the maximum uint32 value",
		},
		"null-value": {
			value:       types.StringNull(),
			expectError: false,
		},
		"unknown-value": {
			value:       types.StringUnknown(),
			expectError: false,
		},
	}

	for name, tc := range testCases {
		name, tc := name, tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			request := validator.StringRequest{
				Path:        path.Root("sid"),
				ConfigValue: tc.value,
			}
			response := validator.StringResponse{}

			ValidateSidRidRange().ValidateString(context.Background(), request, &response)

			if tc.expectError {
				if !response.Diagnostics.HasError() {
					t.Fatal("expected error, got no error")
				}
				if tc.errorContains != "" {
					found := false
					for _, diag := range response.Diagnostics.Errors() {
						if strings.Contains(diag.Detail(), tc.errorContains) {
							found = true
							break
						}
					}
					if !found {
						t.Fatalf("expected error containing '%s', but got: %s", tc.errorContains, response.Diagnostics)
					}
				}
			} else {
				if response.Diagnostics.HasError() {
					t.Fatalf("got unexpected error: %s", response.Diagnostics)
				}
			}
		})
	}
}
