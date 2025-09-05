package graphBetaAppControlForBusinessPolicy

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// MapRemoteResourceStateToTerraform maps the base properties of a DeviceManagementConfigurationPolicy to a Terraform state
func MapRemoteResourceStateToTerraform(ctx context.Context, data *AppControlForBusinessPolicyResourceModel, remoteResource graphmodels.DeviceManagementConfigurationPolicyable) {
	if remoteResource == nil {
		tflog.Debug(ctx, "Remote resource is nil")
		return
	}

	tflog.Debug(ctx, "Starting to map remote state to Terraform state", map[string]interface{}{
		"resourceId": convert.GraphToFrameworkString(remoteResource.GetId()).ValueString(),
	})

	data.ID = convert.GraphToFrameworkString(remoteResource.GetId())
	data.Name = convert.GraphToFrameworkString(remoteResource.GetName())
	data.Description = convert.GraphToFrameworkString(remoteResource.GetDescription())
	data.RoleScopeTagIds = convert.GraphToFrameworkStringSet(ctx, remoteResource.GetRoleScopeTagIds())

	tflog.Debug(ctx, fmt.Sprintf("Finished mapping resource %s with id %s", ResourceName, data.ID.ValueString()))
}

// EnumBitmaskToTypeStringSlice converts technologies bitmask to string slice
func EnumBitmaskToTypeStringSlice(technologies graphmodels.DeviceManagementConfigurationTechnologies) []types.String {
	var values []types.String

	if technologies&graphmodels.NONE_DEVICEMANAGEMENTCONFIGURATIONTECHNOLOGIES != 0 {
		values = append(values, types.StringValue("none"))
	}
	if technologies&graphmodels.MDM_DEVICEMANAGEMENTCONFIGURATIONTECHNOLOGIES != 0 {
		values = append(values, types.StringValue("mdm"))
	}
	if technologies&graphmodels.WINDOWS10XMANAGEMENT_DEVICEMANAGEMENTCONFIGURATIONTECHNOLOGIES != 0 {
		values = append(values, types.StringValue("windows10XManagement"))
	}
	if technologies&graphmodels.CONFIGMANAGER_DEVICEMANAGEMENTCONFIGURATIONTECHNOLOGIES != 0 {
		values = append(values, types.StringValue("configManager"))
	}
	if technologies&graphmodels.APPLEREMOTEMANAGEMENT_DEVICEMANAGEMENTCONFIGURATIONTECHNOLOGIES != 0 {
		values = append(values, types.StringValue("appleRemoteManagement"))
	}
	if technologies&graphmodels.MICROSOFTSENSE_DEVICEMANAGEMENTCONFIGURATIONTECHNOLOGIES != 0 {
		values = append(values, types.StringValue("microsoftSense"))
	}
	if technologies&graphmodels.EXCHANGEONLINE_DEVICEMANAGEMENTCONFIGURATIONTECHNOLOGIES != 0 {
		values = append(values, types.StringValue("exchangeOnline"))
	}
	if technologies&graphmodels.LINUXMDM_DEVICEMANAGEMENTCONFIGURATIONTECHNOLOGIES != 0 {
		values = append(values, types.StringValue("linuxMdm"))
	}
	if technologies&graphmodels.ENROLLMENT_DEVICEMANAGEMENTCONFIGURATIONTECHNOLOGIES != 0 {
		values = append(values, types.StringValue("enrollment"))
	}
	return values
}