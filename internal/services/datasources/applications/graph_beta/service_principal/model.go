// REF: https://learn.microsoft.com/en-us/graph/api/resources/serviceprincipal?view=graph-rest-beta

package graphBetaServicePrincipal

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/datasource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ServicePrincipalDataSourceModel defines the data source model
type ServicePrincipalDataSourceModel struct {
	ID                                  types.String   `tfsdk:"id"`
	ObjectId                            types.String   `tfsdk:"object_id"`
	AppId                               types.String   `tfsdk:"app_id"`
	DisplayName                         types.String   `tfsdk:"display_name"`
	ODataQuery                          types.String   `tfsdk:"odata_query"`
	AppDisplayName                      types.String   `tfsdk:"app_display_name"`
	DeletedDateTime                     types.String   `tfsdk:"deleted_date_time"`
	ApplicationTemplateID               types.String   `tfsdk:"application_template_id"`
	AccountEnabled                      types.Bool     `tfsdk:"account_enabled"`
	AppRoleAssignmentRequired           types.Bool     `tfsdk:"app_role_assignment_required"`
	ServicePrincipalType                types.String   `tfsdk:"service_principal_type"`
	SignInAudience                      types.String   `tfsdk:"sign_in_audience"`
	PreferredSingleSignOnMode           types.String   `tfsdk:"preferred_single_sign_on_mode"`
	Homepage                            types.String   `tfsdk:"homepage"`
	ErrorUrl                            types.String   `tfsdk:"error_url"`
	PublisherName                       types.String   `tfsdk:"publisher_name"`
	ReplyUrls                           types.Set      `tfsdk:"reply_urls"`
	ServicePrincipalNames               types.Set      `tfsdk:"service_principal_names"`
	Tags                                types.Set      `tfsdk:"tags"`
	DisabledByMicrosoftStatus           types.String   `tfsdk:"disabled_by_microsoft_status"`
	AppOwnerOrganizationID              types.String   `tfsdk:"app_owner_organization_id"`
	LoginUrl                            types.String   `tfsdk:"login_url"`
	LogoutUrl                           types.String   `tfsdk:"logout_url"`
	Notes                               types.String   `tfsdk:"notes"`
	NotificationEmailAddresses          types.Set      `tfsdk:"notification_email_addresses"`
	SamlMetadataUrl                     types.String   `tfsdk:"saml_metadata_url"`
	PreferredTokenSigningKeyEndDateTime types.String   `tfsdk:"preferred_token_signing_key_end_date_time"`
	PreferredTokenSigningKeyThumbprint  types.String   `tfsdk:"preferred_token_signing_key_thumbprint"`
	SamlSingleSignOnSettings            types.Object   `tfsdk:"saml_single_sign_on_settings"`
	VerifiedPublisher                   types.Object   `tfsdk:"verified_publisher"`
	Info                                types.Object   `tfsdk:"info"`
	Timeouts                            timeouts.Value `tfsdk:"timeouts"`
}

// Nested object models
type SamlSingleSignOnSettings struct {
	RelayState types.String `tfsdk:"relay_state"`
}

type VerifiedPublisher struct {
	DisplayName         types.String `tfsdk:"display_name"`
	VerifiedPublisherID types.String `tfsdk:"verified_publisher_id"`
	AddedDateTime       types.String `tfsdk:"added_date_time"`
}

type InformationalUrl struct {
	TermsOfServiceUrl   types.String `tfsdk:"terms_of_service_url"`
	SupportUrl          types.String `tfsdk:"support_url"`
	PrivacyStatementUrl types.String `tfsdk:"privacy_statement_url"`
	MarketingUrl        types.String `tfsdk:"marketing_url"`
	LogoUrl             types.String `tfsdk:"logo_url"`
}
