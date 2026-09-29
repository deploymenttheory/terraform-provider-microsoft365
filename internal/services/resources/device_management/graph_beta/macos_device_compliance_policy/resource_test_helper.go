package graphBetaMacosDeviceCompliancePolicy

import (
	"context"

	"github.com/hashicorp/terraform-plugin-testing/terraform"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/exists"
)

// MacosDeviceCompliancePolicyTestResource implements the types.TestResource interface for macOS device compliance policies
type MacosDeviceCompliancePolicyTestResource struct{}

// Exists checks whether the macOS device compliance policy exists in Microsoft Graph
func (r MacosDeviceCompliancePolicyTestResource) Exists(
	ctx context.Context,
	_ any,
	state *terraform.InstanceState,
) (*bool, error) {
	//nolint:wrapcheck // Direct pass-through to generic helper with contextual errors
	return exists.CheckResourceExists(
		ctx,
		state,
		func(client *msgraphbetasdk.GraphServiceClient, ctx context.Context, state *terraform.InstanceState) error {
			_, err := client.DeviceManagement().
				DeviceCompliancePolicies().
				ByDeviceCompliancePolicyId(state.ID).
				Get(ctx, nil)
			return err
		},
	)
}
