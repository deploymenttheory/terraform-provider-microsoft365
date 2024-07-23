package provider

import (
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
)

// setCloudConstants returns the OAuth authority URL, Graph API scope, and Graph API service root based on the provided cloud type.
func setCloudConstants(cloud string) (string, string, string, string, error) {
	switch cloud {
	case "public":
		return constants.PUBLIC_OAUTH_AUTHORITY_URL,
			constants.PUBLIC_GRAPH_API_SCOPE,
			constants.PUBLIC_GRAPH_API_SERVICE_ROOT,
			constants.PUBLIC_GRAPH_BETA_API_SERVICE_ROOT, nil
	case "dod":
		return constants.USDOD_OAUTH_AUTHORITY_URL,
			constants.USDOD_GRAPH_API_SCOPE,
			constants.USDOD_GRAPH_API_SERVICE_ROOT,
			constants.USDOD_GRAPH_BETA_API_SERVICE_ROOT, nil
	case "gcc":
		return constants.USGOV_OAUTH_AUTHORITY_URL,
			constants.USGOV_GRAPH_API_SCOPE,
			constants.USGOV_GRAPH_API_SERVICE_ROOT,
			constants.USGOV_GRAPH_BETA_API_SERVICE_ROOT, nil
	case "gcchigh":
		return constants.USGOVHIGH_OAUTH_AUTHORITY_URL,
			constants.USGOVHIGH_GRAPH_API_SCOPE,
			constants.USGOVHIGH_GRAPH_API_SERVICE_ROOT,
			constants.USGOVHIGH_GRAPH_BETA_API_SERVICE_ROOT, nil
	case "china":
		return constants.CHINA_OAUTH_AUTHORITY_URL,
			constants.CHINA_GRAPH_API_SCOPE,
			constants.CHINA_GRAPH_API_SERVICE_ROOT,
			constants.CHINA_GRAPH_BETA_API_SERVICE_ROOT, nil
	case "ex":
		return constants.EX_OAUTH_AUTHORITY_URL,
			constants.EX_GRAPH_API_SCOPE,
			constants.EX_GRAPH_API_SERVICE_ROOT,
			constants.EX_GRAPH_BETA_API_SERVICE_ROOT, nil
	case "rx":
		return constants.RX_OAUTH_AUTHORITY_URL,
			constants.RX_GRAPH_API_SCOPE,
			constants.RX_GRAPH_API_SERVICE_ROOT,
			constants.RX_GRAPH_BETA_API_SERVICE_ROOT, nil
	default:
		return "", "", "", "", fmt.Errorf("unsupported microsoft cloud type '%s'", cloud)
	}
}
