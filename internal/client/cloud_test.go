package client

import (
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/stretchr/testify/assert"
)

func TestSetCloudConstants(t *testing.T) {
	testCases := []struct {
		name                    string
		cloud                   string
		expectedAuthority       string
		expectedScope           string
		expectedServiceRoot     string
		expectedBetaServiceRoot string
		expectedError           string
	}{
		{
			name:                    "Public Cloud",
			cloud:                   "public",
			expectedAuthority:       constants.PUBLIC_OAUTH_AUTHORITY_URL,
			expectedScope:           constants.PUBLIC_GRAPH_API_SCOPE,
			expectedServiceRoot:     constants.PUBLIC_GRAPH_API_SERVICE_ROOT,
			expectedBetaServiceRoot: constants.PUBLIC_GRAPH_BETA_API_SERVICE_ROOT,
		},
		{
			name:                    "DoD Cloud",
			cloud:                   "dod",
			expectedAuthority:       constants.USDOD_OAUTH_AUTHORITY_URL,
			expectedScope:           constants.USDOD_GRAPH_API_SCOPE,
			expectedServiceRoot:     constants.USDOD_GRAPH_API_SERVICE_ROOT,
			expectedBetaServiceRoot: constants.USDOD_GRAPH_BETA_API_SERVICE_ROOT,
		},
		{
			name:                    "GCC Cloud",
			cloud:                   "gcc",
			expectedAuthority:       constants.USGOV_OAUTH_AUTHORITY_URL,
			expectedScope:           constants.USGOV_GRAPH_API_SCOPE,
			expectedServiceRoot:     constants.USGOV_GRAPH_API_SERVICE_ROOT,
			expectedBetaServiceRoot: constants.USGOV_GRAPH_BETA_API_SERVICE_ROOT,
		},
		{
			name:                    "GCC High Cloud",
			cloud:                   "gcchigh",
			expectedAuthority:       constants.USGOVHIGH_OAUTH_AUTHORITY_URL,
			expectedScope:           constants.USGOVHIGH_GRAPH_API_SCOPE,
			expectedServiceRoot:     constants.USGOVHIGH_GRAPH_API_SERVICE_ROOT,
			expectedBetaServiceRoot: constants.USGOVHIGH_GRAPH_BETA_API_SERVICE_ROOT,
		},
		{
			name:                    "China Cloud",
			cloud:                   "china",
			expectedAuthority:       constants.CHINA_OAUTH_AUTHORITY_URL,
			expectedScope:           constants.CHINA_GRAPH_API_SCOPE,
			expectedServiceRoot:     constants.CHINA_GRAPH_API_SERVICE_ROOT,
			expectedBetaServiceRoot: constants.CHINA_GRAPH_BETA_API_SERVICE_ROOT,
		},
		{
			name:                    "EagleX Cloud",
			cloud:                   "ex",
			expectedAuthority:       constants.EX_OAUTH_AUTHORITY_URL,
			expectedScope:           constants.EX_GRAPH_API_SCOPE,
			expectedServiceRoot:     constants.EX_GRAPH_API_SERVICE_ROOT,
			expectedBetaServiceRoot: constants.EX_GRAPH_BETA_API_SERVICE_ROOT,
		},
		{
			name:                    "Secure Cloud (RX)",
			cloud:                   "rx",
			expectedAuthority:       constants.RX_OAUTH_AUTHORITY_URL,
			expectedScope:           constants.RX_GRAPH_API_SCOPE,
			expectedServiceRoot:     constants.RX_GRAPH_API_SERVICE_ROOT,
			expectedBetaServiceRoot: constants.RX_GRAPH_BETA_API_SERVICE_ROOT,
		},
		{
			name:          "Unsupported Cloud",
			cloud:         "unsupported",
			expectedError: "unsupported microsoft cloud type 'unsupported'",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			authority, scope, serviceRoot, betaServiceRoot, err := SetCloudConstants(tc.cloud)

			if tc.expectedError != "" {
				assert.EqualError(t, err, tc.expectedError)
				assert.Empty(t, authority)
				assert.Empty(t, scope)
				assert.Empty(t, serviceRoot)
				assert.Empty(t, betaServiceRoot)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedAuthority, authority)
				assert.Equal(t, tc.expectedScope, scope)
				assert.Equal(t, tc.expectedServiceRoot, serviceRoot)
				assert.Equal(t, tc.expectedBetaServiceRoot, betaServiceRoot)
			}
		})
	}
}
