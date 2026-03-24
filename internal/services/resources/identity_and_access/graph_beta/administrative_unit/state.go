package graphBetaAdministrativeUnit

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// MapRemoteResourceStateToTerraform maps the remote administrative unit from Kiota SDK to Terraform state
// Note: HardDelete is preserved from the existing state as it's an HCL-only field not returned by the API
func MapRemoteResourceStateToTerraform(ctx context.Context, data AdministrativeUnitResourceModel, remoteResource graphmodels.AdministrativeUnitable) AdministrativeUnitResourceModel {
	if remoteResource == nil {
		tflog.Debug(ctx, "Remote resource is nil")
		return data
	}

	tflog.Debug(ctx, "Starting to map remote state to Terraform state", map[string]any{
		"resourceName": remoteResource.GetDisplayName(),
		"resourceId":   remoteResource.GetId(),
	})

	// Preserve HardDelete from existing state (HCL-only field, not returned by API)
	existingHardDelete := data.HardDelete

	data.ID = convert.GraphToFrameworkString(remoteResource.GetId())
	data.DisplayName = convert.GraphToFrameworkString(remoteResource.GetDisplayName())
	data.Description = convert.GraphToFrameworkString(remoteResource.GetDescription())
	data.IsMemberManagementRestricted = convert.GraphToFrameworkBool(remoteResource.GetIsMemberManagementRestricted())
	data.MembershipRule = convert.GraphToFrameworkString(remoteResource.GetMembershipRule())
	data.MembershipRuleProcessingState = convert.GraphToFrameworkString(remoteResource.GetMembershipRuleProcessingState())
	data.MembershipType = convert.GraphToFrameworkString(remoteResource.GetMembershipType())
	data.Visibility = convert.GraphToFrameworkString(remoteResource.GetVisibility())

	// Restore HardDelete from existing state
	data.HardDelete = existingHardDelete

	tflog.Debug(ctx, fmt.Sprintf("Finished mapping resource %s with id %s", ResourceName, data.ID.ValueString()))
	return data
}
