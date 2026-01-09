package graphBetaAgentIdentityBlueprintCertificateCredential

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance"
	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/generic_client"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// AgentIdentityBlueprintCertificateCredentialTestResource implements the types.TestResource interface
type AgentIdentityBlueprintCertificateCredentialTestResource struct{}

// Exists checks whether the certificate credential exists in the application's keyCredentials
func (r AgentIdentityBlueprintCertificateCredentialTestResource) Exists(ctx context.Context, _ any, state *terraform.InstanceState) (*bool, error) {
	httpClient, err := acceptance.TestHTTPClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get HTTP client: %w", err)
	}

	blueprintID := state.Attributes["blueprint_id"]
	keyID := state.Attributes["key_id"]

	if blueprintID == "" || keyID == "" {
		exists := false
		return &exists, nil
	}

	url := httpClient.GetBaseURL() + "/applications/" + blueprintID

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

		return nil, fmt.Errorf("unexpected error checking certificate credential existence: %s (status: %d)", errorInfo.ErrorCode, errorInfo.StatusCode)
	}

	// Parse the response to check for the keyCredential
	var application map[string]any
	if err := json.NewDecoder(httpResp.Body).Decode(&application); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	keyCredentials, ok := application["keyCredentials"].([]any)
	if !ok {
		exists := false
		return &exists, nil
	}

	// Check if the keyId exists in keyCredentials
	for _, cred := range keyCredentials {
		if credMap, ok := cred.(map[string]any); ok {
			if credKeyID, ok := credMap["keyId"].(string); ok && credKeyID == keyID {
				exists := true
				return &exists, nil
			}
		}
	}

	exists := false
	return &exists, nil
}
