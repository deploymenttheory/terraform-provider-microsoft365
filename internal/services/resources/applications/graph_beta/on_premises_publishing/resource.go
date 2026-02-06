package graphBetaApplicationsOnPremisesPublishing

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

const (
	ResourceName  = "microsoft365_graph_beta_applications_on_premises_publishing"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180
)

var (
	_ resource.Resource                = &OnPremisesPublishingResource{}
	_ resource.ResourceWithConfigure   = &OnPremisesPublishingResource{}
	_ resource.ResourceWithImportState = &OnPremisesPublishingResource{}
)

func NewOnPremisesPublishingResource() resource.Resource {
	return &OnPremisesPublishingResource{
		ReadPermissions: []string{
			"Application.Read.All",
			"Directory.Read.All",
		},
		WritePermissions: []string{
			"Application.ReadWrite.All",
			"Directory.ReadWrite.All",
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

func (r *OnPremisesPublishingResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages on-premises publishing configuration for an Application. " +
			"This resource configures Application Proxy or Global Secure Access (Private Access) settings for an application.\n\n" +
			"For more information, see the [Microsoft Graph API documentation](https://learn.microsoft.com/en-us/graph/api/resources/onpremisespublishing?view=graph-rest-beta).",
		Attributes: map[string]schema.Attribute{
			"application_id": schema.StringAttribute{
				MarkdownDescription: "The unique identifier (object ID) of the application. This is used as the resource ID for import.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"application_type": schema.StringAttribute{
				MarkdownDescription: "The type of application being published. " +
					"Possible values are: `webapp` (for Application Proxy web apps), `nonwebapp` (for Private Access apps). Default is `webapp`.",
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					stringvalidator.OneOf("webapp", "nonwebapp"),
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
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"external_url": schema.StringAttribute{
				MarkdownDescription: "The external URL for accessing the application through Application Proxy.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"is_accessible_via_ztna_client": schema.BoolAttribute{
				MarkdownDescription: "Whether the application is accessible via the Zero Trust Network Access (ZTNA) client. " +
					"Set to `true` for Private Access applications. Default is `false`.",
				Optional: true,
				Computed: true,
				Default:  booldefault.StaticBool(false),
			},
			"is_http_only_cookie_enabled": schema.BoolAttribute{
				MarkdownDescription: "Indicates if HTTP only cookies should be enabled for Application Proxy. Default is `false`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"is_on_prem_publishing_enabled": schema.BoolAttribute{
				MarkdownDescription: "Indicates if on-premises publishing is enabled for this application. Default is `true`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
			},
			"is_persistent_cookie_enabled": schema.BoolAttribute{
				MarkdownDescription: "Indicates if persistent cookies should be enabled for Application Proxy. Default is `false`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"is_secure_cookie_enabled": schema.BoolAttribute{
				MarkdownDescription: "Indicates if secure cookies should be enabled for Application Proxy. Default is `false`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"is_state_session_enabled": schema.BoolAttribute{
				MarkdownDescription: "Indicates if state session is enabled for Application Proxy. Default is `false`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"is_translate_host_header_enabled": schema.BoolAttribute{
				MarkdownDescription: "Indicates if the host header should be translated for Application Proxy. Default is `false`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"is_translate_links_in_body_enabled": schema.BoolAttribute{
				MarkdownDescription: "Indicates if links in the body should be translated for Application Proxy. Default is `false`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"timeouts": commonschema.ResourceTimeouts(ctx),
		},
	}
}
