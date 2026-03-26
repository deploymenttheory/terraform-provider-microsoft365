package graphBetaAdministrativeUnitDirectoryRoleAssignment

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// MapRemoteResourceStateToTerraform maps the remote scoped role membership from the Kiota SDK to Terraform state
func MapRemoteResourceStateToTerraform(ctx context.Context, data AdministrativeUnitDirectoryRoleAssignmentResourceModel, remoteResource graphmodels.ScopedRoleMembershipable) AdministrativeUnitDirectoryRoleAssignmentResourceModel {
	if remoteResource == nil {
		tflog.Debug(ctx, "Remote resource is nil")
		return data
	}

	tflog.Debug(ctx, "Starting to map remote state to Terraform state", map[string]any{
		"resourceId": remoteResource.GetId(),
	})

	data.ID = convert.GraphToFrameworkString(remoteResource.GetId())
	data.AdministrativeUnitID = convert.GraphToFrameworkString(remoteResource.GetAdministrativeUnitId())
	data.DirectoryRoleID = convert.GraphToFrameworkString(remoteResource.GetRoleId())

	if memberInfo := remoteResource.GetRoleMemberInfo(); memberInfo != nil {
		data.RoleMemberID = convert.GraphToFrameworkString(memberInfo.GetId())
		data.RoleMemberDisplayName = convert.GraphToFrameworkString(memberInfo.GetDisplayName())

		// userPrincipalName is not part of the Identityable interface; extract from AdditionalData
		data.RoleMemberUserPrincipalName = types.StringNull()
		if additionalData := memberInfo.GetAdditionalData(); additionalData != nil {
			if upn, ok := additionalData["userPrincipalName"].(string); ok {
				data.RoleMemberUserPrincipalName = types.StringValue(upn)
			}
		}
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished mapping resource %s with id %s", ResourceName, data.ID.ValueString()))
	return data
}
