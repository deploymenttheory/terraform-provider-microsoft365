package graphBetaReuseablePolicySettings

import (
	"context"
	"encoding/json"

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
		tflog.Error(ctx, "Setting instance is nil")
		return
	}

	// Debug log the additional data
	tflog.Debug(ctx, "Additional data from settingInstance", map[string]interface{}{
		"data": settingInstance.GetAdditionalData(),
	})

	// Create base structure
	settingInst := sharedmodels.SettingInstance{}

	// Map basic fields
	if odataType := settingInstance.GetOdataType(); odataType != nil {
		settingInst.ODataType = *odataType
	}
	if settingDefId := settingInstance.GetSettingDefinitionId(); settingDefId != nil {
		settingInst.SettingDefinitionId = *settingDefId
	}

	// Map template reference
	if templateRef := settingInstance.GetSettingInstanceTemplateReference(); templateRef != nil {
		if templateId := templateRef.GetSettingInstanceTemplateId(); templateId != nil {
			settingInst.SettingInstanceTemplateReference = &sharedmodels.SettingInstanceTemplateReference{
				SettingInstanceTemplateId: *templateId,
			}
		}
	}

	// Map simpleSettingValue from additional data
	if additionalData := settingInstance.GetAdditionalData(); additionalData != nil {
		if simpleValue, ok := additionalData["simpleSettingValue"].(map[string]interface{}); ok {
			simpleSettingVal := &sharedmodels.SimpleSettingStruct{}

			if odataType, ok := simpleValue["@odata.type"].(string); ok {
				simpleSettingVal.ODataType = odataType
			}

			// Handle value as an interface{} since it could be string, bool, etc
			if val, exists := simpleValue["value"]; exists {
				simpleSettingVal.Value = val
			}

			settingInst.SimpleSettingValue = simpleSettingVal
		}
	}

	// Create the full content structure
	content := sharedmodels.DeviceConfigV2GraphServiceResourceModel{
		SettingsDetails: []sharedmodels.SettingDetail{
			{
				ID:              "0",
				SettingInstance: settingInst,
			},
		},
	}

	// Convert to string for Terraform state
	jsonBytes, err := json.Marshal(content)
	if err != nil {
		tflog.Error(ctx, "Error marshaling content", map[string]interface{}{"error": err.Error()})
		return
	}

	tflog.Debug(ctx, "Final structured content", map[string]interface{}{
		"content": string(jsonBytes),
	})

	*data = types.StringValue(string(jsonBytes))
}
