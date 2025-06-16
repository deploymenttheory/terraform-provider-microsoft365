package validators

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestStringListAllowedValues(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		values        []string
		list          []attr.Value
		expectedError bool
	}{
		"valid-single-value": {
			values: []string{"mdm", "windows10XManagement"},
			list: []attr.Value{
				types.StringValue("mdm"),
			},
			expectedError: false,
		},
		"valid-multiple-values": {
			values: []string{"mdm", "windows10XManagement", "configManager"},
			list: []attr.Value{
				types.StringValue("mdm"),
				types.StringValue("configManager"),
			},
			expectedError: false,
		},
		"invalid-single-value": {
			values: []string{"mdm", "windows10XManagement"},
			list: []attr.Value{
				types.StringValue("invalid"),
			},
			expectedError: true,
		},
		"invalid-one-of-multiple-values": {
			values: []string{"mdm", "windows10XManagement"},
			list: []attr.Value{
				types.StringValue("mdm"),
				types.StringValue("invalid"),
			},
			expectedError: true,
		},
		"empty-list": {
			values:        []string{"mdm", "windows10XManagement"},
			list:          []attr.Value{},
			expectedError: false,
		},
		"null-element": {
			values: []string{"mdm", "windows10XManagement"},
			list: []attr.Value{
				types.StringNull(),
			},
			expectedError: false,
		},
		"unknown-element": {
			values: []string{"mdm", "windows10XManagement"},
			list: []attr.Value{
				types.StringUnknown(),
			},
			expectedError: false,
		},
		"technology-values": {
			values: []string{
				"none", "mdm", "windows10XManagement", "configManager",
				"intuneManagementExtension", "thirdParty", "documentGateway",
				"appleRemoteManagement", "microsoftSense", "exchangeOnline",
				"mobileApplicationManagement", "linuxMdm", "enrollment",
				"endpointPrivilegeManagement", "unknownFutureValue",
				"windowsOsRecovery", "android",
			},
			list: []attr.Value{
				types.StringValue("mdm"),
				types.StringValue("android"),
				types.StringValue("linuxMdm"),
			},
			expectedError: false,
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			request := validator.ListRequest{
				ConfigValue: types.ListValueMust(types.StringType, testCase.list),
				Path:        path.Root("test"),
			}

			response := validator.ListResponse{}

			StringListAllowedValues(testCase.values...).ValidateList(context.Background(), request, &response)

			if !response.Diagnostics.HasError() && testCase.expectedError {
				t.Fatal("expected error, got no error")
			}

			if response.Diagnostics.HasError() && !testCase.expectedError {
				t.Fatalf("got unexpected error: %s", response.Diagnostics)
			}
		})
	}
}

func TestAllowedValuesListValidator_Description(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		values   []string
		expected string
		markdown bool
	}{
		"basic": {
			values:   []string{"value1", "value2"},
			expected: `each value in the list must be one of: ["value1" "value2"]`,
			markdown: false,
		},
		"markdown": {
			values:   []string{"value1", "value2"},
			expected: `each value in the list must be one of: ["value1" "value2"]`,
			markdown: true,
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			validator := allowedValuesListValidator{
				allowedValues: testCase.values,
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
