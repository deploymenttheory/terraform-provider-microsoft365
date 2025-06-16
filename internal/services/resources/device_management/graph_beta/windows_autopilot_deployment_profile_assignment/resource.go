package graphBetaWindowsAutopilotDeploymentProfileAssignment

import (
	"context"
	"regexp"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

const (
	ResourceName  = "graph_beta_device_management_windows_autopilot_deployment_profile_assignment"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180
)

var (
	// Basic resource interface (CRUD operations)
	_ resource.Resource = &WindowsAutopilotDeploymentProfileAssignmentResource{}

	// Allows the resource to be configured with the provider client
	_ resource.ResourceWithConfigure = &WindowsAutopilotDeploymentProfileAssignmentResource{}

	// Enables import functionality
	_ resource.ResourceWithImportState = &WindowsAutopilotDeploymentProfileAssignmentResource{}

	// Enables plan modification/diff suppression
	_ resource.ResourceWithModifyPlan = &WindowsAutopilotDeploymentProfileAssignmentResource{}
)

func NewWindowsAutopilotDeploymentProfileAssignmentResource() resource.Resource {
	return &WindowsAutopilotDeploymentProfileAssignmentResource{
		ReadPermissions: []string{
			"DeviceManagementServiceConfig.Read.All",
			"DeviceManagementConfiguration.Read.All",
		},
		WritePermissions: []string{
			"DeviceManagementServiceConfig.ReadWrite.All",
			"DeviceManagementConfiguration.ReadWrite.All",
		},
		ResourcePath: "/deviceManagement/windowsAutopilotDeploymentProfiles",
	}
}

type WindowsAutopilotDeploymentProfileAssignmentResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ProviderTypeName string
	TypeName         string
	ReadPermissions  []string
	WritePermissions []string
	ResourcePath     string
}

// Metadata returns the resource type name.
func (r *WindowsAutopilotDeploymentProfileAssignmentResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	r.ProviderTypeName = req.ProviderTypeName
	r.TypeName = ResourceName
	resp.TypeName = r.FullTypeName()
}

// FullTypeName returns the full type name of the resource for logging purposes.
func (r *WindowsAutopilotDeploymentProfileAssignmentResource) FullTypeName() string {
	return r.ProviderTypeName + "_" + r.TypeName
}

// Configure sets the client for the resource.
func (r *WindowsAutopilotDeploymentProfileAssignmentResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = common.SetGraphBetaClientForResource(ctx, req, resp, constants.PROVIDER_NAME+"_"+ResourceName)
}

// ImportState imports the resource state.
func (r *WindowsAutopilotDeploymentProfileAssignmentResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Schema returns the schema for the resource.
func (r *WindowsAutopilotDeploymentProfileAssignmentResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages Windows Autopilot Deployment Profile Assignments in Microsoft Intune.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The ID of the Windows Autopilot Deployment Profile assignment.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"windows_autopilot_deployment_profile_id": schema.StringAttribute{
				Required:    true,
				Description: "The ID of the Windows Autopilot Deployment Profile to attach the assignment to.",
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(constants.GuidRegex),
						"must be a valid GUID in the format 00000000-0000-0000-0000-000000000000",
					),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"source": schema.StringAttribute{
				MarkdownDescription: "Type of resource used for deployment to a group. Possible values are: direct, policySets.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf(
						"direct",
						"policySets",
					),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"source_id": schema.StringAttribute{
				MarkdownDescription: "Identifier for resource used for deployment to a group.",
				Optional:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"target": schema.SingleNestedAttribute{
				Required: true,
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.RequiresReplace(),
				},
				Attributes: map[string]schema.Attribute{
					"target_type": schema.StringAttribute{
						Required: true,
						MarkdownDescription: "The target group type for the Windows Autopilot Deployment Profile assignment. Possible values are:\n\n" +
							"- **allDevices**: Target all devices in the tenant\n" +
							"- **allLicensedUsers**: Target all licensed users in the tenant\n" +
							"- **exclusionGroupAssignment**: Target a specific Entra ID group for exclusion\n" +
							"- **groupAssignment**: Target a specific Entra ID group",
						Validators: []validator.String{
							stringvalidator.OneOf(
								"allDevices",
								"allLicensedUsers",
								"exclusionGroupAssignment",
								"groupAssignment",
							),
						},
					},
					"group_id": schema.StringAttribute{
						MarkdownDescription: "The Entra ID group ID for the Windows Autopilot Deployment Profile assignment target. Required when target_type is 'groupAssignment' or 'exclusionGroupAssignment'.",
						Optional:            true,
						Validators: []validator.String{
							stringvalidator.RegexMatches(
								regexp.MustCompile(constants.GuidRegex),
								"must be a valid GUID in the format 00000000-0000-0000-0000-000000000000",
							),
						},
					},
					"device_and_app_management_assignment_filter_id": schema.StringAttribute{
						MarkdownDescription: "The Id of the scope filter applied to the target assignment.",
						Optional:            true,
						Validators: []validator.String{
							stringvalidator.RegexMatches(
								regexp.MustCompile(constants.GuidRegex),
								"must be a valid GUID in the format 00000000-0000-0000-0000-000000000000",
							),
						},
					},
					"device_and_app_management_assignment_filter_type": schema.StringAttribute{
						Optional: true,
						Computed: true,
						Default:  stringdefault.StaticString("none"),
						MarkdownDescription: "The type of scope filter for the target assignment. Defaults to 'none'. Possible values are:\n\n" +
							"- **include**: Only include devices or users matching the filter\n" +
							"- **exclude**: Exclude devices or users matching the filter\n" +
							"- **none**: No assignment filter applied",
						Validators: []validator.String{
							stringvalidator.OneOf(
								"include",
								"exclude",
								"none",
							),
						},
					},
				},
			},
			"timeouts": commonschema.Timeouts(ctx),
		},
	}
}
