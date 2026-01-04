package graphBetaInitiateOnDemandProactiveRemediationManagedDevice

import (
	"context"
	"regexp"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/action/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

const (
	ActionName = "microsoft365_graph_beta_device_management_managed_device_initiate_on_demand_proactive_remediation"
	InvokeTimeout = 60
)

var (
	_ action.Action                   = &InitiateOnDemandProactiveRemediationManagedDeviceAction{}
	_ action.ActionWithConfigure      = &InitiateOnDemandProactiveRemediationManagedDeviceAction{}
	_ action.ActionWithValidateConfig = &InitiateOnDemandProactiveRemediationManagedDeviceAction{}
)

func NewInitiateOnDemandProactiveRemediationManagedDeviceAction() action.Action {
	return &InitiateOnDemandProactiveRemediationManagedDeviceAction{
		ReadPermissions: []string{
			"DeviceManagementConfiguration.Read.All",
			"DeviceManagementManagedDevices.Read.All",
		},
		WritePermissions: []string{
			"DeviceManagementConfiguration.Read.All",
			"DeviceManagementManagedDevices.Read.All",
		},
	}
}

type InitiateOnDemandProactiveRemediationManagedDeviceAction struct {
	client           *msgraphbetasdk.GraphServiceClient
	ReadPermissions  []string
	WritePermissions []string
}

func (a *InitiateOnDemandProactiveRemediationManagedDeviceAction) Metadata(ctx context.Context, req action.MetadataRequest, resp *action.MetadataResponse) {
	resp.TypeName = ActionName
}

func (a *InitiateOnDemandProactiveRemediationManagedDeviceAction) Configure(ctx context.Context, req action.ConfigureRequest, resp *action.ConfigureResponse) {
	a.client = client.SetGraphBetaClientForAction(ctx, req, resp, ActionName)
}

func (a *InitiateOnDemandProactiveRemediationManagedDeviceAction) Schema(ctx context.Context, req action.SchemaRequest, resp *action.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Initiates on-demand proactive remediation on managed Windows devices using the " +
			"`/deviceManagement/managedDevices/{managedDeviceId}/initiateOnDemandProactiveRemediation` and " +
			"`/deviceManagement/comanagedDevices/{managedDeviceId}/initiateOnDemandProactiveRemediation` endpoints. " +
			"Proactive remediations (also called remediations or health scripts) are PowerShell scripts that detect and " +
			"automatically fix common support issues on Windows devices. This action triggers immediate execution of a " +
			"specified remediation script on selected devices, rather than waiting for the scheduled run. This is useful " +
			"for urgent fixes, troubleshooting, or validating remediation effectiveness.\n\n" +
			"**Important Notes:**\n" +
			"- Only works on Windows 10/11 devices\n" +
			"- Requires script policy ID (remediation script GUID)\n" +
			"- Script executes immediately on device check-in\n" +
			"- Runs with SYSTEM privileges\n" +
			"- Results available in Intune portal and reports\n" +
			"- Script must be already deployed to the device\n" +
			"- Does not create new script deployment\n\n" +
			"**Use Cases:**\n" +
			"- Urgent issue remediation outside scheduled runs\n" +
			"- Troubleshooting and validation\n" +
			"- Post-incident recovery actions\n" +
			"- Ad-hoc compliance fixes\n" +
			"- Testing new remediation scripts\n" +
			"- End-user requested fixes\n\n" +
			"**Platform Support:**\n" +
			"- **Windows**: Windows 10/11 with Intune management extension\n" +
			"- **Other Platforms**: Not supported (Windows-specific feature)\n\n" +
			"**Reference:** [Microsoft Graph API - Initiate On Demand Proactive Remediation](https://learn.microsoft.com/en-us/graph/api/intune-devices-manageddevice-initiateondemandproactiveremediation?view=graph-rest-beta)",
		Attributes: map[string]schema.Attribute{
			"managed_devices": schema.ListNestedAttribute{
				Optional: true,
				MarkdownDescription: "List of managed devices to initiate proactive remediation for. " +
					"Each entry specifies a device and the remediation script to run.\n\n" +
					"**Note:** At least one of `managed_devices` or `comanaged_devices` must be provided. " +
					"Each device can have a different script policy executed.\n\n" +
					"**Important:** The script policy must already be deployed to the device. This action triggers " +
					"immediate execution but does not create a new deployment.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"device_id": schema.StringAttribute{
							Required: true,
							MarkdownDescription: "The unique identifier (GUID) of the managed device to run the remediation script on.\n\n" +
								"**Example**: `\"12345678-1234-1234-1234-123456789abc\"`",
							Validators: []validator.String{
								stringvalidator.RegexMatches(
									regexp.MustCompile(constants.GuidRegex),
									"device_id must be a valid GUID format (e.g., 12345678-1234-1234-1234-123456789abc)",
								),
							},
						},
						"script_policy_id": schema.StringAttribute{
							Required: true,
							MarkdownDescription: "The unique identifier (GUID) of the proactive remediation script policy to execute.\n\n" +
								"**How to find**: Azure Portal → Intune → Devices → Remediations → Select script → Copy GUID from URL or Properties.\n\n" +
								"**Note**: The script must already be assigned/deployed to the device.\n\n" +
								"**Example**: `\"87654321-4321-4321-4321-ba9876543210\"`",
							Validators: []validator.String{
								stringvalidator.LengthAtLeast(1),
							},
						},
					},
				},
			},
			"comanaged_devices": schema.ListNestedAttribute{
				Optional: true,
				MarkdownDescription: "List of co-managed devices to initiate proactive remediation for. " +
					"These are devices managed by both Intune and Configuration Manager (SCCM).\n\n" +
					"**Note:** At least one of `managed_devices` or `comanaged_devices` must be provided.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"device_id": schema.StringAttribute{
							Required: true,
							MarkdownDescription: "The unique identifier (GUID) of the co-managed device to run the remediation script on.\n\n" +
								"**Example**: `\"12345678-1234-1234-1234-123456789abc\"`",
							Validators: []validator.String{
								stringvalidator.RegexMatches(
									regexp.MustCompile(constants.GuidRegex),
									"device_id must be a valid GUID format (e.g., 12345678-1234-1234-1234-123456789abc)",
								),
							},
						},
						"script_policy_id": schema.StringAttribute{
							Required: true,
							MarkdownDescription: "The unique identifier (GUID) of the proactive remediation script policy to execute.\n\n" +
								"**Note**: The script must already be assigned/deployed to the device.\n\n" +
								"**Example**: `\"87654321-4321-4321-4321-ba9876543210\"`",
							Validators: []validator.String{
								stringvalidator.LengthAtLeast(1),
							},
						},
					},
				},
			},
			"ignore_partial_failures": schema.BoolAttribute{
				Optional: true,
				MarkdownDescription: "If set to `true`, the action will succeed even if some operations fail. " +
					"Failed operations will be reported as warnings instead of errors. " +
					"Default: `false` (action fails if any operation fails).",
			},
			"validate_device_exists": schema.BoolAttribute{
				Optional: true,
				MarkdownDescription: "Whether to validate that devices exist and are Windows devices before attempting remediation. " +
					"Disabling this can speed up planning but may result in runtime errors for non-existent or unsupported devices. " +
					"Default: `true`.",
			},
			"timeouts": commonschema.ActionTimeouts(ctx),
		},
	}
}
