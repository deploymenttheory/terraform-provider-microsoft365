package graphBetaGroupLicenseAssignment

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance"
	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/generic_client"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// GroupLicenseAssignmentTestResource implements the types.TestResource interface for group license assignments
type GroupLicenseAssignmentTestResource struct{}

// Exists checks whether the specific group license assignment exists in Microsoft Graph
func (r GroupLicenseAssignmentTestResource) Exists(ctx context.Context, _ any, state *terraform.InstanceState) (*bool, error) {
	httpClient, err := acceptance.TestHTTPClient()
	if err != nil {
		return nil, fmt.Errorf("failed to get HTTP client: %w", err)
	}

	// Extract group_id and sku_id from state
	groupID := state.Attributes["group_id"]
	skuID := state.Attributes["sku_id"]

	url := httpClient.GetBaseURL() + "/groups/" + groupID

	httpReq, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	// Add query parameter to select assignedLicenses
	q := httpReq.URL.Query()
	q.Add("$select", "id,assignedLicenses")
	httpReq.URL.RawQuery = q.Encode()

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

		return nil, fmt.Errorf("unexpected error checking group license existence: %s (status: %d)", errorInfo.ErrorCode, errorInfo.StatusCode)
	}

	var response map[string]any
	if err := json.NewDecoder(httpResp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// Check if the group has the specific license
	assignedLicensesRaw, ok := response["assignedLicenses"]
	if !ok {
		exists := false
		return &exists, nil
	}

	assignedLicenses, ok := assignedLicensesRaw.([]any)
	if !ok || len(assignedLicenses) == 0 {
		exists := false
		return &exists, nil
	}

	// Look for the specific SKU
	for _, license := range assignedLicenses {
		licenseMap, ok := license.(map[string]any)
		if !ok {
			continue
		}

		if licenseSkuID, ok := licenseMap["skuId"].(string); ok && licenseSkuID == skuID {
			exists := true
			return &exists, nil
		}
	}

	exists := false
	return &exists, nil
}
