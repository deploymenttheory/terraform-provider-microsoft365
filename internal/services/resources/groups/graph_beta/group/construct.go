package graphBetaGroup

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// constructResource maps the Terraform schema to the SDK model
// isUpdate indicates whether this is an update operation (PATCH) or create operation (POST)
func constructResource(ctx context.Context, data *GroupResourceModel, isUpdate bool) (graphmodels.Groupable, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing %s resource from model (isUpdate: %t)", ResourceName, isUpdate))

	requestBody := graphmodels.NewGroup()

	convert.FrameworkToGraphString(data.DisplayName, requestBody.SetDisplayName)
	convert.FrameworkToGraphString(data.Description, requestBody.SetDescription)
	convert.FrameworkToGraphString(data.MailNickname, requestBody.SetMailNickname)
	convert.FrameworkToGraphBool(data.MailEnabled, requestBody.SetMailEnabled)
	convert.FrameworkToGraphBool(data.SecurityEnabled, requestBody.SetSecurityEnabled)

	if err := convert.FrameworkToGraphStringSet(ctx, data.GroupTypes, requestBody.SetGroupTypes); err != nil {
		return nil, fmt.Errorf("failed to set group types: %s", err)
	}

	convert.FrameworkToGraphString(data.Visibility, requestBody.SetVisibility)
	convert.FrameworkToGraphBool(data.IsAssignableToRole, requestBody.SetIsAssignableToRole)
	convert.FrameworkToGraphString(data.MembershipRule, requestBody.SetMembershipRule)
	convert.FrameworkToGraphString(data.MembershipRuleProcessingState, requestBody.SetMembershipRuleProcessingState)

	// Set owners and members using additionalData with OData bind (only during creation)
	// A maximum of 20 relationships (owners + members) can be added during group creation
	if !isUpdate {
		additionalData := make(map[string]any)

		// Add owners if specified
		if !data.GroupOwners.IsNull() && !data.GroupOwners.IsUnknown() {
			var owners []string
			diags := data.GroupOwners.ElementsAs(ctx, &owners, false)
			if diags.HasError() {
				return nil, fmt.Errorf("failed to extract group_owners: %v", diags.Errors())
			}
			if len(owners) > 0 {
				ownerUrls := make([]string, len(owners))
				for i, ownerId := range owners {
					ownerUrls[i] = fmt.Sprintf("https://graph.microsoft.com/beta/users/%s", ownerId)
				}
				additionalData["owners@odata.bind"] = ownerUrls
				tflog.Debug(ctx, fmt.Sprintf("Adding %d owners to group", len(owners)))
			}
		}

		// Add members if specified
		if !data.GroupMembers.IsNull() && !data.GroupMembers.IsUnknown() {
			var members []string
			diags := data.GroupMembers.ElementsAs(ctx, &members, false)
			if diags.HasError() {
				return nil, fmt.Errorf("failed to extract group_members: %v", diags.Errors())
			}
			if len(members) > 0 {
				memberUrls := make([]string, len(members))
				for i, memberId := range members {
					memberUrls[i] = fmt.Sprintf("https://graph.microsoft.com/beta/users/%s", memberId)
				}
				additionalData["members@odata.bind"] = memberUrls
				tflog.Debug(ctx, fmt.Sprintf("Adding %d members to group", len(members)))
			}
		}

		// Set additional data if we have any
		if len(additionalData) > 0 {
			requestBody.SetAdditionalData(additionalData)
		}
	}

	if err := constructors.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for resource %s", ResourceName), requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]any{
			"error": err.Error(),
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished constructing %s resource", ResourceName))

	return requestBody, nil
}
