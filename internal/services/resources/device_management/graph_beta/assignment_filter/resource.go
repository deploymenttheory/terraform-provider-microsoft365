package graphBetaAssignmentFilter

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	planmodifiers "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/plan_modifiers"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
	validate "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/validate/attribute"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
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
	ResourceName  = "microsoft365_graph_beta_device_management_assignment_filter"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180
)

var (
	// Basic resource interface (CRUD operations)
	_ resource.Resource = &AssignmentFilterResource{}

	// Allows the resource to be configured with the provider client
	_ resource.ResourceWithConfigure = &AssignmentFilterResource{}

	// Enables import functionality
	_ resource.ResourceWithImportState = &AssignmentFilterResource{}

	// Enables plan modification/diff suppression
	_ resource.ResourceWithModifyPlan = &AssignmentFilterResource{}
)

func NewAssignmentFilterResource() resource.Resource {
	return &AssignmentFilterResource{
		ReadPermissions: []string{
			"DeviceManagementConfiguration.Read.All",
		},
		WritePermissions: []string{
			"DeviceManagementConfiguration.ReadWrite.All",
		},
		ResourcePath: "deviceManagement/assignmentFilters",
	}
}

type AssignmentFilterResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ReadPermissions  []string
	WritePermissions []string
	ResourcePath     string
}

// Metadata returns the resource type name.
func (r *AssignmentFilterResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = ResourceName
}

// Configure sets the client for the resource.
func (r *AssignmentFilterResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = client.SetGraphBetaClientForResource(ctx, req, resp, ResourceName)
}

// ImportState imports the resource state.
func (r *AssignmentFilterResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Schema returns the schema for the resource.
func (r *AssignmentFilterResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages assignment filters in Microsoft Intune using the `/deviceManagement/assignmentFilters` endpoint." +
			" Assignment filters enable granular targeting of policies and applications based on device properties like OS version, manufacturer, " +
			"device name, or custom attributes, allowing more precise deployment control beyond basic group membership. You can learn more about assignment filters " +
			"[here](https://learn.microsoft.com/en-us/intune/intune-service/fundamentals/filters-device-properties).",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					planmodifiers.UseStateForUnknownString(),
				},
				MarkdownDescription: "The unique identifier of the assignment filter.",
			},
			"display_name": schema.StringAttribute{
				Required:    true,
				Description: "The display name of the assignment filter.",
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Description: "The optional description of the assignment filter.",
			},
			"platform": schema.StringAttribute{
				Required: true,
				Description: "The Intune device management type (platform) for the assignment filter. " +
					"This specifies the OS platform type for which the assignment filter will be applied." +
					"Must be one of the following values: android, androidForWork, iOS, macOS, windows10AndLater," +
					"androidAOSP, androidMobileApplicationManagement, iOSMobileApplicationManagement, " +
					"windowsMobileApplicationManagement.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf(
						"android",
						"androidForWork",
						"iOS",
						"macOS",
						//"windowsPhone81", causes a 500 error in acc tests
						//"windows81AndLater", causes a 500 error in acc tests
						"windows10AndLater",
						// "androidWorkProfile", causes a 500 error in acc tests
						//"unknown", causes a 500 error in acc tests
						"androidAOSP",
						"androidMobileApplicationManagement",
						"iOSMobileApplicationManagement",
						"windowsMobileApplicationManagement"),
					validate.RequiredOneOfWhen(
						"assignment_filter_management_type",
						"apps",
						"androidMobileApplicationManagement",
						"iOSMobileApplicationManagement",
						"windowsMobileApplicationManagement"),
				},
			},
			"rule": schema.StringAttribute{
				Required:    true,
				Description: "Rule definition of the assignment filter.",
			},
			"assignment_filter_management_type": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Indicates filter is applied to either 'devices' or 'apps' management type. Possible values are: devices, apps, unknownFutureValue. Default filter will be applied to 'devices'.",
				Validators: []validator.String{
					stringvalidator.OneOf(
						"devices",
						"apps",
					),
				},
			},
			"created_date_time": schema.StringAttribute{
				Computed:    true,
				Description: "The creation time of the assignment filter.",
			},
			"last_modified_date_time": schema.StringAttribute{
				Computed:    true,
				Description: "Last modified time of the assignment filter.",
			},
			"role_scope_tags": schema.SetAttribute{
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
			"timeouts": commonschema.ResourceTimeouts(ctx),
		},
	}
}
