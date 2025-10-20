package graphBetaUpdateWindowsDeviceAccount

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

	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

const (
	ActionName = "graph_beta_device_management_managed_device_update_windows_device_account"
)

var (
	_ action.Action                   = &UpdateWindowsDeviceAccountAction{}
	_ action.ActionWithConfigure      = &UpdateWindowsDeviceAccountAction{}
	_ action.ActionWithValidateConfig = &UpdateWindowsDeviceAccountAction{}
)

func NewUpdateWindowsDeviceAccountAction() action.Action {
	return &UpdateWindowsDeviceAccountAction{
		ReadPermissions: []string{
			"DeviceManagementManagedDevices.PrivilegedOperations.All",
		},
		WritePermissions: []string{
			"DeviceManagementManagedDevices.PrivilegedOperations.All",
		},
	}
}

type UpdateWindowsDeviceAccountAction struct {
	client           *msgraphbetasdk.GraphServiceClient
	ProviderTypeName string
	TypeName         string
	ReadPermissions  []string
	WritePermissions []string
}

func (a *UpdateWindowsDeviceAccountAction) Metadata(ctx context.Context, req action.MetadataRequest, resp *action.MetadataResponse) {
	a.ProviderTypeName = req.ProviderTypeName
	a.TypeName = ActionName
	resp.TypeName = a.FullTypeName()
}

func (a *UpdateWindowsDeviceAccountAction) FullTypeName() string {
	return a.ProviderTypeName + "_" + ActionName
}

func (a *UpdateWindowsDeviceAccountAction) Configure(ctx context.Context, req action.ConfigureRequest, resp *action.ConfigureResponse) {
	a.client = client.SetGraphBetaClientForAction(ctx, req, resp, constants.PROVIDER_NAME+"_"+ActionName)
}

