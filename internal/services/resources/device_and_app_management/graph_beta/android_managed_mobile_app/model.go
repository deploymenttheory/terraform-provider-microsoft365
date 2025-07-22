// REF: https://learn.microsoft.com/en-us/graph/api/resources/intune-mam-managedmobileapp?view=graph-rest-beta
package graphBetaDeviceAndAppManagementAndroidManagedMobileApp

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// AndroidManagedMobileAppResourceModel represents the Terraform resource model for an Android managed mobile app.
type AndroidManagedMobileAppResourceModel struct {
	ID                         types.String                               `tfsdk:"id"`
	Version                    types.String                               `tfsdk:"version"`
	MobileAppIdentifier        *AndroidMobileAppIdentifierResourceModel  `tfsdk:"mobile_app_identifier"`
	ManagedAppProtectionId     types.String                               `tfsdk:"managed_app_protection_id"`
	Timeouts                   timeouts.Value                             `tfsdk:"timeouts"`
}

// AndroidMobileAppIdentifierResourceModel represents the Android mobile app identifier.
type AndroidMobileAppIdentifierResourceModel struct {
	PackageId types.String `tfsdk:"package_id"`
}