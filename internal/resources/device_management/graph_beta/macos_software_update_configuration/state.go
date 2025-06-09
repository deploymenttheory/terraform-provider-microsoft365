package graphBetaMacOSSoftwareUpdateConfiguration

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/state"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// mapRemoteResourceStateToTerraform maps the Graph API model to the Terraform state model.
func mapRemoteResourceStateToTerraform(ctx context.Context, data *MacOSSoftwareUpdateConfigurationResourceModel, remote graphmodels.DeviceConfigurationable) {
	if remote == nil {
		return
	}

	resource, ok := remote.(graphmodels.MacOSSoftwareUpdateConfigurationable)
	if !ok {
		return
	}

	data.ID = state.StringPointerValue(resource.GetId())
	data.DisplayName = state.StringPointerValue(resource.GetDisplayName())
	data.Description = state.StringPointerValue(resource.GetDescription())
	data.RoleScopeTagIds = state.StringSliceToSet(ctx, resource.GetRoleScopeTagIds())
	data.CriticalUpdateBehavior = state.EnumPtrToTypeString(resource.GetCriticalUpdateBehavior())
	data.ConfigDataUpdateBehavior = state.EnumPtrToTypeString(resource.GetConfigDataUpdateBehavior())
	data.FirmwareUpdateBehavior = state.EnumPtrToTypeString(resource.GetFirmwareUpdateBehavior())
	data.AllOtherUpdateBehavior = state.EnumPtrToTypeString(resource.GetAllOtherUpdateBehavior())
	data.UpdateScheduleType = state.EnumPtrToTypeString(resource.GetUpdateScheduleType())
	data.Priority = state.EnumPtrToTypeString(resource.GetPriority())
	data.UpdateTimeWindowUtcOffsetInMinutes = state.Int32PointerValue(resource.GetUpdateTimeWindowUtcOffsetInMinutes())
	data.MaxUserDeferralsCount = state.Int32PointerValue(resource.GetMaxUserDeferralsCount())

	// Custom update time windows
	objType := types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"start_day":  types.StringType,
			"end_day":    types.StringType,
			"start_time": types.StringType,
			"end_time":   types.StringType,
		},
	}
	var windows []attr.Value
	if remoteWindows := resource.GetCustomUpdateTimeWindows(); remoteWindows != nil {
		for _, win := range remoteWindows {
			if win == nil {
				continue
			}
			obj, _ := types.ObjectValue(
				objType.AttrTypes,
				map[string]attr.Value{
					"start_day":  state.EnumPtrToTypeString(win.GetStartDay()),
					"end_day":    state.EnumPtrToTypeString(win.GetEndDay()),
					"start_time": state.TimeOnlyPtrToString(win.GetStartTime()),
					"end_time":   state.TimeOnlyPtrToString(win.GetEndTime()),
				},
			)
			windows = append(windows, obj)
		}
	}
	if len(windows) == 0 {
		data.CustomUpdateTimeWindows = types.ListNull(objType)
	} else {
		data.CustomUpdateTimeWindows = types.ListValueMust(objType, windows)
	}
}
