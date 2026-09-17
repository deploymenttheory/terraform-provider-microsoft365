package graphBetaLinuxDeviceCompliancePolicy

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/exists"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

type LinuxDeviceCompliancePolicyTestResource struct{}

func (r LinuxDeviceCompliancePolicyTestResource) Exists(ctx context.Context, _ any, state *terraform.InstanceState) (*bool, error) {
	//nolint:wrapcheck // Direct pass-through to the shared acceptance helper.
	return exists.CheckResourceExists(ctx, state, func(client *msgraphbetasdk.GraphServiceClient, ctx context.Context, state *terraform.InstanceState) error {
		_, err := client.DeviceManagement().CompliancePolicies().
			ByDeviceManagementCompliancePolicyId(state.ID).Get(ctx, nil)
		return err
	})
}
