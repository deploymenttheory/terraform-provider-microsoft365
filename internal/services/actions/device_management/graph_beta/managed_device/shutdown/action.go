package graphBetaShutdownManagedDevice

import (
	"context"
	"regexp"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/action/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

const (
	ActionName    = "microsoft365_graph_beta_device_management_managed_device_shutdown"
	InvokeTimeout = 60
)

var (
	_ action.Action                   = &ShutdownManagedDeviceAction{}
	_ action.ActionWithConfigure      = &ShutdownManagedDeviceAction{}
	_ action.ActionWithValidateConfig = &ShutdownManagedDeviceAction{}
)

func NewShutdownManagedDeviceAction() action.Action {
	return &ShutdownManagedDeviceAction{
		ReadPermissions: []string{
			"DeviceManagementManagedDevices.PrivilegedOperations.All",
		},
		WritePermissions: []string{
			"DeviceManagementManagedDevices.PrivilegedOperations.All",
		},
	}
}

type ShutdownManagedDeviceAction struct {
	client           *msgraphbetasdk.GraphServiceClient
	ReadPermissions  []string
	WritePermissions []string
}

func (a *ShutdownManagedDeviceAction) Metadata(ctx context.Context, req action.MetadataRequest, resp *action.MetadataResponse) {
	resp.TypeName = ActionName
}

func (a *ShutdownManagedDeviceAction) Configure(ctx context.Context, req action.ConfigureRequest, resp *action.ConfigureResponse) {
	a.client = client.SetGraphBetaClientForAction(ctx, req, resp, ActionName)
}

func (a *ShutdownManagedDeviceAction) Schema(ctx context.Context, req action.SchemaRequest, resp *action.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Shuts down managed devices in Microsoft Intune using the " +
			"`/deviceManagement/managedDevices/{managedDeviceId}/shutDown` endpoint. " +
			"This action is used to power off devices completely, which is useful for energy conservation, " +
			"maintenance operations, or security scenarios. Unlike reboot, shutdown powers the device " +
			"off completely and requires manual intervention to power it back on.\n\n" +
			"**Important Notes:**\n" +
			"- Device shuts down completely (powers off)\n" +
			"- Device requires manual power-on to restart\n" +
			"- Any unsaved work on the device will be lost\n" +
			"- Users receive minimal or no warning before shutdown\n" +
			"- Shutdown is forceful and does not wait for user interaction\n" +
			"- Command is queued if device is offline\n" +
			"- Use with extreme caution - device will be completely offline\n\n" +
			"**Use Cases:**\n" +
			"- Energy conservation during extended non-use periods\n" +
			"- Security incident response (isolate compromised device)\n" +
			"- Hardware maintenance requiring full power-off\n" +
			"- Decommissioning devices before storage or shipment\n" +
			"- Emergency response to prevent data exfiltration\n" +
			"- Scheduled shutdowns for lab or classroom devices\n" +
			"- Reducing power consumption in device fleets\n" +
			"- Preparing devices for physical relocation\n\n" +
			"**Platform Support:**\n" +
			"- **Windows**: Fully supported (Windows 10/11, including Home edition)\n" +
			"- **macOS**: Supported (requires user-approved MDM or supervised)\n" +
			"- **iOS/iPadOS**: Limited support (supervised devices only, rare use case)\n" +
			"- **Android**: Not supported for shutdown action\n\n" +
			"**Shutdown vs Reboot:**\n" +
			"- **Shutdown**: Device powers off completely, requires manual restart\n" +
			"- **Reboot**: Device automatically restarts, comes back online\n" +
			"- Use shutdown for: Long-term offline, security incidents, energy savings\n" +
			"- Use reboot for: Updates, troubleshooting, configuration changes\n\n" +
			"**Best Practices:**\n" +
			"- Only use when device needs to remain offline\n" +
			"- Ensure physical access is available to power device back on\n" +
			"- Notify users before shutdown (device will be offline)\n" +
			"- Schedule for end of business day or weekends\n" +
			"- Document reason for shutdown in change management\n" +
			"- Verify device location before shutdown (ensure accessibility)\n" +
			"- Consider reboot instead if device needs to come back online\n" +
			"- Test with small groups before bulk operations\n\n" +
			"**User Impact:**\n" +
			"- Users lose all unsaved work\n" +
			"- Device becomes completely unavailable\n" +
			"- Active sessions are terminated\n" +
			"- Physical access required to power device back on\n" +
			"- May cause significant disruption to user productivity\n" +
			"- Users cannot access device remotely after shutdown\n\n" +
			"**Reference:** [Microsoft Graph API - Shutdown](https://learn.microsoft.com/en-us/graph/api/intune-devices-manageddevice-shutdown?view=graph-rest-beta)",
		Attributes: map[string]schema.Attribute{
			"device_ids": schema.ListAttribute{
				ElementType: types.StringType,
				Required:    true,
				MarkdownDescription: "List of managed device IDs to shut down. " +
					"Each ID must be a valid GUID format. Multiple devices can be shut down in a single action. " +
					"Example: `[\"12345678-1234-1234-1234-123456789abc\", \"87654321-4321-4321-4321-ba9876543210\"]`\n\n" +
					"**Critical Warning:** Devices will power off completely when they receive this command. " +
					"Physical access will be required to power devices back on. Any unsaved work will be lost. " +
					"Use this action only when devices need to remain powered off.",
				Validators: []validator.List{
					listvalidator.SizeAtLeast(1),
					listvalidator.ValueStringsAre(
						stringvalidator.RegexMatches(
							regexp.MustCompile(constants.GuidRegex),
							"each device ID must be a valid GUID format (e.g., 12345678-1234-1234-1234-123456789abc)",
						),
					),
				},
			},
			"ignore_partial_failures": schema.BoolAttribute{
				Optional: true,
				MarkdownDescription: "When set to `true`, the action will complete successfully even if some devices fail to shut down. " +
					"When `false` (default), the action will fail if any device shutdown fails. " +
					"Use this flag when shutting down multiple devices and you want the action to succeed even if some shutdowns fail.",
			},
			"validate_device_exists": schema.BoolAttribute{
				Optional: true,
				MarkdownDescription: "When set to `true` (default), the action will validate that all specified devices exist " +
					"and support shutdown before attempting to shut them down. " +
					"When `false`, device validation is skipped and the action will attempt to shut down devices directly. " +
					"Disabling validation can improve performance but may result in errors if devices don't exist or are unsupported.",
			},
			"timeouts": commonschema.ActionTimeouts(ctx),
		},
	}
}
