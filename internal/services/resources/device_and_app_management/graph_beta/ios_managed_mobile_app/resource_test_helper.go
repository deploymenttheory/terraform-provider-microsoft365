package graphBetaDeviceAndAppManagementIOSManagedMobileApp

import (
	"context"

	"github.com/hashicorp/terraform-plugin-testing/terraform"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/exists"
)

// IOSManagedMobileAppTestResource lets acceptance tests check whether the app is
// still targeted by its policy. A deleted policy counts as the app not existing.
type IOSManagedMobileAppTestResource struct{}

func (r IOSManagedMobileAppTestResource) Exists(
	ctx context.Context,
	_ any,
	state *terraform.InstanceState,
) (*bool, error) {
	//nolint:wrapcheck // Direct pass-through to generic helper with contextual errors
	return exists.CheckResourceExistsByArrayMembership(
		ctx,
		state,
		"managed_app_protection_id",
		"apps",
		"id",
		"id",
		func(client *msgraphbetasdk.GraphServiceClient, ctx context.Context, policyID string) (any, error) {
			response, err := client.
				DeviceAppManagement().
				IosManagedAppProtections().
				ByIosManagedAppProtectionId(policyID).
				Apps().
				Get(ctx, nil)
			if err != nil {
				return nil, err
			}
			apps := make([]map[string]any, 0, len(response.GetValue()))
			for _, app := range response.GetValue() {
				if app != nil && app.GetId() != nil {
					apps = append(apps, map[string]any{"id": *app.GetId()})
				}
			}
			return map[string]any{"apps": apps}, nil
		},
	)
}
