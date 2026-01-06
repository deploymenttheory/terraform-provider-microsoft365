package graphBetaWindowsFeatureUpdatePolicy

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	planmodifiers "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/plan_modifiers"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
	commonschemagraphbeta "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema/graph_beta/device_management"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

const (
	ResourceName  = "microsoft365_graph_beta_device_management_windows_feature_update_policy"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180
)

var (
	// Basic resource interface (CRUD operations)
	_ resource.Resource = &WindowsFeatureUpdatePolicyResource{}

	// Allows the resource to be configured with the provider client
	_ resource.ResourceWithConfigure = &WindowsFeatureUpdatePolicyResource{}

	// Enables import functionality
	_ resource.ResourceWithImportState = &WindowsFeatureUpdatePolicyResource{}

	// Enables plan modification/diff suppression
	_ resource.ResourceWithModifyPlan = &WindowsFeatureUpdatePolicyResource{}
)

func NewWindowsFeatureUpdatePolicyResource() resource.Resource {
	return &WindowsFeatureUpdatePolicyResource{
		ReadPermissions: []string{
			"DeviceManagementConfiguration.Read.All",
		},
		WritePermissions: []string{
			"DeviceManagementConfiguration.ReadWrite.All",
		},
		ResourcePath: "/deviceManagement/windowsFeatureUpdateProfiles",
	}
}

type WindowsFeatureUpdatePolicyResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ReadPermissions  []string
	WritePermissions []string
	ResourcePath     string
}

// Metadata returns the resource type name.
func (r *WindowsFeatureUpdatePolicyResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = ResourceName
}

func (r *WindowsFeatureUpdatePolicyResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = client.SetGraphBetaClientForResource(ctx, req, resp, ResourceName)
}

// ImportState imports the resource state.
func (r *WindowsFeatureUpdatePolicyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Schema defines the schema for the resource.
func (r *WindowsFeatureUpdatePolicyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages Windows feature update profiles using the `/deviceManagement/windowsFeatureUpdateProfiles` endpoint. Feature update profiles control major Windows version deployments (like Windows 11 24H2) with rollout scheduling, device eligibility rules, and deployment timing to ensure controlled OS upgrades across managed devices.",
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
				Computed:            true,
				Default:             stringdefault.StaticString(""),
				MarkdownDescription: "Optional description of the resource. Maximum length is 1500 characters.",
				Validators: []validator.String{
					stringvalidator.LengthAtMost(1500),
				},
			},
			"feature_update_version": schema.StringAttribute{
				Required: true,
				MarkdownDescription: "The feature update version that will be deployed to the devices targeted by this profile. " +
					"Valid values are: \"Windows 11, version 25H2\", \"Windows 11, version 24H2\", \"Windows 11, version 23H2\". " +
					"By selecting this Feature update to deploy you are agreeing that when applying this operating system to a " +
					"device either (1) the applicable Windows license was purchased though volume licensing, or (2) that you are " +
					"authorized to bind your organization and are accepting on its behalf the relevant Microsoft Software License " +
					"Terms to be found here https://go.microsoft.com/fwlink/?linkid=2171206.",
				Validators: []validator.String{
					stringvalidator.OneOf(
						"Windows 11, version 25H2",
						"Windows 11, version 24H2",
						"Windows 11, version 23H2",
					),
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
				Optional: true,
				MarkdownDescription: "Specifies whether Windows 10 devices that are not eligible for Windows 11 are offered the latest Windows 10" +
					" feature updates. Changes to this field require the resource to be replaced.",
				PlanModifiers: []planmodifier.Bool{
					planmodifiers.NewRequiresReplaceIfChangedBool(),
				},
			},
			"install_feature_updates_optional": schema.BoolAttribute{
				Optional: true,
				Computed: true,
				MarkdownDescription: "If true, the Windows 11 update will become available to users as an optional update. " +
					"If false, the Windows 11 update will become available to users as a required update",
			},
			"rollout_settings": schema.SingleNestedAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "The windows update rollout settings, including offer start date time, offer end date time, and days between each set of offers.",
				Attributes: map[string]schema.Attribute{
					"offer_start_date_time_in_utc": schema.StringAttribute{
						Optional:            true,
						Computed:            true,
						MarkdownDescription: "The UTC offer start date time of the rollout. Must be in RFC3339 format (e.g., '2025-05-01T00:00:00Z').",
					},
					"offer_end_date_time_in_utc": schema.StringAttribute{
						Optional:            true,
						Computed:            true,
						MarkdownDescription: "The UTC offer end date time of the rollout.",
					},
					"offer_interval_in_days": schema.Int32Attribute{
						Optional:            true,
						Computed:            true,
						MarkdownDescription: "The number of days between each set of offers. The value must be between 1 and 14.",
						Validators: []validator.Int32{
							int32validator.Between(1, 14),
						},
					},
				},
			},
			"assignments": commonschemagraphbeta.WindowsSoftwareUpdateAssignmentsSchema(),
			"timeouts":    commonschema.ResourceTimeouts(ctx),
		},
	}
}
