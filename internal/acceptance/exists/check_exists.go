package exists

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance"
	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

// GraphClientCheckFunc is a function type that performs a resource existence check using the Graph SDK client.
// It should return an error if the resource doesn't exist (404) or if there's an API error.
type GraphClientCheckFunc func(client *msgraphbetasdk.GraphServiceClient, ctx context.Context, state *terraform.InstanceState) error

// CheckResourceExists checks if a resource exists using the Microsoft Graph SDK (Kiota).
//
// Parameters:
//   - ctx: Context for the request
//   - state: Terraform instance state containing the resource ID
//   - checkFunc: A function that performs the specific SDK API call to check existence
//
// Returns:
//   - *bool: Pointer to boolean indicating if the resource exists
//   - error: Error if the check fails (other than "not found")
//
// Example usage:
//
//	return exists.CheckResourceExists(ctx, state, func(client *msgraphbetasdk.GraphServiceClient, ctx context.Context, state *terraform.InstanceState) error {
//	    _, err := client.DeviceManagement().WindowsQualityUpdatePolicies().ByWindowsQualityUpdatePolicyId(state.ID).Get(ctx, nil)
//	    return err
//	})
func CheckResourceExists(
	ctx context.Context,
	state *terraform.InstanceState,
	checkFunc GraphClientCheckFunc,
) (*bool, error) {
	graphClient, err := acceptance.TestGraphClient()
	if err != nil {
		return nil, err
	}

	err = checkFunc(graphClient, ctx, state)

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
