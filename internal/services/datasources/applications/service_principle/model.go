// REF: https://learn.microsoft.com/en-us/graph/api/resources/serviceprincipal?view=graph-rest-beta

package graphBetaServicePrincipal

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ServicePrincipalDataSourceModel defines the data source model
type ServicePrincipalDataSourceModel struct {
	FilterType     types.String                 `tfsdk:"filter_type"`     // Required field to specify how to filter
	FilterValue    types.String                 `tfsdk:"filter_value"`    // Value to filter by (not used for "all" or "odata")
	ODataFilter    types.String                 `tfsdk:"odata_filter"`    // OData filter parameter
	ODataTop       types.Int32                  `tfsdk:"odata_top"`       // OData top parameter for limiting results
	ODataSkip      types.Int32                  `tfsdk:"odata_skip"`      // OData skip parameter for pagination
	ODataSelect    types.String                 `tfsdk:"odata_select"`    // OData select parameter for field selection
	ODataOrderBy   types.String                 `tfsdk:"odata_orderby"`   // OData orderby parameter for sorting
	ODataCount     types.Bool                   `tfsdk:"odata_count"`     // OData count parameter
	ODataSearch    types.String                 `tfsdk:"odata_search"`    // OData search parameter
	ODataExpand    types.String                 `tfsdk:"odata_expand"`    // OData expand parameter
	Items          []ServicePrincipalModel      `tfsdk:"items"`           // List of service principals that match the filters
	Timeouts       timeouts.Value               `tfsdk:"timeouts"`
}

// ServicePrincipalModel represents a single service principal
type ServicePrincipalModel struct {
	ID                                        types.String                         `tfsdk:"id"`
	AppID                                     types.String                         `tfsdk:"app_id"`
	AppDisplayName                            types.String                         `tfsdk:"app_display_name"`
	DisplayName                               types.String                         `tfsdk:"display_name"`
	DeletedDateTime                           types.String                         `tfsdk:"deleted_date_time"`
	CreatedDateTime                           types.String                         `tfsdk:"created_date_time"`
	ApplicationTemplateID                     types.String                         `tfsdk:"application_template_id"`
	AccountEnabled                            types.Bool                           `tfsdk:"account_enabled"`
	AppRoleAssignmentRequired                 types.Bool                           `tfsdk:"app_role_assignment_required"`
	ServicePrincipalType                      types.String                         `tfsdk:"service_principal_type"`
	SignInAudience                            types.String                         `tfsdk:"sign_in_audience"`
	PreferredSingleSignOnMode                 types.String                         `tfsdk:"preferred_single_sign_on_mode"`
	Homepage                                  types.String                         `tfsdk:"homepage"`
	ErrorUrl                                  types.String                         `tfsdk:"error_url"`
	PublisherName                             types.String                         `tfsdk:"publisher_name"`
	ReplyUrls                                 []types.String                       `tfsdk:"reply_urls"`
	ServicePrincipalNames                     []types.String                       `tfsdk:"service_principal_names"`
	Tags                                      []types.String                       `tfsdk:"tags"`
	DisabledByMicrosoftStatus                 types.String                         `tfsdk:"disabled_by_microsoft_status"`
	AppOwnerOrganizationID                    types.String                         `tfsdk:"app_owner_organization_id"`
	LoginUrl                                  types.String                         `tfsdk:"login_url"`
	LogoutUrl                                 types.String                         `tfsdk:"logout_url"`
	Notes                                     types.String                         `tfsdk:"notes"`
	NotificationEmailAddresses                []types.String                       `tfsdk:"notification_email_addresses"`
	SamlMetadataUrl                           types.String                         `tfsdk:"saml_metadata_url"`
	PreferredTokenSigningKeyEndDateTime       types.String                         `tfsdk:"preferred_token_signing_key_end_date_time"`
	PreferredTokenSigningKeyThumbprint        types.String                         `tfsdk:"preferred_token_signing_key_thumbprint"`
	SamlSingleSignOnSettings                  *SamlSingleSignOnSettingsModel       `tfsdk:"saml_single_sign_on_settings"`
	VerifiedPublisher                         *VerifiedPublisherModel              `tfsdk:"verified_publisher"`
	Info                                      *InformationalUrlModel               `tfsdk:"info"`
}

// SamlSingleSignOnSettingsModel represents SAML SSO settings
type SamlSingleSignOnSettingsModel struct {
	RelayState types.String `tfsdk:"relay_state"`
}

// VerifiedPublisherModel represents verified publisher information
type VerifiedPublisherModel struct {
	DisplayName          types.String `tfsdk:"display_name"`
	VerifiedPublisherID  types.String `tfsdk:"verified_publisher_id"`
	AddedDateTime        types.String `tfsdk:"added_date_time"`
}

// InformationalUrlModel represents application information URLs
type InformationalUrlModel struct {
	TermsOfServiceUrl      types.String `tfsdk:"terms_of_service_url"`
	SupportUrl             types.String `tfsdk:"support_url"`
	PrivacyStatementUrl    types.String `tfsdk:"privacy_statement_url"`
	MarketingUrl           types.String `tfsdk:"marketing_url"`
	LogoUrl                types.String `tfsdk:"logo_url"`
}