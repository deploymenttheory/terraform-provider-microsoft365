package graphBetaServicePrincipalOwner

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/exists"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

// ServicePrincipalOwnerTestResource implements the types.TestResource interface
type ServicePrincipalOwnerTestResource struct{}

// Exists checks whether the service principal owner assignment exists in Microsoft Graph
func (r ServicePrincipalOwnerTestResource) Exists(ctx context.Context, _ any, state *terraform.InstanceState) (*bool, error) {
	//nolint:wrapcheck // Direct pass-through to generic helper with contextual errors
	return exists.CheckResourceExistsByCompositeID(
		ctx,
		state,
		"service_principal_id",
		func(client *msgraphbetasdk.GraphServiceClient, ctx context.Context, attributeValue string, resourceID string) error {
			// Get all owners and check if our specific owner exists
			owners, err := client.ServicePrincipals().ByServicePrincipalId(attributeValue).Owners().Get(ctx, nil)
			if err != nil {
				return err
			}

			if owners != nil && owners.GetValue() != nil {
				for _, owner := range owners.GetValue() {
					if owner.GetId() != nil && *owner.GetId() == resourceID {
						return nil // Owner found
					}
				}
			}

			return nil // Owner not found, but no error
		},
	)
}
