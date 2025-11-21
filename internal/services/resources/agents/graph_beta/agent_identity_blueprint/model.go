// REF: https://learn.microsoft.com/en-us/graph/api/resources/agentidentityblueprint?view=graph-rest-beta
package graphBetaAgentsAgentIdentityBlueprint

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// AgentIdentityBlueprintResourceModel describes the Terraform resource data model for an agent identity blueprint.
// Agent identity blueprints serve as templates for creating agent identities within the Microsoft Entra ID ecosystem.
type AgentIdentityBlueprintResourceModel struct {
	ID                         types.String         `tfsdk:"id"`
	AppId                      types.String         `tfsdk:"app_id"`
	IdentifierUris             types.Set            `tfsdk:"identifier_uris"`
	CreatedByAppId             types.String         `tfsdk:"created_by_app_id"`
	CreatedDateTime            types.String         `tfsdk:"created_date_time"`
	Description                types.String         `tfsdk:"description"`
	DisabledByMicrosoftStatus  types.String         `tfsdk:"disabled_by_microsoft_status"`
	DisplayName                types.String         `tfsdk:"display_name"`
	GroupMembershipClaims      types.String         `tfsdk:"group_membership_claims"`
	PublisherDomain            types.String         `tfsdk:"publisher_domain"`
	SignInAudience             types.String         `tfsdk:"sign_in_audience"`
	Tags                       types.Set            `tfsdk:"tags"`
	TokenEncryptionKeyId       types.String         `tfsdk:"token_encryption_key_id"`
	UniqueName                 types.String         `tfsdk:"unique_name"`
	ServiceManagementReference types.String         `tfsdk:"service_management_reference"`
	Info                       *InformationalUrl    `tfsdk:"info"`
	KeyCredentials             []KeyCredential      `tfsdk:"key_credentials"`
	PasswordCredentials        []PasswordCredential `tfsdk:"password_credentials"`
	AppRoles                   []AppRole            `tfsdk:"app_roles"`
	Api                        *ApiApplication      `tfsdk:"api"`
	OptionalClaims             *OptionalClaims      `tfsdk:"optional_claims"`
	VerifiedPublisher          *VerifiedPublisher   `tfsdk:"verified_publisher"`
	Web                        *WebApplication      `tfsdk:"web"`
	Certification              *Certification       `tfsdk:"certification"`
	Timeouts                   timeouts.Value       `tfsdk:"timeouts"`
}

// InformationalUrl represents basic profile information of the agent identity blueprint
type InformationalUrl struct {
	LogoUrl             types.String `tfsdk:"logo_url"`
	MarketingUrl        types.String `tfsdk:"marketing_url"`
	PrivacyStatementUrl types.String `tfsdk:"privacy_statement_url"`
	SupportUrl          types.String `tfsdk:"support_url"`
	TermsOfServiceUrl   types.String `tfsdk:"terms_of_service_url"`
}

// KeyCredential represents a key credential associated with the agent identity blueprint
type KeyCredential struct {
	CustomKeyIdentifier types.String `tfsdk:"custom_key_identifier"`
	DisplayName         types.String `tfsdk:"display_name"`
	EndDateTime         types.String `tfsdk:"end_date_time"`
	Key                 types.String `tfsdk:"key"`
	KeyId               types.String `tfsdk:"key_id"`
	StartDateTime       types.String `tfsdk:"start_date_time"`
	Type                types.String `tfsdk:"type"`
	Usage               types.String `tfsdk:"usage"`
}

// PasswordCredential represents a password credential associated with the agent identity blueprint
type PasswordCredential struct {
	CustomKeyIdentifier types.String `tfsdk:"custom_key_identifier"`
	DisplayName         types.String `tfsdk:"display_name"`
	EndDateTime         types.String `tfsdk:"end_date_time"`
	Hint                types.String `tfsdk:"hint"`
	KeyId               types.String `tfsdk:"key_id"`
	SecretText          types.String `tfsdk:"secret_text"`
	StartDateTime       types.String `tfsdk:"start_date_time"`
}

