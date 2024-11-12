package graphBetaSettingsCatalog

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/state"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

func MapRemoteResourceStateToTerraform(ctx context.Context, data *SettingsCatalogProfileResourceModel, remoteResource graphmodels.DeviceManagementConfigurationPolicyable) {
	if remoteResource == nil {
		tflog.Debug(ctx, "Remote resource is nil")
		return
	}

	tflog.Debug(ctx, "Starting to map remote resource state to Terraform state", map[string]interface{}{
		"resourceId": state.StringPtrToString(remoteResource.GetId()),
	})

	// Map basic properties
	data.ID = types.StringValue(state.StringPtrToString(remoteResource.GetId()))
	data.Name = types.StringValue(state.StringPtrToString(remoteResource.GetName()))
	data.Description = types.StringValue(state.StringPtrToString(remoteResource.GetDescription()))
	data.CreatedDateTime = state.TimeToString(remoteResource.GetCreatedDateTime())
	data.LastModifiedDateTime = state.TimeToString(remoteResource.GetLastModifiedDateTime())
	data.SettingsCount = state.Int32PtrToTypeInt64(remoteResource.GetSettingCount())
	data.RoleScopeTagIds = state.SliceToTypeStringSlice(remoteResource.GetRoleScopeTagIds())
	data.IsAssigned = state.BoolPtrToTypeBool(remoteResource.GetIsAssigned())

	// Map enum values
	if platforms := remoteResource.GetPlatforms(); platforms != nil {
		data.Platforms = state.EnumPtrToTypeString(platforms)
	}
	if technologies := remoteResource.GetTechnologies(); technologies != nil {
		data.Technologies = state.EnumPtrToTypeString(technologies)
	}

	tflog.Debug(ctx, "Finished mapping remote resource state to Terraform state", map[string]interface{}{
		"resourceId": data.ID.ValueString(),
	})
}
