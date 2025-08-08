package graphBetaLinuxPlatformScript

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// MapRemoteResourceStateToTerraform states the base properties of a SettingsCatalogProfileResourceModel to a Terraform state
func MapRemoteResourceStateToTerraform(ctx context.Context, data *LinuxPlatformScriptResourceModel, remoteResource graphmodels.DeviceManagementConfigurationPolicyable) {
	if remoteResource == nil {
		tflog.Debug(ctx, "Remote resource is nil")
		return
	}

	tflog.Debug(ctx, "Starting to map remote state to Terraform state", map[string]interface{}{
		"resourceId": convert.GraphToFrameworkString(remoteResource.GetId()).ValueString(),
	})

	data.ID = convert.GraphToFrameworkString(remoteResource.GetId())
	data.Name = convert.GraphToFrameworkString(remoteResource.GetName())
	data.Description = convert.GraphToFrameworkString(remoteResource.GetDescription())

	data.RoleScopeTagIds = convert.GraphToFrameworkStringSet(ctx, remoteResource.GetRoleScopeTagIds())

	if platforms := remoteResource.GetPlatforms(); platforms != nil {
		data.Platforms = convert.GraphToFrameworkEnum(platforms)
	}

	if technologies := remoteResource.GetTechnologies(); technologies != nil {
		data.Technologies = EnumBitmaskToTypeStringSlice(*technologies)
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished mapping resource %s with id %s", ResourceName, data.ID.ValueString()))
}

func EnumBitmaskToTypeStringSlice(technologies graphmodels.DeviceManagementConfigurationTechnologies) []types.String {
	var values []types.String

	if technologies&graphmodels.NONE_DEVICEMANAGEMENTCONFIGURATIONTECHNOLOGIES != 0 {
		values = append(values, types.StringValue("none"))
	}
	if technologies&graphmodels.LINUXMDM_DEVICEMANAGEMENTCONFIGURATIONTECHNOLOGIES != 0 {
		values = append(values, types.StringValue("linuxMdm"))
	}
	return values
}
