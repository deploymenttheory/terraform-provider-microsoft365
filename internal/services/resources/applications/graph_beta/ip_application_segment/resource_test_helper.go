package graphBetaApplicationsIpApplicationSegment

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/exists"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	abstractions "github.com/microsoft/kiota-abstractions-go"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/models/odataerrors"
)

// IpApplicationSegmentTestResource implements the types.TestResource interface for IP application segments
type IpApplicationSegmentTestResource struct{}

// buildSegmentItemPath constructs the URL path for a specific application segment.
func buildSegmentItemPath(applicationID, segmentID string) string {
	return fmt.Sprintf("applications/%s/onPremisesPublishing/segmentsConfiguration/microsoft.graph.ipSegmentConfiguration/applicationSegments/%s", applicationID, segmentID)
}

// Exists checks whether the IP application segment exists in Microsoft Graph
func (r IpApplicationSegmentTestResource) Exists(ctx context.Context, _ any, state *terraform.InstanceState) (*bool, error) {
	//nolint:wrapcheck // Direct pass-through to generic helper with contextual errors
	return exists.CheckResourceExists(ctx, state, func(client *msgraphbetasdk.GraphServiceClient, ctx context.Context, state *terraform.InstanceState) error {
		applicationID := state.Attributes["application_object_id"]
		segmentID := state.ID

		// Build request using SDK's RequestAdapter with manual path construction
		requestInfo := abstractions.NewRequestInformation()
		requestInfo.Method = abstractions.GET
		requestInfo.UrlTemplate = "{+baseurl}/" + buildSegmentItemPath(applicationID, segmentID)
		requestInfo.PathParameters = map[string]string{
			"baseurl": "https://graph.microsoft.com/beta",
		}

		errorMapping := abstractions.ErrorMappings{
			"XXX": odataerrors.CreateODataErrorFromDiscriminatorValue,
		}

		_, err := client.GetAdapter().Send(
			ctx,
			requestInfo,
			graphmodels.CreateIpApplicationSegmentFromDiscriminatorValue,
			errorMapping,
		)
		return err
	})
}
