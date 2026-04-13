package graphBetaWindowsUpdatesAutopatchDeploymentState

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	planmodifiers "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/plan_modifiers"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

const (
	ResourceName  = "microsoft365_graph_beta_windows_updates_autopatch_deployment_state"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180
)

var (
	// Basic resource interface (CRUD operations)
	_ resource.Resource = &WindowsUpdatesAutopatchDeploymentStateResource{}
	// Allows the resource to be configured with the provider client
	_ resource.ResourceWithConfigure = &WindowsUpdatesAutopatchDeploymentStateResource{}
	// Enables import functionality
	_ resource.ResourceWithImportState = &WindowsUpdatesAutopatchDeploymentStateResource{}
)

func NewWindowsUpdatesAutopatchDeploymentStateResource() resource.Resource {
	return &WindowsUpdatesAutopatchDeploymentStateResource{
		ReadPermissions: []string{
			"WindowsUpdates.Read.All",
		},
		WritePermissions: []string{
			"WindowsUpdates.ReadWrite.All",
		},
		ResourcePath: "/admin/windows/updates/deployments",
	}
}

type WindowsUpdatesAutopatchDeploymentStateResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ReadPermissions  []string
	WritePermissions []string
	ResourcePath     string
}

func (r *WindowsUpdatesAutopatchDeploymentStateResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = ResourceName
}

func (r *WindowsUpdatesAutopatchDeploymentStateResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = client.SetGraphBetaClientForResource(ctx, req, resp, ResourceName)
}

func (r *WindowsUpdatesAutopatchDeploymentStateResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("deployment_id"), req.ID)...)
}

func (r *WindowsUpdatesAutopatchDeploymentStateResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages the lifecycle state of a Windows Update autopatch deployment. Using the `/admin/windows/updates/deployments/{deploymentId}/update` endpoint. " +
			"This resource allows pausing, resuming, or archiving a deployment independently from its configuration. " +
			"The deployment must be created first using the `microsoft365_graph_beta_windows_updates_autopatch_deployment` resource. ",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					planmodifiers.UseStateForUnknownString(),
				},
				MarkdownDescription: "The ID of the deployment (same as deployment_id).",
			},
			"deployment_id": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					planmodifiers.RequiresReplaceString(),
				},
				MarkdownDescription: "The ID of the deployment whose state is being managed.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"requested_value": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The requested state of the deployment. Valid values are: `none` (active/offering), `paused`, `archived`.",
				Validators: []validator.String{
					stringvalidator.OneOf("none", "paused", "archived"),
				},
			},
			"effective_value": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The effective state value of the deployment. Possible values: `scheduled`, `offering`, `paused`, `faulted`, `archived` (read-only).",
			},
			"timeouts": commonschema.ResourceTimeouts(ctx),
		},
	}
}
