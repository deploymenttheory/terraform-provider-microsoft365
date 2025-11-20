package graphBetaWindowsAutopilotDeviceIdentity

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
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
	ResourceName  = "microsoft365_graph_beta_device_management_windows_autopilot_device_identity"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180
)

var (
	// Basic resource interface (CRUD operations)
	_ resource.Resource = &WindowsAutopilotDeviceIdentityResource{}

	// Allows the resource to be configured with the provider client
	_ resource.ResourceWithConfigure = &WindowsAutopilotDeviceIdentityResource{}

	// Enables import functionality
	_ resource.ResourceWithImportState = &WindowsAutopilotDeviceIdentityResource{}

	// Enables plan modification/diff suppression
	_ resource.ResourceWithModifyPlan = &WindowsAutopilotDeviceIdentityResource{}
)

func NewWindowsAutopilotDeviceIdentityResource() resource.Resource {
	return &WindowsAutopilotDeviceIdentityResource{
		ReadPermissions: []string{
			"DeviceManagementServiceConfig.Read.All",
			"DeviceManagementConfiguration.Read.All",
		},
		WritePermissions: []string{
			"DeviceManagementServiceConfig.ReadWrite.All",
			"DeviceManagementConfiguration.ReadWrite.All",
		},
		ResourcePath: "/deviceManagement/windowsAutopilotDeviceIdentities",
	}
}

type WindowsAutopilotDeviceIdentityResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ReadPermissions  []string
	WritePermissions []string
	ResourcePath     string
}

// Metadata returns the resource type name.
func (r *WindowsAutopilotDeviceIdentityResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = ResourceName
}

// Configure sets the client for the resource.
func (r *WindowsAutopilotDeviceIdentityResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = client.SetGraphBetaClientForResource(ctx, req, resp, ResourceName)
}

// ImportState imports the resource state.
func (r *WindowsAutopilotDeviceIdentityResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Schema defines the schema for the resource.
func (r *WindowsAutopilotDeviceIdentityResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages Windows Autopilot device identities using the `/deviceManagement/windowsAutopilotDeviceIdentities` endpoint. Windows Autopilot device identities represent devices registered with the Windows Autopilot service for automated device provisioning.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					planmodifiers.UseStateForUnknownString(),
				},
				MarkdownDescription: "The unique identifier for the Windows Autopilot device identity.",
			},
			"group_tag": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Group Tag of the Windows Autopilot device. This tag can be used to add devices to a logical group and can be leveraged by deployment profiles for group-specific configurations.",
			},
			"purchase_order_identifier": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Purchase Order Identifier of the Windows autopilot device.",
			},
			"serial_number": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Serial number of the Windows autopilot device.",
			},
			"product_key": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Product Key of the Windows autopilot device.",
			},
			"manufacturer": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Oem manufacturer of the Windows autopilot device.",
				PlanModifiers: []planmodifier.String{
					planmodifiers.UseStateForUnknownString(),
				},
			},
			"model": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Model name of the Windows autopilot device.",
				PlanModifiers: []planmodifier.String{
					planmodifiers.UseStateForUnknownString(),
				},
			},
			"enrollment_state": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Intune enrollment state of the Windows autopilot device. Possible values are: unknown, enrolled, pendingReset, failed, notContacted, blocked.",
				Validators: []validator.String{
					stringvalidator.OneOf(
						"unknown",
						"enrolled",
						"pendingReset",
						"failed",
						"notContacted",
						"blocked",
					),
				},
			},
			"last_contacted_date_time": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Intune Last Contacted Date Time of the Windows autopilot device.",
			},
			"addressable_user_name": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Addressable user name.",
			},
			"user_principal_name": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "User Principal Name.",
			},
			"resource_name": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Resource Name.",
			},
			"sku_number": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "SKU Number",
			},
			"system_family": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "System Family",
			},
			"azure_active_directory_device_id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "AAD Device ID - to be deprecated",
			},
			"azure_ad_device_id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "AAD Device ID",
			},
			"managed_device_id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Managed Device ID",
			},
			"display_name": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Display Name",
			},
			"deployment_profile_assignment_status": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Profile assignment status of the Windows autopilot device. Possible values are: unknown, assignedInSync, assignedOutOfSync, assignedUnkownSyncState, notAssigned, pending, failed.",
				Validators: []validator.String{
					stringvalidator.OneOf(
						"unknown",
						"assignedInSync",
						"assignedOutOfSync",
						"assignedUnkownSyncState",
						"notAssigned",
						"pending",
						"failed",
					),
				},
			},
			"deployment_profile_assignment_detailed_status": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Profile assignment detailed status of the Windows autopilot device. Possible values are: none, hardwareRequirementsNotMet, surfaceHubProfileNotSupported, holoLensProfileNotSupported, windowsPcProfileNotSupported, surfaceHub2SProfileNotSupported, unknownFutureValue.",
				Validators: []validator.String{
					stringvalidator.OneOf(
						"none",
						"hardwareRequirementsNotMet",
						"surfaceHubProfileNotSupported",
						"holoLensProfileNotSupported",
						"windowsPcProfileNotSupported",
						"surfaceHub2SProfileNotSupported",
						"unknownFutureValue",
					),
				},
			},
			"deployment_profile_assigned_date_time": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Profile set time of the Windows autopilot device.",
			},
			"remediation_state": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Device Remediation State. Possible values are: unknown, noRemediationRequired, automaticRemediationRequired, manualRemediationRequired, unknownFutureValue.",
				Validators: []validator.String{
					stringvalidator.OneOf(
						"unknown",
						"noRemediationRequired",
						"automaticRemediationRequired",
						"manualRemediationRequired",
						"unknownFutureValue",
					),
				},
			},
			"remediation_state_last_modified_date_time": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "RemediationState set time of Autopilot device.",
			},
			"userless_enrollment_status": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Enrollment status for userless enrollments. Possible values are: unknown, allowed, blocked, unknownFutureValue.",
				Validators: []validator.String{
					stringvalidator.OneOf(
						"unknown",
						"allowed",
						"blocked",
						"unknownFutureValue",
					),
				},
			},
			"user_assignment": schema.SingleNestedAttribute{
				Optional:            true,
				MarkdownDescription: "User assignment configuration for the Windows Autopilot device. This block allows you to assign or unassign users to the device.",
				Attributes: map[string]schema.Attribute{
					"user_principal_name": schema.StringAttribute{
						Required:            true,
						MarkdownDescription: "User Principal Name (UPN) of the user to be assigned to the device.",
						Validators: []validator.String{
							stringvalidator.LengthAtLeast(1),
						},
					},
					"addressable_user_name": schema.StringAttribute{
						Computed:            true,
						MarkdownDescription: "Addressable user name that gets set on the device.",
					},
				},
			},
			"timeouts": commonschema.Timeouts(ctx),
		},
	}
}
