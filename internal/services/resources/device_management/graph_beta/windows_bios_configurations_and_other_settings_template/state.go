package graphBetaWindowsBiosConfigurationsAndOtherSettingsTemplate

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
)

func MapRemoteResourceStateToTerraform(
	ctx context.Context,
	data *WindowsBiosConfigurationsAndOtherSettingsTemplateResourceModel,
	remoteResource graphmodels.HardwareConfigurationable,
) {
	if remoteResource == nil {
		tflog.Debug(ctx, "Remote resource is nil")
		return
	}

	tflog.Debug(ctx, "Starting to map remote state to Terraform state", map[string]any{
		"resourceId": convert.GraphToFrameworkString(remoteResource.GetId()).ValueString(),
	})

	data.ID = convert.GraphToFrameworkString(remoteResource.GetId())
	data.DisplayName = convert.GraphToFrameworkString(remoteResource.GetDisplayName())
	data.Description = convert.GraphToFrameworkString(remoteResource.GetDescription())
	data.FileName = convert.GraphToFrameworkString(remoteResource.GetFileName())
	data.HardwareConfigurationFormat = convert.GraphToFrameworkEnum(
		remoteResource.GetHardwareConfigurationFormat(),
	)
	data.PerDevicePasswordDisabled = convert.GraphToFrameworkBool(
		remoteResource.GetPerDevicePasswordDisabled(),
	)
	data.RoleScopeTagIds = convert.GraphToFrameworkStringSet(
		ctx,
		remoteResource.GetRoleScopeTagIds(),
	)
	data.Version = convert.GraphToFrameworkInt32(remoteResource.GetVersion())
	data.CreatedDateTime = convert.GraphToFrameworkTime(remoteResource.GetCreatedDateTime())
	data.LastModifiedDateTime = convert.GraphToFrameworkTime(
		remoteResource.GetLastModifiedDateTime(),
	)

	// Kiota base64 decodes configurationFileContent during deserialization, so re-encode it to match
	// the base64 the schema exposes. The API omits the field on some responses; keep the configured
	// value in that case rather than nulling an attribute the practitioner set.
	if content := remoteResource.GetConfigurationFileContent(); content != nil {
		data.ConfigurationFileContent = types.StringValue(helpers.ByteStringToBase64(content))
	} else {
		tflog.Debug(
			ctx,
			"Remote resource did not return configurationFileContent, retaining the value already in state",
			map[string]any{
				"resourceId": data.ID.ValueString(),
			},
		)
	}

	assignments := remoteResource.GetAssignments()
	tflog.Debug(ctx, "Retrieved assignments from remote resource", map[string]any{
		"assignmentCount": len(assignments),
		"resourceId":      data.ID.ValueString(),
	})

	MapAssignmentsToTerraform(ctx, data, assignments)

	tflog.Debug(
		ctx,
		fmt.Sprintf("Finished mapping resource %s with id %s", ResourceName, data.ID.ValueString()),
	)
}
