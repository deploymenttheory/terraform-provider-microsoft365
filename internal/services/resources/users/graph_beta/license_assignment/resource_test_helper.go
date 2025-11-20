package graphBetaUserLicenseAssignment

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance"
	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/generic_client"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// UserLicenseAssignmentTestResource implements the types.TestResource interface for user license assignments
type UserLicenseAssignmentTestResource struct{}

// Exists checks whether the user license assignment exists in Microsoft Graph
func (r UserLicenseAssignmentTestResource) Exists(ctx context.Context, _ any, state *terraform.InstanceState) (*bool, error) {
	httpClient, err := acceptance.TestHTTPClient()
	if err != nil {
		return nil, fmt.Errorf("failed to get HTTP client: %w", err)
	}

	userID := state.ID

	url := httpClient.GetBaseURL() + "/users/" + userID + "/licenseDetails"

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

		return nil, fmt.Errorf("unexpected error checking user license existence: %s (status: %d)", errorInfo.ErrorCode, errorInfo.StatusCode)
	}

	var response map[string]any
	if err := json.NewDecoder(httpResp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// Check if the user has any license details
	value, ok := response["value"]
	if !ok {
		exists := false
		return &exists, nil
	}

	licenses, ok := value.([]any)
	if !ok || len(licenses) == 0 {
		exists := false
		return &exists, nil
	}

	exists := true
	return &exists, nil
}
