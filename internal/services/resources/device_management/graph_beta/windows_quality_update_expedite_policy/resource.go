package graphBetaWindowsQualityUpdateExpeditePolicy

import (
	"context"
	"regexp"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	planmodifiers "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/plan_modifiers"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
	commonschemagraphbeta "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema/graph_beta/device_management"
	attributevalidate "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/validate/attribute"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
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
	ResourceName  = "microsoft365_graph_beta_device_management_windows_quality_update_expedite_policy"
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
	ReadPermissions  []string
	WritePermissions []string
	ResourcePath     string
}

// Metadata returns the resource type name.
func (r *WindowsQualityUpdateExpeditePolicyResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = ResourceName
}

func (r *WindowsQualityUpdateExpeditePolicyResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = client.SetGraphBetaClientForResource(ctx, req, resp, ResourceName)
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
				Computed:            true,
				MarkdownDescription: "The Intune Windows Quality Update Profile (expedite policy) profile id.",
				PlanModifiers: []planmodifier.String{
					planmodifiers.UseStateForUnknownString(),
				},
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
				Computed:            true,
				MarkdownDescription: "The description of the profile which is specified by the user. Max allowed length is 1500 chars.",
				Validators: []validator.String{
					stringvalidator.LengthAtMost(1500),
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
			"release_date_display_name": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Friendly release date to display for a Quality Update release",
			},
			"deployable_content_display_name": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Friendly display name of the quality update profile deployable content",
			},
			"expedited_update_settings": schema.SingleNestedAttribute{
				Required:            true,
				MarkdownDescription: "Expedited Quality update settings.",
				Attributes: map[string]schema.Attribute{
					"quality_update_release": schema.StringAttribute{
						Required: true,
						MarkdownDescription: "Expedite installation of quality updates if device OS version less than the quality update release identifier. " +
							"Value must be a valid ISO 8601 datetime format (e.g., 2025-12-09T00:00:00Z). Valid values as of December 2025: 2025-12-09T00:00:00Z, 2025-11-20T00:00:00Z",
						Validators: []validator.String{
							stringvalidator.OneOf("2025-12-09T00:00:00Z", "2025-11-20T00:00:00Z"),
							attributevalidate.RegexMatches(
								regexp.MustCompile(constants.ISO8601DateTimeRegex),
								"value must be a valid ISO 8601 datetime format (e.g., 2025-12-09T00:00:00Z)",
							),
						},
					},
					"days_until_forced_reboot": schema.Int32Attribute{
						Required:            true,
						MarkdownDescription: "If a reboot is required, select the number of days before it's enforced. Valid values are: 0, 1, and 2.",
						Validators: []validator.Int32{
							int32validator.OneOf(0, 1, 2),
						},
					},
				},
			},
			"assignments": commonschemagraphbeta.WindowsSoftwareUpdateAssignmentsSchema(),
			"timeouts":    commonschema.Timeouts(ctx),
		},
	}
}
