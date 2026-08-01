package graphBetaServicePrincipal

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ServicePrincipalResourceModel describes the Terraform resource data model for a Service Principal
type ServicePrincipalResourceModel struct {
	ID                         types.String   `tfsdk:"id"`
	AppID                      types.String   `tfsdk:"app_id"`
	DisplayName                types.String   `tfsdk:"display_name"`
	AccountEnabled             types.Bool     `tfsdk:"account_enabled"`
	AlternativeNames           types.Set      `tfsdk:"alternative_names"`
	AppRoleAssignmentRequired  types.Bool     `tfsdk:"app_role_assignment_required"`
	Description                types.String   `tfsdk:"description"`
	Homepage                   types.String   `tfsdk:"homepage"`
	LoginURL                   types.String   `tfsdk:"login_url"`
	LogoutURL                  types.String   `tfsdk:"logout_url"`
	Notes                      types.String   `tfsdk:"notes"`
	NotificationEmailAddresses types.Set      `tfsdk:"notification_email_addresses"`
	PreferredSingleSignOnMode  types.String   `tfsdk:"preferred_single_sign_on_mode"`
	SamlSingleSignOnSettings   types.Object   `tfsdk:"saml_single_sign_on_settings"`
	ServicePrincipalType       types.String   `tfsdk:"service_principal_type"`
	ServicePrincipalNames      types.Set      `tfsdk:"service_principal_names"`
	SignInAudience             types.String   `tfsdk:"sign_in_audience"`
	Tags                       types.Set      `tfsdk:"tags"`
	HardDelete                 types.Bool     `tfsdk:"hard_delete"`
	Timeouts                   timeouts.Value `tfsdk:"timeouts"`
}

// SamlSingleSignOnSettingsModel describes the nested saml_single_sign_on_settings object
type SamlSingleSignOnSettingsModel struct {
	RelayState types.String `tfsdk:"relay_state"`
}
