package graphBetaSyncManagedDevice

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
	ActionName    = "microsoft365_graph_beta_device_management_managed_device_sync_device"
	InvokeTimeout = 60
)

var (
	_ action.Action                   = &SyncManagedDeviceAction{}
	_ action.ActionWithConfigure      = &SyncManagedDeviceAction{}
	_ action.ActionWithValidateConfig = &SyncManagedDeviceAction{}
)

func NewSyncManagedDeviceAction() action.Action {
	return &SyncManagedDeviceAction{
		ReadPermissions: []string{
			"DeviceManagementManagedDevices.PrivilegedOperations.All",
		},
		WritePermissions: []string{
			"DeviceManagementManagedDevices.PrivilegedOperations.All",
		},
	}
}

type SyncManagedDeviceAction struct {
	client           *msgraphbetasdk.GraphServiceClient
	ReadPermissions  []string
	WritePermissions []string
}

func (a *SyncManagedDeviceAction) Metadata(ctx context.Context, req action.MetadataRequest, resp *action.MetadataResponse) {
	resp.TypeName = ActionName
}

func (a *SyncManagedDeviceAction) Configure(ctx context.Context, req action.ConfigureRequest, resp *action.ConfigureResponse) {
	a.client = client.SetGraphBetaClientForAction(ctx, req, resp, ActionName)
}

func (a *SyncManagedDeviceAction) Schema(ctx context.Context, req action.SchemaRequest, resp *action.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Forces managed and co-managed devices to immediately check in with Intune using the " +
			"`/deviceManagement/managedDevices/{managedDeviceId}/syncDevice` and " +
			"`/deviceManagement/comanagedDevices/{managedDeviceId}/syncDevice` endpoints. " +
			"This action triggers an immediate synchronization, causing devices to apply the latest policies, " +
			"configurations, and updates from Intune without waiting for the standard check-in interval.\n\n" +
			"**What This Action Does:**\n" +
			"- Forces immediate check-in with Intune\n" +
			"- Applies latest policies and configurations\n" +
			"- Downloads pending applications\n" +
			"- Reports updated device inventory\n" +
			"- Enforces compliance evaluation\n" +
			"- Processes queued remote actions\n" +
			"- Updates device status in console\n\n" +
			"**Managed vs Co-Managed Devices:**\n" +
			"- **Managed Devices**: Fully managed by Intune only\n" +
			"- **Co-Managed Devices**: Managed by both Intune and Configuration Manager (SCCM)\n" +
			"- This action supports both types independently or together\n\n" +
			"**Platform Support:**\n" +
			"- **Windows**: Full support (managed and co-managed)\n" +
			"- **macOS**: Full support (managed only)\n" +
			"- **iOS/iPadOS**: Full support (managed only)\n" +
			"- **Android**: Full support (managed only)\n" +
			"- **ChromeOS**: Limited support\n\n" +
			"**Common Use Cases:**\n" +
			"- Apply new policies immediately\n" +
			"- Force app installation/updates\n" +
			"- Trigger compliance re-evaluation\n" +
			"- Update device inventory quickly\n" +
			"- Verify policy deployment\n" +
			"- Troubleshoot deployment issues\n" +
			"- Emergency configuration changes\n\n" +
			"**Check-In Behavior:**\n" +
			"- Normal interval: Every 8 hours (varies by platform)\n" +
			"- This action: Immediate (within 1-5 minutes)\n" +
			"- Device must be online and powered on\n" +
			"- Network connectivity required\n" +
			"- Results visible in Intune admin center\n\n" +
			"**Important Considerations:**\n" +
			"- Device must be online to receive command\n" +
			"- Command queued if device is offline\n" +
			"- Sync completes when device comes online\n" +
			"- Multiple syncs in short period may delay each other\n" +
			"- No user disruption (background operation)\n\n" +
			"**Reference:** [Microsoft Graph API - Sync Device](https://learn.microsoft.com/en-us/graph/api/intune-devices-manageddevice-syncdevice?view=graph-rest-beta)",
		Attributes: map[string]schema.Attribute{
			"managed_device_ids": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
				MarkdownDescription: "List of managed device IDs to sync. These are devices fully managed by Intune only. " +
					"Each ID must be a valid GUID format. Multiple devices can be synced in a single action. " +
					"Example: `[\"12345678-1234-1234-1234-123456789abc\", \"87654321-4321-4321-4321-ba9876543210\"]`\n\n" +
					"**Note:** At least one of `managed_device_ids` or `comanaged_device_ids` must be provided. " +
					"You can provide both to sync different types of devices in one action.",
				Validators: []validator.List{
					listvalidator.ValueStringsAre(
						stringvalidator.RegexMatches(
							regexp.MustCompile(constants.GuidRegex),
							"each device ID must be a valid GUID format (e.g., 12345678-1234-1234-1234-123456789abc)",
						),
					),
				},
			},
			"comanaged_device_ids": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
				MarkdownDescription: "List of co-managed device IDs to sync. These are devices managed by both Intune and " +
					"Configuration Manager (SCCM). Each ID must be a valid GUID format. " +
					"Example: `[\"12345678-1234-1234-1234-123456789abc\"]`\n\n" +
					"**Co-Management Context:**\n" +
					"- Devices managed by both Intune and Configuration Manager\n" +
					"- Typically Windows 10/11 enterprise devices\n" +
					"- Workloads split between Intune and ConfigMgr\n" +
					"- Sync affects Intune-managed workloads only\n\n" +
					"**Note:** At least one of `managed_device_ids` or `comanaged_device_ids` must be provided.",
				Validators: []validator.List{
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
				MarkdownDescription: "When set to `true`, the action will complete successfully even if some devices fail to sync. " +
					"When `false` (default), the action will fail if any device sync fails. " +
					"Use this flag when syncing multiple devices and you want the action to succeed even if some syncs fail.",
			},
			"validate_device_exists": schema.BoolAttribute{
				Optional: true,
				MarkdownDescription: "When set to `true` (default), the action will validate that all specified devices exist " +
					"and support sync before attempting to sync them. " +
					"When `false`, device validation is skipped and the action will attempt to sync devices directly. " +
					"Disabling validation can improve performance but may result in errors if devices don't exist or are unsupported.",
			},
			"timeouts": commonschema.ActionTimeouts(ctx),
		},
	}
}
