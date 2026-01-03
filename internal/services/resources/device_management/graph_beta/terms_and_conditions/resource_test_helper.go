package graphBetaTermsAndConditions

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance"
	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// TermsAndConditionsTestResource implements the types.TestResource interface for terms and conditions
type TermsAndConditionsTestResource struct{}

// Exists checks whether the terms and conditions exists in Microsoft Graph
func (r TermsAndConditionsTestResource) Exists(ctx context.Context, _ any, state *terraform.InstanceState) (*bool, error) {
	graphClient, err := acceptance.TestGraphClient()
	if err != nil {
		return nil, err
	}

	termsAndConditionsID := state.ID

	result, err := graphClient.
		DeviceManagement().
		TermsAndConditions().
		ByTermsAndConditionsId(termsAndConditionsID).
		Get(ctx, nil)

	if err != nil {
		errorInfo := errors.GraphError(ctx, err)
		if errorInfo.StatusCode == 404 ||
			errorInfo.ErrorCode == "ResourceNotFound" ||
			errorInfo.ErrorCode == "ItemNotFound" ||
			errorInfo.ErrorCode == "Request_ResourceNotFound" {
			exists := false
			return &exists, nil
		}
		return nil, err
	}

	exists := result != nil
	return &exists, nil
}
