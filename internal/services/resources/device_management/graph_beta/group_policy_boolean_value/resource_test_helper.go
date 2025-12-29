package graphBetaGroupPolicyBooleanValue

import (
	"context"
	"fmt"
	"net/http"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance"
	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/generic_client"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// GroupPolicyBooleanValueTestResource implements the types.TestResource interface for Group Policy Boolean Values
type GroupPolicyBooleanValueTestResource struct{}

// Exists checks whether the Group Policy Boolean Value exists in Microsoft Graph
func (r GroupPolicyBooleanValueTestResource) Exists(ctx context.Context, _ any, state *terraform.InstanceState) (*bool, error) {
	httpClient, err := acceptance.TestHTTPClient()
	if err != nil {
		return nil, fmt.Errorf("failed to get HTTP client: %w", err)
	}

	groupPolicyConfigurationID := state.Attributes["group_policy_configuration_id"]
	definitionValueID := state.Attributes["group_policy_definition_value_id"]

	if groupPolicyConfigurationID == "" || definitionValueID == "" {
		exists := false
		return &exists, nil
	}

	url := httpClient.GetBaseURL() + "/deviceManagement/groupPolicyConfigurations/" + groupPolicyConfigurationID + "/definitionValues/" + definitionValueID

	httpReq, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	httpResp, err := httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to execute HTTP request: %w", err)
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode == http.StatusNotFound {
		exists := false
		return &exists, nil
	}

	if httpResp.StatusCode != http.StatusOK {
		errorInfo := errors.ExtractHTTPGraphError(ctx, httpResp)
		if errorInfo.ErrorCode == "ResourceNotFound" ||
			errorInfo.ErrorCode == "Request_ResourceNotFound" ||
			errorInfo.ErrorCode == "ItemNotFound" {
			exists := false
			return &exists, nil
		}

		return nil, fmt.Errorf("unexpected error checking group policy boolean value existence: %s (status: %d)", errorInfo.ErrorCode, errorInfo.StatusCode)
	}

	exists := true
	return &exists, nil
}

