package graphVersionResourceTemplate

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

var _ resource.Resource = &ResourceTemplateResource{}
var _ resource.ResourceWithConfigure = &ResourceTemplateResource{}
var _ resource.ResourceWithImportState = &ResourceTemplateResource{}

func NewResourceTemplateResource() resource.Resource {
	return &ResourceTemplateResource{
		ReadPermissions: []string{
			"DeviceManagementConfiguration.Read.All",
		},
		WritePermissions: []string{
			"DeviceManagementConfiguration.ReadWrite.All",
		},
	}
}

type ResourceTemplateResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ProviderTypeName string
	TypeName         string
	ReadPermissions  []string
	WritePermissions []string
}

// Metadata returns the resource type name.
func (r *ResourceTemplateResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_graph_apitype_resource_type_resource_name"
}

// Configure sets the client for the resource.
func (r *ResourceTemplateResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = common.SetGraphBetaClientForResource(ctx, req, resp, r.TypeName)
}

// ImportState imports the resource state.
func (r *ResourceTemplateResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

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
