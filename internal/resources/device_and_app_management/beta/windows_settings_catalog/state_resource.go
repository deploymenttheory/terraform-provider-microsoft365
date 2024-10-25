package graphBetaWindowsSettingsCatalog

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/state"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

func MapRemoteResourceStateToTerraform(ctx context.Context, data *WindowsSettingsCatalogProfileResourceModel, remoteResource graphmodels.DeviceManagementConfigurationPolicyable) {
	if remoteResource == nil {
		tflog.Debug(ctx, "Remote resource is nil")
		return
	}

	tflog.Debug(ctx, "Starting to map remote resource state to Terraform state", map[string]interface{}{
		"resourceId": state.StringPtrToString(remoteResource.GetId()),
	})

	// Map basic properties
	data.ID = types.StringValue(state.StringPtrToString(remoteResource.GetId()))
	data.DisplayName = types.StringValue(state.StringPtrToString(remoteResource.GetName()))
	data.Description = types.StringValue(state.StringPtrToString(remoteResource.GetDescription()))
	data.CreationSource = types.StringValue(state.StringPtrToString(remoteResource.GetCreationSource()))
	data.CreatedDateTime = state.TimeToString(remoteResource.GetCreatedDateTime())
	data.LastModifiedDateTime = state.TimeToString(remoteResource.GetLastModifiedDateTime())
	data.IsAssigned = state.BoolPtrToTypeBool(remoteResource.GetIsAssigned())
	data.SettingsCount = state.Int32PtrToTypeInt32(remoteResource.GetSettingCount())

	// Map enum values
	if platforms := remoteResource.GetPlatforms(); platforms != nil {
		data.Platforms = state.EnumPtrToTypeString(platforms)
	}
	if technologies := remoteResource.GetTechnologies(); technologies != nil {
		data.Technologies = state.EnumPtrToTypeString(technologies)
	}

	// Map role scope tag IDs
	data.RoleScopeTagIds = state.SliceToTypeStringSlice(remoteResource.GetRoleScopeTagIds())

	// Map template reference if present
	// if template := remoteResource.GetTemplateReference(); template != nil {
	// 	data.TemplateReference = TemplateReference{
	// 		TemplateID:             types.StringValue(state.StringPtrToString(template.GetTemplateId())),
	// 		TemplateFamily:         state.EnumPtrToTypeString(template.GetTemplateFamily()).ValueString(),
	// 		TemplateDisplayName:    types.StringValue(state.StringPtrToString(template.GetTemplateDisplayName())),
	// 		TemplateDisplayVersion: types.StringValue(state.StringPtrToString(template.GetTemplateDisplayVersion())),
	// 	}
	// }

	tflog.Debug(ctx, "Finished mapping remote resource state to Terraform state", map[string]interface{}{
		"resourceId": data.ID.ValueString(),
	})
}
