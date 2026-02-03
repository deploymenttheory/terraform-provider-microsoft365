package graphBetaServicePrincipal

import (
	"context"
	"regexp"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

const (
	DataSourceName = "microsoft365_graph_beta_applications_service_principal"
	ReadTimeout    = 180
)

var (
	// Basic datasource interface (CRUD operations)
	_ datasource.DataSource = &ServicePrincipalDataSource{}

	// Allows the datasource to be configured with the provider client
	_ datasource.DataSourceWithConfigure = &ServicePrincipalDataSource{}
)

func NewServicePrincipalDataSource() datasource.DataSource {
	return &ServicePrincipalDataSource{
		ReadPermissions: []string{
			"Application.ReadWrite.OwnedBy",
			"CustomSecAttributeAssignment.Reade.All",
		},
	}
}

type ServicePrincipalDataSource struct {
	client *msgraphbetasdk.GraphServiceClient

	ReadPermissions []string
}

func (d *ServicePrincipalDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = DataSourceName
}

func (d *ServicePrincipalDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.client = client.SetGraphBetaClientForDataSource(ctx, req, resp, DataSourceName)
}

// Schema defines the schema for the data source
func (d *ServicePrincipalDataSource) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Retrieves information about a Microsoft Entra ID service principal using the `/servicePrincipals` endpoint. This data source is used to query enterprise applications and managed identities by ID, app ID, display name, or advanced OData filtering.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Unique identifier for the service principal object. This property is referred to as **Object ID** in the Microsoft Entra admin center. Key. Not nullable. Read-only. Supports `$filter` (`eq`, `ne`, `not`, `in`).",
			},
			"object_id": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "The object ID of the service principal. One of `object_id`, `app_id`, `display_name`, or `odata_query` must be specified.",
				Validators: []validator.String{
					stringvalidator.AtLeastOneOf(
						path.MatchRoot("app_id"),
						path.MatchRoot("display_name"),
						path.MatchRoot("odata_query"),
					),
					stringvalidator.ConflictsWith(
						path.MatchRoot("app_id"),
						path.MatchRoot("display_name"),
						path.MatchRoot("odata_query"),
					),
					stringvalidator.RegexMatches(
						regexp.MustCompile(constants.GuidRegex),
						"must be a valid GUID in the format 00000000-0000-0000-0000-000000000000",
					),
				},
			},
			"app_id": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "The unique identifier for the associated application (client ID). Not nullable. Read-only. Alternate key. Supports `$filter` (`eq`). One of `object_id`, `app_id`, `display_name`, or `odata_query` must be specified.",
				Validators: []validator.String{
					stringvalidator.AtLeastOneOf(
						path.MatchRoot("object_id"),
						path.MatchRoot("display_name"),
						path.MatchRoot("odata_query"),
					),
					stringvalidator.ConflictsWith(
						path.MatchRoot("object_id"),
						path.MatchRoot("display_name"),
						path.MatchRoot("odata_query"),
					),
					stringvalidator.RegexMatches(
						regexp.MustCompile(constants.GuidRegex),
						"must be a valid GUID in the format 00000000-0000-0000-0000-000000000000",
					),
				},
			},
			"display_name": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "The display name for the service principal. Supports `$filter` (`eq`, `ne`, `not`, `ge`, `le`, `in`, `startsWith`, and `eq` on null values), `$search`, and `$orderby`. One of `object_id`, `app_id`, `display_name`, or `odata_query` must be specified.",
				Validators: []validator.String{
					stringvalidator.AtLeastOneOf(
						path.MatchRoot("object_id"),
						path.MatchRoot("app_id"),
						path.MatchRoot("odata_query"),
					),
					stringvalidator.ConflictsWith(
						path.MatchRoot("object_id"),
						path.MatchRoot("app_id"),
						path.MatchRoot("odata_query"),
					),
				},
			},
			"odata_query": schema.StringAttribute{
				Optional: true,
				MarkdownDescription: "Custom OData filter query. " +
					"Use this for advanced filtering when the standard lookup attributes don't meet your needs. " +
					"Cannot be combined with `object_id`, `app_id`, or `display_name`. " +
					"Example: `displayName eq 'My Service Principal' and servicePrincipalType eq 'Application'`",
				Validators: []validator.String{
					stringvalidator.AtLeastOneOf(
						path.MatchRoot("object_id"),
						path.MatchRoot("app_id"),
						path.MatchRoot("display_name"),
					),
					stringvalidator.ConflictsWith(
						path.MatchRoot("object_id"),
						path.MatchRoot("app_id"),
						path.MatchRoot("display_name"),
					),
				},
			},
			"app_display_name": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The display name exposed by the associated application.",
			},
			"deleted_date_time": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The date and time the service principal was deleted. The DateTimeOffset type represents date and time information using ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z. Read-only.",
			},
			"application_template_id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Unique identifier of the applicationTemplate that the servicePrincipal was created from. Read-only. Supports `$filter` (`eq`, `ne`, `NOT`, `startsWith`).",
			},
			"account_enabled": schema.BoolAttribute{
				Computed:            true,
				MarkdownDescription: "true if the service principal account is enabled; otherwise, false. If set to false, then no users are able to sign in to this app, even if they're assigned to it. Supports `$filter` (`eq`, `ne`, `not`, `in`).",
			},
			"app_role_assignment_required": schema.BoolAttribute{
				Computed:            true,
				MarkdownDescription: "Specifies whether users or other service principals need to be granted an app role assignment for this service principal before users can sign in or apps can get tokens. The default value is false. Not nullable. Supports `$filter` (`eq`, `ne`, `NOT`).",
			},
			"service_principal_type": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Identifies if the service principal represents an Application, a ManagedIdentity, or a legacy application (socialIdp). This is set by Azure AD internally. For a service principal that represents an Application this is set as Application. For a service principal that represent a managed identity this is set as ManagedIdentity. For a service principal representing a legacy app this is set as SocialIdp. Supports `$filter` (`eq`, `ne`, `NOT`, `in`).",
			},
			"sign_in_audience": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Specifies the Microsoft accounts that are supported for the current application. Supported values are `AzureADMyOrg`, `AzureADMultipleOrgs`, `AzureADandPersonalMicrosoftAccount`, `PersonalMicrosoftAccount`. Read-only. Supports `$filter` (`eq`, `ne`, `NOT`, `startsWith`).",
			},
			"preferred_single_sign_on_mode": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Specifies the single sign-on mode configured for this application. Azure AD uses the preferred single sign-on mode to launch the application from Microsoft 365 or the Azure AD My Apps. The supported values are `password`, `saml`, `notSupported`, and `oidc`.",
			},
			"homepage": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Home page or landing page of the application.",
			},
			"error_url": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Deprecated. Do not use.",
			},
			"publisher_name": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The display name of the tenant in which the associated application is registered. Provided only when the application publisher is from a different tenant. Supports `$filter` (`eq`, `ne`, `ge`, `le`, `startsWith`).",
			},
			"reply_urls": schema.SetAttribute{
				Computed:            true,
				ElementType:         types.StringType,
				MarkdownDescription: "The URLs that user tokens are sent to for sign in with the associated application, or the redirect URIs that OAuth 2.0 authorization codes and access tokens are sent to for the associated application. Not nullable.",
			},
			"service_principal_names": schema.SetAttribute{
				Computed:            true,
				ElementType:         types.StringType,
				MarkdownDescription: "Contains the list of identifiersUris, copied over from the associated application. More values can be added to hybrid applications. These values can be used to identify the permissions exposed by this app within Azure AD. Supports `$filter` (`eq`, `ne`, `ge`, `le`, `startsWith`).",
			},
			"tags": schema.SetAttribute{
				Computed:            true,
				ElementType:         types.StringType,
				MarkdownDescription: "Custom strings that can be used to categorize and identify the service principal. Not nullable. Supports `$filter` (`eq`, `ne`, `NOT`, `ge`, `le`, `startsWith`).",
			},
			"disabled_by_microsoft_status": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Specifies whether Microsoft has disabled the registered application. Possible values are: `null` (default value), `NotDisabled`, and `DisabledDueToViolationOfServicesAgreement` (reasons include suspicious, abusive, or malicious activity, or a violation of the Microsoft Services Agreement). Supports `$filter` (`eq`, `ne`, `NOT`).",
			},
			"app_owner_organization_id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Contains the tenant ID where the application is registered. This is applicable only to service principals backed by applications. Supports `$filter` (`eq`, `ne`, `NOT`, `ge`, `le`).",
			},
			"login_url": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Specifies the URL where the service provider redirects the user to Azure AD to authenticate. Azure AD uses the URL to launch the application from Microsoft 365 or the Azure AD My Apps. When blank, Azure AD performs IdP-initiated sign-on for applications configured with SAML-based single sign-on.",
			},
			"logout_url": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Specifies the URL that the Microsoft's authorization service uses to sign out a user using OpenId Connect front-channel, back-channel, or SAML sign out protocols.",
			},
			"notes": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Free text field to capture information about the service principal, typically used for operational purposes. Maximum allowed size is 1,024 characters.",
			},
			"notification_email_addresses": schema.SetAttribute{
				Computed:            true,
				ElementType:         types.StringType,
				MarkdownDescription: "Specifies the list of email addresses where Azure AD sends a notification when the active certificate is near the expiration date. This is only for the certificates used to sign the SAML token issued for Azure AD Gallery applications.",
			},
			"saml_metadata_url": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The URL where the service exposes SAML metadata for federation.",
			},
			"preferred_token_signing_key_end_date_time": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Specifies the expiration date of the keyCredential used for token signing, marked by preferredTokenSigningKeyThumbprint. Updating this attribute isn't currently supported. For details, see ServicePrincipal property differences. The Timestamp type represents date and time information using ISO 8601 format and is always in UTC time.",
			},
			"preferred_token_signing_key_thumbprint": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "This property can be used on SAML applications (apps that have preferredSingleSignOnMode set to saml) to control which certificate is used to sign the SAML responses. For applications that aren't SAML, don't write or otherwise rely on this property.",
			},
			"saml_single_sign_on_settings": schema.SingleNestedAttribute{
				Computed:            true,
				MarkdownDescription: "The collection for settings related to saml single sign-on.",
				Attributes: map[string]schema.Attribute{
					"relay_state": schema.StringAttribute{
						Computed:            true,
						MarkdownDescription: "The relative URI the service provider would redirect to after completion of the single sign-on flow.",
					},
				},
			},
			"verified_publisher": schema.SingleNestedAttribute{
				Computed:            true,
				MarkdownDescription: "Specifies the verified publisher of the application which this service principal represents.",
				Attributes: map[string]schema.Attribute{
					"display_name": schema.StringAttribute{
						Computed:            true,
						MarkdownDescription: "The verified publisher name from the app publisher's Microsoft Partner Network (MPN) account.",
					},
					"verified_publisher_id": schema.StringAttribute{
						Computed:            true,
						MarkdownDescription: "The ID of the verified publisher from the app publisher's Partner Center account.",
					},
					"added_date_time": schema.StringAttribute{
						Computed:            true,
						MarkdownDescription: "The timestamp when the verified publisher was first added or most recently updated.",
					},
				},
			},
			"info": schema.SingleNestedAttribute{
				Computed:            true,
				MarkdownDescription: "Basic profile information of the acquired application such as app's marketing, support, terms of service and privacy statement URLs. The terms of service and privacy statement are surfaced to users through the user consent experience. For more information, see How to: Add Terms of service and privacy statement for registered Azure AD apps.",
				Attributes: map[string]schema.Attribute{
					"terms_of_service_url": schema.StringAttribute{
						Computed:            true,
						MarkdownDescription: "Link to the application's terms of service statement.",
					},
					"support_url": schema.StringAttribute{
						Computed:            true,
						MarkdownDescription: "Link to the application's support page.",
					},
					"privacy_statement_url": schema.StringAttribute{
						Computed:            true,
						MarkdownDescription: "Link to the application's privacy statement.",
					},
					"marketing_url": schema.StringAttribute{
						Computed:            true,
						MarkdownDescription: "Link to the application's marketing page.",
					},
					"logo_url": schema.StringAttribute{
						Computed:            true,
						MarkdownDescription: "CDN URL to the application's logo.",
					},
				},
			},
			"timeouts": commonschema.DatasourceTimeouts(ctx),
		},
	}
}
