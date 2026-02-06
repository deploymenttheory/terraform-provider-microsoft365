package graphBetaMacOSSoftwareUpdateConfiguration

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// constructResource maps the Terraform resource model to the Graph API request model.
func constructResource(ctx context.Context, data *MacOSSoftwareUpdateConfigurationResourceModel) (graphmodels.DeviceConfigurationable, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing %s resource from model", ResourceName))

	requestBody := graphmodels.NewMacOSSoftwareUpdateConfiguration()

	convert.FrameworkToGraphString(data.DisplayName, requestBody.SetDisplayName)
	convert.FrameworkToGraphString(data.Description, requestBody.SetDescription)
	if err := convert.FrameworkToGraphStringSet(ctx, data.RoleScopeTagIds, requestBody.SetRoleScopeTagIds); err != nil {
		return nil, fmt.Errorf("failed to set role_scope_tag_ids: %w", err)
	}

	if err := convert.FrameworkToGraphEnum(
		data.CriticalUpdateBehavior,
		graphmodels.ParseMacOSSoftwareUpdateBehavior,
		requestBody.SetCriticalUpdateBehavior,
	); err != nil {
		return nil, fmt.Errorf("failed to set critical_update_behavior: %w", err)
	}

	if err := convert.FrameworkToGraphEnum(
		data.ConfigDataUpdateBehavior,
		graphmodels.ParseMacOSSoftwareUpdateBehavior,
		requestBody.SetConfigDataUpdateBehavior,
	); err != nil {
		return nil, fmt.Errorf("failed to set config_data_update_behavior: %w", err)
	}

	if err := convert.FrameworkToGraphEnum(
		data.FirmwareUpdateBehavior,
		graphmodels.ParseMacOSSoftwareUpdateBehavior,
		requestBody.SetFirmwareUpdateBehavior,
	); err != nil {
		return nil, fmt.Errorf("failed to set firmware_update_behavior: %w", err)
	}

	if err := convert.FrameworkToGraphEnum(
		data.AllOtherUpdateBehavior,
		graphmodels.ParseMacOSSoftwareUpdateBehavior,
		requestBody.SetAllOtherUpdateBehavior,
	); err != nil {
		return nil, fmt.Errorf("failed to set all_other_update_behavior: %w", err)
	}

	if err := convert.FrameworkToGraphEnum(
		data.UpdateScheduleType,
		graphmodels.ParseMacOSSoftwareUpdateScheduleType,
		requestBody.SetUpdateScheduleType,
	); err != nil {
		return nil, fmt.Errorf("failed to set update_schedule_type: %w", err)
	}

	// Custom update time windows
	if !data.CustomUpdateTimeWindows.IsNull() && !data.CustomUpdateTimeWindows.IsUnknown() {
		var timeWindows []graphmodels.CustomUpdateTimeWindowable
		for _, v := range data.CustomUpdateTimeWindows.Elements() {
			if v.IsNull() || v.IsUnknown() {
				continue
			}
			m := v.(types.Object)
			win := graphmodels.NewCustomUpdateTimeWindow()
			if err := convert.FrameworkToGraphEnum(
				m.Attributes()["start_day"].(types.String),
				graphmodels.ParseDayOfWeek,
				win.SetStartDay,
			); err != nil {
				return nil, fmt.Errorf("invalid start_day: %w", err)
			}
			if err := convert.FrameworkToGraphEnum(
				m.Attributes()["end_day"].(types.String),
				graphmodels.ParseDayOfWeek,
				win.SetEndDay,
			); err != nil {
				return nil, fmt.Errorf("invalid end_day: %w", err)
			}
			_ = convert.FrameworkToGraphTimeOnly(m.Attributes()["start_time"].(types.String), win.SetStartTime)
			_ = convert.FrameworkToGraphTimeOnly(m.Attributes()["end_time"].(types.String), win.SetEndTime)
			timeWindows = append(timeWindows, win)
		}
		requestBody.SetCustomUpdateTimeWindows(timeWindows)
	}

	convert.FrameworkToGraphInt32(data.MaxUserDeferralsCount, requestBody.SetMaxUserDeferralsCount)

	if err := convert.FrameworkToGraphEnum(
		data.Priority,
		graphmodels.ParseMacOSPriority,
		requestBody.SetPriority,
	); err != nil {
		return nil, fmt.Errorf("failed to set priority: %w", err)
	}

	if err := constructors.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for resource %s", ResourceName), requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]any{
			"error": err.Error(),
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished constructing %s resource", ResourceName))

	return requestBody, nil
}
