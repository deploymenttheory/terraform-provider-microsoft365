package graphBetaGroupMemberAssignment

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// MapRemoteStateToTerraform maps the properties of a DirectoryObject member to a Terraform state.
func MapRemoteStateToTerraform(ctx context.Context, data *GroupMemberAssignmentResourceModel, remoteResource graphmodels.DirectoryObjectable) {
	if remoteResource == nil {
		tflog.Debug(ctx, "Remote resource is nil")
		return
	}

	tflog.Debug(ctx, "Starting to map remote state to Terraform state", map[string]interface{}{
		"memberID": types.StringPointerValue(remoteResource.GetId()),
		"groupID":  data.GroupID.ValueString(),
	})

	// Set the member ID
	data.MemberID = types.StringPointerValue(remoteResource.GetId())

	// Create composite ID from group_id and member_id
	compositeID := fmt.Sprintf("%s/%s", data.GroupID.ValueString(), data.MemberID.ValueString())
	data.ID = types.StringValue(compositeID)

	// Set member type from @odata.type
	if odataType := remoteResource.GetOdataType(); odataType != nil {
		// Extract the type from @odata.type (e.g., "#microsoft.graph.user" -> "User")
		memberType := extractMemberTypeFromOdataType(*odataType)
		data.MemberType = types.StringValue(memberType)
		// Also set the MemberObjectType if it's not already set
		if data.MemberObjectType.IsNull() || data.MemberObjectType.IsUnknown() {
			data.MemberObjectType = types.StringValue(memberType)
		}
	}

	// Set display name if available
	// Try to get display name from different member types
	if user, ok := remoteResource.(graphmodels.Userable); ok {
		data.MemberDisplayName = types.StringPointerValue(user.GetDisplayName())
	} else if group, ok := remoteResource.(graphmodels.Groupable); ok {
		data.MemberDisplayName = types.StringPointerValue(group.GetDisplayName())
	} else if servicePrincipal, ok := remoteResource.(graphmodels.ServicePrincipalable); ok {
		data.MemberDisplayName = types.StringPointerValue(servicePrincipal.GetDisplayName())
	} else if device, ok := remoteResource.(graphmodels.Deviceable); ok {
		data.MemberDisplayName = types.StringPointerValue(device.GetDisplayName())
	} else if orgContact, ok := remoteResource.(graphmodels.OrgContactable); ok {
		data.MemberDisplayName = types.StringPointerValue(orgContact.GetDisplayName())
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished mapping remote state to Terraform state %s with id %s", ResourceName, data.ID.ValueString()))
}

// extractMemberTypeFromOdataType extracts the member type from @odata.type
func extractMemberTypeFromOdataType(odataType string) string {
	switch odataType {
	case "#microsoft.graph.user":
		return "User"
	case "#microsoft.graph.group":
		return "Group"
	case "#microsoft.graph.servicePrincipal":
		return "ServicePrincipal"
	case "#microsoft.graph.device":
		return "Device"
	case "#microsoft.graph.orgContact":
		return "OrganizationalContact"
	default:
		return "Unknown"
	}
}
