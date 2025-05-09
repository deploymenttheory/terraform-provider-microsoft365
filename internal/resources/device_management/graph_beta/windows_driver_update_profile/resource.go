package graphBetaWindowsDriverUpdateProfile

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common"
	planmodifiers "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/plan_modifiers"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/schema"
	commonschemagraphbeta "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/schema/graph_beta/device_management"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

const (
	ResourceName  = "graph_beta_device_management_windows_driver_update_profile"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180
)

var (
	// Basic resource interface (CRUD operations)
	_ resource.Resource = &WindowsDriverUpdateProfileResource{}

	// Allows the resource to be configured with the provider client
	_ resource.ResourceWithConfigure = &WindowsDriverUpdateProfileResource{}

	// Enables import functionality
	_ resource.ResourceWithImportState = &WindowsDriverUpdateProfileResource{}

	// Enables plan modification/diff suppression
	_ resource.ResourceWithModifyPlan = &WindowsDriverUpdateProfileResource{}
)

func NewWindowsDriverUpdateProfileResource() resource.Resource {
	return &WindowsDriverUpdateProfileResource{
		ReadPermissions: []string{
			"DeviceManagementConfiguration.Read.All",
		},
		WritePermissions: []string{
			"DeviceManagementConfiguration.ReadWrite.All",
		},
		ResourcePath: "/deviceManagement/windowsDriverUpdateProfiles",
	}
}

type WindowsDriverUpdateProfileResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ProviderTypeName string
	TypeName         string
	ReadPermissions  []string
	WritePermissions []string
	ResourcePath     string
}

// Metadata returns the resource type name.
func (r *WindowsDriverUpdateProfileResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + ResourceName
}

// Configure sets the client for the resource.
func (r *WindowsDriverUpdateProfileResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = common.SetGraphBetaClientForResource(ctx, req, resp, r.TypeName)
}

// ImportState imports the resource state.
func (r *WindowsDriverUpdateProfileResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Schema defines the schema for the resource.
func (r *WindowsDriverUpdateProfileResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a Windows Driver Update Profile in Microsoft Intune. This correlates to the gui location: Devices -> Manage Updates -> Windows Updates -> Driver Updates.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					planmodifiers.UseStateForUnknownString(),
				},
				MarkdownDescription: "The Intune policy id.",
			},
			"display_name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The display name for the profile.",
			},
			"description": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The description of the profile which is specified by the user.",
			},
			"approval_type": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Driver update profile approval type. For example, manual or automatic approval. Possible values are: `manual`, `automatic`.",
				Validators: []validator.String{
					stringvalidator.OneOf("manual", "automatic"),
				},
				PlanModifiers: []planmodifier.String{
					planmodifiers.RequiresReplaceString(),
				},
			},
			"device_reporting": schema.Int32Attribute{
				Computed:            true,
				MarkdownDescription: "Number of devices reporting for this profile",
			},
			"new_updates": schema.Int32Attribute{
				Computed:            true,
				MarkdownDescription: "Number of new driver updates available for this profile.",
			},
			"deployment_deferral_in_days": schema.Int32Attribute{
				Optional:            true,
				MarkdownDescription: "Deployment deferral settings in days, only applicable when ApprovalType is set to automatic approval.",
				PlanModifiers: []planmodifier.Int32{
					planmodifiers.RequiresOtherAttributeValueInt32(path.Root("approval_type"), "automatic"),
				},
			},
			"created_date_time": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The date time that the profile was created.",
			},
			"last_modified_date_time": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The date time that the profile was last modified.",
			},
			"role_scope_tag_ids": schema.SetAttribute{
				ElementType:         types.StringType,
				Optional:            true,
				MarkdownDescription: "List of Scope Tags for this Driver Update entity.",
			},
			"inventory_sync_status": schema.SingleNestedAttribute{
				Computed:            true,
				MarkdownDescription: "Driver inventory sync status for this profile.",
				Attributes: map[string]schema.Attribute{
					"last_successful_sync_date_time": schema.StringAttribute{
						Computed:            true,
						MarkdownDescription: "Last successful sync date time for driver inventory.",
					},
					"driver_inventory_sync_state": schema.StringAttribute{
						Computed:            true,
						MarkdownDescription: "Driver inventory sync state for this profile.",
					},
				},
			},
			"timeouts": commonschema.Timeouts(ctx),
		},
		Blocks: map[string]schema.Block{
			"assignment": commonschemagraphbeta.WindowsUpdateAssignments(),
		},
	}
}
