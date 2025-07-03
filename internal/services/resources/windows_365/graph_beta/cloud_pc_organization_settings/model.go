// REF: https://learn.microsoft.com/en-us/graph/api/resources/cloudpcorganizationsettings?view=graph-rest-beta
package graphBetaCloudPcOrganizationSettings

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type CloudPcOrganizationSettingsResourceModel struct {
	ID                  types.String          `tfsdk:"id"`
	EnableMEMAutoEnroll types.Bool            `tfsdk:"enable_mem_auto_enroll"`
	EnableSingleSignOn  types.Bool            `tfsdk:"enable_single_sign_on"`
	OsVersion           types.String          `tfsdk:"os_version"`
	UserAccountType     types.String          `tfsdk:"user_account_type"`
	WindowsSettings     *WindowsSettingsModel `tfsdk:"windows_settings"`
}

type WindowsSettingsModel struct {
	Language types.String `tfsdk:"language"`
}
