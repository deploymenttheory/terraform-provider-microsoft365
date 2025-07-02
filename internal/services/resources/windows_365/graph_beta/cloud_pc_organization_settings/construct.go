package graphBetaCloudPcOrganizationSettings

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

func constructResource(ctx context.Context, data *CloudPcOrganizationSettingsResourceModel) (*models.CloudPcOrganizationSettings, error) {
	requestBody := models.NewCloudPcOrganizationSettings()

	convert.FrameworkToGraphBool(data.EnableMEMAutoEnroll, requestBody.SetEnableMEMAutoEnroll)
	convert.FrameworkToGraphBool(data.EnableSingleSignOn, requestBody.SetEnableSingleSignOn)

	if err := convert.FrameworkToGraphEnum(
		data.OsVersion,
		models.ParseCloudPcOperatingSystem,
		func(v models.CloudPcOperatingSystem) { requestBody.SetOsVersion(&v) },
	); err != nil {
		return nil, fmt.Errorf("failed to set os_version: %w", err)
	}

	if err := convert.FrameworkToGraphEnum(
		data.UserAccountType,
		models.ParseCloudPcUserAccountType,
		func(v models.CloudPcUserAccountType) { requestBody.SetUserAccountType(&v) },
	); err != nil {
		return nil, fmt.Errorf("failed to set user_account_type: %w", err)
	}

	// WindowsSettings mapping
	if data.WindowsSettings != nil {
		ws := models.NewCloudPcWindowsSettings()
		convert.FrameworkToGraphString(data.WindowsSettings.Language, ws.SetLanguage)
		requestBody.SetWindowsSettings(ws)
	}

	return requestBody, nil
}
