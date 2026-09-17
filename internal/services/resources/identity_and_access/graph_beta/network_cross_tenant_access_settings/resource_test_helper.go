package graphBetaNetworkCrossTenantAccessSettings

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

type NetworkCrossTenantAccessSettingsTestResource struct{}

func (r NetworkCrossTenantAccessSettingsTestResource) Exists(ctx context.Context, _ any, state *terraform.InstanceState) (*bool, error) {
	client, err := acceptance.TestGraphClient()
	if err != nil {
		return nil, fmt.Errorf("create acceptance Graph client: %w", err)
	}
	remote, err := client.NetworkAccess().Settings().CrossTenantAccess().Get(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("read network cross-tenant access settings: %w", err)
	}
	var model NetworkCrossTenantAccessSettingsResourceModel
	if err := mapRemoteState(&model, remote); err != nil {
		return nil, err
	}
	exists := state.ID == singletonID
	return &exists, nil
}
