package graphBetaWindowsUpdatesAutopatchUpdatableAssetGroup

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
	ResourceName  = "microsoft365_graph_beta_windows_updates_autopatch_updatable_asset_group"
	CreateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180
)

var (
	_ resource.Resource                = &WindowsUpdatesAutopatchUpdatableAssetGroupResource{}
	_ resource.ResourceWithConfigure   = &WindowsUpdatesAutopatchUpdatableAssetGroupResource{}
	_ resource.ResourceWithImportState = &WindowsUpdatesAutopatchUpdatableAssetGroupResource{}
)

func NewWindowsUpdatesAutopatchUpdatableAssetGroupResource() resource.Resource {
	return &WindowsUpdatesAutopatchUpdatableAssetGroupResource{
		ReadPermissions: []string{
			"WindowsUpdates.ReadWrite.All",
		},
		WritePermissions: []string{
			"WindowsUpdates.ReadWrite.All",
		},
		ResourcePath: "/admin/windows/updates/updatableAssets",
	}
}

type WindowsUpdatesAutopatchUpdatableAssetGroupResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ReadPermissions  []string
	WritePermissions []string
	ResourcePath     string
}

func (r *WindowsUpdatesAutopatchUpdatableAssetGroupResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = ResourceName
}

func (r *WindowsUpdatesAutopatchUpdatableAssetGroupResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = client.SetGraphBetaClientForResource(ctx, req, resp, ResourceName)
}

func (r *WindowsUpdatesAutopatchUpdatableAssetGroupResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *WindowsUpdatesAutopatchUpdatableAssetGroupResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a Windows Autopatch updatable asset group using the `/admin/windows/updates/updatableAssets` endpoint. " +
			"An updatable asset group is a logical container for grouping Entra ID devices for Windows Autopatch targeting. " +
			"The group itself has no configurable properties — membership is managed separately via the " +
			"`microsoft365_graph_beta_windows_updates_autopatch_updatable_asset_group_assignment` resource. " +
			"Any change requires resource replacement.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					planmodifiers.UseStateForUnknownString(),
				},
				MarkdownDescription: "The unique identifier for the updatable asset group.",
			},
			"timeouts": commonschema.ResourceTimeouts(ctx),
		},
	}
}
