package graphBetaAndroidEnrollmentNotifications

import (
	"context"
	"fmt"
	"strings"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// constructResource maps the Terraform schema to the SDK model for Device Enrollment Notification Configuration.
func constructResource(ctx context.Context, data *AndroidEnrollmentNotificationsResourceModel, isForCreate bool, currentTemplateGUIDs ...[]string) (graphmodels.DeviceEnrollmentConfigurationable, error) {
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

	// api expects email_00000000-0000-0000-0000-000000000000 and push_00000000-0000-0000-0000-000000000000 for
	// create operations. Intune then updates the zero's to the actual template GUIDs.
	// For update operations, you must use the actual template GUIDs else it causes a 400 error.
	if !data.NotificationTemplates.IsNull() && !data.NotificationTemplates.IsUnknown() {
		templateValues := make([]string, 0)

		if isForCreate {
			// For create operations, use blank GUIDs. Intune will update the GUIDs to the actual template GUIDs.
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
		} else if len(currentTemplateGUIDs) > 0 {
			// For update operations, use actual template GUIDs
			planTemplateTypes := make(map[string]bool)
			for _, templateType := range data.NotificationTemplates.Elements() {
				templateTypeStr := templateType.String()
				templateTypeStr = templateTypeStr[1 : len(templateTypeStr)-1] // Remove quotes
				planTemplateTypes[templateTypeStr] = true
			}

			// Include existing template GUIDs for types that are still wanted
			existingTypes := make(map[string]bool)
			for _, currentTemplate := range currentTemplateGUIDs[0] {
				if strings.Contains(strings.ToLower(currentTemplate), "email") && planTemplateTypes["email"] {
					templateValues = append(templateValues, currentTemplate)
					existingTypes["email"] = true
				} else if strings.Contains(strings.ToLower(currentTemplate), "push") && planTemplateTypes["push"] {
					templateValues = append(templateValues, currentTemplate)
					existingTypes["push"] = true
				}
			}

			// Add zero GUIDs for new template types that don't exist yet
			for templateType := range planTemplateTypes {
				if !existingTypes[templateType] {
					switch templateType {
					case "email":
						templateValues = append(templateValues, "email_00000000-0000-0000-0000-000000000000")
					case "push":
						templateValues = append(templateValues, "push_00000000-0000-0000-0000-000000000000")
					}
				}
			}

			tflog.Debug(ctx, fmt.Sprintf("Using template GUIDs for update (existing + new): %v", templateValues))
		}

		if len(templateValues) > 0 {
			requestBody.SetNotificationTemplates(templateValues)
		}
	}

	if err := constructors.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for resource %s", ResourceName), requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]any{
			"error": err.Error(),
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished constructing %s resource", ResourceName))

	return requestBody, nil
}
