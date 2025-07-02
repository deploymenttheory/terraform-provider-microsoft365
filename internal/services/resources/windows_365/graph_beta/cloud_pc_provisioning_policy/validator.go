package graphBetaCloudPcProvisioningPolicy

import (
	"context"
	"fmt"

	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

// validateResource checks that the display_name is unique among all Cloud PC provisioning policies.
// Returns an error if a policy with the same display_name already exists.
// For updates, it excludes the current resource from the uniqueness check.
func validateResource(ctx context.Context, client *msgraphbetasdk.GraphServiceClient, data *CloudPcProvisioningPolicyResourceModel) error {
	if client == nil {
		return fmt.Errorf("microsoft Graph client is not available for uniqueness validation")
	}
	if data.DisplayName.IsNull() || data.DisplayName.IsUnknown() {
		return nil
	}
	displayName := data.DisplayName.ValueString()

	// Skip uniqueness check for updates - if the resource has an ID, it's an update
	isUpdate := !data.ID.IsNull() && !data.ID.IsUnknown() && data.ID.ValueString() != ""
	if isUpdate {
		return nil
	}

	policies, err := client.
		DeviceManagement().
		VirtualEndpoint().
		ProvisioningPolicies().
		Get(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to list Cloud PC provisioning policies for uniqueness check: %v", err)
	}
	if policies == nil || policies.GetValue() == nil {
		return nil
	}
	for _, policy := range policies.GetValue() {
		if policy == nil || policy.GetDisplayName() == nil {
			continue
		}
		if *policy.GetDisplayName() == displayName {
			return fmt.Errorf("a Cloud PC provisioning policy with display_name '%s' already exists. Policy names must be unique", displayName)
		}
	}
	return nil
}