func (a *UpdateWindowsDeviceAccountAction) Schema(ctx context.Context, req action.SchemaRequest, resp *action.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Updates the device account configuration on Windows devices using the " +
			"`/deviceManagement/managedDevices/{managedDeviceId}/updateWindowsDeviceAccount` and " +
			"`/deviceManagement/comanagedDevices/{managedDeviceId}/updateWindowsDeviceAccount` endpoints. " +
			"This action is specifically designed for shared Windows devices like Surface Hub and Microsoft Teams Rooms " +
			"that require device account configuration for Exchange and Skype for Business/Teams integration.\n\n" +
			"**What This Action Does:**\n" +
			"- Updates device account credentials\n" +
			"- Configures Exchange server settings\n" +
			"- Sets up calendar sync\n" +
			"- Configures Teams/SfB settings\n" +
			"- Manages password rotation\n" +
			"- Updates SIP address configuration\n\n" +
			"**Target Devices:**\n" +
			"- **Surface Hub**: Collaboration devices\n" +
			"- **Microsoft Teams Rooms**: Meeting room systems\n" +
			"- **Shared Windows devices**: Kiosk/common area devices\n\n" +
			"**Platform Support:**\n" +
			"- **Windows 10/11**: Surface Hub, Teams Rooms\n" +
			"- **Windows 10 IoT**: Teams Rooms appliances\n\n" +
			"**Common Use Cases:**\n" +
			"- Update device account password\n" +
			"- Reconfigure Exchange server\n" +
			"- Update calendar sync settings\n" +
			"- Change Teams/SfB configuration\n" +
			"- Rotate device credentials\n" +
			"- Fix authentication issues\n\n" +
			"**Important Considerations:**\n" +
			"- Requires device reboot to apply\n" +
			"- Password stored securely\n" +
			"- Exchange connectivity required\n" +
			"- Teams/SfB license needed\n" +
			"- Affects device functionality\n\n" +
			"**Reference:** [Microsoft Graph API - Update Windows Device Account](https://learn.microsoft.com/en-us/graph/api/intune-devices-manageddevice-updatewindowsdeviceaccount?view=graph-rest-beta)",
		Attributes: map[string]schema.Attribute{
			"managed_devices": schema.ListNestedAttribute{
				Optional: true,
				MarkdownDescription: "List of managed Windows devices to update with individual device account configurations. " +
					"Each entry specifies a device ID and its complete device account settings including credentials, " +
					"Exchange server, and synchronization options. These are devices fully managed by Intune only.\n\n" +
					"**Note:** At least one of `managed_devices` or `comanaged_devices` must be provided.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"device_id": schema.StringAttribute{
							Required: true,
							MarkdownDescription: "The managed device ID (GUID) of the Windows device to update. " +
								"Example: `12345678-1234-1234-1234-123456789abc`",
							Validators: []validator.String{
								stringvalidator.RegexMatches(
									regexp.MustCompile(constants.GuidRegex),
									"device_id must be a valid GUID format (e.g., 12345678-1234-1234-1234-123456789abc)",
								),
							},
						},
						"device_account_email": schema.StringAttribute{
							Required: true,
							MarkdownDescription: "The email address of the device account (resource mailbox). " +
								"This is typically a room mailbox in Exchange for Teams Rooms or Surface Hub. " +
								"Example: `conference-room-01@company.com` or `surfacehub-lobby@company.com`\n\n" +
								"**Requirements:**\n" +
								"- Must be a valid email address\n" +
								"- Must exist in Exchange/Microsoft 365\n" +
								"- Should be a room or equipment mailbox\n" +
								"- Requires appropriate licenses",
						},
						"password": schema.StringAttribute{
							Required: true,
							MarkdownDescription: "The password for the device account. This password is used to authenticate " +
								"the device with Exchange and Teams/Skype for Business services. " +
								"The password is transmitted securely and stored encrypted.\n\n" +
								"**Best Practices:**\n" +
								"- Use a strong, complex password\n" +
								"- Consider enabling password rotation\n" +
								"- Store securely (use Terraform sensitive values)\n" +
								"- Rotate regularly for security\n" +
								"- Follow organizational password policies",
						},
						"password_rotation_enabled": schema.BoolAttribute{
							Required: true,
							MarkdownDescription: "Whether automatic password rotation is enabled for the device account. " +
								"When enabled, the device will automatically rotate its password periodically.\n\n" +
								"- **`true`**: Enable automatic password rotation (recommended for security)\n" +
								"- **`false`**: Disable automatic password rotation (manual management required)\n\n" +
								"**Note:** When enabled, ensure the device account has appropriate permissions " +
								"in Active Directory to change its own password.",
						},
						"calendar_sync_enabled": schema.BoolAttribute{
							Required: true,
							MarkdownDescription: "Whether calendar synchronization is enabled for the device. " +
								"This determines if the device will sync its calendar from Exchange.\n\n" +
								"- **`true`**: Enable calendar sync (shows meetings, availability)\n" +
								"- **`false`**: Disable calendar sync (no meeting information displayed)\n\n" +
								"**Use Cases:**\n" +
								"- Teams Rooms: Typically enabled (display meeting schedule)\n" +
								"- Surface Hub: Typically enabled (meeting coordination)\n" +
								"- Kiosk devices: May be disabled (no calendar needed)",
						},
						"exchange_server": schema.StringAttribute{
							Optional: true,
							MarkdownDescription: "The Exchange server address for mailbox connectivity. " +
								"This can be an on-premises Exchange server or Exchange Online (Microsoft 365).\n\n" +
								"**Examples:**\n" +
								"- Exchange Online: `outlook.office365.com`\n" +
								"- On-premises: `mail.company.com` or `exchange.company.local`\n" +
								"- Autodiscover: Leave blank to use autodiscover\n\n" +
								"**Note:** If not specified, the device will attempt to use Exchange autodiscover " +
								"to locate the appropriate Exchange server automatically.",
						},
						"session_initiation_protocol_address": schema.StringAttribute{
							Optional: true,
							MarkdownDescription: "The Session Initiation Protocol (SIP) address for Teams/Skype for Business connectivity. " +
								"This is the SIP URI for the device account, typically matching the email address but with 'sip:' prefix.\n\n" +
								"**Format:** `sip:username@domain.com`\n\n" +
								"**Examples:**\n" +
								"- `sip:conference-room-01@company.com`\n" +
								"- `sip:surfacehub-lobby@company.com`\n\n" +
								"**Requirements:**\n" +
								"- Required for Teams/Skype for Business functionality\n" +
								"- Must match the device account UPN or email\n" +
								"- Account must be enabled for Teams/SfB\n" +
								"- Requires appropriate licensing",
						},
					},
				},
				Validators: []validator.List{
					listvalidator.SizeAtLeast(1),
				},
			},
			"comanaged_devices": schema.ListNestedAttribute{
				Optional: true,
				MarkdownDescription: "List of co-managed Windows devices to update with individual device account configurations. " +
					"These are devices managed by both Intune and Configuration Manager (SCCM). " +
					"Configuration is identical to managed_devices.\n\n" +
					"**Note:** At least one of `managed_devices` or `comanaged_devices` must be provided.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"device_id": schema.StringAttribute{
							Required: true,
							MarkdownDescription: "The co-managed device ID (GUID). " +
								"Example: `12345678-1234-1234-1234-123456789abc`",
							Validators: []validator.String{
								stringvalidator.RegexMatches(
									regexp.MustCompile(constants.GuidRegex),
									"device_id must be a valid GUID format",
								),
							},
						},
						"device_account_email": schema.StringAttribute{
							Required:            true,
							MarkdownDescription: "The email address of the device account. See managed_devices.device_account_email for details.",
						},
						"password": schema.StringAttribute{
							Required:            true,
							MarkdownDescription: "The password for the device account. See managed_devices.password for details.",
						},
						"password_rotation_enabled": schema.BoolAttribute{
							Required:            true,
							MarkdownDescription: "Whether automatic password rotation is enabled. See managed_devices.password_rotation_enabled for details.",
						},
						"calendar_sync_enabled": schema.BoolAttribute{
							Required:            true,
							MarkdownDescription: "Whether calendar synchronization is enabled. See managed_devices.calendar_sync_enabled for details.",
						},
						"exchange_server": schema.StringAttribute{
							Optional:            true,
							MarkdownDescription: "The Exchange server address. See managed_devices.exchange_server for details.",
						},
						"session_initiation_protocol_address": schema.StringAttribute{
							Optional:            true,
							MarkdownDescription: "The SIP address for Teams/SfB. See managed_devices.session_initiation_protocol_address for details.",
						},
					},
				},
				Validators: []validator.List{
					listvalidator.SizeAtLeast(1),
				},
			},
			"timeouts": commonschema.Timeouts(ctx),
		},
	}
}
