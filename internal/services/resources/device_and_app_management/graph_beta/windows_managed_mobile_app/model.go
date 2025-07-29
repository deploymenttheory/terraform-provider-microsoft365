// REF: https://learn.microsoft.com/en-us/graph/api/resources/intune-mam-managedmobileapp?view=graph-rest-beta
package graphBetaDeviceAndAppManagementWindowsManagedMobileApp

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// WindowsManagedMobileAppResourceModel represents the Terraform resource model for a Windows managed mobile app.
type WindowsManagedMobileAppResourceModel struct {
	ID                     types.String                             `tfsdk:"id"`
	Version                types.String                             `tfsdk:"version"`
	MobileAppIdentifier    *WindowsMobileAppIdentifierResourceModel `tfsdk:"mobile_app_identifier"`
	ManagedAppProtectionId types.String                             `tfsdk:"managed_app_protection_id"`
	Timeouts               timeouts.Value                           `tfsdk:"timeouts"`
}

// WindowsMobileAppIdentifierResourceModel represents the Windows mobile app identifier.
type WindowsMobileAppIdentifierResourceModel struct {
	WindowsAppId types.String `tfsdk:"windows_app_id"`
}
