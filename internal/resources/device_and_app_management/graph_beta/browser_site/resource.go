package graphBetaBrowserSite

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/schema"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

const (
	ResourceName  = "graph_beta_device_and_app_management_browser_site"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180
)

var (
	// Basic resource interface (CRUD operations)
	_ resource.Resource = &BrowserSiteResource{}

	// Allows the resource to be configured with the provider client
	_ resource.ResourceWithConfigure = &BrowserSiteResource{}

	// Enables import functionality
	_ resource.ResourceWithImportState = &BrowserSiteResource{}

	// Enables plan modification/diff suppression
	_ resource.ResourceWithModifyPlan = &BrowserSiteResource{}
)

func NewBrowserSiteResource() resource.Resource {
	return &BrowserSiteResource{
		ReadPermissions: []string{
			"BrowserSiteLists.Read.All",
		},
		WritePermissions: []string{
			"BrowserSiteLists.ReadWrite.All",
		},
		ResourcePath: "/admin/edge/internetExplorerMode/siteLists/{browserSiteList-id}/sites",
	}
}

type BrowserSiteResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ProviderTypeName string
	TypeName         string
	ReadPermissions  []string
	WritePermissions []string
	ResourcePath     string
}

// Metadata returns the resource type name.
func (r *BrowserSiteResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + ResourceName
}

// Configure sets the client for the resource.
func (r *BrowserSiteResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = common.SetGraphBetaClientForResource(ctx, req, resp, r.TypeName)
}

// ImportState imports the resource state.
func (r *BrowserSiteResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *BrowserSiteResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a browser site in Microsoft Intune.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The unique identifier for the site.",
				Computed:    true,
			},
			"browser_site_list_assignment_id": schema.StringAttribute{
				Required:    true,
				Description: "The browser site list id to assign this browser site to.",
			},
			"allow_redirect": schema.BoolAttribute{
				Description: "Controls the behavior of redirected sites. If `true`, indicates that the site will open in Internet Explorer 11 or Microsoft Edge even if the site is navigated to as part of a HTTP or meta refresh redirection chain.",
				Required:    true,
			},
			"comment": schema.StringAttribute{
				Description: "The comment for the site.",
				Optional:    true,
			},
			"compatibility_mode": schema.StringAttribute{
				Description: "Controls what compatibility setting is used for specific sites or domains.",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.OneOf(
						"default",
						"internetExplorer8Enterprise",
						"internetExplorer7Enterprise",
						"internetExplorer11",
						"internetExplorer10",
						"internetExplorer9",
						"internetExplorer8",
						"internetExplorer7",
						"internetExplorer5",
						"unknownFutureValue"),
				},
			},
			"created_date_time": schema.StringAttribute{
				Description: "The date and time when the site was created.",
				Computed:    true,
			},
			"deleted_date_time": schema.StringAttribute{
				Description: "The date and time when the site was deleted.",
				Computed:    true,
			},
			"history": schema.ListNestedAttribute{
				Description: "The history of modifications applied to the site.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"allow_redirect": schema.BoolAttribute{
							Description: "Controls the behavior of redirected sites for this history entry.",
							Computed:    true,
						},
						"comment": schema.StringAttribute{
							Description: "The comment for this history entry.",
							Computed:    true,
						},
						"compatibility_mode": schema.StringAttribute{
							Description: "The compatibility mode for this history entry.",
							Computed:    true,
						},
						"merge_type": schema.StringAttribute{
							Description: "The merge type for this history entry.",
							Computed:    true,
						},
						"published_date_time": schema.StringAttribute{
							Description: "The date and time when this history entry was published.",
							Computed:    true,
						},
						"target_environment": schema.StringAttribute{
							Description: "The target environment for this history entry.",
							Computed:    true,
						},
					},
				},
			},
			"last_modified_date_time": schema.StringAttribute{
				Description: "The date and time when the site was last modified.",
				Computed:    true,
			},
			"merge_type": schema.StringAttribute{
				Description: "The merge type of the site.",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.OneOf("noMerge", "default", "unknownFutureValue"),
				},
			},
			"status": schema.StringAttribute{
				Description: "Indicates the status of the site.",
				Computed:    true,
			},
			"target_environment": schema.StringAttribute{
				Description: "The target environment that the site should open in.",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.OneOf("internetExplorerMode", "internetExplorer11", "microsoftEdge", "configurable", "none", "unknownFutureValue"),
				},
			},
			"web_url": schema.StringAttribute{
				Description: "The URL of the site.",
				Required:    true,
			},
			"timeouts": commonschema.Timeouts(ctx),
		},
	}
}