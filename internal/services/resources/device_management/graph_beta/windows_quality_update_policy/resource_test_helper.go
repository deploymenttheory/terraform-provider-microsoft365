package graphBetaWindowsQualityUpdatePolicy

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance"
	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// WindowsQualityUpdatePolicyTestResource implements the types.TestResource interface for Windows Quality Update Policies
type WindowsQualityUpdatePolicyTestResource struct{}

// Exists checks whether the Windows Quality Update Policy exists in Microsoft Graph
func (r WindowsQualityUpdatePolicyTestResource) Exists(ctx context.Context, _ any, state *terraform.InstanceState) (*bool, error) {
	graphClient, err := acceptance.TestGraphClient()
	if err != nil {
		return nil, err
	}

	policyID := state.ID

	_, err = graphClient.
		DeviceManagement().
		WindowsQualityUpdatePolicies().
		ByWindowsQualityUpdatePolicyId(policyID).
		Get(ctx, nil)

	if err != nil {
		errorInfo := errors.GraphError(ctx, err)
		if errorInfo.StatusCode == 404 ||
			errorInfo.ErrorCode == "ResourceNotFound" ||
			errorInfo.ErrorCode == "Request_ResourceNotFound" ||
			errorInfo.ErrorCode == "ItemNotFound" {
			exists := false
			return &exists, nil
		}
		return nil, err
	}

	exists := true
	return &exists, nil
}
