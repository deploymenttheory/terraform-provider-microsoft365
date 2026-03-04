package graphBetaApplicationsOnPremisesPublishing

import (
	"context"
	"regexp"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/validate/attribute"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/identityschema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

const (
	ResourceName  = "microsoft365_graph_beta_applications_application_on_premises_publishing"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180
)

var (
	// Basic resource interface (CRUD operations)
	_ resource.Resource = &OnPremisesPublishingResource{}

	// Allows the resource to be configured with the provider client
	_ resource.ResourceWithConfigure = &OnPremisesPublishingResource{}

	// Enables import functionality
	_ resource.ResourceWithImportState = &OnPremisesPublishingResource{}

	// Enables identity schema for list resource support
	_ resource.ResourceWithIdentity = &OnPremisesPublishingResource{}
)

func NewOnPremisesPublishingResource() resource.Resource {
	return &OnPremisesPublishingResource{
		ReadPermissions: []string{
			"Application.Read.All",
			"Directory.Read.All",
		},
		WritePermissions: []string{
			"Application.ReadWrite.All",
		},
		ResourcePath: "/applications/{application-id}",
	}
}

type OnPremisesPublishingResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ReadPermissions  []string
	WritePermissions []string
	ResourcePath     string
}

func (r *OnPremisesPublishingResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = ResourceName
}

func (r *OnPremisesPublishingResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = client.SetGraphBetaClientForResource(ctx, req, resp, ResourceName)
}

func (r *OnPremisesPublishingResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("application_id"), req, resp)
}

// IdentitySchema defines the identity schema for this resource, used by list operations to uniquely identify instances
func (r *OnPremisesPublishingResource) IdentitySchema(ctx context.Context, req resource.IdentitySchemaRequest, resp *resource.IdentitySchemaResponse) {
	resp.IdentitySchema = identityschema.Schema{
		Attributes: map[string]identityschema.Attribute{
			"id": identityschema.StringAttribute{
				RequiredForImport: true,
			},
		},
	}
}

func (r *OnPremisesPublishingResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages on-premises publishing configuration for a Microsoft Entra Application using the `/applications/{application-id}` endpoint. " +
			"This resource configures the `onPremisesPublishing` property which enables Application Proxy or Global Secure Access (Private Access) settings. " +
			"Use this resource to enable non-web applications for Private Access or to configure web applications for Application Proxy.\n\n" +
			"For more information, see the [Microsoft Graph API documentation](https://learn.microsoft.com/en-us/graph/api/resources/onpremisespublishing?view=graph-rest-beta).",
		Attributes: map[string]schema.Attribute{
			"application_id": schema.StringAttribute{
				MarkdownDescription: "The unique identifier (object ID) of the application. This is used as the resource ID for import.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"alternate_url": schema.StringAttribute{
				MarkdownDescription: "User-friendly URL pointing to traffic manager in front of multiple app proxy applications.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"application_server_timeout": schema.StringAttribute{
				MarkdownDescription: "Duration the connector waits for a response from the backend application before closing the connection. " +
					"Possible values are: `default` (85 seconds), `long` (180 seconds).",
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					stringvalidator.OneOf("default", "long"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"application_type": schema.StringAttribute{
				MarkdownDescription: "System-defined value indicating whether this application is an application proxy configured application. " +
					"Possible values are: `quickaccessapp`, `nonwebapp`.",
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					stringvalidator.OneOf("quickaccessapp", "nonwebapp"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"external_authentication_type": schema.StringAttribute{
				MarkdownDescription: "The external authentication type for Application Proxy. " +
					"Possible values are: `passthru`, `aadPreAuthentication`. This is typically used for webapp type.",
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					stringvalidator.OneOf("passthru", "aadPreAuthentication"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"internal_url": schema.StringAttribute{
				MarkdownDescription: "The internal URL of the application. This is the URL that users are redirected to when accessing the application through Application Proxy.",
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{
					attribute.RegexMatches(regexp.MustCompile(constants.HttpOrHttpsUrlRegex), "must be a valid HTTP or HTTPS URL"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"external_url": schema.StringAttribute{
				MarkdownDescription: "The external URL for accessing the application through Application Proxy.",
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{
					attribute.RegexMatches(regexp.MustCompile(constants.HttpOrHttpsUrlRegex), "must be a valid HTTP or HTTPS URL"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"is_accessible_via_ztna_client": schema.BoolAttribute{
				MarkdownDescription: "Whether the application is accessible via the Zero Trust Network Access (ZTNA) client. " +
					"Set to `true` for Private Access applications.",
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"is_backend_certificate_validation_enabled": schema.BoolAttribute{
				MarkdownDescription: "Indicates whether backend SSL certificate validation is enabled for the application. " +
					"For all new Application Proxy apps, the property is set to `true` by default. For all existing apps, the property is set to `false`.",
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"is_continuous_access_evaluation_enabled": schema.BoolAttribute{
				MarkdownDescription: "Indicates whether continuous access evaluation is enabled for Application Proxy application. " +
					"For all Application Proxy apps, the property is set to `true` by default.",
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"is_dns_resolution_enabled": schema.BoolAttribute{
				MarkdownDescription: "Indicates whether Microsoft Entra Private Access should handle DNS resolution. Default is `false`.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"is_http_only_cookie_enabled": schema.BoolAttribute{
				MarkdownDescription: "Indicates if HTTP only cookies should be enabled for Application Proxy.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"is_on_prem_publishing_enabled": schema.BoolAttribute{
				MarkdownDescription: "Indicates if the application is currently being published via Application Proxy.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"is_persistent_cookie_enabled": schema.BoolAttribute{
				MarkdownDescription: "Indicates if persistent cookies should be enabled for Application Proxy.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"is_secure_cookie_enabled": schema.BoolAttribute{
				MarkdownDescription: "Indicates if secure cookies should be enabled for Application Proxy.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"is_state_session_enabled": schema.BoolAttribute{
				MarkdownDescription: "Indicates if state session is enabled for Application Proxy.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"is_translate_host_header_enabled": schema.BoolAttribute{
				MarkdownDescription: "Indicates if the host header should be translated for Application Proxy.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"is_translate_links_in_body_enabled": schema.BoolAttribute{
				MarkdownDescription: "Indicates if links in the body should be translated for Application Proxy.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"use_alternate_url_for_translation_and_redirect": schema.BoolAttribute{
				MarkdownDescription: "Indicates whether the application should use `alternate_url` instead of `external_url`.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"waf_provider": schema.StringAttribute{
				MarkdownDescription: "Web Application Firewall (WAF) provider for the application.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"timeouts": commonschema.ResourceTimeouts(ctx),
		},
	}
}
