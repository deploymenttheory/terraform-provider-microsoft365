// REF: https://learn.microsoft.com/en-us/graph/api/resources/intune-mam-managedmobileapp?view=graph-rest-beta
package graphBetaDeviceAndAppManagementIOSManagedMobileApp

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// IOSManagedMobileAppResourceModel represents the Terraform resource model for an iOS managed mobile app.
type IOSManagedMobileAppResourceModel struct {
	ID                     types.String                         `tfsdk:"id"`
	Version                types.String                         `tfsdk:"version"`
	MobileAppIdentifier    *IOSMobileAppIdentifierResourceModel `tfsdk:"mobile_app_identifier"`
	ManagedAppProtectionId types.String                         `tfsdk:"managed_app_protection_id"`
	Timeouts               timeouts.Value                       `tfsdk:"timeouts"`
}

// IOSMobileAppIdentifierResourceModel represents the iOS mobile app identifier.
type IOSMobileAppIdentifierResourceModel struct {
	BundleId types.String `tfsdk:"bundle_id"`
}
