package graphBetaBrowserSiteList

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/schema"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

const (
	ResourceName  = "graph_beta_m365_admin_browser_site_list"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180
)

var (
	// Basic resource interface (CRUD operations)
	_ resource.Resource = &BrowserSiteListResource{}

	// Allows the resource to be configured with the provider client
	_ resource.ResourceWithConfigure = &BrowserSiteListResource{}

	// Enables import functionality
	_ resource.ResourceWithImportState = &BrowserSiteListResource{}

	// Enables plan modification/diff suppression
	_ resource.ResourceWithModifyPlan = &BrowserSiteListResource{}
)

func NewBrowserSiteListResource() resource.Resource {
	return &BrowserSiteListResource{
		ReadPermissions: []string{
			"BrowserSiteLists.Read.All",
		},
		WritePermissions: []string{
			"BrowserSiteLists.ReadWrite.All",
		},
		ResourcePath: "/admin/edge/internetExplorerMode/siteLists",
	}
}

type BrowserSiteListResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ProviderTypeName string
	TypeName         string
	ReadPermissions  []string
	WritePermissions []string
	ResourcePath     string
}

// Metadata returns the resource type name.
func (r *BrowserSiteListResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + ResourceName
}

// Configure sets the client for the resource.
func (r *BrowserSiteListResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = common.SetGraphBetaClientForResource(ctx, req, resp, r.TypeName)
}

// ImportState imports the resource state.
func (r *BrowserSiteListResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *BrowserSiteListResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a browser site list in Microsoft 365 Admin Centre. Settings > Org settings > Microsoft Edge site lists.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The unique identifier for the site list.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"description": schema.StringAttribute{
				Description: "The description of the site list.",
				Optional:    true,
			},
			"display_name": schema.StringAttribute{
				Description: "The name of the site list.",
				Required:    true,
			},
			"last_modified_date_time": schema.StringAttribute{
				Description: "The date and time when the site list was last modified.",
				Computed:    true,
			},
			"published_date_time": schema.StringAttribute{
				Description: "The date and time when the site list was published.",
				Computed:    true,
			},
			"revision": schema.StringAttribute{
				Description: "The current revision of the site list.",
				Computed:    true,
			},
			"status": schema.StringAttribute{
				Description: "The current status of the site list.",
				Computed:    true,
				Validators: []validator.String{
					stringvalidator.OneOf(
						"draft",
						"published",
						"pending",
						"unknownFutureValue"),
				},
			},
			"timeouts": commonschema.Timeouts(ctx),
		},
	}
}