// AppRole represents an app role that can be requested by client applications
type AppRole struct {
	AllowedMemberTypes types.Set    `tfsdk:"allowed_member_types"`
	Description        types.String `tfsdk:"description"`
	DisplayName        types.String `tfsdk:"display_name"`
	Id                 types.String `tfsdk:"id"`
	IsEnabled          types.Bool   `tfsdk:"is_enabled"`
	Origin             types.String `tfsdk:"origin"`
	Value              types.String `tfsdk:"value"`
}

// ApiApplication represents API-specific settings for the agent identity blueprint
type ApiApplication struct {
	AcceptMappedClaims          types.Bool                 `tfsdk:"accept_mapped_claims"`
	KnownClientApplications     types.Set                  `tfsdk:"known_client_applications"`
	PreAuthorizedApplications   []PreAuthorizedApplication `tfsdk:"pre_authorized_applications"`
	RequestedAccessTokenVersion types.Int64                `tfsdk:"requested_access_token_version"`
	Oauth2PermissionScopes      []PermissionScope          `tfsdk:"oauth2_permission_scopes"`
}

// PreAuthorizedApplication represents a pre-authorized client application
type PreAuthorizedApplication struct {
	AppId                  types.String `tfsdk:"app_id"`
	DelegatedPermissionIds types.Set    `tfsdk:"delegated_permission_ids"`
}

// PermissionScope represents an OAuth 2.0 permission scope
type PermissionScope struct {
	AdminConsentDescription types.String `tfsdk:"admin_consent_description"`
	AdminConsentDisplayName types.String `tfsdk:"admin_consent_display_name"`
	Id                      types.String `tfsdk:"id"`
	IsEnabled               types.Bool   `tfsdk:"is_enabled"`
	Origin                  types.String `tfsdk:"origin"`
	Type                    types.String `tfsdk:"type"`
	UserConsentDescription  types.String `tfsdk:"user_consent_description"`
	UserConsentDisplayName  types.String `tfsdk:"user_consent_display_name"`
	Value                   types.String `tfsdk:"value"`
}

// OptionalClaims represents optional claims configuration
type OptionalClaims struct {
	AccessToken types.Set `tfsdk:"access_token"`
	IdToken     types.Set `tfsdk:"id_token"`
	Saml2Token  types.Set `tfsdk:"saml2_token"`
}

// VerifiedPublisher represents the verified publisher of the agent identity blueprint
type VerifiedPublisher struct {
	AddedDateTime       types.String `tfsdk:"added_date_time"`
	DisplayName         types.String `tfsdk:"display_name"`
	VerifiedPublisherId types.String `tfsdk:"verified_publisher_id"`
}

// WebApplication represents web application settings
type WebApplication struct {
	HomePageUrl           types.String           `tfsdk:"home_page_url"`
	ImplicitGrantSettings *ImplicitGrantSettings `tfsdk:"implicit_grant_settings"`
	LogoutUrl             types.String           `tfsdk:"logout_url"`
	RedirectUris          types.Set              `tfsdk:"redirect_uris"`
	RedirectUriSettings   []RedirectUriSettings  `tfsdk:"redirect_uri_settings"`
}

// ImplicitGrantSettings represents implicit grant flow settings
type ImplicitGrantSettings struct {
	EnableAccessTokenIssuance types.Bool `tfsdk:"enable_access_token_issuance"`
	EnableIdTokenIssuance     types.Bool `tfsdk:"enable_id_token_issuance"`
}

// RedirectUriSettings represents redirect URI settings
type RedirectUriSettings struct {
	Index types.Int64  `tfsdk:"index"`
	Uri   types.String `tfsdk:"uri"`
}

// Certification represents certification information
type Certification struct {
	CertificationDetailsUrl         types.String `tfsdk:"certification_details_url"`
	CertificationExpirationDateTime types.String `tfsdk:"certification_expiration_date_time"`
	IsCertifiedByMicrosoft          types.Bool   `tfsdk:"is_certified_by_microsoft"`
	IsPublisherAttested             types.Bool   `tfsdk:"is_publisher_attested"`
	LastCertificationDateTime       types.String `tfsdk:"last_certification_date_time"`
}
