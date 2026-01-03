package graphBetaActivateDeviceEsimManagedDevice

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
	ActionName    = "microsoft365_graph_beta_device_management_managed_device_activate_device_esim"
	InvokeTimeout = 60
)

var (
	_ action.ActionWithValidateConfig = (*ActivateDeviceEsimManagedDeviceAction)(nil)
	_ action.ActionWithConfigure      = (*ActivateDeviceEsimManagedDeviceAction)(nil)
)

func NewActivateDeviceEsimManagedDeviceAction() action.Action {
	return &ActivateDeviceEsimManagedDeviceAction{
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

type ActivateDeviceEsimManagedDeviceAction struct {
	client           *msgraphbetasdk.GraphServiceClient
	ReadPermissions  []string
	WritePermissions []string
}

func (a *ActivateDeviceEsimManagedDeviceAction) Metadata(ctx context.Context, req action.MetadataRequest, resp *action.MetadataResponse) {
	resp.TypeName = ActionName
}

func (a *ActivateDeviceEsimManagedDeviceAction) Configure(ctx context.Context, req action.ConfigureRequest, resp *action.ConfigureResponse) {
	a.client = client.SetGraphBetaClientForAction(ctx, req, resp, ActionName)
}

func (a *ActivateDeviceEsimManagedDeviceAction) Schema(ctx context.Context, req action.SchemaRequest, resp *action.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Activates eSIM on managed cellular devices using the " +
			"`/deviceManagement/managedDevices/{managedDeviceId}/activateDeviceEsim` and " +
			"`/deviceManagement/comanagedDevices/{managedDeviceId}/activateDeviceEsim` endpoints. " +
			"This action enables eSIM functionality on compatible devices by providing a carrier activation URL. " +
			"eSIM (embedded SIM) technology allows devices to connect to cellular networks without a physical SIM card, " +
			"providing greater flexibility for device deployment and carrier management. This action supports activating " +
			"eSIM on multiple devices in a single operation with per-device carrier URL configuration.\n\n" +
			"**Important Notes:**\n" +
			"- Only applicable to devices with eSIM hardware capability\n" +
			"- Requires carrier-specific activation URL\n" +
			"- Device must support eSIM technology\n" +
			"- Carrier must support eSIM activation\n" +
			"- Device must be online to receive activation\n" +
			"- Each device requires its own carrier activation URL\n\n" +
			"**Use Cases:**\n" +
			"- Initial eSIM activation on new devices\n" +
			"- Switching carriers on eSIM-capable devices\n" +
			"- Bulk eSIM deployment for corporate devices\n" +
			"- Remote eSIM provisioning for field devices\n" +
			"- International device deployment with local carriers\n\n" +
			"**Platform Support:**\n" +
			"- **iOS/iPadOS**: Supported on eSIM-capable devices (iPhone XS and later, cellular iPads)\n" +
			"- **Windows**: Supported on eSIM-capable Windows devices with cellular modems\n" +
			"- **Android**: Support varies by device manufacturer and Android version\n" +
			"- **Other Platforms**: Not applicable\n\n" +
			"**Reference:** [Microsoft Graph API - Activate Device eSIM](https://learn.microsoft.com/en-us/graph/api/intune-devices-manageddevice-activatedeviceesim?view=graph-rest-beta)",
		Attributes: map[string]schema.Attribute{
			"managed_devices": schema.ListNestedAttribute{
				Optional: true,
				MarkdownDescription: "List of managed devices to activate eSIM on. These are devices fully managed by Intune only. " +
					"Each entry specifies a device ID and the carrier-specific activation URL.\n\n" +
					"**Examples:**\n" +
					"```hcl\n" +
					"managed_devices = [\n" +
					"  {\n" +
					"    device_id   = \"12345678-1234-1234-1234-123456789abc\"\n" +
					"    carrier_url = \"https://carrier.example.com/esim/activate?token=abc123\"\n" +
					"  },\n" +
					"  {\n" +
					"    device_id   = \"87654321-4321-4321-4321-987654321cba\"\n" +
					"    carrier_url = \"https://carrier.example.com/esim/activate?token=def456\"\n" +
					"  }\n" +
					"]\n" +
					"```\n\n" +
					"**Platform Support:** iOS (iPhone XS+), Windows 10/11 with cellular, Android (varies by manufacturer)\n\n" +
					"**Note:** At least one of `managed_devices` or `comanaged_devices` must be provided. " +
					"Device must be online and support eSIM technology.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"device_id": schema.StringAttribute{
							Required: true,
							MarkdownDescription: "The unique identifier (GUID) of the managed device to activate eSIM on. " +
								"Device must have eSIM hardware capability. " +
								"Example: `\"12345678-1234-1234-1234-123456789abc\"`",
							Validators: []validator.String{
								stringvalidator.RegexMatches(
									regexp.MustCompile(constants.GuidRegex),
									"device_id must be a valid GUID format (e.g., 12345678-1234-1234-1234-123456789abc)",
								),
							},
						},
						"carrier_url": schema.StringAttribute{
							Required: true,
							MarkdownDescription: "The carrier-specific activation URL for this device's eSIM. " +
								"This URL is provided by the mobile carrier and contains the activation profile. " +
								"Format varies by carrier. " +
								"Example: `\"https://carrier.example.com/esim/activate?token=abc123\"`",
							Validators: []validator.String{
								stringvalidator.LengthAtLeast(1),
							},
						},
					},
				},
			},
			"comanaged_devices": schema.ListNestedAttribute{
				Optional: true,
				MarkdownDescription: "List of co-managed devices to activate eSIM on. These are devices managed by both Intune and " +
					"Configuration Manager (SCCM). Each entry specifies a device ID and the carrier activation URL.\n\n" +
					"**Examples:**\n" +
					"```hcl\n" +
					"comanaged_devices = [\n" +
					"  {\n" +
					"    device_id   = \"aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee\"\n" +
					"    carrier_url = \"https://carrier.example.com/esim/activate?code=xyz789\"\n" +
					"  }\n" +
					"]\n" +
					"```\n\n" +
					"**Platform Support:** Windows 10/11 with cellular modems (primary), limited iOS/Android support\n\n" +
					"**Note:** At least one of `managed_devices` or `comanaged_devices` must be provided. " +
					"Device must be online and support eSIM technology.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"device_id": schema.StringAttribute{
							Required: true,
							MarkdownDescription: "The unique identifier (GUID) of the co-managed device to activate eSIM on. " +
								"Example: `\"12345678-1234-1234-1234-123456789abc\"`",
							Validators: []validator.String{
								stringvalidator.RegexMatches(
									regexp.MustCompile(constants.GuidRegex),
									"device_id must be a valid GUID format (e.g., 12345678-1234-1234-1234-123456789abc)",
								),
							},
						},
						"carrier_url": schema.StringAttribute{
							Required: true,
							MarkdownDescription: "The carrier activation URL for this co-managed device. " +
								"Example: `\"https://carrier.example.com/esim/activate?code=xyz789\"`",
							Validators: []validator.String{
								stringvalidator.LengthAtLeast(1),
							},
						},
					},
				},
			},
			"ignore_partial_failures": schema.BoolAttribute{
				Optional: true,
				MarkdownDescription: "If set to `true`, the action will succeed even if some devices fail eSIM activation. " +
					"Failed devices will be reported as warnings instead of errors. " +
					"Default: `false` (action fails if any device fails).",
			},
			"validate_device_exists": schema.BoolAttribute{
				Optional: true,
				MarkdownDescription: "Whether to validate that devices exist before attempting activation. " +
					"Disabling this can speed up planning but may result in runtime errors for non-existent devices. " +
					"Default: `true`.",
			},
			"timeouts": commonschema.ActionTimeouts(ctx),
		},
	}
}
