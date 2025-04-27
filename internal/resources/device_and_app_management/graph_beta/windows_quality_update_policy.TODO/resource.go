package graphBetaWindowsQualityUpdateProfile

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common"
	planmodifiers "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/plan_modifiers"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

const (
	ResourceName  = "graph_beta_device_and_app_management_windows_feature_update_profile"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180
)

var (
	// Basic resource interface (CRUD operations)
	_ resource.Resource = &WindowsQualityUpdateProfileResource{}

	// Allows the resource to be configured with the provider client
	_ resource.ResourceWithConfigure = &WindowsQualityUpdateProfileResource{}

	// Enables import functionality
	_ resource.ResourceWithImportState = &WindowsQualityUpdateProfileResource{}

	// Enables plan modification/diff suppression
	_ resource.ResourceWithModifyPlan = &WindowsQualityUpdateProfileResource{}
)

func NewWindowsQualityUpdateProfileResource() resource.Resource {
	return &WindowsQualityUpdateProfileResource{
		ReadPermissions: []string{
			"DeviceManagementConfiguration.Read.All",
		},
		WritePermissions: []string{
			"DeviceManagementConfiguration.ReadWrite.All",
		},
		ResourcePath: "/deviceManagement/WindowsQualityUpdateProfiles",
	}
}

type WindowsQualityUpdateProfileResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ProviderTypeName string
	TypeName         string
	ReadPermissions  []string
	WritePermissions []string
	ResourcePath     string
}

// Metadata returns the resource type name.
func (r *WindowsQualityUpdateProfileResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + ResourceName
}

// Configure sets the client for the resource.
func (r *WindowsQualityUpdateProfileResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = common.SetGraphBetaClientForResource(ctx, req, resp, r.TypeName)
}

// ImportState imports the resource state.
func (r *WindowsQualityUpdateProfileResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Schema defines the schema for the resource.
func (r *WindowsQualityUpdateProfileResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a Windows Quality Update Profile in Microsoft Intune. This correlates to the gui location: Devices -> Manage Updates -> Windows Updates -> Quality Updates.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					planmodifiers.UseStateForUnknownString(),
				},
				MarkdownDescription: "The Intune quality update profile id.",
			},
			"display_name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The display name for the profile.",
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
				MarkdownDescription: "List of Scope Tags for this Quality Update entity.",
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
				MarkdownDescription: "Expedited update settings.",
				Attributes: map[string]schema.Attribute{
					"quality_update_release": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "The quality update release identifier.",
					},
					"days_until_forced_reboot": schema.Int32Attribute{
						Optional:            true,
						MarkdownDescription: "Number of days before a forced reboot occurs.",
					},
				},
			},
			"timeouts": commonschema.Timeouts(ctx),
		},
	}
}
