package graphBetaGroupPolicyDefinition

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance"
	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/generic_client"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// GroupPolicyDefinitionTestResource implements the types.TestResource interface for Group Policy definitions
type GroupPolicyDefinitionTestResource struct{}

// Exists checks whether the Group Policy definition value exists in Microsoft Graph
func (r GroupPolicyDefinitionTestResource) Exists(ctx context.Context, _ any, state *terraform.InstanceState) (*bool, error) {
	httpClient, err := acceptance.TestHTTPClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get HTTP client: %w", err)
	}

	// The ID is in format: configID/definitionValueID
	idParts := strings.Split(state.ID, "/")
	if len(idParts) != 2 {
		return nil, fmt.Errorf("invalid ID format, expected configID/definitionValueID, got: %s", state.ID)
	}

	configID := idParts[0]
	defValueID := idParts[1]

	url := httpClient.GetBaseURL() + "/deviceManagement/groupPolicyConfigurations/" + configID + "/definitionValues/" + defValueID

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

		return nil, fmt.Errorf("unexpected error checking Group Policy definition value existence: %s (status: %d)", errorInfo.ErrorCode, errorInfo.StatusCode)
	}

	exists := true
	return &exists, nil
}
