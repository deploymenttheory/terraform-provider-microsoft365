package graphBetaWindowsQualityUpdateExpeditePolicy

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common"
	planmodifiers "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/plan_modifiers"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/schema"
	commonschemagraphbeta "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/schema/graph_beta/device_management"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
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
	ResourceName  = "graph_beta_device_and_app_management_windows_quality_update_expedite_policy"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180
)

var (
	// Basic resource interface (CRUD operations)
	_ resource.Resource = &WindowsQualityUpdateExpeditePolicyResource{}

	// Allows the resource to be configured with the provider client
	_ resource.ResourceWithConfigure = &WindowsQualityUpdateExpeditePolicyResource{}

	// Enables import functionality
	_ resource.ResourceWithImportState = &WindowsQualityUpdateExpeditePolicyResource{}

	// Enables plan modification/diff suppression
	_ resource.ResourceWithModifyPlan = &WindowsQualityUpdateExpeditePolicyResource{}
)

func NewWindowsQualityUpdateExpeditePolicyResource() resource.Resource {
	return &WindowsQualityUpdateExpeditePolicyResource{
		ReadPermissions: []string{
			"DeviceManagementConfiguration.Read.All",
		},
		WritePermissions: []string{
			"DeviceManagementConfiguration.ReadWrite.All",
		},
		ResourcePath: "/deviceManagement/WindowsQualityUpdateExpeditePolicys",
	}
}

type WindowsQualityUpdateExpeditePolicyResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ProviderTypeName string
	TypeName         string
	ReadPermissions  []string
	WritePermissions []string
	ResourcePath     string
}

// Metadata returns the resource type name.
func (r *WindowsQualityUpdateExpeditePolicyResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + ResourceName
}

// Configure sets the client for the resource.
func (r *WindowsQualityUpdateExpeditePolicyResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = common.SetGraphBetaClientForResource(ctx, req, resp, r.TypeName)
}

// ImportState imports the resource state.
func (r *WindowsQualityUpdateExpeditePolicyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Schema defines the schema for the resource.
func (r *WindowsQualityUpdateExpeditePolicyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a Windows Quality Update Profile (expedite policy) in Microsoft Intune. This correlates to the gui location: Devices -> Manage Updates -> Windows Updates -> Quality Updates.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					planmodifiers.UseStateForUnknownString(),
				},
				MarkdownDescription: "The Intune Windows Quality Update Profile (expedite policy) profile id.",
			},
			"display_name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The display name for the Windows Quality Update Profile (expedite policy).",
				Validators: []validator.String{
					stringvalidator.LengthAtMost(200),
				},
			},
			"description": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The description of the profile which is specified by the user.",
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
				MarkdownDescription: "List of Intune Scope Tags for this Quality Update entity.",
			},
			"release_date_display_name": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Friendly release date to display for a Quality Update release",
			},
			"deployable_content_display_name": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Friendly display name of the quality update profile deployable content",
			},
			"expedited_update_settings": schema.SingleNestedAttribute{
				Optional:            true,
				MarkdownDescription: "Expedited Quality update settings.",
				Attributes: map[string]schema.Attribute{
					"quality_update_release": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "Expedite installation of quality updates if device OS version less than the quality update release identifier. ",
						Validators: []validator.String{
							stringvalidator.OneOf("2025-04-22T00:00:00Z", "2025-04-08T00:00:00Z", "2025-03-11T00:00:00Z"),
						},
					},
					"days_until_forced_reboot": schema.Int32Attribute{
						Optional:            true,
						MarkdownDescription: "Number of days to wait before restart is enforced. Valid values are: 0, 1, and 2.",
						Validators: []validator.Int32{
							int32validator.OneOf(0, 1, 2),
						},
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
