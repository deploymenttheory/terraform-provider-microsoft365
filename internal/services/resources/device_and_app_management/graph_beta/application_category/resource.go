package graphBetaApplicationCategory

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	planmodifiers "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/plan_modifiers"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"

	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

const (
	ResourceName  = "microsoft365_graph_beta_device_and_app_management_application_category"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180
)

var (
	// Basic resource interface (CRUD operations)
	_ resource.Resource = &ApplicationCategoryResource{}

	// Allows the resource to be configured with the provider client
	_ resource.ResourceWithConfigure = &ApplicationCategoryResource{}

	// Enables import functionality
	_ resource.ResourceWithImportState = &ApplicationCategoryResource{}

	// Enables plan modification/diff suppression
	_ resource.ResourceWithModifyPlan = &ApplicationCategoryResource{}
)

func NewApplicationCategoryResource() resource.Resource {
	return &ApplicationCategoryResource{
		ReadPermissions: []string{
			"DeviceManagementConfiguration.Read.All",
			"DeviceManagementApps.Read.All",
		},
		WritePermissions: []string{
			"DeviceManagementConfiguration.ReadWrite.All",
			"DeviceManagementApps.ReadWrite.All",
		},
		ResourcePath: "/deviceAppManagement/mobileAppCategories",
	}
}

type ApplicationCategoryResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ReadPermissions  []string
	WritePermissions []string
	ResourcePath     string
}

// Metadata returns the resource type name.
func (r *ApplicationCategoryResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = ResourceName
}

// Configure sets the client for the resource.
func (r *ApplicationCategoryResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = client.SetGraphBetaClientForResource(ctx, req, resp, ResourceName)
}

// ImportState imports the resource state.
func (r *ApplicationCategoryResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Schema returns the schema for the resource.
func (r *ApplicationCategoryResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages application categories using the `/deviceAppManagement/mobileAppCategories` endpoint. This resource is used to application categories help organize mobile apps in the Company Portal, enabling users to browse and discover apps by category (Business, Productivity, etc.) and allowing administrators to logically group applications for easier management.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					planmodifiers.UseStateForUnknownString(),
				},
				MarkdownDescription: "The unique identifier for this Intune application category",
			},
			"display_name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the Intune application category",
			},
			"last_modified_date_time": schema.StringAttribute{
				Computed:    true,
				Description: "The date and time when the Intune application category was last modified. This property is read-only.",
			},
			"timeouts": commonschema.ResourceTimeouts(ctx),
		},
	}
}
