package graphBetaReuseablePolicySettings

import (
	"context"

	sharedmodels "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/shared_models/graph_beta/device_and_app_management"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/state"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// MapRemoteResourceStateToTerraform states the base properties of a EndpointPrivilegeManagementResourceModel to a Terraform state
func MapRemoteResourceStateToTerraform(ctx context.Context, data *sharedmodels.ReuseablePolicySettingsResourceModel, remoteResource graphmodels.DeviceManagementReusablePolicySettingable) {
	if remoteResource == nil {
		tflog.Debug(ctx, "Remote resource is nil")
		return
	}

	tflog.Debug(ctx, "Starting to map remote resource state to Terraform state", map[string]interface{}{
		"resourceId": state.StringPtrToString(remoteResource.GetId()),
	})

	data.ID = types.StringPointerValue(remoteResource.GetId())
	data.Name = types.StringPointerValue(remoteResource.GetDisplayName())
	data.Description = types.StringPointerValue(remoteResource.GetDescription())
	data.CreatedDateTime = state.TimeToString(remoteResource.GetCreatedDateTime())
	data.LastModifiedDateTime = state.TimeToString(remoteResource.GetLastModifiedDateTime())
	data.Version = types.Int64PointerValue(int64(*remoteResource.GetVersion()))
	data.SettingDefinitionId = types.StringPointerValue(remoteResource.GetSettingDefinitionId())
	data.ReferencingConfigurationPolicyCount = types.Int64PointerValue(int64(*remoteResource.GetReferencingConfigurationPolicyCount()))

	if settingInstance := remoteResource.GetSettingInstance(); settingInstance != nil {
		// Handle setting instance mapping based on your needs
	}

	tflog.Debug(ctx, "Finished mapping remote resource state to Terraform state", map[string]interface{}{
		"resourceId": data.ID.ValueString(),
	})
}
