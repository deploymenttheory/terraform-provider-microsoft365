package graphBetaVisionOSDeviceEnrollmentPolicy

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/exists"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

// VisionOSDeviceEnrollmentPolicyTestResource implements the types.TestResource interface for
// visionOS Automated Device Enrollment (ADE) profiles.
type VisionOSDeviceEnrollmentPolicyTestResource struct{}

// Exists checks whether the visionOS ADE enrollment profile exists in Microsoft Graph.
func (r VisionOSDeviceEnrollmentPolicyTestResource) Exists(ctx context.Context, _ any, state *terraform.InstanceState) (*bool, error) {
	//nolint:wrapcheck // Direct pass-through to generic helper with contextual errors
	return exists.CheckResourceExists(ctx, state, func(client *msgraphbetasdk.GraphServiceClient, ctx context.Context, state *terraform.InstanceState) error {
		_, err := client.DeviceManagement().ConfigurationPolicies().ByDeviceManagementConfigurationPolicyId(state.ID).Get(ctx, nil)
		return err
	})
}
