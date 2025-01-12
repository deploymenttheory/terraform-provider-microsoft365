package graphBetaReuseablePolicySettings

import (
	"context"
	"encoding/json"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/normalize"
	sharedmodels "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/shared_models/graph_beta/device_and_app_management"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/state"
	"github.com/hashicorp/terraform-plugin-framework/attr"
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

	elements := make([]attr.Value, 0)
	for _, policy := range remoteResource.GetReferencingConfigurationPolicies() {
		if id := policy.GetId(); id != nil {
			elements = append(elements, types.StringValue(*id))
		}
	}
	data.ReferencingConfigurationPolicies = types.ListValueMust(types.StringType, elements)

	StateReusablePolicySettings(ctx, &data.Settings, remoteResource.GetSettingInstance())

	tflog.Debug(ctx, "Finished mapping remote resource state to Terraform state", map[string]interface{}{
		"resourceId": data.ID.ValueString(),
	})
}

func StateReusablePolicySettings(ctx context.Context, data *types.String, settingInstance graphmodels.DeviceManagementConfigurationSettingInstanceable) {
	if settingInstance == nil {
		// Handle nil settingInstance case
		emptyContent := map[string]interface{}{
			"settingsDetails": []interface{}{
				map[string]interface{}{
					"settingInstance": map[string]interface{}{},
				},
			},
		}
		jsonBytes, err := json.Marshal(emptyContent)
		if err != nil {
			tflog.Error(ctx, "Failed to marshal empty content", map[string]interface{}{"error": err.Error()})
			return
		}
		*data = types.StringValue(string(jsonBytes))
		return
	}

	// Convert the settingInstance to a map for proper JSON handling
	settingInstanceBytes, err := json.Marshal(settingInstance)
	if err != nil {
		tflog.Error(ctx, "Failed to marshal setting instance", map[string]interface{}{"error": err.Error()})
		return
	}

	var settingInstanceMap map[string]interface{}
	if err := json.Unmarshal(settingInstanceBytes, &settingInstanceMap); err != nil {
		tflog.Error(ctx, "Failed to unmarshal setting instance to map", map[string]interface{}{"error": err.Error()})
		return
	}

	structuredContent := map[string]interface{}{
		"settingsDetails": []interface{}{
			map[string]interface{}{
				"settingInstance": settingInstanceMap,
			},
		},
	}

	jsonBytes, err := json.Marshal(structuredContent)
	if err != nil {
		tflog.Error(ctx, "Failed to marshal structured content", map[string]interface{}{"error": err.Error()})
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
