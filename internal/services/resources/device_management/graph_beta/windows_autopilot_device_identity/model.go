// REF: https://learn.microsoft.com/en-us/graph/api/resources/intune-enrollment-windowsautopilotdeviceidentity?view=graph-rest-beta
package graphBetaWindowsAutopilotDeviceIdentity

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// WindowsAutopilotDeviceIdentityResourceModel represents the values for Windows Autopilot device identities
type WindowsAutopilotDeviceIdentityResourceModel struct {
	ID                                        types.String         `tfsdk:"id"`
	GroupTag                                  types.String         `tfsdk:"group_tag"`
	PurchaseOrderIdentifier                   types.String         `tfsdk:"purchase_order_identifier"`
	SerialNumber                              types.String         `tfsdk:"serial_number"`
	ProductKey                                types.String         `tfsdk:"product_key"`
	Manufacturer                              types.String         `tfsdk:"manufacturer"`
	Model                                     types.String         `tfsdk:"model"`
	EnrollmentState                           types.String         `tfsdk:"enrollment_state"`
	LastContactedDateTime                     types.String         `tfsdk:"last_contacted_date_time"`
	AddressableUserName                       types.String         `tfsdk:"addressable_user_name"`
	UserPrincipalName                         types.String         `tfsdk:"user_principal_name"`
	ResourceName                              types.String         `tfsdk:"resource_name"`
	SkuNumber                                 types.String         `tfsdk:"sku_number"`
	SystemFamily                              types.String         `tfsdk:"system_family"`
	AzureActiveDirectoryDeviceId              types.String         `tfsdk:"azure_active_directory_device_id"`
	AzureAdDeviceId                           types.String         `tfsdk:"azure_ad_device_id"`
	ManagedDeviceId                           types.String         `tfsdk:"managed_device_id"`
	DisplayName                               types.String         `tfsdk:"display_name"`
	DeploymentProfileAssignmentStatus         types.String         `tfsdk:"deployment_profile_assignment_status"`
	DeploymentProfileAssignmentDetailedStatus types.String         `tfsdk:"deployment_profile_assignment_detailed_status"`
	DeploymentProfileAssignedDateTime         types.String         `tfsdk:"deployment_profile_assigned_date_time"`
	RemediationState                          types.String         `tfsdk:"remediation_state"`
	RemediationStateLastModifiedDateTime      types.String         `tfsdk:"remediation_state_last_modified_date_time"`
	UserlessEnrollmentStatus                  types.String         `tfsdk:"userless_enrollment_status"`
	UserAssignment                            *UserAssignmentModel `tfsdk:"user_assignment"`

	Timeouts timeouts.Value `tfsdk:"timeouts"`
}

// UserAssignmentModel represents the user assignment configuration for a Windows Autopilot device
type UserAssignmentModel struct {
	UserPrincipalName   types.String `tfsdk:"user_principal_name"`
	AddressableUserName types.String `tfsdk:"addressable_user_name"`
}

// ImportedWindowsAutopilotDeviceIdentityStateModel represents the state of an imported Windows Autopilot device
type ImportedWindowsAutopilotDeviceIdentityStateModel struct {
	DeviceImportStatus   types.String `tfsdk:"device_import_status"`
	DeviceRegistrationId types.String `tfsdk:"device_registration_id"`
	DeviceErrorCode      types.Int32  `tfsdk:"device_error_code"`
	DeviceErrorName      types.String `tfsdk:"device_error_name"`
}
