package graphBetaSettingsCatalogInventoryPolicy

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

func MapRemoteResourceStateToTerraform(ctx context.Context, data *InventoryPolicyResourceModel, remoteResource graphmodels.DeviceManagementConfigurationPolicyable) {
	if remoteResource == nil {
		tflog.Debug(ctx, "Remote resource is nil")
		return
	}

	data.ID = convert.GraphToFrameworkString(remoteResource.GetId())
	data.Name = convert.GraphToFrameworkString(remoteResource.GetName())
	data.Description = convert.GraphToFrameworkString(remoteResource.GetDescription())
	data.CreatedDateTime = convert.GraphToFrameworkTime(remoteResource.GetCreatedDateTime())
	data.LastModifiedDateTime = convert.GraphToFrameworkTime(remoteResource.GetLastModifiedDateTime())
	data.SettingsCount = convert.GraphToFrameworkInt32(remoteResource.GetSettingCount())

	data.RoleScopeTagIds = convert.GraphToFrameworkStringSet(ctx, remoteResource.GetRoleScopeTagIds())

	if platforms := remoteResource.GetPlatforms(); platforms != nil {
		data.Platforms = convert.GraphToFrameworkEnum(platforms)
	}

	// "extensibility" is not in the SDK enum, so GetTechnologies() returns nil — fall back to AdditionalData
	if technologies := remoteResource.GetTechnologies(); technologies != nil {
		data.Technologies = convert.GraphToFrameworkEnum(technologies)
	} else if additionalData := remoteResource.GetAdditionalData(); additionalData != nil {
		if techVal, ok := additionalData["technologies"]; ok {
			if techStr, ok := techVal.(string); ok {
				data.Technologies = types.StringValue(techStr)
			}
		}
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished mapping %s with id %s", ResourceName, data.ID.ValueString()))
}
