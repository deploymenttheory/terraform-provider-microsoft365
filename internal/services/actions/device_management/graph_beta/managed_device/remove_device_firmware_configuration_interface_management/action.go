package graphBetaRemoveDeviceFirmwareConfigurationInterfaceManagementManagedDevice

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
	ActionName    = "microsoft365_graph_beta_device_management_managed_device_remove_device_firmware_configuration_interface_management"
	InvokeTimeout = 60
)

var (
	_ action.Action                   = &RemoveDeviceFirmwareConfigurationInterfaceManagementManagedDeviceAction{}
	_ action.ActionWithConfigure      = &RemoveDeviceFirmwareConfigurationInterfaceManagementManagedDeviceAction{}
	_ action.ActionWithValidateConfig = &RemoveDeviceFirmwareConfigurationInterfaceManagementManagedDeviceAction{}
)

func NewRemoveDeviceFirmwareConfigurationInterfaceManagementManagedDeviceAction() action.Action {
	return &RemoveDeviceFirmwareConfigurationInterfaceManagementManagedDeviceAction{
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

type RemoveDeviceFirmwareConfigurationInterfaceManagementManagedDeviceAction struct {
	client           *msgraphbetasdk.GraphServiceClient
	ReadPermissions  []string
	WritePermissions []string
}

func (a *RemoveDeviceFirmwareConfigurationInterfaceManagementManagedDeviceAction) Metadata(ctx context.Context, req action.MetadataRequest, resp *action.MetadataResponse) {
	resp.TypeName = ActionName
}

func (a *RemoveDeviceFirmwareConfigurationInterfaceManagementManagedDeviceAction) Configure(ctx context.Context, req action.ConfigureRequest, resp *action.ConfigureResponse) {
	a.client = client.SetGraphBetaClientForAction(ctx, req, resp, ActionName)
}

func (a *RemoveDeviceFirmwareConfigurationInterfaceManagementManagedDeviceAction) Schema(ctx context.Context, req action.SchemaRequest, resp *action.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Removes managed devices from Device Firmware Configuration Interface (DFCI) management in Microsoft Intune using the " +
			"`/deviceManagement/managedDevices/{managedDeviceId}/removeDeviceFirmwareConfigurationInterfaceManagement` and " +
			"`/deviceManagement/comanagedDevices/{managedDeviceId}/removeDeviceFirmwareConfigurationInterfaceManagement` endpoints. " +
			"This action is used to remove the DFCI management capability from devices, reverting them to standard Intune management without " +
			"firmware-level control. DFCI enables Intune to manage UEFI (BIOS) settings on compatible Windows devices, providing low-level security controls. " +
			"After removal, the device's UEFI settings can no longer be managed remotely via Intune.\n\n" +
			"**Important Notes:**\n" +
			"- Only works on Windows devices with DFCI-capable firmware\n" +
			"- Requires devices enrolled with DFCI management enabled\n" +
			"- Removes ability to manage UEFI/BIOS settings remotely\n" +
			"- Does not unenroll device from Intune\n" +
			"- Standard MDM management continues\n" +
			"- Typically used on Surface and compatible OEM devices\n" +
			"- Cannot be easily reversed\n\n" +
			"**Use Cases:**\n" +
			"- Decommissioning devices from DFCI management\n" +
			"- Transitioning to standard management\n" +
			"- Removing firmware-level security controls\n" +
			"- Preparing devices for transfer or resale\n" +
			"- Troubleshooting DFCI-related issues\n" +
			"- Disabling low-level hardware management\n\n" +
			"**Platform Support:**\n" +
			"- **Windows**: DFCI-capable devices only (Surface, select OEM devices)\n" +
			"- **Other Platforms**: Not supported (DFCI is Windows-specific)\n\n" +
			"**Reference:** [Microsoft Graph API - Remove Device Firmware Configuration Interface Management](https://learn.microsoft.com/en-us/graph/api/intune-devices-manageddevice-removedevicefirmwareconfigurationinterfacemanagement?view=graph-rest-beta)",
		Attributes: map[string]schema.Attribute{
			"managed_device_ids": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
				MarkdownDescription: "List of managed device IDs (GUIDs) to remove from DFCI management. These are devices " +
					"fully managed by Intune that currently have DFCI management enabled.\n\n" +
					"**Note:** At least one of `managed_device_ids` or `comanaged_device_ids` must be provided. " +
					"You can provide both to remove DFCI management from different types of devices in one action.\n\n" +
					"**Important:** After removal, these devices will continue standard Intune MDM management but will no longer " +
					"support remote UEFI/BIOS configuration through Intune.\n\n" +
					"Example: `[\"12345678-1234-1234-1234-123456789abc\", \"87654321-4321-4321-4321-ba9876543210\"]`",
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
				MarkdownDescription: "List of co-managed device IDs (GUIDs) to remove from DFCI management. These are devices " +
					"managed by both Intune and Configuration Manager (SCCM) that currently have DFCI management enabled.\n\n" +
					"**Note:** At least one of `managed_device_ids` or `comanaged_device_ids` must be provided.\n\n" +
					"Example: `[\"abcdef12-3456-7890-abcd-ef1234567890\"]`",
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
				MarkdownDescription: "If set to `true`, the action will succeed even if some operations fail. " +
					"Failed operations will be reported as warnings instead of errors. " +
					"Default: `false` (action fails if any operation fails).",
			},
			"validate_device_exists": schema.BoolAttribute{
				Optional: true,
				MarkdownDescription: "Whether to validate that devices exist and are Windows devices before attempting DFCI removal. " +
					"Disabling this can speed up planning but may result in runtime errors for non-existent or non-Windows devices. " +
					"Default: `true`.",
			},
			"timeouts": commonschema.ActionTimeouts(ctx),
		},
	}
}
