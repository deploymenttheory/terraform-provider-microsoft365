package graphBetaSettingsCatalog

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/state"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// MapRemoteResourceStateToTerraform states the base properties of a SettingsCatalogProfileResourceModel to a Terraform state
func MapRemoteResourceStateToTerraform(ctx context.Context, data *SettingsCatalogProfileResourceModel, remoteResource graphmodels.DeviceManagementConfigurationPolicyable) {
	if remoteResource == nil {
		tflog.Debug(ctx, "Remote resource is nil")
		return
	}

	tflog.Debug(ctx, "Starting to map remote resource state to Terraform state", map[string]interface{}{
		"resourceId": state.StringPtrToString(remoteResource.GetId()),
	})

	data.ID = types.StringValue(state.StringPtrToString(remoteResource.GetId()))
	data.Name = types.StringValue(state.StringPtrToString(remoteResource.GetName()))
	data.Description = types.StringValue(state.StringPtrToString(remoteResource.GetDescription()))
	data.RoleScopeTagIds = state.SliceToTypeStringSlice(remoteResource.GetRoleScopeTagIds())
	data.IsAssigned = state.BoolPtrToTypeBool(remoteResource.GetIsAssigned())
	data.CreatedDateTime = state.TimeToString(remoteResource.GetCreatedDateTime())
	data.LastModifiedDateTime = state.TimeToString(remoteResource.GetLastModifiedDateTime())
	data.SettingsCount = state.Int32PtrToTypeInt64(remoteResource.GetSettingCount())

	if platforms := remoteResource.GetPlatforms(); platforms != nil {
		data.Platforms = state.EnumPtrToTypeString(platforms)
	}
	if technologies := remoteResource.GetTechnologies(); technologies != nil {
		data.Technologies = EnumBitmaskToTypeStringSlice(*technologies)
	}

	tflog.Debug(ctx, "Finished mapping remote resource state to Terraform state", map[string]interface{}{
		"resourceId": data.ID.ValueString(),
	})
}

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
	if technologies&graphmodels.MOBILEAPPLICATIONMANAGEMENT_DEVICEMANAGEMENTCONFIGURATIONTECHNOLOGIES != 0 {
		values = append(values, types.StringValue("mobileApplicationManagement"))
	}
	if technologies&graphmodels.LINUXMDM_DEVICEMANAGEMENTCONFIGURATIONTECHNOLOGIES != 0 {
		values = append(values, types.StringValue("linuxMdm"))
	}
	if technologies&graphmodels.ENROLLMENT_DEVICEMANAGEMENTCONFIGURATIONTECHNOLOGIES != 0 {
		values = append(values, types.StringValue("enrollment"))
	}
	if technologies&graphmodels.ENDPOINTPRIVILEGEMANAGEMENT_DEVICEMANAGEMENTCONFIGURATIONTECHNOLOGIES != 0 {
		values = append(values, types.StringValue("endpointPrivilegeManagement"))
	}
	if technologies&graphmodels.UNKNOWNFUTUREVALUE_DEVICEMANAGEMENTCONFIGURATIONTECHNOLOGIES != 0 {
		values = append(values, types.StringValue("unknownFutureValue"))
	}
	if technologies&graphmodels.WINDOWSOSRECOVERY_DEVICEMANAGEMENTCONFIGURATIONTECHNOLOGIES != 0 {
		values = append(values, types.StringValue("windowsOsRecovery"))
	}

	return values
}
