package graphBetaWindowsAutopatchDeploymentAudience

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/identityschema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

const (
	ResourceName  = "microsoft365_graph_beta_device_management_windows_autopatch_deployment_audience"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180
)

var (
	_ resource.Resource                = &WindowsUpdateDeploymentAudienceResource{}
	_ resource.ResourceWithConfigure   = &WindowsUpdateDeploymentAudienceResource{}
	_ resource.ResourceWithImportState = &WindowsUpdateDeploymentAudienceResource{}
	_ resource.ResourceWithIdentity    = &WindowsUpdateDeploymentAudienceResource{}
)

func NewWindowsUpdateDeploymentAudienceResource() resource.Resource {
	return &WindowsUpdateDeploymentAudienceResource{
		ReadPermissions: []string{
			"WindowsUpdates.ReadWrite.All",
		},
		WritePermissions: []string{
			"WindowsUpdates.ReadWrite.All",
		},
		ResourcePath: "/admin/windows/updates/deploymentAudiences",
	}
}

type WindowsUpdateDeploymentAudienceResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ReadPermissions  []string
	WritePermissions []string
	ResourcePath     string
}

func (r *WindowsUpdateDeploymentAudienceResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = ResourceName
}

func (r *WindowsUpdateDeploymentAudienceResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = client.SetGraphBetaClientForResource(ctx, req, resp, ResourceName)
}

func (r *WindowsUpdateDeploymentAudienceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *WindowsUpdateDeploymentAudienceResource) IdentitySchema(ctx context.Context, req resource.IdentitySchemaRequest, resp *resource.IdentitySchemaResponse) {
	resp.IdentitySchema = identityschema.Schema{
		Attributes: map[string]identityschema.Attribute{
			"id": identityschema.StringAttribute{
				RequiredForImport: true,
			},
		},
	}
}

func (r *WindowsUpdateDeploymentAudienceResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a Windows Update deployment audience container in Microsoft 365. " +
			"A deployment audience is a container that can be populated with devices or groups. " +
			"Use the `microsoft365_graph_beta_device_management_windows_autopatch_deployment_audience_members` resource to manage the actual members and exclusions. " +
			"See the [Microsoft Graph API documentation](https://learn.microsoft.com/en-us/graph/api/resources/windowsupdates-deploymentaudience?view=graph-rest-beta) for more information.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The unique identifier for the deployment audience.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"timeouts": commonschema.ResourceTimeouts(ctx),
		},
	}
}
