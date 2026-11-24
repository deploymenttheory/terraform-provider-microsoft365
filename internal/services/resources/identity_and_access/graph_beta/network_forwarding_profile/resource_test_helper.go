package graphBetaNetworkForwardingProfile

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-testing/terraform"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/exists"
	customrequests "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/custom_requests"
)

// NetworkForwardingProfileTestResource provides resource existence checks for acceptance tests.
type NetworkForwardingProfileTestResource struct{}

// Exists checks whether the resource exists in Microsoft Graph.
func (r NetworkForwardingProfileTestResource) Exists(
	ctx context.Context,
	_ any,
	state *terraform.InstanceState,
) (*bool, error) {
	//nolint:wrapcheck // The shared existence helper supplies operation context.
	return exists.CheckResourceExists(
		ctx,
		state,
		func(client *msgraphbetasdk.GraphServiceClient, ctx context.Context, state *terraform.InstanceState) error {
			r := NewNetworkForwardingProfileResource().(*NetworkForwardingProfileResource)
			getConfig := customrequests.GetRequestConfig{
				Endpoint:          r.ResourcePath,
				ResourceIDPattern: customrequests.ResourceIDPatternSegment,
				ResourceID:        state.ID,
			}
			if _, err := customrequests.GetRequestByResourceId(ctx, client.GetAdapter(), getConfig); err != nil {
				return fmt.Errorf("get forwarding profile %s: %w", state.ID, err)
			}
			return nil
		},
	)
}
