package graphBetaServicePrincipalAppRoleAssignedTo

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/exists"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

// ServicePrincipalAppRoleAssignedToTestResource implements the types.TestResource interface
type ServicePrincipalAppRoleAssignedToTestResource struct{}

// Exists checks whether the app role assignment exists in Microsoft Graph
func (r ServicePrincipalAppRoleAssignedToTestResource) Exists(ctx context.Context, _ any, state *terraform.InstanceState) (*bool, error) {
	//nolint:wrapcheck // Direct pass-through to generic helper with contextual errors
	return exists.CheckResourceExistsByCompositeID(
		ctx,
		state,
		"resource_object_id",
		func(client *msgraphbetasdk.GraphServiceClient, ctx context.Context, attributeValue string, resourceID string) error {
			_, err := client.ServicePrincipals().ByServicePrincipalId(attributeValue).AppRoleAssignedTo().ByAppRoleAssignmentId(resourceID).Get(ctx, nil)
			return err
		},
	)
}
