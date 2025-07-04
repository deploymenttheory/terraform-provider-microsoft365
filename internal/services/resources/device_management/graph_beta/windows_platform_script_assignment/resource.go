package graphBetaWindowsPlatformScriptAssignment

import (
	"context"
	"regexp"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	planmodifiers "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/plan_modifiers"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

const (
	ResourceName  = "graph_beta_device_management_windows_platform_script_assignment"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180
)

var (
	// Basic resource interface (CRUD operations)
	_ resource.Resource = &WindowsPlatformScriptAssignmentResource{}

	// Allows the resource to be configured with the provider client
	_ resource.ResourceWithConfigure = &WindowsPlatformScriptAssignmentResource{}

	// Enables import functionality
	_ resource.ResourceWithImportState = &WindowsPlatformScriptAssignmentResource{}

	// Enables plan modification/diff suppression
	_ resource.ResourceWithModifyPlan = &WindowsPlatformScriptAssignmentResource{}
)

func NewWindowsPlatformScriptAssignmentResource() resource.Resource {
	return &WindowsPlatformScriptAssignmentResource{
		ReadPermissions: []string{
			"DeviceManagementConfiguration.Read.All",
			"DeviceManagementConfiguration.ReadWrite.All",
			"DeviceManagementManagedDevices.Read.All",
		},
		WritePermissions: []string{
			"DeviceManagementConfiguration.ReadWrite.All",
			"DeviceManagementManagedDevices.ReadWrite.All",
		},
		ResourcePath: "deviceManagement/deviceShellScripts",
	}
}

type WindowsPlatformScriptAssignmentResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ProviderTypeName string
	TypeName         string
	ReadPermissions  []string
	WritePermissions []string
	ResourcePath     string
}

func (r *WindowsPlatformScriptAssignmentResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	r.ProviderTypeName = req.ProviderTypeName
	r.TypeName = ResourceName
	resp.TypeName = r.FullTypeName()
}

func (r *WindowsPlatformScriptAssignmentResource) FullTypeName() string {
	return r.ProviderTypeName + "_" + r.TypeName
}

func (r *WindowsPlatformScriptAssignmentResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = client.SetGraphBetaClientForResource(ctx, req, resp, constants.PROVIDER_NAME+"_"+ResourceName)
}

func (r *WindowsPlatformScriptAssignmentResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *WindowsPlatformScriptAssignmentResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages Device Management Script Assignments in Microsoft Intune using the" +
			"`/deviceManagement/deviceManagementScripts/{deviceManagementScriptId}/assignments`" +
			"See [MacosPlatformScriptAssignment resource type](https://learn.microsoft.com/en-us/graph/api/resources/intune-devices-MacosPlatformScriptAssignment?view=graph-rest-beta).",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					planmodifiers.UseStateForUnknownString(),
				},
				MarkdownDescription: "The unique identifier for the custom attribute shell script.",
			},
			"windows_platform_script_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The ID of the device management script to assign.",
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(constants.GuidRegex),
						"must be a valid GUID in the format 00000000-0000-0000-0000-000000000000",
					),
				},
			},
			"target": schema.SingleNestedAttribute{
				Required:            true,
				MarkdownDescription: "The assignment target. See [deviceAndAppManagementAssignmentTarget](https://learn.microsoft.com/en-us/graph/api/resources/intune-shared-deviceandappmanagementassignmenttarget?view=graph-rest-beta).",
				Attributes: map[string]schema.Attribute{
					"target_type": schema.StringAttribute{
						Required:            true,
						MarkdownDescription: "The target group type for the script assignment. Possible values: `allDevices`, `allLicensedUsers`, `groupAssignment`.",
						Validators: []validator.String{
							stringvalidator.OneOf(
								"allDevices",
								"allLicensedUsers",
								"groupAssignment",
							),
						},
					},
					"group_id": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "The Entra ID group ID for the assignment target. Required when target_type is 'groupAssignment'.",
						Validators: []validator.String{
							stringvalidator.RegexMatches(
								regexp.MustCompile(constants.GuidRegex),
								"must be a valid GUID in the format 00000000-0000-0000-0000-000000000000",
							),
						},
					},
					"device_and_app_management_assignment_filter_id": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "The Id of the scope filter applied to the target assignment.",
						Validators: []validator.String{
							stringvalidator.RegexMatches(
								regexp.MustCompile(constants.GuidRegex),
								"must be a valid GUID in the format 00000000-0000-0000-0000-000000000000",
							),
						},
					},
					"device_and_app_management_assignment_filter_type": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "The type of scope filter for the target assignment. Possible values: `include`, `exclude`, `none`.",
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
