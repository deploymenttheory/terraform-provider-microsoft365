package graphBetaCloudPcOrganizationSettings

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

func constructResource(ctx context.Context, data *CloudPcOrganizationSettingsResourceModel) (*models.CloudPcOrganizationSettings, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing %s resource from model", ResourceName))

	requestBody := models.NewCloudPcOrganizationSettings()

	convert.FrameworkToGraphBool(data.EnableMEMAutoEnroll, requestBody.SetEnableMEMAutoEnroll)
	convert.FrameworkToGraphBool(data.EnableSingleSignOn, requestBody.SetEnableSingleSignOn)

	if err := convert.FrameworkToGraphEnum(data.UserAccountType, models.ParseCloudPcUserAccountType, requestBody.SetUserAccountType); err != nil {
		return nil, fmt.Errorf("invalid user account type: %s", err)
	}

	if err := convert.FrameworkToGraphEnum(data.OsVersion, models.ParseCloudPcOperatingSystem, requestBody.SetOsVersion); err != nil {
		return nil, fmt.Errorf("invalid run as account type: %s", err)
	}

	// WindowsSettings mapping
	if data.WindowsSettings != nil {
		ws := models.NewCloudPcWindowsSettings()
		convert.FrameworkToGraphString(data.WindowsSettings.Language, ws.SetLanguage)
		requestBody.SetWindowsSettings(ws)
	}

	if err := constructors.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for resource %s", ResourceName), requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]any{
			"error": err.Error(),
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished constructing %s resource", ResourceName))

	return requestBody, nil
}
