package graphBetaWindowsFeatureUpdateProfile

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
	ResourceName  = "graph_beta_device_management_windows_feature_update_profile"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180
)

var (
	// Basic resource interface (CRUD operations)
	_ resource.Resource = &WindowsFeatureUpdateProfileResource{}

	// Allows the resource to be configured with the provider client
	_ resource.ResourceWithConfigure = &WindowsFeatureUpdateProfileResource{}

	// Enables import functionality
	_ resource.ResourceWithImportState = &WindowsFeatureUpdateProfileResource{}

	// Enables plan modification/diff suppression
	_ resource.ResourceWithModifyPlan = &WindowsFeatureUpdateProfileResource{}
)

func NewWindowsFeatureUpdateProfileResource() resource.Resource {
	return &WindowsFeatureUpdateProfileResource{
		ReadPermissions: []string{
			"DeviceManagementConfiguration.Read.All",
		},
		WritePermissions: []string{
			"DeviceManagementConfiguration.ReadWrite.All",
		},
		ResourcePath: "/deviceManagement/windowsFeatureUpdateProfiles",
	}
}

type WindowsFeatureUpdateProfileResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ProviderTypeName string
	TypeName         string
	ReadPermissions  []string
	WritePermissions []string
	ResourcePath     string
}

// Metadata returns the resource type name.
func (r *WindowsFeatureUpdateProfileResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + ResourceName
}

// Configure sets the client for the resource.
func (r *WindowsFeatureUpdateProfileResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = common.SetGraphBetaClientForResource(ctx, req, resp, r.TypeName)
}

// ImportState imports the resource state.
func (r *WindowsFeatureUpdateProfileResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Schema defines the schema for the resource.
func (r *WindowsFeatureUpdateProfileResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a Windows Feature Update Profile in Microsoft Intune. This correlates to the gui location: Devices -> Manage Updates -> Windows Updates -> Feature Updates.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					planmodifiers.UseStateForUnknownString(),
				},
				MarkdownDescription: "The Identifier of the entity.",
			},
			"display_name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The display name of the profile.",
			},
			"description": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The description of the profile which is specified by the user.",
			},
			"feature_update_version": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The feature update version that will be deployed to the devices targeted by this profile. Valid values are: \"Windows 11, version 24H2\", \"Windows 11, version 23H2\", \"Windows 11, version 22H2\", \"Windows 10, version 22H2\". By selecting this Feature update to deploy you are agreeing that when applying this operating system to a device either (1) the applicable Windows license was purchased though volume licensing, or (2) that you are authorized to bind your organization and are accepting on its behalf the relevant Microsoft Software License Terms to be found here https://go.microsoft.com/fwlink/?linkid=2171206.",
				Validators: []validator.String{
					stringvalidator.OneOf(
						"Windows 11, version 24H2",
						"Windows 11, version 23H2",
						"Windows 11, version 22H2",
						"Windows 10, version 22H2"),
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
				Computed:            true,
				MarkdownDescription: "Set of scope tag IDs for this Settings Catalog template profile.",
				PlanModifiers: []planmodifier.Set{
					planmodifiers.DefaultSetValue(
						[]attr.Value{types.StringValue("0")},
					),
				},
			},
			"deployable_content_display_name": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Friendly display name of the quality update profile deployable content",
			},
			"end_of_support_date": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The last supported date for a feature update",
			},
			"install_latest_windows10_on_windows11_ineligible_device": schema.BoolAttribute{
				Optional:            true,
				MarkdownDescription: "Specifies whether Windows 10 devices that are not eligible for Windows 11 are offered the latest Windows 10 feature updates. Changes to this field require the resource to be replaced.",
				PlanModifiers: []planmodifier.Bool{
					planmodifiers.NewRequiresReplaceIfChangedBool(),
				},
			},
			"install_feature_updates_optional": schema.BoolAttribute{
				Optional:            true,
				MarkdownDescription: "If true, the Windows 11 update will become optional",
			},
			"rollout_settings": schema.SingleNestedAttribute{
				Optional:            true,
				MarkdownDescription: "The windows update rollout settings, including offer start date time, offer end date time, and days between each set of offers.",
				Attributes: map[string]schema.Attribute{
					"offer_start_date_time_in_utc": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "The UTC offer start date time of the rollout.",
					},
					"offer_end_date_time_in_utc": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "The UTC offer end date time of the rollout.",
					},
					"offer_interval_in_days": schema.Int32Attribute{
						Optional:            true,
						MarkdownDescription: "The number of days between each set of offers.",
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
