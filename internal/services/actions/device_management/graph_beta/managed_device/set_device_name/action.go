package graphBetaSetDeviceNameManagedDevice

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
	ActionName    = "microsoft365_graph_beta_device_management_managed_device_set_device_name"
	InvokeTimeout = 60
)

var (
	_ action.Action                   = &SetDeviceNameManagedDeviceAction{}
	_ action.ActionWithConfigure      = &SetDeviceNameManagedDeviceAction{}
	_ action.ActionWithValidateConfig = &SetDeviceNameManagedDeviceAction{}
)

func NewSetDeviceNameManagedDeviceAction() action.Action {
	return &SetDeviceNameManagedDeviceAction{
		ReadPermissions: []string{
			"DeviceManagementManagedDevices.PrivilegedOperations.All",
		},
		WritePermissions: []string{
			"DeviceManagementManagedDevices.PrivilegedOperations.All",
		},
	}
}

type SetDeviceNameManagedDeviceAction struct {
	client           *msgraphbetasdk.GraphServiceClient
	ReadPermissions  []string
	WritePermissions []string
}

func (a *SetDeviceNameManagedDeviceAction) Metadata(ctx context.Context, req action.MetadataRequest, resp *action.MetadataResponse) {
	resp.TypeName = ActionName
}

func (a *SetDeviceNameManagedDeviceAction) Configure(ctx context.Context, req action.ConfigureRequest, resp *action.ConfigureResponse) {
	a.client = client.SetGraphBetaClientForAction(ctx, req, resp, ActionName)
}

func (a *SetDeviceNameManagedDeviceAction) Schema(ctx context.Context, req action.SchemaRequest, resp *action.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Sets a custom device name for managed devices in Microsoft Intune using the " +
			"`/deviceManagement/managedDevices/{managedDeviceId}/setDeviceName` and " +
			"`/deviceManagement/comanagedDevices/{managedDeviceId}/setDeviceName` endpoints. " +
			"This action is used to assign meaningful, custom names to devices for easier identification " +
			"and management in the Intune console. Device names can be used to reflect location, user, function, or " +
			"organizational naming conventions. This action supports setting names on multiple devices in a single operation " +
			"with per-device name customization.\n\n" +
			"**Important Notes:**\n" +
			"- Device name length and character restrictions vary by platform\n" +
			"- Some platforms may have specific naming conventions or limitations\n" +
			"- Device must be online to receive the name change command\n" +
			"- Name changes may take time to reflect after device check-in\n" +
			"- Each device can have its own unique custom name\n\n" +
			"**Use Cases:**\n" +
			"- Implementing organizational naming conventions\n" +
			"- Identifying devices by location (e.g., 'NYC-Floor3-Conf-01')\n" +
			"- Associating devices with users or departments\n" +
			"- Standardizing device names across the organization\n" +
			"- Renaming devices after reassignment or relocation\n\n" +
			"**Platform Support:**\n" +
			"- **Windows**: Fully supported with various name length restrictions\n" +
			"- **iOS/iPadOS**: Supported for supervised devices\n" +
			"- **macOS**: Supported for managed devices\n" +
			"- **Android**: Support varies by management mode\n\n" +
			"**Reference:** [Microsoft Graph API - Set Device Name](https://learn.microsoft.com/en-us/graph/api/intune-devices-manageddevice-setdevicename?view=graph-rest-beta)",
		Attributes: map[string]schema.Attribute{
			"managed_devices": schema.ListNestedAttribute{
				Optional: true,
				MarkdownDescription: "List of managed devices to set custom names for. These are devices fully managed by Intune only. " +
					"Each entry specifies a device ID and the new name to assign to that device.\n\n" +
					"**Note:** At least one of `managed_devices` or `comanaged_devices` must be provided. " +
					"You can provide both to set names on different types of devices in one action.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"device_id": schema.StringAttribute{
							Required: true,
							MarkdownDescription: "The unique identifier (GUID) of the managed device to rename. " +
								"Example: `\"12345678-1234-1234-1234-123456789abc\"`",
							Validators: []validator.String{
								stringvalidator.RegexMatches(
									regexp.MustCompile(constants.GuidRegex),
									"device_id must be a valid GUID format (e.g., 12345678-1234-1234-1234-123456789abc)",
								),
							},
						},
						"device_name": schema.StringAttribute{
							Required: true,
							MarkdownDescription: "The new name to assign to this device. Device naming requirements vary by platform. " +
								"Consult platform-specific documentation for character and length restrictions. " +
								"Example: `\"NYC-Marketing-Laptop-01\"`",
							Validators: []validator.String{
								stringvalidator.LengthAtLeast(1),
								stringvalidator.LengthAtMost(255),
							},
						},
					},
				},
			},
			"comanaged_devices": schema.ListNestedAttribute{
				Optional: true,
				MarkdownDescription: "List of co-managed devices to set custom names for. These are devices managed by both Intune and " +
					"Configuration Manager (SCCM). Each entry specifies a device ID and the new name to assign.\n\n" +
					"**Note:** At least one of `managed_devices` or `comanaged_devices` must be provided.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"device_id": schema.StringAttribute{
							Required: true,
							MarkdownDescription: "The unique identifier (GUID) of the co-managed device to rename. " +
								"Example: `\"12345678-1234-1234-1234-123456789abc\"`",
							Validators: []validator.String{
								stringvalidator.RegexMatches(
									regexp.MustCompile(constants.GuidRegex),
									"device_id must be a valid GUID format (e.g., 12345678-1234-1234-1234-123456789abc)",
								),
							},
						},
						"device_name": schema.StringAttribute{
							Required: true,
							MarkdownDescription: "The new name to assign to this co-managed device. " +
								"Example: `\"NYC-IT-Desktop-05\"`",
							Validators: []validator.String{
								stringvalidator.LengthAtLeast(1),
								stringvalidator.LengthAtMost(255),
							},
						},
					},
				},
			},
			"ignore_partial_failures": schema.BoolAttribute{
				Optional: true,
				MarkdownDescription: "When set to `true`, the action will complete successfully even if some devices fail to be renamed. " +
					"When `false` (default), the action will fail if any device rename operation fails. " +
					"Use this flag when renaming multiple devices and you want the action to succeed even if some renames fail.",
			},
			"validate_device_exists": schema.BoolAttribute{
				Optional: true,
				MarkdownDescription: "When set to `true` (default), the action will validate that all specified devices exist " +
					"before attempting to rename them. When `false`, device validation is skipped and the action will attempt to rename devices directly. " +
					"Disabling validation can improve performance but may result in errors if devices don't exist.",
			},
			"timeouts": commonschema.ActionTimeouts(ctx),
		},
	}
}
