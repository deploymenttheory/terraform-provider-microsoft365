package graphBetaAndroidEnrollmentNotifications

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// MapRemoteStateToTerraform maps the remote resource state to the Terraform resource model.
func MapRemoteStateToTerraform(ctx context.Context, data *AndroidEnrollmentNotificationsResourceModel, remoteResource graphmodels.DeviceEnrollmentConfigurationable) {
	if remoteResource == nil {
		tflog.Debug(ctx, "Remote resource is nil")
		return
	}

	tflog.Debug(ctx, "Starting to map remote state to Terraform state", map[string]interface{}{
		"resourceId": convert.GraphToFrameworkString(remoteResource.GetId()).ValueString(),
	})

	data.ID = convert.GraphToFrameworkString(remoteResource.GetId())
	data.DisplayName = convert.GraphToFrameworkString(remoteResource.GetDisplayName())
	data.Description = convert.GraphToFrameworkString(remoteResource.GetDescription())
	data.RoleScopeTagIds = convert.GraphToFrameworkStringSet(ctx, remoteResource.GetRoleScopeTagIds())

	// For DeviceEnrollmentNotificationConfiguration, we set a fixed platform type
	data.PlatformType = types.StringValue("androidForWork")

	// Set default locale - this is a fixed value for this resource type
	data.DefaultLocale = types.StringValue("en-US")

	if enrollmentNotificationConfig, ok := remoteResource.(graphmodels.DeviceEnrollmentNotificationConfigurationable); ok {
		templates := enrollmentNotificationConfig.GetNotificationTemplates()
		if len(templates) > 0 {
			templateValues := make([]attr.Value, 0, len(templates))
			for _, template := range templates {
				// Transform from API format back to user-friendly format
				switch template {
				case "email_00000000-0000-0000-0000-000000000000":
					templateValues = append(templateValues, types.StringValue("email"))
				case "push_00000000-0000-0000-0000-000000000000":
					templateValues = append(templateValues, types.StringValue("push"))
				}
			}

			if len(templateValues) > 0 {
				setVal, diags := types.SetValue(types.StringType, templateValues)
				if !diags.HasError() {
					data.NotificationTemplates = setVal
				}
			}
		}
	} else {
		// Fallback to default values if type assertion fails
		templates := []attr.Value{
			types.StringValue("email"),
			types.StringValue("push"),
		}
		setVal, diags := types.SetValue(types.StringType, templates)
		if !diags.HasError() {
			data.NotificationTemplates = setVal
		}
	}

	data.Priority = convert.GraphToFrameworkInt32(remoteResource.GetPriority())
	data.CreatedDateTime = convert.GraphToFrameworkTime(remoteResource.GetCreatedDateTime())
	data.LastModifiedDateTime = convert.GraphToFrameworkTime(remoteResource.GetLastModifiedDateTime())
	data.Version = convert.GraphToFrameworkInt32(remoteResource.GetVersion())
	data.DeviceEnrollmentConfigurationType = convert.GraphToFrameworkEnum(remoteResource.GetDeviceEnrollmentConfigurationType())

	// Note: Assignments are handled separately in the Read function via dedicated API call

	tflog.Debug(ctx, fmt.Sprintf("Finished mapping resource %s with id %s", ResourceName, data.ID.ValueString()))
}
