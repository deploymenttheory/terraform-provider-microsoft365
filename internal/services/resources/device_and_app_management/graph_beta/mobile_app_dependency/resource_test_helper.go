package graphBetaMobileAppDependency

import (
	"context"

	"github.com/hashicorp/terraform-plugin-testing/terraform"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/exists"
)

// MobileAppDependencyTestResource implements the shared acceptance test interface.
type MobileAppDependencyTestResource struct{}

// Exists checks the actual relationship ID, including ID-only imports.
func (r MobileAppDependencyTestResource) Exists(
	ctx context.Context,
	_ any,
	state *terraform.InstanceState,
) (*bool, error) {
	//nolint:wrapcheck // The shared helper adds context to Graph errors.
	return exists.CheckResourceExists(
		ctx,
		state,
		func(client *msgraphbetasdk.GraphServiceClient, ctx context.Context, state *terraform.InstanceState) error {
			_, err := client.DeviceAppManagement().
				MobileAppRelationships().
				ByMobileAppRelationshipId(state.ID).
				Get(ctx, nil)
			return err
		},
	)
}
