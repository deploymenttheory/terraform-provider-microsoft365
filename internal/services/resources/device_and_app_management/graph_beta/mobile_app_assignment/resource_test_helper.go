package graphBetaDeviceAndAppManagementAppAssignment

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/exists"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

// MobileAppAssignmentTestResource implements the types.TestResource interface for mobile app assignments
type MobileAppAssignmentTestResource struct{}

// Exists checks whether the mobile app assignment exists in Microsoft Graph.
//
// An assignment is addressed by the app that owns it as well as its own id, so the check
// reads mobile_app_id from state alongside the resource id.
func (r MobileAppAssignmentTestResource) Exists(ctx context.Context, _ any, state *terraform.InstanceState) (*bool, error) {
	//nolint:wrapcheck // Direct pass-through to generic helper with contextual errors
	return exists.CheckResourceExistsByCompositeID(ctx, state, "mobile_app_id", func(client *msgraphbetasdk.GraphServiceClient, ctx context.Context, mobileAppID string, assignmentID string) error {
		_, err := client.
			DeviceAppManagement().
			MobileApps().
			ByMobileAppId(mobileAppID).
			Assignments().
			ByMobileAppAssignmentId(assignmentID).
			Get(ctx, nil)

		return err
	})
}
