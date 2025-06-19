package graphBetaReuseablePolicySettings

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	sharedmodels "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/shared_models/graph_beta/device_management"
	sharedStater "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/state/graph_beta/device_management"
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
		"resourceId": convert.GraphToFrameworkString(remoteResource.GetId()),
	})

	// Add debug logs to trace the ID
	id := remoteResource.GetId()
	tflog.Debug(ctx, "Remote resource ID value", map[string]interface{}{
		"id":    id,
		"isNil": id == nil,
	})

	// Check Entity interface implementation
	if entity, ok := remoteResource.(graphmodels.Entityable); ok {
		tflog.Debug(ctx, "Entity ID value", map[string]interface{}{
			"id": entity.GetId(),
		})
	}

	data.ID = convert.GraphToFrameworkString(remoteResource.GetId())
	data.DisplayName = convert.GraphToFrameworkString(remoteResource.GetDisplayName())
	data.Description = convert.GraphToFrameworkString(remoteResource.GetDescription())
	data.CreatedDateTime = convert.GraphToFrameworkTime(remoteResource.GetCreatedDateTime())
	data.LastModifiedDateTime = convert.GraphToFrameworkTime(remoteResource.GetLastModifiedDateTime())
	data.Version = convert.GraphToFrameworkInt32(remoteResource.GetVersion())
	data.ReferencingConfigurationPolicyCount = convert.GraphToFrameworkInt32(remoteResource.GetReferencingConfigurationPolicyCount())

	elements := make([]attr.Value, 0)
	for _, policy := range remoteResource.GetReferencingConfigurationPolicies() {
		elements = append(elements, convert.GraphToFrameworkString(policy.GetId()))
	}
	data.ReferencingConfigurationPolicies = types.ListValueMust(types.StringType, elements)

	// State the settings catalog fields
	if settingInstance := remoteResource.GetSettingInstance(); settingInstance != nil {
		// Create a wrapper to match the expected format
		wrappedSettings := map[string]interface{}{
			"settings": []interface{}{
				map[string]interface{}{
					"id":              "0", // Single setting always has ID 0
					"settingInstance": settingInstance,
				},
			},
		}

		// Convert to JSON
		settingsJson, err := json.Marshal(wrappedSettings)
		if err != nil {
			tflog.Error(ctx, "Failed to marshal settings", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		// Use the shared stater to handle the settings
		sharedStater.StateReusablePolicySettings(ctx, data, settingsJson)
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished stating resource %s with id %s", ResourceName, data.ID.ValueString()))
}
