package graphBetaGroupLicenseAssignment

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/exists"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/groups"
)

// GroupLicenseAssignmentTestResource implements the types.TestResource interface for group license assignments
type GroupLicenseAssignmentTestResource struct{}

// Exists checks whether the specific group license assignment exists in Microsoft Graph
func (r GroupLicenseAssignmentTestResource) Exists(ctx context.Context, _ any, state *terraform.InstanceState) (*bool, error) {
	//nolint:wrapcheck // Direct pass-through to generic helper with contextual errors
	return exists.CheckResourceExistsByArrayMembership(
		ctx,
		state,
		"group_id",
		"assignedLicenses",
		"skuId",
		"sku_id",
		func(client *msgraphbetasdk.GraphServiceClient, ctx context.Context, parentID string) (any, error) {
			return client.Groups().ByGroupId(parentID).Get(ctx, &groups.GroupItemRequestBuilderGetRequestConfiguration{
				QueryParameters: &groups.GroupItemRequestBuilderGetQueryParameters{
					Select: []string{"id", "assignedLicenses"},
				},
			})
		},
	)
}
