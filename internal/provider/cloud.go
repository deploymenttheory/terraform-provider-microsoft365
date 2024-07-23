package provider

import (
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
)

// setCloudConstants returns the OAuth authority URL and Graph API scope based on the provided cloud type.
func setCloudConstants(cloud string) (string, string, error) {
	switch cloud {
	case "public":
		return constants.PUBLIC_OAUTH_AUTHORITY_URL, constants.PUBLIC_GRAPH_API_SCOPE, nil
	case "dod":
		return constants.USDOD_OAUTH_AUTHORITY_URL, constants.USDOD_GRAPH_API_SCOPE, nil
	case "gcc":
		return constants.USGOV_OAUTH_AUTHORITY_URL, constants.USGOV_GRAPH_API_SCOPE, nil
	case "gcchigh":
		return constants.USGOVHIGH_OAUTH_AUTHORITY_URL, constants.USGOVHIGH_GRAPH_API_SCOPE, nil
	case "china":
		return constants.CHINA_OAUTH_AUTHORITY_URL, constants.CHINA_GRAPH_API_SCOPE, nil
	case "ex":
		return constants.EX_OAUTH_AUTHORITY_URL, constants.EX_GRAPH_API_SCOPE, nil
	case "rx":
		return constants.RX_OAUTH_AUTHORITY_URL, constants.RX_GRAPH_API_SCOPE, nil
	default:
		return "", "", fmt.Errorf("unsupported microsoft cloud type '%s'", cloud)
	}
}
