package graphBetaWindowsUpdatesAutopatchUpdatableAssetGroupAssignment

import (
	"context"
	"regexp"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	planmodifiers "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/plan_modifiers"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

const (
	ResourceName  = "microsoft365_graph_beta_windows_updates_autopatch_updatable_asset_group_assignment"
	CreateTimeout = 180
	ReadTimeout   = 180
	UpdateTimeout = 180
	DeleteTimeout = 180
	ResourcePath  = "/admin/windows/updates/updatableAssets"
)

var (
	// Basic resource interface (CRUD operations)
	_ resource.Resource = &WindowsUpdatesAutopatchUpdatableAssetGroupAssignmentResource{}
	// Allows the resource to be configured with the provider client
	_ resource.ResourceWithConfigure = &WindowsUpdatesAutopatchUpdatableAssetGroupAssignmentResource{}
	// Enables import functionality
	_ resource.ResourceWithImportState = &WindowsUpdatesAutopatchUpdatableAssetGroupAssignmentResource{}
)

func NewWindowsUpdatesAutopatchUpdatableAssetGroupAssignmentResource() resource.Resource {
	return &WindowsUpdatesAutopatchUpdatableAssetGroupAssignmentResource{
		ReadPermissions: []string{
			"WindowsUpdates.ReadWrite.All",
		},
		WritePermissions: []string{
			"WindowsUpdates.ReadWrite.All",
		},
		ResourcePath: ResourcePath,
	}
}

type WindowsUpdatesAutopatchUpdatableAssetGroupAssignmentResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ReadPermissions  []string
	WritePermissions []string
	ResourcePath     string
}

func (r *WindowsUpdatesAutopatchUpdatableAssetGroupAssignmentResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = ResourceName
}

func (r *WindowsUpdatesAutopatchUpdatableAssetGroupAssignmentResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = client.SetGraphBetaClientForResource(ctx, req, resp, ResourceName)
}

func (r *WindowsUpdatesAutopatchUpdatableAssetGroupAssignmentResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("updatable_asset_group_id"), req.ID)...)
}

func (r *WindowsUpdatesAutopatchUpdatableAssetGroupAssignmentResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages membership of a Windows Autopatch updatable asset group using the `/admin/windows/updates/updatableAssets/{groupId}/addMembersById` and `removeMembersById` endpoints. " +
			"This resource manages which Entra ID devices are members of a given updatable asset group. " +
			"When the resource is created, devices are added to the group. " +
			"When updated, the diff between plan and state is applied. " +
			"When deleted, all managed devices are removed from the group. " +
			"The updatable asset group must already exist, managed by the `microsoft365_graph_beta_windows_updates_autopatch_updatable_asset_group` resource.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The composite identifier. Format: `{updatable_asset_group_id}`.",
				PlanModifiers: []planmodifier.String{
					planmodifiers.UseStateForUnknownString(),
				},
			},
			"updatable_asset_group_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The ID of the Windows Autopatch updatable asset group to manage membership for. " +
					"Must reference an existing updatable asset group managed by the `microsoft365_graph_beta_windows_updates_autopatch_updatable_asset_group` resource.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.RegexMatches(regexp.MustCompile(constants.GuidRegex), "Must be a valid UUID format (xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx)"),
				},
			},
			"entra_device_ids": schema.SetAttribute{
				ElementType:         types.StringType,
				Required:            true,
				MarkdownDescription: "Set of Entra ID device registration IDs (`azure_ad_device_id`) to add as members of the updatable asset group.",
				Validators: []validator.Set{
					setvalidator.SizeAtLeast(1),
					setvalidator.ValueStringsAre(
						stringvalidator.RegexMatches(regexp.MustCompile(constants.GuidRegex), "Must be a valid UUID format (xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx)"),
					),
				},
			},
			"timeouts": commonschema.ResourceTimeouts(ctx),
		},
	}
}
