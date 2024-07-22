package assignmentFilter

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// platformValidator is the custom validator type
type platformValidator struct{}

// ValidateString performs the validation.
func (v platformValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsUnknown() || req.ConfigValue.IsNull() {
		return
	}

	_, err := models.ParseDevicePlatformType(req.ConfigValue.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid Device Platform Type",
			fmt.Sprintf("The platform type '%s' is not valid. Supported types: %v", req.ConfigValue.ValueString(), getAllPlatformStrings()),
		)
	}
}

// Description describes the validation in plain text.
func (v platformValidator) Description(ctx context.Context) string {
	return "must be a valid device platform type"
}

// MarkdownDescription describes the validation in Markdown.
func (v platformValidator) MarkdownDescription(ctx context.Context) string {
	return "must be a valid device platform type"
}

// getAllPlatformStrings returns all the valid platform strings
func getAllPlatformStrings() []string {
	platformTypes := []models.DevicePlatformType{
		models.ANDROID_DEVICEPLATFORMTYPE,
		models.ANDROIDFORWORK_DEVICEPLATFORMTYPE,
		models.IOS_DEVICEPLATFORMTYPE,
		models.MACOS_DEVICEPLATFORMTYPE,
		models.WINDOWSPHONE81_DEVICEPLATFORMTYPE,
		models.WINDOWS81ANDLATER_DEVICEPLATFORMTYPE,
		models.WINDOWS10ANDLATER_DEVICEPLATFORMTYPE,
		models.ANDROIDWORKPROFILE_DEVICEPLATFORMTYPE,
		models.UNKNOWN_DEVICEPLATFORMTYPE,
		models.ANDROIDAOSP_DEVICEPLATFORMTYPE,
		models.ANDROIDMOBILEAPPLICATIONMANAGEMENT_DEVICEPLATFORMTYPE,
		models.IOSMOBILEAPPLICATIONMANAGEMENT_DEVICEPLATFORMTYPE,
		models.UNKNOWNFUTUREVALUE_DEVICEPLATFORMTYPE,
		models.WINDOWSMOBILEAPPLICATIONMANAGEMENT_DEVICEPLATFORMTYPE,
	}

	var platformStrings []string
	for _, platform := range platformTypes {
		platformStrings = append(platformStrings, platform.String())
	}
	return platformStrings
}

// assignmentFilterManagementTypeValidator is the custom validator type
type assignmentFilterManagementTypeValidator struct{}

// Description describes the validation in plain text.
func (v assignmentFilterManagementTypeValidator) Description(ctx context.Context) string {
	return "must be a valid assignment filter management type"
}

// MarkdownDescription describes the validation in Markdown.
func (v assignmentFilterManagementTypeValidator) MarkdownDescription(ctx context.Context) string {
	return "must be a valid assignment filter management type"
}

// ValidateString performs the validation.
func (v assignmentFilterManagementTypeValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsUnknown() || req.ConfigValue.IsNull() {
		return
	}

	validTypes := getAllManagementTypeStrings()
	value := req.ConfigValue.ValueString()
	for _, validType := range validTypes {
		if value == validType {
			return
		}
	}

	resp.Diagnostics.AddError(
		"Invalid Assignment Filter Management Type",
		fmt.Sprintf("The management type '%s' is not valid. Supported types: %v", value, validTypes),
	)
}

// getAllManagementTypeStrings returns all the valid management type strings
func getAllManagementTypeStrings() []string {
	return []string{"devices", "apps", "unknownFutureValue"}
}

// assignmentFilterTypeValidator is the custom validator type
type assignmentFilterTypeValidator struct{}

// ValidateString performs the validation.
func (v assignmentFilterTypeValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsUnknown() || req.ConfigValue.IsNull() {
		return
	}

	validTypes := getAllAssignmentFilterTypes()
	value := req.ConfigValue.ValueString()
	for _, validType := range validTypes {
		if value == validType {
			return
		}
	}

	resp.Diagnostics.AddError(
		"Invalid Assignment Filter Type",
		fmt.Sprintf("The assignment filter type '%s' is not valid. Supported types: %v", value, validTypes),
	)
}

// Description describes the validation in plain text.
func (v assignmentFilterTypeValidator) Description(ctx context.Context) string {
	return "must be a valid assignment filter type"
}

// MarkdownDescription describes the validation in Markdown.
func (v assignmentFilterTypeValidator) MarkdownDescription(ctx context.Context) string {
	return "must be a valid assignment filter type"
}

// getAllAssignmentFilterTypes returns all the valid assignment filter type strings
func getAllAssignmentFilterTypes() []string {
	types := []models.AssociatedAssignmentPayloadType{
		models.UNKNOWN_ASSOCIATEDASSIGNMENTPAYLOADTYPE,
		models.DEVICECONFIGURATIONANDCOMPLIANCE_ASSOCIATEDASSIGNMENTPAYLOADTYPE,
		models.APPLICATION_ASSOCIATEDASSIGNMENTPAYLOADTYPE,
		models.ANDROIDENTERPRISEAPP_ASSOCIATEDASSIGNMENTPAYLOADTYPE,
		models.ENROLLMENTCONFIGURATION_ASSOCIATEDASSIGNMENTPAYLOADTYPE,
		models.GROUPPOLICYCONFIGURATION_ASSOCIATEDASSIGNMENTPAYLOADTYPE,
		models.ZEROTOUCHDEPLOYMENTDEVICECONFIGPROFILE_ASSOCIATEDASSIGNMENTPAYLOADTYPE,
		models.ANDROIDENTERPRISECONFIGURATION_ASSOCIATEDASSIGNMENTPAYLOADTYPE,
		models.DEVICEFIRMWARECONFIGURATIONINTERFACEPOLICY_ASSOCIATEDASSIGNMENTPAYLOADTYPE,
		models.RESOURCEACCESSPOLICY_ASSOCIATEDASSIGNMENTPAYLOADTYPE,
		models.WIN32APP_ASSOCIATEDASSIGNMENTPAYLOADTYPE,
		models.DEVICEMANAGMENTCONFIGURATIONANDCOMPLIANCEPOLICY_ASSOCIATEDASSIGNMENTPAYLOADTYPE,
	}

	var typeStrings []string
	for _, t := range types {
		typeStrings = append(typeStrings, t.String())
	}
	return typeStrings
}
