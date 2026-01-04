package graphBetaWindowsDefenderUpdateSignatures

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
	ActionName = "microsoft365_graph_beta_device_management_managed_device_windows_defender_update_signatures"
	InvokeTimeout = 60
)

var (
	_ action.Action                   = &WindowsDefenderUpdateSignaturesAction{}
	_ action.ActionWithConfigure      = &WindowsDefenderUpdateSignaturesAction{}
	_ action.ActionWithValidateConfig = &WindowsDefenderUpdateSignaturesAction{}
)

func NewWindowsDefenderUpdateSignaturesAction() action.Action {
	return &WindowsDefenderUpdateSignaturesAction{
		ReadPermissions: []string{
			"DeviceManagementManagedDevices.PrivilegedOperations.All",
		},
		WritePermissions: []string{
			"DeviceManagementManagedDevices.PrivilegedOperations.All",
		},
	}
}

type WindowsDefenderUpdateSignaturesAction struct {
	client           *msgraphbetasdk.GraphServiceClient
	ReadPermissions  []string
	WritePermissions []string
}

func (a *WindowsDefenderUpdateSignaturesAction) Metadata(ctx context.Context, req action.MetadataRequest, resp *action.MetadataResponse) {
	resp.TypeName = ActionName
}

func (a *WindowsDefenderUpdateSignaturesAction) Configure(ctx context.Context, req action.ConfigureRequest, resp *action.ConfigureResponse) {
	a.client = client.SetGraphBetaClientForAction(ctx, req, resp, ActionName)
}

func (a *WindowsDefenderUpdateSignaturesAction) Schema(ctx context.Context, req action.SchemaRequest, resp *action.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Forces Windows devices to immediately update Windows Defender (Microsoft Defender Antivirus) signatures using the " +
			"`/deviceManagement/managedDevices/{managedDeviceId}/windowsDefenderUpdateSignatures` and " +
			"`/deviceManagement/comanagedDevices/{managedDeviceId}/windowsDefenderUpdateSignatures` endpoints. " +
			"This action triggers an immediate update of antivirus definitions without waiting for the standard update schedule.\n\n" +
			"**What This Action Does:**\n" +
			"- Forces immediate signature update\n" +
			"- Downloads latest threat definitions\n" +
			"- Updates malware detection database\n" +
			"- Ensures current threat protection\n" +
			"- Works on managed and co-managed devices\n" +
			"- No device reboot required\n" +
			"- Completes in 1-5 minutes\n\n" +
			"**When to Use:**\n" +
			"- Zero-day threat emergence\n" +
			"- Critical security updates\n" +
			"- Before antivirus scans\n" +
			"- After new threat intel\n" +
			"- Compliance requirements\n" +
			"- Outdated definitions detected\n" +
			"- Emergency response scenarios\n\n" +
			"**Platform Support:**\n" +
			"- **Windows 10/11**: Full support (managed and co-managed)\n" +
			"- **Windows Server**: Full support (if Defender enabled)\n" +
			"- **Other platforms**: Not supported (Windows Defender only)\n\n" +
			"**Update Process:**\n" +
			"- Device receives update command\n" +
			"- Connects to Microsoft Update servers\n" +
			"- Downloads latest signatures\n" +
			"- Applies updates automatically\n" +
			"- Reports completion to Intune\n" +
			"- No user interaction required\n\n" +
			"**Important Considerations:**\n" +
			"- Device must be online\n" +
			"- Internet connectivity required\n" +
			"- Minimal performance impact\n" +
			"- Updates in background\n" +
			"- No device reboot needed\n" +
			"- Automatic threat protection\n\n" +
			"**Reference:** [Microsoft Graph API - Windows Defender Update Signatures](https://learn.microsoft.com/en-us/graph/api/intune-devices-manageddevice-windowsdefenderupdatesignatures?view=graph-rest-beta)",
		Attributes: map[string]schema.Attribute{
			"managed_device_ids": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
				MarkdownDescription: "List of managed device IDs to update Windows Defender signatures. These are devices fully managed by Intune only. " +
					"Each ID must be a valid GUID format. Multiple devices can be updated in a single action. " +
					"Example: `[\"12345678-1234-1234-1234-123456789abc\", \"87654321-4321-4321-4321-ba9876543210\"]`\n\n" +
					"**Note:** At least one of `managed_device_ids` or `comanaged_device_ids` must be provided. " +
					"You can provide both to update different types of devices in one action.",
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
				MarkdownDescription: "List of co-managed device IDs to update Windows Defender signatures. These are devices managed by both Intune and " +
					"Configuration Manager (SCCM). Each ID must be a valid GUID format. " +
					"Example: `[\"12345678-1234-1234-1234-123456789abc\"]`\n\n" +
					"**Co-Management Context:**\n" +
					"- Devices managed by both Intune and Configuration Manager\n" +
					"- Typically Windows 10/11 enterprise devices\n" +
					"- This action updates signatures via Intune endpoint\n" +
					"- ConfigMgr can also manage definition updates independently\n" +
					"- No conflict between systems\n\n" +
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
				MarkdownDescription: "When set to `true`, the action will complete successfully even if some devices fail to update signatures. " +
					"When `false` (default), the action will fail if any device update fails. " +
					"Use this flag when updating multiple devices and you want the action to succeed even if some updates fail.",
			},
			"validate_device_exists": schema.BoolAttribute{
				Optional: true,
				MarkdownDescription: "When set to `true` (default), the action will validate that all specified devices exist " +
					"and are Windows devices before attempting to update signatures. " +
					"When `false`, device validation is skipped and the action will attempt to update signatures directly. " +
					"Disabling validation can improve performance but may result in errors if devices don't exist or are not Windows devices.",
			},
			"timeouts": commonschema.ActionTimeouts(ctx),
		},
	}
}
