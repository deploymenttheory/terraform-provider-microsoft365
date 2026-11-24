package graphBetaApplicationsOnPremisesIpApplicationSegment

import (
	"context"

	"github.com/hashicorp/terraform-plugin-testing/terraform"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/exists"
	customrequests "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/custom_requests"
)

// OnPremisesIpApplicationSegmentTestResource implements the types.TestResource interface for IP application segments
type OnPremisesIpApplicationSegmentTestResource struct{}

// Exists checks whether the IP application segment exists in Microsoft Graph
func (r OnPremisesIpApplicationSegmentTestResource) Exists(ctx context.Context, _ any, state *terraform.InstanceState) (*bool, error) {
	//nolint:wrapcheck // Direct pass-through to generic helper with contextual errors
	return exists.CheckResourceExists(ctx, state, func(client *msgraphbetasdk.GraphServiceClient, ctx context.Context, state *terraform.InstanceState) error {
		getConfig := customrequests.GetRequestConfig{
			APIVersion:        customrequests.GraphAPIBeta,
			Endpoint:          "applications/" + state.Attributes["application_object_id"] + segmentsEndpointSuffix,
			ResourceIDPattern: customrequests.ResourceIDPatternSegment,
			ResourceID:        state.ID,
		}
		_, err := customrequests.GetRequestByResourceId(ctx, client.GetAdapter(), getConfig)
		return err
	})
}
