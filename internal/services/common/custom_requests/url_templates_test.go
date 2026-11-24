package customrequests

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
			expected: "{+baseurl}/deviceManagement/configurationPolicies('12345')",
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
		{
			name: "slash brace pattern writes the ID literally",
			config: GetRequestConfig{
				Endpoint:          "/applications",
				ResourceIDPattern: "/{id}",
				ResourceID:        "d557c813-b8e5-4efc-b00e-9c0bd5fd10df",
			},
			expected: "{+baseurl}/applications/d557c813-b8e5-4efc-b00e-9c0bd5fd10df",
		},
		{
			name: "slash brace pattern with suffix",
			config: GetRequestConfig{
				Endpoint:          "applications/abc/owners",
				ResourceIDPattern: "/{id}",
				ResourceID:        "d557c813-b8e5-4efc-b00e-9c0bd5fd10df",
				EndpointSuffix:    "/$ref",
			},
			expected: "{+baseurl}/applications/abc/owners/d557c813-b8e5-4efc-b00e-9c0bd5fd10df/$ref",
		},
		{
			name: "quoted brace pattern",
			config: GetRequestConfig{
				Endpoint:          "/deviceManagement/configurationPolicies",
				ResourceIDPattern: "('{id}')",
				ResourceID:        "12345",
			},
			expected: "{+baseurl}/deviceManagement/configurationPolicies('12345')",
		},
		{
			name: "no pattern keeps the endpoint",
			config: GetRequestConfig{
				Endpoint: "/applications/abc/microsoft.graph.agentIdentityBlueprint/sponsors",
			},
			expected: "{+baseurl}/applications/abc/microsoft.graph.agentIdentityBlueprint/sponsors",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ByIDRequestUrlTemplate(tt.config)
			assert.Equal(t, tt.expected, result, "URL template should match expected value")
			assert.Equal(t, tt.expected, ByIDRequestUrlTemplate(DeleteRequestConfig{Endpoint: tt.config.Endpoint, ResourceIDPattern: tt.config.ResourceIDPattern, ResourceID: tt.config.ResourceID, EndpointSuffix: tt.config.EndpointSuffix}), "DELETE must build the same template")
			assert.Equal(t, tt.expected, ByIDRequestUrlTemplate(PatchRequestConfig{Endpoint: tt.config.Endpoint, ResourceIDPattern: tt.config.ResourceIDPattern, ResourceID: tt.config.ResourceID, EndpointSuffix: tt.config.EndpointSuffix}), "PATCH must build the same template")
		})
	}
}
