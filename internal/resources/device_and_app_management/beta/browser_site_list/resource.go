package graphbetabrowsersite

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

// Ensure provider defined types fully satisfy framework interfaces
var _ resource.Resource = &BrowserSiteListResource{}
var _ resource.ResourceWithConfigure = &BrowserSiteListResource{}
var _ resource.ResourceWithImportState = &BrowserSiteListResource{}

func NewBrowserSiteListResource() resource.Resource {
	return &BrowserSiteListResource{}
}

type BrowserSiteListResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ProviderTypeName string
	TypeName         string
}

// GetID returns the ID of a resource from the state model.
func (s *BrowserSiteListResourceModel) GetID() string {
	return s.ID.ValueString()
}

// GetTypeName returns the type name of the resource from the state model.
func (r *BrowserSiteListResource) GetTypeName() string {
	return r.TypeName
}

// Metadata returns the resource type name.
func (r *BrowserSiteListResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_graph_beta_device_and_app_management_browser_site_list"
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
		Description: "Manages a browser site list in Microsoft Intune.",
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
					stringvalidator.OneOf("draft", "published", "pending", "unknownFutureValue"),
				},
			},
			"timeouts": commonschema.Timeouts(ctx),
		},
	}
}
