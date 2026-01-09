package graphBetaUserLicenseAssignment

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/exists"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/users"
)

// UserLicenseAssignmentTestResource implements the types.TestResource interface for user license assignments
type UserLicenseAssignmentTestResource struct{}

// Exists checks whether the specific user license assignment exists in Microsoft Graph
func (r UserLicenseAssignmentTestResource) Exists(ctx context.Context, _ any, state *terraform.InstanceState) (*bool, error) {
	//nolint:wrapcheck // Direct pass-through to generic helper with contextual errors
	return exists.CheckResourceExistsByArrayMembership(
		ctx,
		state,
		"user_id",
		"assignedLicenses",
		"skuId",
		"sku_id",
		func(client *msgraphbetasdk.GraphServiceClient, ctx context.Context, parentID string) (any, error) {
			return client.Users().ByUserId(parentID).Get(ctx, &users.UserItemRequestBuilderGetRequestConfiguration{
				QueryParameters: &users.UserItemRequestBuilderGetQueryParameters{
					Select: []string{"id", "assignedLicenses"},
				},
			})
		},
	)
}
