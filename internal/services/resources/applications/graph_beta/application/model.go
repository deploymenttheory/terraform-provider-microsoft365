// REF: https://learn.microsoft.com/en-us/graph/api/application-post-applications?view=graph-rest-beta&tabs=go
// REF: https://learn.microsoft.com/en-us/graph/api/resources/application?view=graph-rest-beta
package graphBetaApplication

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ApplicationResourceModel represents the Terraform resource model for Microsoft Entra Applications
type ApplicationResourceModel struct {
	ID                             types.String   `tfsdk:"id"`
	AppId                          types.String   `tfsdk:"app_id"`
	DisplayName                    types.String   `tfsdk:"display_name"`
	Description                    types.String   `tfsdk:"description"`
	SignInAudience                 types.String   `tfsdk:"sign_in_audience"`
	SignInAudienceRestrictions     types.Object   `tfsdk:"sign_in_audience_restrictions"` // Optional+Computed must be types.Object
	IdentifierUris                 types.Set      `tfsdk:"identifier_uris"`
	GroupMembershipClaims          types.Set      `tfsdk:"group_membership_claims"`
	Notes                          types.String   `tfsdk:"notes"`
	IsDeviceOnlyAuthSupported      types.Bool     `tfsdk:"is_device_only_auth_supported"`
	IsFallbackPublicClient         types.Bool     `tfsdk:"is_fallback_public_client"`
	// OAuth2RequirePostResponse      types.Bool     `tfsdk:"oauth2_require_post_response"` // Field doesn't exist in SDK
	ServiceManagementReference     types.String   `tfsdk:"service_management_reference"`
	Tags                           types.Set      `tfsdk:"tags"`
	DisabledByMicrosoftStatus      types.String   `tfsdk:"disabled_by_microsoft_status"`
	PublisherDomain                types.String   `tfsdk:"publisher_domain"`
	CreatedDateTime                types.String   `tfsdk:"created_date_time"`
	DeletedDateTime                types.String   `tfsdk:"deleted_date_time"`
	Api                            types.Object   `tfsdk:"api"`                     // Optional+Computed must be types.Object
	AppRoles                       types.Set      `tfsdk:"app_roles"`
	Info                           types.Object   `tfsdk:"info"`                    // Optional+Computed must be types.Object
	KeyCredentials                 types.Set      `tfsdk:"key_credentials"`
	PasswordCredentials            types.Set      `tfsdk:"password_credentials"`
	OptionalClaims                 types.Object   `tfsdk:"optional_claims"`         // Optional+Computed must be types.Object
	ParentalControlSettings        types.Object   `tfsdk:"parental_control_settings"` // Optional+Computed must be types.Object
	PublicClient                   types.Object   `tfsdk:"public_client"`           // Optional+Computed must be types.Object
	RequiredResourceAccess         types.Set      `tfsdk:"required_resource_access"`
	Spa                            types.Object   `tfsdk:"spa"`                     // Optional+Computed must be types.Object
	Web                            types.Object   `tfsdk:"web"`                     // Optional+Computed must be types.Object
	OwnerUserIds                   types.Set      `tfsdk:"owner_user_ids"`
	PreventDuplicateNames          types.Bool     `tfsdk:"prevent_duplicate_names"`
	HardDelete                     types.Bool     `tfsdk:"hard_delete"`
	Timeouts                       timeouts.Value `tfsdk:"timeouts"`
}

// SignInAudienceRestrictions represents the sign-in audience restrictions for multitenant applications
type SignInAudienceRestrictions struct {
	ODataType           types.String `tfsdk:"odata_type"`
	IsHomeTenantAllowed types.Bool   `tfsdk:"is_home_tenant_allowed"`
	AllowedTenantIds    types.Set    `tfsdk:"allowed_tenant_ids"`
}

// ApplicationApi represents the API configuration for the application
type ApplicationApi struct {
	AcceptMappedClaims          types.Bool  `tfsdk:"accept_mapped_claims"`
	KnownClientApplications     types.Set   `tfsdk:"known_client_applications"`
	OAuth2PermissionScopes      types.Set   `tfsdk:"oauth2_permission_scopes"`
	PreAuthorizedApplications   types.Set   `tfsdk:"pre_authorized_applications"`
	RequestedAccessTokenVersion types.Int32 `tfsdk:"requested_access_token_version"`
}

// ApplicationApiPermissionScope represents an OAuth2 permission scope exposed by the application
type ApplicationApiPermissionScope struct {
	ID                      types.String `tfsdk:"id"`
	AdminConsentDescription types.String `tfsdk:"admin_consent_description"`
	AdminConsentDisplayName types.String `tfsdk:"admin_consent_display_name"`
	IsEnabled               types.Bool   `tfsdk:"is_enabled"`
	Type                    types.String `tfsdk:"type"`
	UserConsentDescription  types.String `tfsdk:"user_consent_description"`
	UserConsentDisplayName  types.String `tfsdk:"user_consent_display_name"`
	Value                   types.String `tfsdk:"value"`
}

