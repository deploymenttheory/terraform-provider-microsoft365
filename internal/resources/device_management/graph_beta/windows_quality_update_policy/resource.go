package graphBetaWindowsQualityUpdatePolicy

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common"
	planmodifiers "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/plan_modifiers"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/schema"
	commonschemagraphbeta "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/schema/graph_beta/device_management"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

const (
	ResourceName  = "graph_beta_device_management_windows_quality_update_policy"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180
)

var (
	// Basic resource interface (CRUD operations)
	_ resource.Resource = &WindowsQualityUpdatePolicyResource{}

	// Allows the resource to be configured with the provider client
	_ resource.ResourceWithConfigure = &WindowsQualityUpdatePolicyResource{}

	// Enables import functionality
	_ resource.ResourceWithImportState = &WindowsQualityUpdatePolicyResource{}

	// Enables plan modification/diff suppression
	_ resource.ResourceWithModifyPlan = &WindowsQualityUpdatePolicyResource{}
)

func NewWindowsQualityUpdatePolicyResource() resource.Resource {
	return &WindowsQualityUpdatePolicyResource{
		ReadPermissions: []string{
			"DeviceManagementConfiguration.Read.All",
		},
		WritePermissions: []string{
			"DeviceManagementConfiguration.ReadWrite.All",
		},
		ResourcePath: "/deviceManagement/WindowsQualityUpdatePolicies",
	}
}

type WindowsQualityUpdatePolicyResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ProviderTypeName string
	TypeName         string
	ReadPermissions  []string
	WritePermissions []string
	ResourcePath     string
}

// Metadata returns the resource type name.
func (r *WindowsQualityUpdatePolicyResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + ResourceName
}

// Configure sets the client for the resource.
func (r *WindowsQualityUpdatePolicyResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = common.SetGraphBetaClientForResource(ctx, req, resp, r.TypeName)
}

// ImportState imports the resource state.
func (r *WindowsQualityUpdatePolicyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Schema defines the schema for the resource.
func (r *WindowsQualityUpdatePolicyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a Windows Quality Update Policy in Microsoft Intune. This correlates to the gui location: Devices -> Manage Updates -> Windows Updates -> Quality Updates.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					planmodifiers.UseStateForUnknownString(),
				},
				MarkdownDescription: "The Intune quality update policy id.",
			},
			"display_name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The display name for the policy.",
				Validators: []validator.String{
					stringvalidator.LengthAtMost(200),
				},
			},
			"description": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The description of the policy which is specified by the user. Max allowed length is 1500 chars.",
			},
			"hotpatch_enabled": schema.BoolAttribute{
				Optional:            true,
				MarkdownDescription: "Indicates if hotpatch is enabled for the tenants. When 'true', tenant can apply quality updates without rebooting their devices. When 'false', tenant devices will receive cold patch associated with Windows quality updates.",
			},
			"created_date_time": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Timestamp of when the profile was created. The value cannot be modified and is automatically populated when the profile is created.",
			},
			"last_modified_date_time": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Timestamp of when the profile was modified. The value cannot be modified and is automatically populated when the profile is modified.",
			},
			"role_scope_tag_ids": schema.SetAttribute{
				ElementType:         types.StringType,
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Set of scope tag IDs for this Settings Catalog template profile.",
				PlanModifiers: []planmodifier.Set{
					planmodifiers.DefaultSetValue(
						[]attr.Value{types.StringValue("0")},
					),
				},
			},
			"timeouts": commonschema.Timeouts(ctx),
		},
		Blocks: map[string]schema.Block{
			"assignment": commonschemagraphbeta.WindowsUpdateAssignments(),
		},
	}
}
