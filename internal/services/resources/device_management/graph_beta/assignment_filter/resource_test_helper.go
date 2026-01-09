package graphBetaAssignmentFilter

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance"
	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/generic_client"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// AssignmentFilterTestResource implements the types.TestResource interface for assignment filters
type AssignmentFilterTestResource struct{}

// Exists checks whether the assignment filter exists in Microsoft Graph
func (r AssignmentFilterTestResource) Exists(ctx context.Context, _ any, state *terraform.InstanceState) (*bool, error) {
	httpClient, err := acceptance.TestHTTPClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get HTTP client: %w", err)
	}

	filterID := state.ID

	url := httpClient.GetBaseURL() + "/deviceManagement/assignmentFilters/" + filterID

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

		return nil, fmt.Errorf("unexpected error checking filter existence: %s (status: %d)", errorInfo.ErrorCode, errorInfo.StatusCode)
	}

	var filter map[string]any
	if err := json.NewDecoder(httpResp.Body).Decode(&filter); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if filter == nil {
		exists := false
		return &exists, nil
	}

	exists := true
	return &exists, nil
}
