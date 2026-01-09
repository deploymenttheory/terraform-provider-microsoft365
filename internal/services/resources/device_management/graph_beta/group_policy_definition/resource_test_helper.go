package graphBetaGroupPolicyDefinition

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/exists"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

// GroupPolicyDefinitionTestResource implements the types.TestResource interface for Group Policy definitions
type GroupPolicyDefinitionTestResource struct{}

// Exists checks whether the Group Policy definition value exists in Microsoft Graph
func (r GroupPolicyDefinitionTestResource) Exists(ctx context.Context, _ any, state *terraform.InstanceState) (*bool, error) {
	//nolint:wrapcheck // Direct pass-through to generic helper with contextual errors
	return exists.CheckResourceExistsBySplitID(
		ctx,
		state,
		"/",
		2,
		func(client *msgraphbetasdk.GraphServiceClient, ctx context.Context, idParts []string) error {
			_, err := client.DeviceManagement().GroupPolicyConfigurations().ByGroupPolicyConfigurationId(idParts[0]).DefinitionValues().ByGroupPolicyDefinitionValueId(idParts[1]).Get(ctx, nil)
			return err
		},
	)
}
