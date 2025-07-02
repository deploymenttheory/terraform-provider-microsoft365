package graphBetaCloudPcOrganizationSettings

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-framework/types"
	msgraphbetamodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

func MapRemoteStateToTerraform(ctx context.Context, data *CloudPcOrganizationSettingsResourceModel, remote msgraphbetamodels.CloudPcOrganizationSettingsable) {
	if remote == nil {
		return
	}
	data.ID = types.StringValue("singleton")
	data.EnableMEMAutoEnroll = convert.GraphToFrameworkBool(remote.GetEnableMEMAutoEnroll())
	data.EnableSingleSignOn = convert.GraphToFrameworkBool(remote.GetEnableSingleSignOn())
	data.OsVersion = convert.GraphToFrameworkEnum(remote.GetOsVersion())
	data.UserAccountType = convert.GraphToFrameworkEnum(remote.GetUserAccountType())

	if ws := remote.GetWindowsSettings(); ws != nil {
		data.WindowsSettings = &WindowsSettingsModel{
			Language: convert.GraphToFrameworkString(ws.GetLanguage()),
		}
	} else {
		data.WindowsSettings = nil
	}
}
