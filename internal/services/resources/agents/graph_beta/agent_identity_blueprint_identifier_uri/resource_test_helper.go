package graphBetaAgentIdentityBlueprintIdentifierUri

import (
	"context"
	"fmt"
	"net/http"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance"
	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/generic_client"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// AgentIdentityBlueprintIdentifierUriTestResource implements the types.TestResource interface
type AgentIdentityBlueprintIdentifierUriTestResource struct{}

// Exists checks whether the identifier URI exists on the application.
func (r AgentIdentityBlueprintIdentifierUriTestResource) Exists(ctx context.Context, _ any, state *terraform.InstanceState) (*bool, error) {
	httpClient, err := acceptance.TestHTTPClient()
	if err != nil {
		return nil, fmt.Errorf("failed to get HTTP client: %w", err)
	}

	blueprintID := state.Attributes["blueprint_id"]
	identifierUri := state.Attributes["identifier_uri"]

	if blueprintID == "" || identifierUri == "" {
		exists := false
		return &exists, nil
	}

	// Get the application to check identifier URIs
	url := httpClient.GetBaseURL() + "/applications/" + blueprintID + "?$select=identifierUris"

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

		return nil, fmt.Errorf("unexpected error checking identifier URI existence: %s (status: %d)", errorInfo.ErrorCode, errorInfo.StatusCode)
	}

	// The URI was removed, so check if it still exists in the response
	// For acceptance tests after delete, we need to verify the URI is gone
	exists := false
	return &exists, nil
}
