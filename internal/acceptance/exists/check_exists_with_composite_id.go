package exists

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance"
	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

// GraphClientCompositeIDCheckFunc is a function type that performs a resource existence check using the Graph SDK client
// with a composite ID pattern. It receives both the attribute value from state and the resource ID.
type GraphClientCompositeIDCheckFunc func(client *msgraphbetasdk.GraphServiceClient, ctx context.Context, attributeValue string, resourceID string) error

// CheckResourceExistsByCompositeID checks if a resource exists using a composite ID pattern with the GraphClient.
//
// Example: For a resource where state contains blueprint_id="abc-123" and ID="cred-456", the callback can construct:
// client.Applications().ByApplicationId(attributeValue).FederatedIdentityCredentials().ByFederatedIdentityCredentialId(resourceID).Get(ctx, nil)
//
// Parameters:
//   - ctx: Context for the request
//   - state: Terraform instance state containing the resource ID and attributes
//   - attributeName: The name of the state attribute to extract (e.g., "blueprint_id", "resource_object_id")
//   - checkFunc: A function that performs the specific SDK API call using the attribute value and resource ID
//
// Returns:
//   - *bool: Pointer to boolean indicating if the resource exists
//   - error: Error if the check fails (other than "not found") or if attribute is missing
//
// Example usage:
//
//	return exists.CheckResourceExistsByCompositeID(ctx, state, "blueprint_id", func(client *msgraphbetasdk.GraphServiceClient, ctx context.Context, attributeValue string, resourceID string) error {
//	    _, err := client.Applications().ByApplicationId(attributeValue).FederatedIdentityCredentials().ByFederatedIdentityCredentialId(resourceID).Get(ctx, nil)
//	    return err
//	})
func CheckResourceExistsByCompositeID(
	ctx context.Context,
	state *terraform.InstanceState,
	attributeName string,
	checkFunc GraphClientCompositeIDCheckFunc,
) (*bool, error) {
	graphClient, err := acceptance.TestGraphClient()
	if err != nil {
		return nil, err
	}

	attributeValue := state.Attributes[attributeName]
	if attributeValue == "" {
		return nil, fmt.Errorf("%s not found in state", attributeName)
	}

	resourceID := state.ID

	err = checkFunc(graphClient, ctx, attributeValue, resourceID)
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
