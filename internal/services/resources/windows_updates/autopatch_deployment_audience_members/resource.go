package graphBetaWindowsUpdatesAutopatchDeploymentAudienceMembers

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	planmodifiers "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/plan_modifiers"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/identityschema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

const (
	ResourceName  = "microsoft365_graph_beta_windows_updates_autopatch_deployment_audience_members"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180
)

var (
	// Basic resource interface (CRUD operations)
	_ resource.Resource = &WindowsUpdatesAutopatchDeploymentAudienceMembersResource{}
	// Allows the resource to be configured with the provider client
	_ resource.ResourceWithConfigure = &WindowsUpdatesAutopatchDeploymentAudienceMembersResource{}
	// Enables import functionality
	_ resource.ResourceWithImportState = &WindowsUpdatesAutopatchDeploymentAudienceMembersResource{}
	// Enables identity schema for list resource support
	_ resource.ResourceWithIdentity = &WindowsUpdatesAutopatchDeploymentAudienceMembersResource{}
)

func NewWindowsUpdatesAutopatchDeploymentAudienceMembersResource() resource.Resource {
	return &WindowsUpdatesAutopatchDeploymentAudienceMembersResource{
		ReadPermissions: []string{
			"WindowsUpdates.ReadWrite.All",
		},
		WritePermissions: []string{
			"WindowsUpdates.ReadWrite.All",
		},
		ResourcePath: "/admin/windows/updates/deploymentAudiences",
	}
}

type WindowsUpdatesAutopatchDeploymentAudienceMembersResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ReadPermissions  []string
	WritePermissions []string
	ResourcePath     string
}

func (r *WindowsUpdatesAutopatchDeploymentAudienceMembersResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = ResourceName
}

func (r *WindowsUpdatesAutopatchDeploymentAudienceMembersResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = client.SetGraphBetaClientForResource(ctx, req, resp, ResourceName)
}

func (r *WindowsUpdatesAutopatchDeploymentAudienceMembersResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Expected format: {audience_id}_{member_type}
	// Example: 12345678-1234-1234-1234-123456789012_azureADDevice
	id := req.ID
	separatorIndex := len(id) - len("_azureADDevice")
	if separatorIndex > 0 && id[separatorIndex:] == "_azureADDevice" {
		audienceID := id[:separatorIndex]
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("audience_id"), audienceID)...)
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("member_type"), "azureADDevice")...)
		return
	}

	separatorIndex = len(id) - len("_updatableAssetGroup")
	if separatorIndex > 0 && id[separatorIndex:] == "_updatableAssetGroup" {
		audienceID := id[:separatorIndex]
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("audience_id"), audienceID)...)
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("member_type"), "updatableAssetGroup")...)
		return
	}

	resp.Diagnostics.AddError(
		"Invalid Import ID",
		fmt.Sprintf("Expected import ID format: {audience_id}_{member_type}, got: %s", req.ID),
	)
}

func (r *WindowsUpdatesAutopatchDeploymentAudienceMembersResource) IdentitySchema(ctx context.Context, req resource.IdentitySchemaRequest, resp *resource.IdentitySchemaResponse) {
	resp.IdentitySchema = identityschema.Schema{
		Attributes: map[string]identityschema.Attribute{
			"id": identityschema.StringAttribute{
				RequiredForImport: true,
			},
		},
	}
}

func (r *WindowsUpdatesAutopatchDeploymentAudienceMembersResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages members and exclusions for a Windows Update deployment audience in Microsoft 365. Using the `/admin/windows/updates/deploymentAudiences/{audienceId}/updateAudienceById` endpoint. " +
			"This resource uses the `updateAudienceById` action to add or remove devices and groups from an audience. " +
			"All members must be of the same type (either all `azureADDevice` or all `updatableAssetGroup`). " +
			"The audience container must be created first using the `microsoft365_graph_beta_windows_updates_autopatch_deployment_audience` resource. ",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The composite identifier for this resource. Format: `{audience_id}_{member_type}`.",
				PlanModifiers: []planmodifier.String{
					planmodifiers.UseStateForUnknownString(),
				},
			},
			"audience_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The unique identifier of the deployment audience to manage members for.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"member_type": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The type of members in this audience. All members and exclusions must be of the same type. Valid values are: `azureADDevice`, `updatableAssetGroup`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf("azureADDevice", "updatableAssetGroup"),
				},
			},
			"members": schema.SetAttribute{
				ElementType:         types.StringType,
				Optional:            true,
				MarkdownDescription: "Set of device or Entra group IDs to include in the deployment audience.",
			},
			"exclusions": schema.SetAttribute{
				ElementType:         types.StringType,
				Optional:            true,
				MarkdownDescription: "Set of device or Entra group IDs to exclude from the deployment audience.",
			},
			"timeouts": commonschema.ResourceTimeouts(ctx),
		},
	}
}
