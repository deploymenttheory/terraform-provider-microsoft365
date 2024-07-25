package graphBetaAssignmentFilter

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

/* platform type validator */
type platformTypeValidator struct{}

func (v platformTypeValidator) Description(ctx context.Context) string {
	return fmt.Sprintf("platform must be one of: %s", strings.Join(validPlatformTypes, ", "))
}

func (v platformTypeValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v platformTypeValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	value := req.ConfigValue.ValueString()
	for _, validType := range validPlatformTypes {
		if value == validType {
			return
		}
	}

	resp.Diagnostics.AddAttributeError(
		req.Path,
		"Invalid Platform Type",
		fmt.Sprintf("Platform must be one of: %s", strings.Join(validPlatformTypes, ", ")),
	)
}

var validPlatformTypes = []string{
	"android",
	"androidForWork",
	"iOS",
	"macOS",
	"windowsPhone81",
	"windows81AndLater",
	"windows10AndLater",
	"androidWorkProfile",
	"unknown",
	"androidAOSP",
	"androidMobileApplicationManagement",
	"iOSMobileApplicationManagement",
	"windowsMobileApplicationManagement",
}

/* assignmentFilterManagement Type validator */

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
