package graphBetaPlayLostModeSoundManagedDevice

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
	ActionName = "graph_beta_device_management_managed_device_play_lost_mode_sound"
)

var (
	_ action.Action                   = &PlayLostModeSoundManagedDeviceAction{}
	_ action.ActionWithConfigure      = &PlayLostModeSoundManagedDeviceAction{}
	_ action.ActionWithValidateConfig = &PlayLostModeSoundManagedDeviceAction{}
)

func NewPlayLostModeSoundManagedDeviceAction() action.Action {
	return &PlayLostModeSoundManagedDeviceAction{
		ReadPermissions: []string{
			"DeviceManagementManagedDevices.PrivilegedOperations.All",
		},
		WritePermissions: []string{
			"DeviceManagementManagedDevices.PrivilegedOperations.All",
		},
	}
}

type PlayLostModeSoundManagedDeviceAction struct {
	client           *msgraphbetasdk.GraphServiceClient
	ProviderTypeName string
	TypeName         string
	ReadPermissions  []string
	WritePermissions []string
}

func (a *PlayLostModeSoundManagedDeviceAction) Metadata(ctx context.Context, req action.MetadataRequest, resp *action.MetadataResponse) {
	a.ProviderTypeName = req.ProviderTypeName
	a.TypeName = ActionName
	resp.TypeName = a.FullTypeName()
}

func (a *PlayLostModeSoundManagedDeviceAction) FullTypeName() string {
	return a.ProviderTypeName + "_" + ActionName
}

func (a *PlayLostModeSoundManagedDeviceAction) Configure(ctx context.Context, req action.ConfigureRequest, resp *action.ConfigureResponse) {
	a.client = client.SetGraphBetaClientForAction(ctx, req, resp, constants.PROVIDER_NAME+"_"+ActionName)
}

func (a *PlayLostModeSoundManagedDeviceAction) Schema(ctx context.Context, req action.SchemaRequest, resp *action.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Plays a sound on iOS/iPadOS managed devices in lost mode using the " +
			"`/deviceManagement/managedDevices/{managedDeviceId}/playLostModeSound` and " +
			"`/deviceManagement/comanagedDevices/{managedDeviceId}/playLostModeSound` endpoints. " +
			"This action helps locate lost devices by triggering an audible alert that plays even if the device is in silent mode. " +
			"The sound plays for a specified duration to assist in physically locating the device. " +
			"This action supports playing sounds on multiple devices in a single operation with per-device duration settings.\n\n" +
			"**Important Notes:**\n" +
			"- Only applicable to iOS and iPadOS devices in lost mode\n" +
			"- Device must be supervised\n" +
			"- Device must currently be in lost mode\n" +
			"- Sound plays even if device is in silent mode\n" +
			"- Requires device to be online to receive command\n" +
			"- Each device can have its own sound duration\n\n" +
			"**Use Cases:**\n" +
			"- Device is nearby but cannot be visually located\n" +
			"- Device is in lost mode and needs audible alert\n" +
			"- Assisting user in finding device in office or home\n" +
			"- Confirming device location before recovery\n\n" +
			"**Platform Support:**\n" +
			"- **iOS/iPadOS**: Fully supported (supervised devices in lost mode only)\n" +
			"- **Other Platforms**: Not applicable - lost mode is iOS/iPadOS only\n\n" +
			"**Reference:** [Microsoft Graph API - Play Lost Mode Sound](https://learn.microsoft.com/en-us/graph/api/intune-devices-manageddevice-playlostmodesound?view=graph-rest-beta)",
		Attributes: map[string]schema.Attribute{
			"timeouts": commonschema.Timeouts(ctx),
		},
		Blocks: map[string]schema.Block{
			"managed_devices": schema.ListNestedBlock{
				MarkdownDescription: "List of managed devices to play lost mode sound on. These are iOS/iPadOS devices fully managed by Intune only. " +
					"Each entry specifies a device ID and the duration to play the sound.\n\n" +
					"**Note:** At least one of `managed_devices` or `comanaged_devices` must be provided. " +
					"You can provide both to play sounds on different types of devices in one action.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"device_id": schema.StringAttribute{
							Required: true,
							MarkdownDescription: "The unique identifier (GUID) of the managed device to play sound on. " +
								"Example: `\"12345678-1234-1234-1234-123456789abc\"`",
							Validators: []validator.String{
								stringvalidator.RegexMatches(
									regexp.MustCompile(constants.GuidRegex),
									"device_id must be a valid GUID format (e.g., 12345678-1234-1234-1234-123456789abc)",
								),
							},
						},
						"duration_in_minutes": schema.StringAttribute{
							Optional: true,
							MarkdownDescription: "The duration in minutes to play the lost mode sound. " +
								"If not specified, the sound will play for the default duration. " +
								"Example: `\"5\"` for 5 minutes",
						},
					},
				},
			},
			"comanaged_devices": schema.ListNestedBlock{
				MarkdownDescription: "List of co-managed devices to play lost mode sound on. These are iOS/iPadOS devices managed by both Intune and " +
					"Configuration Manager (SCCM). Each entry specifies a device ID and the duration to play the sound.\n\n" +
					"**Note:** At least one of `managed_devices` or `comanaged_devices` must be provided.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"device_id": schema.StringAttribute{
							Required: true,
							MarkdownDescription: "The unique identifier (GUID) of the co-managed device to play sound on. " +
								"Example: `\"12345678-1234-1234-1234-123456789abc\"`",
							Validators: []validator.String{
								stringvalidator.RegexMatches(
									regexp.MustCompile(constants.GuidRegex),
									"device_id must be a valid GUID format (e.g., 12345678-1234-1234-1234-123456789abc)",
								),
							},
						},
						"duration_in_minutes": schema.StringAttribute{
							Optional: true,
							MarkdownDescription: "The duration in minutes to play the lost mode sound. " +
								"Example: `\"5\"`",
						},
					},
				},
			},
		},
	}
}
