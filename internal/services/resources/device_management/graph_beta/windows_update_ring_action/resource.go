package graphBetaWindowsUpdateRingAction

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	planmodifiers "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/plan_modifiers"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

const (
	ResourceName  = "graph_beta_device_management_windows_update_ring_action"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180
)

var (
	// Basic resource interface (CRUD operations)
	_ resource.Resource = &WindowsUpdateRingActionResource{}

	// Allows the resource to be configured with the provider client
	_ resource.ResourceWithConfigure = &WindowsUpdateRingActionResource{}

	// Enables import functionality
	_ resource.ResourceWithImportState = &WindowsUpdateRingActionResource{}
)

func NewWindowsUpdateRingActionResource() resource.Resource {
	return &WindowsUpdateRingActionResource{
		ReadPermissions: []string{
			"DeviceManagementConfiguration.Read.All",
		},
		WritePermissions: []string{
			"DeviceManagementConfiguration.ReadWrite.All",
		},
		ResourcePath: "/deviceManagement/deviceConfigurations",
	}
}

type WindowsUpdateRingActionResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ProviderTypeName string
	TypeName         string
	ReadPermissions  []string
	WritePermissions []string
	ResourcePath     string
}

// Metadata returns the resource type name.
func (r *WindowsUpdateRingActionResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	r.ProviderTypeName = req.ProviderTypeName
	r.TypeName = ResourceName
	resp.TypeName = r.FullTypeName()
}

// FullTypeName returns the full resource type name in the format "providername_resourcename".
func (r *WindowsUpdateRingActionResource) FullTypeName() string {
	return r.ProviderTypeName + "_" + r.TypeName
}

// Configure sets the client for the resource.
func (r *WindowsUpdateRingActionResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = client.SetGraphBetaClientForResource(ctx, req, resp, constants.PROVIDER_NAME+"_"+ResourceName)
}

// ImportState imports the resource state.
func (r *WindowsUpdateRingActionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Schema defines the schema for the resource.
func (r *WindowsUpdateRingActionResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages Windows Update for Business configuration policy actions using the `/deviceManagement/deviceConfigurations` endpoint. This resource allows performing actions on Windows Update Rings such as pausing, resuming, extending pause periods, and rolling back updates for managed Windows 10/11 devices.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					planmodifiers.UseStateForUnknownString(),
				},
				MarkdownDescription: "Unique identifier for this action resource.",
			},
			"update_ring_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The ID of the Windows Update Ring to perform actions on.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"pause_feature_updates": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				MarkdownDescription: "Set to `true` to pause feature updates for the specified Windows Update Ring. This will set `featureUpdatesPaused` to `true` via PATCH API call.",
			},
			"resume_feature_updates": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				MarkdownDescription: "Set to `true` to resume feature updates for the specified Windows Update Ring. This will set `featureUpdatesPaused` to `false` via PATCH API call.",
			},
			"extend_feature_updates_pause": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				MarkdownDescription: "Set to `true` to extend the feature updates pause by the maximum allowed period (35 days). This will make a POST API call to the `extendFeatureUpdatesPause` endpoint.",
			},
			"rollback_feature_updates": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				MarkdownDescription: "Set to `true` to rollback feature updates for the specified Windows Update Ring. This will set `featureUpdatesWillBeRolledBack` to `true` via PATCH API call.",
			},
			"pause_quality_updates": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				MarkdownDescription: "Set to `true` to pause quality updates for the specified Windows Update Ring. This will set `qualityUpdatesPaused` to `true` via PATCH API call.",
			},
			"resume_quality_updates": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				MarkdownDescription: "Set to `true` to resume quality updates for the specified Windows Update Ring. This will set `qualityUpdatesPaused` to `false` via PATCH API call.",
			},
			"rollback_quality_updates": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				MarkdownDescription: "Set to `true` to rollback quality updates for the specified Windows Update Ring. This will set `qualityUpdatesWillBeRolledBack` to `true` via PATCH API call.",
			},
			"description": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Description of the actions being performed on the Windows Update Ring.",
			},
			"last_action_performed": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The last action that was successfully performed on the Windows Update Ring.",
			},
			"last_action_timestamp": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Timestamp when the last action was performed.",
			},
			"timeouts": commonschema.Timeouts(ctx),
		},
	}
}
