package graphBetaServicePrincipal

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ServicePrincipalResourceModel describes the Terraform resource data model for a Service Principal.
//
// Properties that Microsoft Entra reflects from the associated application (its display name,
// home page, logout URL, reply URLs, app roles, permission scopes, informational URLs and so on)
// are deliberately absent: they are owned by the application, change whenever it changes, and are
// available from the application resource or the service principal data source. Modelling them
// here would make every update plan churn with "known after apply" for values this resource does
// not own.
type ServicePrincipalResourceModel struct {
	ID                                  types.String   `tfsdk:"id"`
	AppID                               types.String   `tfsdk:"app_id"`
	AccountEnabled                      types.Bool     `tfsdk:"account_enabled"`
	AlternativeNames                    types.Set      `tfsdk:"alternative_names"`
	AppOwnerOrganizationID              types.String   `tfsdk:"app_owner_organization_id"`
	AppRoleAssignmentRequired           types.Bool     `tfsdk:"app_role_assignment_required"`
	ApplicationTemplateID               types.String   `tfsdk:"application_template_id"`
	CreatedByAppID                      types.String   `tfsdk:"created_by_app_id"`
	Description                         types.String   `tfsdk:"description"`
	IsDisabled                          types.Bool     `tfsdk:"is_disabled"`
	KeyCredentials                      types.Set      `tfsdk:"key_credentials"`
	LoginURL                            types.String   `tfsdk:"login_url"`
	Notes                               types.String   `tfsdk:"notes"`
	NotificationEmailAddresses          types.Set      `tfsdk:"notification_email_addresses"`
	PasswordCredentials                 types.Set      `tfsdk:"password_credentials"`
	PreferredSingleSignOnMode           types.String   `tfsdk:"preferred_single_sign_on_mode"`
	PreferredTokenSigningKeyEndDateTime types.String   `tfsdk:"preferred_token_signing_key_end_date_time"`
	PreferredTokenSigningKeyThumbprint  types.String   `tfsdk:"preferred_token_signing_key_thumbprint"`
	SamlSingleSignOnSettings            types.Object   `tfsdk:"saml_single_sign_on_settings"`
	ServicePrincipalType                types.String   `tfsdk:"service_principal_type"`
	Tags                                types.Set      `tfsdk:"tags"`
	TokenEncryptionKeyID                types.String   `tfsdk:"token_encryption_key_id"`
	HardDelete                          types.Bool     `tfsdk:"hard_delete"`
	Timeouts                            timeouts.Value `tfsdk:"timeouts"`
}

// SamlSingleSignOnSettingsModel describes the nested saml_single_sign_on_settings object
type SamlSingleSignOnSettingsModel struct {
	RelayState types.String `tfsdk:"relay_state"`
}
