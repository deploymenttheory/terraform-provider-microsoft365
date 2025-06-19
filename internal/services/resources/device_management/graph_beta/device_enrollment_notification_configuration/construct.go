package graphBetaDeviceEnrollmentNotificationConfiguration

import (
	"context"
	"fmt"
	"strings"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/devicemanagement"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// constructResource maps the Terraform resource model to the Graph API request model
func constructResource(ctx context.Context, data *DeviceEnrollmentNotificationConfigurationResourceModel) (models.DeviceEnrollmentNotificationConfigurationable, error) {
	requestBody := models.NewDeviceEnrollmentNotificationConfiguration()

	convert.FrameworkToGraphString(data.DisplayName, requestBody.SetDisplayName)
	convert.FrameworkToGraphString(data.Description, requestBody.SetDescription)
	convert.FrameworkToGraphInt32(data.Priority, requestBody.SetPriority)

	// Platform type is always Windows for enrollment notifications
	platformType := models.WINDOWS_ENROLLMENTRESTRICTIONPLATFORMTYPE
	requestBody.SetPlatformType(&platformType)

	// Set branding options from the set
	if !data.BrandingOptions.IsNull() && !data.BrandingOptions.IsUnknown() {
		// Convert set to comma-separated string for the helper
		elements := data.BrandingOptions.Elements()
		var opts []string
		for _, elem := range elements {
			if strVal, ok := elem.(types.String); ok && !strVal.IsNull() && !strVal.IsUnknown() {
				opts = append(opts, strVal.ValueString())
			}
		}
		joined := strings.Join(opts, ",")
		stringVal := types.StringValue(joined)
		err := convert.FrameworkToGraphBitmaskEnum(
			stringVal,
			models.ParseEnrollmentNotificationBrandingOptions,
			requestBody.SetBrandingOptions,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to set branding options: %s", err)
		}
	}

	var notificationTemplates []string
	if !data.TemplateTypes.IsNull() && !data.TemplateTypes.IsUnknown() {
		for _, element := range data.TemplateTypes.Elements() {
			if stringVal, ok := element.(types.String); ok && !stringVal.IsNull() {
				templateType := stringVal.ValueString()

				notificationTemplates = append(notificationTemplates, fmt.Sprintf("%s_00000000-0000-0000-0000-000000000000", templateType))
			}
		}
	}
	if len(notificationTemplates) > 0 {
		requestBody.SetNotificationTemplates(notificationTemplates)
	}

	if !data.NotificationMessageTemplateId.IsNull() && !data.NotificationMessageTemplateId.IsUnknown() {
		if templateId, err := uuid.Parse(data.NotificationMessageTemplateId.ValueString()); err == nil {
			requestBody.SetNotificationMessageTemplateId(&templateId)
		}
	}

	if err := convert.FrameworkToGraphStringSet(ctx, data.RoleScopeTagIds, requestBody.SetRoleScopeTagIds); err != nil {
		return nil, fmt.Errorf("failed to set role scope tags: %s", err)
	}

	if err := constructors.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for resource %s", ResourceName), requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]interface{}{
			"error": err.Error(),
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished constructing %s resource", ResourceName))

	return requestBody, nil
}

// constructLocalizedMessage constructs a localized message request body
func constructLocalizedMessage(ctx context.Context, message *LocalizedNotificationMessageModel) *models.LocalizedNotificationMessage {
	if message == nil {
		return nil
	}

	requestBody := models.NewLocalizedNotificationMessage()

	// Set default locale if not specified
	locale := "en-us"
	if !message.Locale.IsNull() && !message.Locale.IsUnknown() {
		locale = message.Locale.ValueString()
	}
	requestBody.SetLocale(&locale)

	if !message.Subject.IsNull() && !message.Subject.IsUnknown() {
		subject := message.Subject.ValueString()
		requestBody.SetSubject(&subject)
	}

	if !message.MessageTemplate.IsNull() && !message.MessageTemplate.IsUnknown() {
		messageTemplate := message.MessageTemplate.ValueString()
		requestBody.SetMessageTemplate(&messageTemplate)
	}

	// Set default to true if not specified
	isDefault := true
	if !message.IsDefault.IsNull() && !message.IsDefault.IsUnknown() {
		isDefault = message.IsDefault.ValueBool()
	}
	requestBody.SetIsDefault(&isDefault)

	if err := constructors.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for resource %s", ResourceName), requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]interface{}{
			"error": err.Error(),
		})
	}

	return requestBody
}

// constructAssignmentsRequestBody creates the request body for the assign action.
func constructAssignmentsRequestBody(ctx context.Context, assignments []AssignmentModel) (devicemanagement.DeviceEnrollmentConfigurationsItemAssignPostRequestBodyable, error) {
	requestBody := devicemanagement.NewDeviceEnrollmentConfigurationsItemAssignPostRequestBody()

	assignmentList := make([]models.EnrollmentConfigurationAssignmentable, 0, len(assignments))

	tflog.Debug(ctx, "Starting Device Enrollment Notification Configuration assignment construction")

	for _, assignmentData := range assignments {
		if assignmentData.Target == nil {
			continue
		}

		assignment := models.NewEnrollmentConfigurationAssignment()
		var target models.DeviceAndAppManagementAssignmentTargetable

		targetType := assignmentData.Target.TargetType.ValueString()
		switch targetType {
		case "group":
			groupTarget := models.NewGroupAssignmentTarget()
			convert.FrameworkToGraphString(assignmentData.Target.GroupId, groupTarget.SetGroupId)
			target = groupTarget
		case "allDevices":
			target = models.NewAllDevicesAssignmentTarget()
		case "allLicensedUsers":
			target = models.NewAllLicensedUsersAssignmentTarget()
		default:
			return nil, fmt.Errorf("unsupported target type: %s", targetType)
		}

		convert.FrameworkToGraphString(assignmentData.Target.DeviceAndAppManagementAssignmentFilterId, target.SetDeviceAndAppManagementAssignmentFilterId)
		err := convert.FrameworkToGraphEnum(
			assignmentData.Target.DeviceAndAppManagementAssignmentFilterType,
			models.ParseDeviceAndAppManagementAssignmentFilterType,
			target.SetDeviceAndAppManagementAssignmentFilterType,
		)
		if err != nil {
			return nil, err
		}

		assignment.SetTarget(target)
		assignmentList = append(assignmentList, assignment)
	}

	requestBody.SetEnrollmentConfigurationAssignments(assignmentList)

	if err := constructors.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for resource %s", ResourceName), requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]interface{}{
			"error": err.Error(),
		})
	}

	return requestBody, nil
}
