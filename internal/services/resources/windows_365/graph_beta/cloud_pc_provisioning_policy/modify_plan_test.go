package graphBetaCloudPcProvisioningPolicy

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/stretchr/testify/require"
)

func TestUnitResourceCloudPcProvisioningPolicy_ComputedOnlyPlan(t *testing.T) {
	persistenceType := tftypes.Object{AttributeTypes: map[string]tftypes.Type{"enabled": tftypes.Bool, "storage_size_category": tftypes.String}}
	objectType := tftypes.Object{AttributeTypes: map[string]tftypes.Type{"display_name": tftypes.String, "image_id": tftypes.String, "image_display_name": tftypes.String, "description": tftypes.String, "alternate_resource_url": tftypes.String, "cloud_pc_group_display_name": tftypes.String, "grace_period_in_hours": tftypes.Number, "user_settings_persistence_configuration": persistenceType}}
	for _, scenario := range []string{"no_effective_change", "rename", "new_image", "enable_persistence", "unknown_description", "unknown_persistence"} {
		t.Run(scenario, func(t *testing.T) {
			previous := map[string]tftypes.Value{"display_name": tftypes.NewValue(tftypes.String, "Policy"), "image_id": tftypes.NewValue(tftypes.String, "image-id"), "image_display_name": tftypes.NewValue(tftypes.String, "Image name"), "description": tftypes.NewValue(tftypes.String, nil), "alternate_resource_url": tftypes.NewValue(tftypes.String, nil), "cloud_pc_group_display_name": tftypes.NewValue(tftypes.String, nil), "grace_period_in_hours": tftypes.NewValue(tftypes.Number, 4), "user_settings_persistence_configuration": tftypes.NewValue(persistenceType, map[string]tftypes.Value{"enabled": tftypes.NewValue(tftypes.Bool, false), "storage_size_category": tftypes.NewValue(tftypes.String, "eightGB")})}
			planned := map[string]tftypes.Value{}
			configured := map[string]tftypes.Value{}
			for name, value := range previous {
				planned[name] = value
				configured[name] = value
			}
			for _, name := range []string{"image_display_name", "description", "alternate_resource_url", "cloud_pc_group_display_name", "grace_period_in_hours"} {
				planned[name] = tftypes.NewValue(objectType.AttributeTypes[name], tftypes.UnknownValue)
				configured[name] = tftypes.NewValue(objectType.AttributeTypes[name], nil)
			}
			configured["user_settings_persistence_configuration"] = tftypes.NewValue(persistenceType, nil)
			switch scenario {
			case "rename":
				planned["display_name"] = tftypes.NewValue(tftypes.String, "Renamed")
				configured["display_name"] = planned["display_name"]
			case "new_image":
				planned["image_id"] = tftypes.NewValue(tftypes.String, "new-image-id")
				configured["image_id"] = planned["image_id"]
			case "enable_persistence":
				planned["user_settings_persistence_configuration"] = tftypes.NewValue(persistenceType, map[string]tftypes.Value{"enabled": tftypes.NewValue(tftypes.Bool, true), "storage_size_category": tftypes.NewValue(tftypes.String, "eightGB")})
				configured["user_settings_persistence_configuration"] = planned["user_settings_persistence_configuration"]
			case "unknown_persistence":
				planned["user_settings_persistence_configuration"] = tftypes.NewValue(persistenceType, tftypes.UnknownValue)
				configured["user_settings_persistence_configuration"] = planned["user_settings_persistence_configuration"]
			case "unknown_description":
				configured["description"] = tftypes.NewValue(tftypes.String, tftypes.UnknownValue)
			}
			request := resource.ModifyPlanRequest{State: tfsdk.State{Raw: tftypes.NewValue(objectType, previous)}, Plan: tfsdk.Plan{Raw: tftypes.NewValue(objectType, planned)}, Config: tfsdk.Config{Raw: tftypes.NewValue(objectType, configured)}}
			response := resource.ModifyPlanResponse{Plan: request.Plan}
			(&CloudPcProvisioningPolicyResource{}).ModifyPlan(context.Background(), request, &response)
			require.False(t, response.Diagnostics.HasError(), response.Diagnostics)
			if scenario == "no_effective_change" {
				require.True(t, response.Plan.Raw.Equal(request.State.Raw))
			} else {
				require.True(t, response.Plan.Raw.Equal(request.Plan.Raw))
				var attributes map[string]tftypes.Value
				require.NoError(t, response.Plan.Raw.As(&attributes))
				require.False(t, attributes["image_display_name"].IsKnown())
			}
		})
	}
}
