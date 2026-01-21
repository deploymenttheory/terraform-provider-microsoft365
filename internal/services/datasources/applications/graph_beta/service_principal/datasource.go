package graphBetaServicePrincipal

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

const (
	datasourceName = "graph_beta_applications_service_principal"
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
	client           *msgraphbetasdk.GraphServiceClient
	ProviderTypeName string
	TypeName         string
	ReadPermissions  []string
}

// Metadata returns the datasource type name.
func (r *ServicePrincipalDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + datasourceName
}

// Configure sets the client for the data source
func (d *ServicePrincipalDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.client = client.SetGraphBetaClientForDataSource(ctx, req, resp, d.TypeName)
}

// Schema defines the schema for the data source
func (d *ServicePrincipalDataSource) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Retrieves service principals from Microsoft Entra ID using the `/servicePrincipals` endpoint. This data source is used to query enterprise applications and managed identities with advanced filtering capabilities.",
		Attributes: map[string]schema.Attribute{
			"filter_type": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Type of filter to apply. Valid values are: `all`, `id`, `app_id`, `display_name`, `odata`.",
				Validators: []validator.String{
					stringvalidator.OneOf("all", "id", "app_id", "display_name", "odata"),
				},
			},
			"filter_value": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Value to filter by. Not required when filter_type is 'all' or 'odata'.",
			},
			"odata_filter": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "OData $filter parameter for filtering results. Only used when filter_type is 'odata'. Example: preferredSingleSignOnMode ne 'notSupported'.",
			},
			"odata_top": schema.Int32Attribute{
				Optional:            true,
				MarkdownDescription: "OData $top parameter to limit the number of results. Only used when filter_type is 'odata'.",
			},
			"odata_skip": schema.Int32Attribute{
				Optional:            true,
				MarkdownDescription: "OData $skip parameter for pagination. Only used when filter_type is 'odata'.",
			},
			"odata_select": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "OData $select parameter to specify which fields to include. Only used when filter_type is 'odata'.",
			},
			"odata_orderby": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "OData $orderby parameter to sort results. Only used when filter_type is 'odata'. Example: displayName.",
			},
			"odata_count": schema.BoolAttribute{
				Optional:            true,
				MarkdownDescription: "OData $count parameter to include count of total results. Only used when filter_type is 'odata'.",
			},
			"odata_search": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "OData $search parameter for full-text search. Only used when filter_type is 'odata'.",
			},
			"odata_expand": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "OData $expand parameter to include related entities. Only used when filter_type is 'odata'.",
			},
			"items": schema.ListNestedAttribute{
				Computed:            true,
				MarkdownDescription: "The list of service principals that match the filter criteria.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The unique identifier for the service principal.",
						},
						"app_id": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The unique identifier for the associated application.",
						},
						"app_display_name": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The display name exposed by the associated application.",
						},
						"display_name": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The display name for the service principal.",
						},
						"deleted_date_time": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The date and time the service principal was deleted.",
						},
						"application_template_id": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "Unique identifier of the applicationTemplate that the servicePrincipal was created from.",
						},
						"account_enabled": schema.BoolAttribute{
							Computed:            true,
							MarkdownDescription: "true if the service principal account is enabled; otherwise, false.",
						},
						"app_role_assignment_required": schema.BoolAttribute{
							Computed:            true,
							MarkdownDescription: "Specifies whether users or other service principals need to be granted an app role assignment for this service principal before users can sign in or apps can get tokens.",
						},
						"service_principal_type": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "Identifies if the service principal represents an Application, a ManagedIdentity, or a legacy application.",
						},
						"sign_in_audience": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "Specifies the Microsoft accounts that are supported for the current application.",
						},
						"preferred_single_sign_on_mode": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "Specifies the single sign-on mode configured for this application.",
						},
						"homepage": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "Home page or landing page of the application.",
						},
						"publisher_name": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The display name of the tenant in which the associated application is registered.",
						},
						"reply_urls": schema.ListAttribute{
							Computed:            true,
							ElementType:         types.StringType,
							MarkdownDescription: "The URLs that user tokens are sent to for sign in with the associated application, or the redirect URIs that OAuth 2.0 authorization codes and access tokens are sent to for the associated application.",
						},
						"tags": schema.ListAttribute{
							Computed:            true,
							ElementType:         types.StringType,
							MarkdownDescription: "Custom strings that can be used to categorize and identify the service principal.",
						},
						"created_date_time": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The date and time the service principal was created.",
						},
						"disabled_by_microsoft_status": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "Specifies whether Microsoft has disabled the registered application.",
						},
						"app_owner_organization_id": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "Contains the tenant id where the application is registered.",
						},
						"login_url": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "Specifies the URL where the service provider redirects the user to Azure AD to authenticate.",
						},
						"logout_url": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "Specifies the URL that will be used by Microsoft's authorization service to logout an user using OpenId Connect front-channel, back-channel or SAML logout protocols.",
						},
						"notes": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "Free text field to capture information about the service principal, typically used for operational purposes.",
						},
						"notification_email_addresses": schema.ListAttribute{
							Computed:            true,
							ElementType:         types.StringType,
							MarkdownDescription: "Specifies the list of email addresses where Azure AD sends a notification when the active certificate is near the expiration date.",
						},
						"error_url": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "Deprecated. Do not use.",
						},
						"service_principal_names": schema.ListAttribute{
							Computed:            true,
							ElementType:         types.StringType,
							MarkdownDescription: "Contains the list of identifiersUris, copied over from the associated application.",
						},
						"saml_metadata_url": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The URL where the service exposes SAML metadata for federation.",
						},
						"preferred_token_signing_key_end_date_time": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "Specifies the expiration date of the keyCredential used for token signing, marked by preferredTokenSigningKeyThumbprint.",
						},
						"preferred_token_signing_key_thumbprint": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "Reserved for internal use only. Do not write or otherwise rely on this property.",
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
							MarkdownDescription: "Basic profile information of the acquired application such as app's marketing, support, terms of service and privacy statement URLs.",
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
					},
				},
			},
			"timeouts": commonschema.ResourceTimeouts(ctx),
		},
	}
}
