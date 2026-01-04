package graphBetaRebootNowManagedDevice

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
	ActionName = "microsoft365_graph_beta_device_management_managed_device_reboot_now"
	InvokeTimeout = 60
)

var (
	_ action.Action                   = &RebootNowManagedDeviceAction{}
	_ action.ActionWithConfigure      = &RebootNowManagedDeviceAction{}
	_ action.ActionWithValidateConfig = &RebootNowManagedDeviceAction{}
)

func NewRebootNowManagedDeviceAction() action.Action {
	return &RebootNowManagedDeviceAction{
		ReadPermissions: []string{
			"DeviceManagementManagedDevices.PrivilegedOperations.All",
		},
		WritePermissions: []string{
			"DeviceManagementManagedDevices.PrivilegedOperations.All",
		},
	}
}

type RebootNowManagedDeviceAction struct {
	client           *msgraphbetasdk.GraphServiceClient
	ReadPermissions  []string
	WritePermissions []string
}

func (a *RebootNowManagedDeviceAction) Metadata(ctx context.Context, req action.MetadataRequest, resp *action.MetadataResponse) {
	resp.TypeName = ActionName
}

func (a *RebootNowManagedDeviceAction) Configure(ctx context.Context, req action.ConfigureRequest, resp *action.ConfigureResponse) {
	a.client = client.SetGraphBetaClientForAction(ctx, req, resp, ActionName)
}

func (a *RebootNowManagedDeviceAction) Schema(ctx context.Context, req action.SchemaRequest, resp *action.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Remotely reboots managed devices using the " +
			"`/deviceManagement/managedDevices/{managedDeviceId}/rebootNow` endpoint. " +
			"This action immediately restarts devices, which is essential for applying updates, " +
			"troubleshooting system issues, or ensuring configuration changes take effect. " +
			"The reboot command is sent to devices immediately if online, or queued for execution " +
			"when the device next checks in with Intune.\n\n" +
			"**Important Notes:**\n" +
			"- Device reboots immediately upon receiving command (if online)\n" +
			"- Any unsaved work on the device will be lost\n" +
			"- Users receive minimal or no warning before reboot\n" +
			"- Reboot is forceful and does not wait for user interaction\n" +
			"- Command is queued if device is offline\n" +
			"- Use with caution during business hours\n\n" +
			"**Use Cases:**\n" +
			"- Applying Windows updates that require restart\n" +
			"- Installing software that requires system reboot\n" +
			"- Troubleshooting devices with performance issues\n" +
			"- Forcing application of configuration profiles\n" +
			"- Clearing temporary system issues\n" +
			"- Maintenance windows for device refresh\n" +
			"- Resolving frozen or unresponsive remote devices\n" +
			"- Completing BitLocker encryption setup\n\n" +
			"**Platform Support:**\n" +
			"- **Windows**: Fully supported (Windows 10/11, including Home edition)\n" +
			"- **macOS**: Supported (requires user-approved MDM or supervised)\n" +
			"- **iOS/iPadOS**: Limited support (supervised devices only)\n" +
			"- **Android**: Not supported for reboot action\n\n" +
			"**Best Practices:**\n" +
			"- Schedule reboots during maintenance windows or off-hours\n" +
			"- Notify users in advance when possible\n" +
			"- Use for non-interactive devices (kiosks, shared devices)\n" +
			"- Consider user impact before rebooting during business hours\n" +
			"- Test with small device groups before bulk operations\n" +
			"- Document reason for reboot in change management system\n" +
			"- Combine with compliance policies for automated maintenance\n\n" +
			"**User Impact:**\n" +
			"- Users may lose unsaved work\n" +
			"- Active sessions are terminated\n" +
			"- Video calls and presentations are interrupted\n" +
			"- File transfers may be interrupted\n" +
			"- Users may not receive advance warning\n" +
			"- Device is unavailable for 2-5 minutes during restart\n\n" +
			"**Reference:** [Microsoft Graph API - Reboot Now](https://learn.microsoft.com/en-us/graph/api/intune-devices-manageddevice-rebootnow?view=graph-rest-beta)",
		Attributes: map[string]schema.Attribute{
			"device_ids": schema.ListAttribute{
				ElementType: types.StringType,
				Required:    true,
				MarkdownDescription: "List of managed device IDs to reboot. " +
					"Each ID must be a valid GUID format. Multiple devices can be rebooted in a single action. " +
					"Example: `[\"12345678-1234-1234-1234-123456789abc\", \"87654321-4321-4321-4321-ba9876543210\"]`\n\n" +
					"**Important:** Devices will reboot immediately when they receive this command. " +
					"Any unsaved work will be lost. Use with caution, especially during business hours.",
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
				MarkdownDescription: "If set to `true`, the action will succeed even if some operations fail. " +
					"Failed operations will be reported as warnings instead of errors. " +
					"Default: `false` (action fails if any operation fails).",
			},
			"validate_device_exists": schema.BoolAttribute{
				Optional: true,
				MarkdownDescription: "Whether to validate that devices exist and support remote reboot before attempting to send the reboot command. " +
					"Disabling this can speed up planning but may result in runtime errors for non-existent or unsupported devices. " +
					"Default: `true`.",
			},
			"timeouts": commonschema.ActionTimeouts(ctx),
		},
	}
}
