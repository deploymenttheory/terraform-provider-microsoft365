package graphBetaReuseablePolicySettings

import (
	"context"
	"encoding/json"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/normalize"
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
	data.DisplayName = types.StringPointerValue(remoteResource.GetDisplayName())
	data.Description = types.StringPointerValue(remoteResource.GetDescription())
	data.CreatedDateTime = state.TimeToString(remoteResource.GetCreatedDateTime())
	data.LastModifiedDateTime = state.TimeToString(remoteResource.GetLastModifiedDateTime())
	data.Version = types.Int32PointerValue(remoteResource.GetVersion())
	data.ReferencingConfigurationPolicyCount = types.Int32PointerValue(remoteResource.GetReferencingConfigurationPolicyCount())

	StateReusablePolicySettings(ctx, &data.Settings, remoteResource.GetSettingInstance())

	tflog.Debug(ctx, "Finished mapping remote resource state to Terraform state", map[string]interface{}{
		"resourceId": data.ID.ValueString(),
	})
}

func StateReusablePolicySettings(ctx context.Context, data *types.String, settingInstance graphmodels.DeviceManagementConfigurationSettingInstanceable) {
	var configSettings map[string]interface{}
	if err := json.Unmarshal([]byte(data.ValueString()), &configSettings); err != nil {
		tflog.Error(ctx, "Failed to unmarshal config settings", map[string]interface{}{"error": err.Error()})
		return
	}

	// Structure the single setting instance under settingsDetails
	structuredContent := map[string]interface{}{
		"settingsDetails": []interface{}{
			map[string]interface{}{
				"settingInstance": settingInstance,
			},
		},
	}

	if err := normalize.PreserveSecretSettings(configSettings, structuredContent); err != nil {
		tflog.Error(ctx, "Error stating secret settings from HCL", map[string]interface{}{"error": err.Error()})
		return
	}

	jsonBytes, err := json.Marshal(structuredContent)
	if err != nil {
		tflog.Error(ctx, "Failed to marshal JSON structured content during preparation for normalization", map[string]interface{}{"error": err.Error()})
		return
	}

	normalizedJSON, err := normalize.JSONAlphabetically(string(jsonBytes))
	if err != nil {
		tflog.Error(ctx, "Failed to normalize settings JSON alphabetically", map[string]interface{}{"error": err.Error()})
		return
	}

	tflog.Debug(ctx, "Normalized settings", map[string]interface{}{"settings": normalizedJSON})

	*data = types.StringValue(normalizedJSON)
}
