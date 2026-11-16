package graphBetaIdentityAndAccessMobilityManagementPolicy

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
)

// MapRemoteResourceStateToTerraform maps the remote mobility management policy and its included group IDs to Terraform state
func MapRemoteResourceStateToTerraform(
	ctx context.Context,
	data MobilityManagementPolicyResourceModel,
	remoteResource graphmodels.MobilityManagementPolicyable,
	includedGroupIDs []string,
) MobilityManagementPolicyResourceModel {
	if remoteResource == nil {
		tflog.Debug(ctx, "Remote resource is nil")
		return data
	}

	tflog.Debug(ctx, "Starting to map remote state to Terraform state", map[string]any{
		"resourceName": remoteResource.GetDisplayName(),
		"resourceId":   remoteResource.GetId(),
	})

	if remoteResource.GetId() != nil {
		data.PolicyID = types.StringValue(*remoteResource.GetId())
	}
	data.ID = types.StringValue(buildResourceID(data.PolicyType.ValueString(), data.PolicyID.ValueString()))
	data.DisplayName = convert.GraphToFrameworkString(remoteResource.GetDisplayName())
	data.Description = convert.GraphToFrameworkString(remoteResource.GetDescription())
	data.IsValid = convert.GraphToFrameworkBool(remoteResource.GetIsValid())
	data.ComplianceUrl = convert.GraphToFrameworkString(remoteResource.GetComplianceUrl())
	data.DiscoveryUrl = convert.GraphToFrameworkString(remoteResource.GetDiscoveryUrl())
	data.TermsOfUseUrl = convert.GraphToFrameworkString(remoteResource.GetTermsOfUseUrl())

	if appliesTo := remoteResource.GetAppliesTo(); appliesTo != nil {
		data.AppliesTo = types.StringValue(appliesTo.String())
	} else {
		data.AppliesTo = types.StringNull()
	}

	data.IncludedGroups = convert.GraphToFrameworkStringSet(ctx, includedGroupIDs)

	if deviceResource, ok := remoteResource.(graphmodels.MobileDeviceManagementPolicyable); ok {
		data.IsMdmEnrollmentDuringRegistrationDisabled = convert.GraphToFrameworkBool(deviceResource.GetIsMdmEnrollmentDuringRegistrationDisabled())
	} else {
		data.IsMdmEnrollmentDuringRegistrationDisabled = types.BoolNull()
	}

	tflog.Debug(ctx, "Finished mapping remote state to Terraform state", map[string]any{
		"resourceId": data.ID.ValueString(),
	})

	return data
}
