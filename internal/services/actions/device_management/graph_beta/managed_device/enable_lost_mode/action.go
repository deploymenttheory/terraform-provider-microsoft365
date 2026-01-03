package graphBetaEnableLostModeManagedDevice

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
	ActionName = "microsoft365_graph_beta_device_management_managed_device_enable_lost_mode"
	InvokeTimeout = 60
)

var (
	_ action.Action                   = &EnableLostModeManagedDeviceAction{}
	_ action.ActionWithConfigure      = &EnableLostModeManagedDeviceAction{}
	_ action.ActionWithValidateConfig = &EnableLostModeManagedDeviceAction{}
)

func NewEnableLostModeManagedDeviceAction() action.Action {
	return &EnableLostModeManagedDeviceAction{
		ReadPermissions: []string{
			"DeviceManagementManagedDevices.PrivilegedOperations.All",
		},
		WritePermissions: []string{
			"DeviceManagementManagedDevices.PrivilegedOperations.All",
		},
	}
}

type EnableLostModeManagedDeviceAction struct {
	client           *msgraphbetasdk.GraphServiceClient
	ReadPermissions  []string
	WritePermissions []string
}

func (a *EnableLostModeManagedDeviceAction) Metadata(ctx context.Context, req action.MetadataRequest, resp *action.MetadataResponse) {
	resp.TypeName = ActionName
}

func (a *EnableLostModeManagedDeviceAction) Configure(ctx context.Context, req action.ConfigureRequest, resp *action.ConfigureResponse) {
	a.client = client.SetGraphBetaClientForAction(ctx, req, resp, ActionName)
}

func (a *EnableLostModeManagedDeviceAction) Schema(ctx context.Context, req action.SchemaRequest, resp *action.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Enables lost mode on iOS/iPadOS managed devices using the " +
			"`/deviceManagement/managedDevices/{managedDeviceId}/enableLostMode` and " +
			"`/deviceManagement/comanagedDevices/{managedDeviceId}/enableLostMode` endpoints. " +
			"This action locks the device and displays a custom message with contact information on the lock screen. " +
			"Lost mode is a feature that helps locate and secure lost iOS/iPadOS devices by locking them " +
			"and enabling device location tracking. " +
			"This action supports enabling lost mode on multiple devices in a single operation with per-device messages.\n\n" +
			"**Important Notes:**\n" +
			"- Only applicable to iOS and iPadOS devices (iOS 9.3+)\n" +
			"- Device must be supervised\n" +
			"- Requires device to be online to receive command\n" +
			"- Locks device and displays custom message with contact information\n" +
			"- Enables device location tracking\n" +
			"- Each device can have its own custom message, phone number, and footnote\n\n" +
			"**Use Cases:**\n" +
			"- Device has been reported lost or stolen\n" +
			"- Need to lock device and display recovery contact information\n" +
			"- Need to track device location for recovery\n" +
			"- Prevent unauthorized access to corporate data\n\n" +
			"**Platform Support:**\n" +
			"- **iOS/iPadOS**: Fully supported (iOS 9.3+, supervised devices only)\n" +
			"- **Other Platforms**: Not applicable - lost mode is iOS/iPadOS only\n\n" +
			"**Reference:** [Microsoft Graph API - Enable Lost Mode](https://learn.microsoft.com/en-us/graph/api/intune-devices-manageddevice-enablelostmode?view=graph-rest-beta)",
		Attributes: map[string]schema.Attribute{
			"managed_devices": schema.ListNestedAttribute{
				Optional: true,
				MarkdownDescription: "List of managed devices to enable lost mode for. These are iOS/iPadOS devices fully managed by Intune only. " +
					"Each entry specifies a device ID and the custom lost mode configuration for that device.\n\n" +
					"**Note:** At least one of `managed_devices` or `comanaged_devices` must be provided. " +
					"You can provide both to enable lost mode on different types of devices in one action.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"device_id": schema.StringAttribute{
							Required: true,
							MarkdownDescription: "The unique identifier (GUID) of the managed device to enable lost mode for. " +
								"Example: `\"12345678-1234-1234-1234-123456789abc\"`",
							Validators: []validator.String{
								stringvalidator.RegexMatches(
									regexp.MustCompile(constants.GuidRegex),
									"device_id must be a valid GUID format (e.g., 12345678-1234-1234-1234-123456789abc)",
								),
							},
						},
						"message": schema.StringAttribute{
							Required: true,
							MarkdownDescription: "The message to display on this device's lock screen. This message should provide information on how to return the device. " +
								"Example: `\"This device has been lost. Please contact IT at 555-0123 to return.\"`\n\n" +
								"**Requirements:**\n" +
								"- Must not be empty\n" +
								"- Should include clear instructions for device return\n" +
								"- Recommended: Include contact information and identification details",
							Validators: []validator.String{
								stringvalidator.LengthAtLeast(1),
							},
						},
						"phone_number": schema.StringAttribute{
							Required: true,
							MarkdownDescription: "The phone number to display on this device's lock screen. This should be a contact number for returning the device. " +
								"Example: `\"555-0123\"` or `\"+1-555-0123\"`\n\n" +
								"**Requirements:**\n" +
								"- Must not be empty\n" +
								"- Should be a valid phone number format\n" +
								"- Can include international dialing codes",
							Validators: []validator.String{
								stringvalidator.LengthAtLeast(1),
							},
						},
						"footnote": schema.StringAttribute{
							Optional: true,
							MarkdownDescription: "An optional footnote to display below the message on this device's lock screen. " +
								"This can be used for additional instructions or legal information. " +
								"Example: `\"Property of Contoso Corporation\"`",
						},
					},
				},
			},
			"comanaged_devices": schema.ListNestedAttribute{
				Optional: true,
				MarkdownDescription: "List of co-managed devices to enable lost mode for. These are iOS/iPadOS devices managed by both Intune and " +
					"Configuration Manager (SCCM). Each entry specifies a device ID and the custom lost mode configuration.\n\n" +
					"**Note:** At least one of `managed_devices` or `comanaged_devices` must be provided.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"device_id": schema.StringAttribute{
							Required: true,
							MarkdownDescription: "The unique identifier (GUID) of the co-managed device to enable lost mode for. " +
								"Example: `\"12345678-1234-1234-1234-123456789abc\"`",
							Validators: []validator.String{
								stringvalidator.RegexMatches(
									regexp.MustCompile(constants.GuidRegex),
									"device_id must be a valid GUID format (e.g., 12345678-1234-1234-1234-123456789abc)",
								),
							},
						},
						"message": schema.StringAttribute{
							Required: true,
							MarkdownDescription: "The message to display on this device's lock screen. This message should provide information on how to return the device. " +
								"Example: `\"This device has been lost. Please contact IT at 555-0123 to return.\"`",
							Validators: []validator.String{
								stringvalidator.LengthAtLeast(1),
							},
						},
						"phone_number": schema.StringAttribute{
							Required: true,
							MarkdownDescription: "The phone number to display on this device's lock screen. " +
								"Example: `\"555-0123\"`",
							Validators: []validator.String{
								stringvalidator.LengthAtLeast(1),
							},
						},
						"footer": schema.StringAttribute{
							Optional:            true,
							MarkdownDescription: "An optional footer to display below the message on this device's lock screen.",
						},
					},
				},
			},
			"timeouts": commonschema.ActionTimeouts(ctx),
		},
	}
}
