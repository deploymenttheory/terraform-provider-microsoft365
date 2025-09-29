package graphBetaGroupMemberAssignment

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// constructResource constructs a reference object for adding a member to a group
// It first validates that the member type is compatible with the target group type
func constructResource(ctx context.Context, data *GroupMemberAssignmentResourceModel, client *msgraphbetasdk.GraphServiceClient) (graphmodels.ReferenceCreateable, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing %s resource from model", ResourceName))

	if err := ValidateGroupMemberAssignment(ctx, client, *data, false); err != nil {
		return nil, err
	}

	memberObjectType := data.MemberObjectType.ValueString()
	memberId := data.MemberID.ValueString()

	requestBody := graphmodels.NewReferenceCreate()

	// Create the @odata.id reference URL based on the member object type
	var odataId string

	switch memberObjectType {
	case "User":
		// For users, we can use either directoryObjects or users endpoint
		// Using directoryObjects as it's more general and works for both security and M365 groups
		odataId = fmt.Sprintf("https://graph.microsoft.com/beta/directoryObjects/%s", memberId)
	case "Group":
		// For groups, use the groups endpoint
		odataId = fmt.Sprintf("https://graph.microsoft.com/beta/groups/%s", memberId)
	case "Device":
		// For devices, use the devices endpoint
		odataId = fmt.Sprintf("https://graph.microsoft.com/beta/devices/%s", memberId)
	case "ServicePrincipal":
		// For service principals, use the servicePrincipals endpoint
		odataId = fmt.Sprintf("https://graph.microsoft.com/beta/servicePrincipals/%s", memberId)
	case "OrganizationalContact":
		// For organizational contacts, use the contacts endpoint
		odataId = fmt.Sprintf("https://graph.microsoft.com/beta/contacts/%s", memberId)
	default:
		return nil, fmt.Errorf("unsupported member object type: %s", memberObjectType)
	}

	requestBody.SetOdataId(&odataId)

	tflog.Debug(ctx, fmt.Sprintf("Constructed member reference with @odata.id: %s for member type: %s", odataId, memberObjectType))

	if err := constructors.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for resource %s", ResourceName), requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]any{
			"error": err.Error(),
		})
	}
	tflog.Debug(ctx, fmt.Sprintf("Finished constructing %s resource", ResourceName))

	return requestBody, nil
}
