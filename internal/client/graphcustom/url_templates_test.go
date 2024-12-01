package graphcustom

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestByIDRequestUrlTemplate(t *testing.T) {
	tests := []struct {
		name     string
		config   GetRequestConfig
		expected string
	}{
		{
			name: "basic endpoint with ID pattern",
			config: GetRequestConfig{
				Endpoint:          "/deviceManagement/configurationPolicies",
				ResourceIDPattern: "('id')",
				ResourceID:        "12345",
			},
			expected: "{+baseurl}/deviceManagement/configurationPolicies('12345')",
		},
		{
			name: "endpoint with ID pattern and suffix",
			config: GetRequestConfig{
				Endpoint:          "/deviceManagement/configurationPolicies",
				ResourceIDPattern: "('id')",
				ResourceID:        "12345",
				EndpointSuffix:    "/settings",
			},
			expected: "{+baseurl}/deviceManagement/configurationPolicies('12345')/settings",
		},
		{
			name: "endpoint without leading slash",
			config: GetRequestConfig{
				Endpoint:          "deviceManagement/configurationPolicies",
				ResourceIDPattern: "('id')",
				ResourceID:        "12345",
			},
			expected: "{+baseurl}deviceManagement/configurationPolicies('12345')",
		},
		{
			name: "endpoint with different ID pattern format",
			config: GetRequestConfig{
				Endpoint:          "/users",
				ResourceIDPattern: "(id)",
				ResourceID:        "user@contoso.com",
			},
			expected: "{+baseurl}/users(user@contoso.com)",
		},
		{
			name: "complex ID with special characters",
			config: GetRequestConfig{
				Endpoint:          "/deviceManagement/configurationPolicies",
				ResourceIDPattern: "('id')",
				ResourceID:        "12345-67890-abcdef",
			},
			expected: "{+baseurl}/deviceManagement/configurationPolicies('12345-67890-abcdef')",
		},
		{
			name: "endpoint with multiple path segments and suffix",
			config: GetRequestConfig{
				Endpoint:          "/users/mailFolders/messages",
				ResourceIDPattern: "('id')",
				ResourceID:        "ABC123",
				EndpointSuffix:    "/attachments",
			},
			expected: "{+baseurl}/users/mailFolders/messages('ABC123')/attachments",
		},
		{
			name: "empty suffix",
			config: GetRequestConfig{
				Endpoint:          "/deviceManagement/configurationPolicies",
				ResourceIDPattern: "('id')",
				ResourceID:        "12345",
				EndpointSuffix:    "",
			},
			expected: "{+baseurl}/deviceManagement/configurationPolicies('12345')",
		},
		{
			name: "suffix without leading slash",
			config: GetRequestConfig{
				Endpoint:          "/deviceManagement/configurationPolicies",
				ResourceIDPattern: "('id')",
				ResourceID:        "12345",
				EndpointSuffix:    "settings",
			},
			expected: "{+baseurl}/deviceManagement/configurationPolicies('12345')settings",
		},
		{
			name: "GUID in resource ID",
			config: GetRequestConfig{
				Endpoint:          "/deviceManagement/configurationPolicies",
				ResourceIDPattern: "('id')",
				ResourceID:        "d557c813-b8e5-4efc-b00e-9c0bd5fd10df",
			},
			expected: "{+baseurl}/deviceManagement/configurationPolicies('d557c813-b8e5-4efc-b00e-9c0bd5fd10df')",
		},
		{
			name: "empty resource ID",
			config: GetRequestConfig{
				Endpoint:          "/deviceManagement/configurationPolicies",
				ResourceIDPattern: "('id')",
				ResourceID:        "",
			},
			expected: "{+baseurl}/deviceManagement/configurationPolicies('')",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ByIDRequestUrlTemplate(tt.config)
			assert.Equal(t, tt.expected, result, "URL template should match expected value")
		})
	}
}
