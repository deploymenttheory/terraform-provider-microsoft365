package graphVersionResourceTemplate

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

const (
	ResourceName  = "microsoft365_graph_api_service_type_resource_name"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180
)

var (
	// Basic resource interface (CRUD operations)
	_ resource.Resource = &ResourceTemplateResource{}

	// Allows the resource to be configured with the provider client
	_ resource.ResourceWithConfigure = &ResourceTemplateResource{}

	// Enables import functionality
	_ resource.ResourceWithImportState = &ResourceTemplateResource{}

	// Enables plan modification/diff suppression
	_ resource.ResourceWithModifyPlan = &ResourceTemplateResource{}
)

func NewResourceTemplateResource() resource.Resource {
	return &ResourceTemplateResource{
		ReadPermissions: []string{
			"ReadPermission.Read.All",
		},
		WritePermissions: []string{
			"WritePermission.ReadWrite.All",
		},
		ResourcePath: "/serviceType/resourceName",
	}
}

type ResourceTemplateResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ReadPermissions  []string
	WritePermissions []string
	ResourcePath     string
}

// Metadata returns the resource type name.
func (r *ResourceTemplateResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = ResourceName
}

// Configure sets the client for the resource.
func (r *ResourceTemplateResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = client.SetGraphBetaClientForResource(ctx, req, resp, ResourceName)
}

// ImportState imports the resource state.
func (r *ResourceTemplateResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Function to create the full device management win32 lob app schema
func (r *ResourceTemplateResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "The resource `resource_name` manages a graph api resource of type `resource_name`",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Unique Identifier for the resource.",
				Computed:    true,
			},
			"etc": schema.StringAttribute{
				Description: "Add schema from here.",
				Computed:    true,
			},
			"timeouts": commonschema.Timeouts(ctx),
		},
	}
}