// ApplicationApiPreAuthorizedApplication represents a pre-authorized client application
type ApplicationApiPreAuthorizedApplication struct {
	AppId               types.String `tfsdk:"app_id"`
	DelegatedPermission types.Set    `tfsdk:"delegated_permission_ids"`
}

// ApplicationAppRole represents an app role that can be assigned to users, groups, or service principals
type ApplicationAppRole struct {
	ID                 types.String `tfsdk:"id"`
	AllowedMemberTypes types.Set    `tfsdk:"allowed_member_types"`
	Description        types.String `tfsdk:"description"`
	DisplayName        types.String `tfsdk:"display_name"`
	IsEnabled          types.Bool   `tfsdk:"is_enabled"`
	Origin             types.String `tfsdk:"origin"`
	Value              types.String `tfsdk:"value"`
}

// ApplicationInformationalUrl represents URLs for the application's info
type ApplicationInformationalUrl struct {
	LogoUrl             types.String `tfsdk:"logo_url"`
	MarketingUrl        types.String `tfsdk:"marketing_url"`
	PrivacyStatementUrl types.String `tfsdk:"privacy_statement_url"`
	SupportUrl          types.String `tfsdk:"support_url"`
	TermsOfServiceUrl   types.String `tfsdk:"terms_of_service_url"`
}

// ApplicationKeyCredential represents a certificate credential
type ApplicationKeyCredential struct {
	CustomKeyIdentifier types.String `tfsdk:"custom_key_identifier"`
	DisplayName         types.String `tfsdk:"display_name"`
	EndDateTime         types.String `tfsdk:"end_date_time"`
	Key                 types.String `tfsdk:"key"`
	KeyId               types.String `tfsdk:"key_id"`
	StartDateTime       types.String `tfsdk:"start_date_time"`
	Type                types.String `tfsdk:"type"`
	Usage               types.String `tfsdk:"usage"`
}

// ApplicationPasswordCredential represents a password credential
type ApplicationPasswordCredential struct {
	CustomKeyIdentifier types.String `tfsdk:"custom_key_identifier"`
	DisplayName         types.String `tfsdk:"display_name"`
	EndDateTime         types.String `tfsdk:"end_date_time"`
	Hint                types.String `tfsdk:"hint"`
	KeyId               types.String `tfsdk:"key_id"`
	SecretText          types.String `tfsdk:"secret_text"`
	StartDateTime       types.String `tfsdk:"start_date_time"`
}

// ApplicationOptionalClaims represents optional claims configuration
type ApplicationOptionalClaims struct {
	AccessToken types.Set `tfsdk:"access_token"`
	IdToken     types.Set `tfsdk:"id_token"`
	Saml2Token  types.Set `tfsdk:"saml2_token"`
}

// ApplicationOptionalClaim represents a single optional claim
type ApplicationOptionalClaim struct {
	Name                 types.String `tfsdk:"name"`
	Source               types.String `tfsdk:"source"`
	Essential            types.Bool   `tfsdk:"essential"`
	AdditionalProperties types.Set    `tfsdk:"additional_properties"`
}

// ApplicationParentalControlSettings represents parental control settings
type ApplicationParentalControlSettings struct {
	CountriesBlockedForMinors types.Set    `tfsdk:"countries_blocked_for_minors"`
	LegalAgeGroupRule         types.String `tfsdk:"legal_age_group_rule"`
}

// ApplicationPublicClient represents public client configuration
type ApplicationPublicClient struct {
	RedirectUris types.Set `tfsdk:"redirect_uris"`
}

// ApplicationRequiredResourceAccess represents required resource access (API permissions)
type ApplicationRequiredResourceAccess struct {
	ResourceAppId  types.String `tfsdk:"resource_app_id"`
	ResourceAccess types.Set    `tfsdk:"resource_access"`
}

// ApplicationResourceAccess represents a single resource access (permission)
type ApplicationResourceAccess struct {
	ID   types.String `tfsdk:"id"`
	Type types.String `tfsdk:"type"`
}

// ApplicationSpa represents single page application configuration
type ApplicationSpa struct {
	RedirectUris types.Set `tfsdk:"redirect_uris"`
}

// ApplicationWeb represents web application configuration
type ApplicationWeb struct {
	HomePageUrl          types.String `tfsdk:"home_page_url"`
	LogoutUrl            types.String `tfsdk:"logout_url"`
	RedirectUris         types.Set    `tfsdk:"redirect_uris"`
	ImplicitGrantSettings types.Object `tfsdk:"implicit_grant_settings"` // Optional+Computed must be types.Object
	RedirectUriSettings   types.Set    `tfsdk:"redirect_uri_settings"`
}

// ApplicationWebImplicitGrantSettings represents implicit grant settings
type ApplicationWebImplicitGrantSettings struct {
	EnableAccessTokenIssuance types.Bool `tfsdk:"enable_access_token_issuance"`
	EnableIdTokenIssuance     types.Bool `tfsdk:"enable_id_token_issuance"`
}

// ApplicationWebRedirectUriSettings represents redirect URI settings
type ApplicationWebRedirectUriSettings struct {
	Uri   types.String `tfsdk:"uri"`
	Index types.Int32  `tfsdk:"index"`
}
