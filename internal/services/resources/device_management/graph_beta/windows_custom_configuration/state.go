package graphBetaWindowsCustomConfiguration

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// OmaSettingType returns the object type for OmaSettingResourceModel
func OmaSettingType() types.ObjectType {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"odata_type":   types.StringType,
			"display_name": types.StringType,
			"description":  types.StringType,
			"oma_uri":      types.StringType,
			"value":        types.StringType,
			"file_name":    types.StringType,
		},
	}
}

func MapRemoteResourceStateToTerraform(ctx context.Context, data *WindowsCustomConfigurationResourceModel, remoteResource graphmodels.DeviceConfigurationable) {
	if remoteResource == nil {
		tflog.Debug(ctx, "Remote resource is nil")
		return
	}

	tflog.Debug(ctx, "Starting to map remote state to Terraform state", map[string]any{
		"resourceId": convert.GraphToFrameworkString(remoteResource.GetId()).ValueString(),
	})

	// Map common properties
	data.ID = convert.GraphToFrameworkString(remoteResource.GetId())
	data.DisplayName = convert.GraphToFrameworkString(remoteResource.GetDisplayName())
	data.Description = convert.GraphToFrameworkString(remoteResource.GetDescription())
	data.RoleScopeTagIds = convert.GraphToFrameworkStringSet(ctx, remoteResource.GetRoleScopeTagIds())

	if config, ok := remoteResource.(graphmodels.Windows10CustomConfigurationable); ok {
		mapOmaSettings(ctx, data, config)
	} else {
		tflog.Error(ctx, "Remote resource is not a Windows10CustomConfiguration", map[string]any{
			"type": fmt.Sprintf("%T", remoteResource),
		})
	}

	// Map assignments
	assignments := remoteResource.GetAssignments()
	tflog.Debug(ctx, "Retrieved assignments from remote resource", map[string]any{
		"assignmentCount": len(assignments),
		"resourceId":      data.ID.ValueString(),
	})

	if len(assignments) == 0 {
		data.Assignments = types.SetNull(WindowsCustomConfigurationAssignmentType())
	} else {
		mapAssignmentsToTerraform(ctx, data, assignments)
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished mapping resource %s with id %s", ResourceName, data.ID.ValueString()))
}

func mapOmaSettings(ctx context.Context, data *WindowsCustomConfigurationResourceModel, config graphmodels.Windows10CustomConfigurationable) {
	remoteSettings := config.GetOmaSettings()
	if len(remoteSettings) == 0 {
		data.OmaSettings = types.ListNull(OmaSettingType())
		return
	}

	settingModels := make([]OmaSettingResourceModel, 0, len(remoteSettings))
	for _, remoteSetting := range remoteSettings {
		if remoteSetting == nil {
			continue
		}
		settingModels = append(settingModels, mapOmaSetting(ctx, remoteSetting))
	}

	listValue, diags := types.ListValueFrom(ctx, OmaSettingType(), settingModels)
	if diags.HasError() {
		tflog.Error(ctx, "Failed to create oma settings list", map[string]any{
			"errors": diags.Errors(),
		})
		data.OmaSettings = types.ListNull(OmaSettingType())
		return
	}

	data.OmaSettings = listValue
}

// mapOmaSetting maps a microsoft.graph.omaSetting subtype to the Terraform model, converting
// the typed value back to its string representation.
func mapOmaSetting(ctx context.Context, remoteSetting graphmodels.OmaSettingable) OmaSettingResourceModel {
	settingModel := OmaSettingResourceModel{
		OdataType:   convert.GraphToFrameworkString(remoteSetting.GetOdataType()),
		DisplayName: convert.GraphToFrameworkString(remoteSetting.GetDisplayName()),
		Description: convert.GraphToFrameworkString(remoteSetting.GetDescription()),
		OmaUri:      convert.GraphToFrameworkString(remoteSetting.GetOmaUri()),
		Value:       types.StringNull(),
		FileName:    types.StringNull(),
	}

	switch setting := remoteSetting.(type) {
	case *graphmodels.OmaSettingString:
		settingModel.Value = convert.GraphToFrameworkString(setting.GetValue())
	case *graphmodels.OmaSettingInteger:
		if value := setting.GetValue(); value != nil {
			settingModel.Value = types.StringValue(strconv.FormatInt(int64(*value), 10))
		}
	case *graphmodels.OmaSettingBoolean:
		if value := setting.GetValue(); value != nil {
			settingModel.Value = types.StringValue(strconv.FormatBool(*value))
		}
	case *graphmodels.OmaSettingBase64:
		settingModel.Value = convert.GraphToFrameworkString(setting.GetValue())
		settingModel.FileName = convert.GraphToFrameworkString(setting.GetFileName())
	case *graphmodels.OmaSettingDateTime:
		if value := setting.GetValue(); value != nil {
			settingModel.Value = types.StringValue(value.UTC().Format(time.RFC3339))
		}
	case *graphmodels.OmaSettingFloatingPoint:
		if value := setting.GetValue(); value != nil {
			settingModel.Value = types.StringValue(strconv.FormatFloat(float64(*value), 'f', -1, 32))
		}
	case *graphmodels.OmaSettingStringXml:
		settingModel.Value = convert.GraphToFrameworkBytes(setting.GetValue())
		settingModel.FileName = convert.GraphToFrameworkString(setting.GetFileName())
	default:
		tflog.Error(ctx, "Unknown oma setting type", map[string]any{
			"type": fmt.Sprintf("%T", remoteSetting),
		})
	}

	return settingModel
}
