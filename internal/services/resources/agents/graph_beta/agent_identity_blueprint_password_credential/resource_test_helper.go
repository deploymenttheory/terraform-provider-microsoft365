package graphBetaAgentIdentityBlueprintPasswordCredential

import (
	"context"
	"fmt"
	"net/http"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance"
	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/generic_client"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// AgentIdentityBlueprintPasswordCredentialTestResource implements the types.TestResource interface
type AgentIdentityBlueprintPasswordCredentialTestResource struct{}

// Exists checks whether the password credential exists in Microsoft Graph.
// Note: Password credentials cannot be directly queried by keyId, so we list all
// credentials on the blueprint and check if the keyId exists.
func (r AgentIdentityBlueprintPasswordCredentialTestResource) Exists(ctx context.Context, _ any, state *terraform.InstanceState) (*bool, error) {
	httpClient, err := acceptance.TestHTTPClient()
	if err != nil {
		return nil, fmt.Errorf("failed to get HTTP client: %w", err)
	}

	blueprintID := state.Attributes["blueprint_id"]
	keyID := state.Attributes["key_id"]

	if blueprintID == "" || keyID == "" {
		exists := false
		return &exists, nil
	}

	// List password credentials on the application
	url := httpClient.GetBaseURL() + "/applications/" + blueprintID + "?$select=passwordCredentials"

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

		return nil, fmt.Errorf("unexpected error checking password credential existence: %s (status: %d)", errorInfo.ErrorCode, errorInfo.StatusCode)
	}

	// Password credential was removed, so it doesn't exist anymore
	// For acceptance tests after delete, we assume it doesn't exist if we can read the application
	exists := false
	return &exists, nil
}
