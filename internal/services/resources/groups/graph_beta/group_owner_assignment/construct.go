package graphBetaGroupOwnerAssignment

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// constructResource constructs a reference object for adding an owner to a group
func constructResource(ctx context.Context, data *GroupOwnerAssignmentResourceModel, client *msgraphbetasdk.GraphServiceClient) (graphmodels.ReferenceCreateable, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing %s resource from model", ResourceName))

	if err := ValidateGroupOwnerAssignment(ctx, client, *data, false); err != nil {
		return nil, err
	}

	ownerObjectType := data.OwnerObjectType.ValueString()
	ownerId := data.OwnerID.ValueString()

	requestBody := graphmodels.NewReferenceCreate()

	// Create the @odata.id reference URL based on the owner object type
	var odataId string

	switch ownerObjectType {
	case "User":
		// For users, we use the users endpoint
		odataId = fmt.Sprintf("https://graph.microsoft.com/beta/users/%s", ownerId)
	case "ServicePrincipal":
		// For service principals, use the servicePrincipals endpoint
		odataId = fmt.Sprintf("https://graph.microsoft.com/beta/servicePrincipals/%s", ownerId)
	default:
		return nil, fmt.Errorf("unsupported owner object type: %s", ownerObjectType)
	}

	requestBody.SetOdataId(&odataId)

	tflog.Debug(ctx, fmt.Sprintf("Constructed owner reference with @odata.id: %s for owner type: %s", odataId, ownerObjectType))

	if err := constructors.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for resource %s", ResourceName), requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]any{
			"error": err.Error(),
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished constructing %s resource", ResourceName))

	return requestBody, nil
}
