package graphBetaAndroidEnrollmentNotifications

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// constructResource maps the Terraform schema to the SDK model for Device Enrollment Notification Configuration.
func constructResource(ctx context.Context, data *AndroidEnrollmentNotificationsResourceModel, isForCreate bool) (graphmodels.DeviceEnrollmentConfigurationable, error) {
	tflog.Debug(ctx, "Constructing DeviceEnrollmentNotificationConfiguration resource from Terraform state")

	requestBody := graphmodels.NewDeviceEnrollmentNotificationConfiguration()

	convert.FrameworkToGraphString(data.DisplayName, requestBody.SetDisplayName)
	convert.FrameworkToGraphString(data.Description, requestBody.SetDescription)
	convert.FrameworkToGraphString(data.DefaultLocale, requestBody.SetDefaultLocale)

	if !data.PlatformType.IsNull() && !data.PlatformType.IsUnknown() {
		platformTypeStr := data.PlatformType.ValueString()
		switch platformTypeStr {
		case "androidForWork":
			platformTypeValue := graphmodels.ANDROIDFORWORK_ENROLLMENTRESTRICTIONPLATFORMTYPE
			requestBody.SetPlatformType(&platformTypeValue)
		case "android":
			platformTypeValue := graphmodels.ANDROID_ENROLLMENTRESTRICTIONPLATFORMTYPE
			requestBody.SetPlatformType(&platformTypeValue)
		default:
			return nil, fmt.Errorf("unsupported platform type: %s", platformTypeStr)
		}
	}

	if err := convert.FrameworkToGraphStringSet(ctx, data.RoleScopeTagIds, requestBody.SetRoleScopeTagIds); err != nil {
		return nil, fmt.Errorf("failed to set role scope tags: %s", err)
	}

	// Transform notification template schema values to expected api format
	if !data.NotificationTemplates.IsNull() && !data.NotificationTemplates.IsUnknown() {
		templateValues := make([]string, 0)
		for _, templateType := range data.NotificationTemplates.Elements() {
			templateTypeStr := templateType.String()
			templateTypeStr = templateTypeStr[1 : len(templateTypeStr)-1]

			switch templateTypeStr {
			case "email":
				templateValues = append(templateValues, "email_00000000-0000-0000-0000-000000000000")
			case "push":
				templateValues = append(templateValues, "push_00000000-0000-0000-0000-000000000000")
			default:
				return nil, fmt.Errorf("unsupported notification template type: %s", templateTypeStr)
			}
		}
		requestBody.SetNotificationTemplates(templateValues)
	}

	if err := constructors.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for resource %s", ResourceName), requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]any{
			"error": err.Error(),
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished constructing %s resource", ResourceName))

	return requestBody, nil
}
