package exists

import (
	"context"
	"fmt"
	"strings"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance"
	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

// GraphClientSplitIDCheckFunc is a function type that performs a resource existence check using the Graph SDK client
// with split ID parts. It receives the split ID parts as a slice and should construct the appropriate SDK call.
type GraphClientSplitIDCheckFunc func(client *msgraphbetasdk.GraphServiceClient, ctx context.Context, idParts []string) error

// CheckResourceExistsBySplitID checks if a resource exists by splitting a composite ID and using the GraphClient.
//
// Example: state.ID = "configID/definitionValueID" can be used to call:
// client.DeviceManagement().GroupPolicyConfigurations().ByGroupPolicyConfigurationId(idParts[0]).DefinitionValues().ByGroupPolicyDefinitionValueId(idParts[1]).Get(ctx, nil)
//
// Parameters:
//   - ctx: Context for the request
//   - state: Terraform instance state with composite ID
//   - separator: The separator used in the ID (e.g., "/")
//   - expectedParts: Expected number of parts after splitting
//   - checkFunc: A function that performs the specific SDK API call using the split ID parts
//
// Returns:
//   - *bool: Pointer to boolean indicating if the resource exists
//   - error: Error if the check fails (other than "not found") or if ID format is invalid
//
// Example usage:
//
//	return exists.CheckResourceExistsBySplitID(ctx, state, "/", 2, func(client *msgraphbetasdk.GraphServiceClient, ctx context.Context, idParts []string) error {
//	    _, err := client.DeviceManagement().GroupPolicyConfigurations().ByGroupPolicyConfigurationId(idParts[0]).DefinitionValues().ByGroupPolicyDefinitionValueId(idParts[1]).Get(ctx, nil)
//	    return err
//	})
func CheckResourceExistsBySplitID(
	ctx context.Context,
	state *terraform.InstanceState,
	separator string,
	expectedParts int,
	checkFunc GraphClientSplitIDCheckFunc,
) (*bool, error) {
	graphClient, err := acceptance.TestGraphClient()
	if err != nil {
		return nil, err
	}

	idParts := strings.Split(state.ID, separator)
	if len(idParts) != expectedParts {
		return nil, fmt.Errorf("invalid ID format, expected %d parts separated by '%s', got: %s", expectedParts, separator, state.ID)
	}

	err = checkFunc(graphClient, ctx, idParts)
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
