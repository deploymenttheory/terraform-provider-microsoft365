package graphDeviceConfigurationAssignment

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	planmodifiers "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/plan_modifiers"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	msgraphsdk "github.com/microsoftgraph/msgraph-sdk-go"
)

const (
	ResourceName  = "microsoft365_graph_device_management_device_configuration_assignment"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180
)

var (
	// Basic resource interface (CRUD operations)
	_ resource.Resource = &DeviceConfigurationAssignmentResource{}

	// Allows the resource to be configured with the provider client
	_ resource.ResourceWithConfigure = &DeviceConfigurationAssignmentResource{}

	// Enables import functionality
	_ resource.ResourceWithImportState = &DeviceConfigurationAssignmentResource{}

	// Enables plan modification/diff suppression
	_ resource.ResourceWithModifyPlan = &DeviceConfigurationAssignmentResource{}
)

func NewDeviceConfigurationAssignmentResource() resource.Resource {
	return &DeviceConfigurationAssignmentResource{
		ReadPermissions: []string{
			"DeviceManagementConfiguration.Read.All",
		},
		WritePermissions: []string{
			"DeviceManagementConfiguration.ReadWrite.All",
		},
		ResourcePath: "/deviceManagement/deviceConfigurations/{deviceConfigurationId}/assignments",
	}
}

type DeviceConfigurationAssignmentResource struct {
	client           *msgraphsdk.GraphServiceClient
	ReadPermissions  []string
	WritePermissions []string
	ResourcePath     string
}

// Metadata returns the resource type name.
func (r *DeviceConfigurationAssignmentResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = ResourceName
}

// Configure sets the client for the resource.
func (r *DeviceConfigurationAssignmentResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = client.SetGraphStableClientForResource(ctx, req, resp, ResourceName)
}

// ImportState imports the resource state.
func (r *DeviceConfigurationAssignmentResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// The import ID is expected to be in the format "deviceConfigId:assignmentId"
	idParts := strings.Split(req.ID, ":")
	if len(idParts) != 2 {
		resp.Diagnostics.AddError(
			"Invalid Import ID",
			fmt.Sprintf("Expected import ID in format 'deviceConfigId:assignmentId', got: %s", req.ID),
		)
		return
	}

	deviceConfigId := idParts[0]
	assignmentId := idParts[1]

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("device_configuration_id"), deviceConfigId)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), assignmentId)...)
}

// Schema returns the schema for the resource.
func (r *DeviceConfigurationAssignmentResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages device configuration assignments using the `/deviceManagement/deviceConfigurations/{deviceConfigurationId}/assignments` endpoint. This resource is used to controls assignments for device configurations.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The unique identifier of the device configuration assignment.",
				PlanModifiers: []planmodifier.String{
					planmodifiers.UseStateForUnknownString(),
				},
			},
			"device_configuration_id": schema.StringAttribute{
				Required:    true,
				Description: "The ID of the device configuration policy to assign.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(constants.GuidRegex),
						"must be a valid GUID in the format 00000000-0000-0000-0000-000000000000",
					),
				},
			},
			"target_type": schema.StringAttribute{
				Required: true,
				Description: "The target group type for the device configuration assignment. Possible values are:\n\n" +
					"- **allDevices**: Target all devices in the tenant\n" +
					"- **allLicensedUsers**: Target all licensed users in the tenant\n" +
					"- **configurationManagerCollection**: Target System Centre Configuration Manager collection\n" +
					"- **exclusionGroupAssignment**: Target a specific Entra ID group for exclusion\n" +
					"- **groupAssignment**: Target a specific Entra ID group",
				Validators: []validator.String{
					stringvalidator.OneOf(
						"allDevices",
						"allLicensedUsers",
						"configurationManagerCollection",
						"exclusionGroupAssignment",
						"groupAssignment",
					),
				},
			},
			"group_id": schema.StringAttribute{
				Optional: true,
				Computed: true,
				MarkdownDescription: "The ID of the Azure AD group to assign the device configuration to.\n\n" +
					"**Required when:**\n" +
					"- `target_type` is `groupAssignment` or `exclusionGroupAssignment`\n" +
					"- `target_type` is `configurationManagerCollection` (represents the collection ID)\n\n" +
					"**Not used when:**\n" +
					"- `target_type` is `allDevices` or `allLicensedUsers`",
				Default: stringdefault.StaticString(""),
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(constants.GuidRegex),
						"must be a valid GUID in the format 00000000-0000-0000-0000-000000000000",
					),
				},
			},
			"timeouts": commonschema.ResourceTimeouts(ctx),
		},
	}
}
