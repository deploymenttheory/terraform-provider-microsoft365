package graphBetaLocateManagedDevice

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
	ActionName    = "microsoft365_graph_beta_device_management_managed_device_locate_device"
	InvokeTimeout = 60
)

var (
	_ action.Action                   = &LocateManagedDeviceAction{}
	_ action.ActionWithConfigure      = &LocateManagedDeviceAction{}
	_ action.ActionWithValidateConfig = &LocateManagedDeviceAction{}
)

func NewLocateManagedDeviceAction() action.Action {
	return &LocateManagedDeviceAction{
		ReadPermissions: []string{
			"DeviceManagementManagedDevices.PrivilegedOperations.All",
		},
		WritePermissions: []string{
			"DeviceManagementManagedDevices.PrivilegedOperations.All",
		},
	}
}

type LocateManagedDeviceAction struct {
	client           *msgraphbetasdk.GraphServiceClient
	ReadPermissions  []string
	WritePermissions []string
}

func (a *LocateManagedDeviceAction) Metadata(ctx context.Context, req action.MetadataRequest, resp *action.MetadataResponse) {
	resp.TypeName = ActionName
}

func (a *LocateManagedDeviceAction) Configure(ctx context.Context, req action.ConfigureRequest, resp *action.ConfigureResponse) {
	a.client = client.SetGraphBetaClientForAction(ctx, req, resp, ActionName)
}

func (a *LocateManagedDeviceAction) Schema(ctx context.Context, req action.SchemaRequest, resp *action.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Triggers device location for one or more managed devices using the " +
			"`/deviceManagement/managedDevices/{managedDeviceId}/locateDevice` endpoint. " +
			"This action requests the device to report its current geographic location, which is then viewable " +
			"in the Microsoft Intune admin center. The locate device feature is essential for finding lost or " +
			"stolen devices and is commonly used in conjunction with lost mode.\n\n" +
			"**Important Notes:**\n" +
			"- Device must be online to receive and respond to the locate command\n" +
			"- Location services must be enabled on the device\n" +
			"- Device must have GPS/location hardware capability\n" +
			"- Location data is displayed in the Intune portal, not returned via API\n" +
			"- Multiple location requests can be sent over time to track device movement\n" +
			"- Location accuracy depends on device capabilities (GPS, WiFi, cellular triangulation)\n\n" +
			"**Use Cases:**\n" +
			"- Locating lost or stolen devices\n" +
			"- Tracking devices in lost mode\n" +
			"- Verifying device location for security/compliance\n" +
			"- Finding devices before performing remote wipe\n" +
			"- Assisting users who have misplaced their devices\n" +
			"- Asset tracking and recovery operations\n\n" +
			"**Platform Support:**\n" +
			"- **iOS/iPadOS**: Fully supported (iOS 9.3+, supervised devices)\n" +
			"- **Android Enterprise**: Supported (fully managed, dedicated, and work profile devices)\n" +
			"- **macOS**: Supported (macOS 10.13+, supervised devices or user-approved MDM)\n" +
			"- **Windows**: Limited support (Windows 10/11 with location services enabled)\n\n" +
			"**Location Data:**\n" +
			"- Location data is displayed in the Intune admin center under device properties\n" +
			"- Shows latitude, longitude, altitude (if available)\n" +
			"- Includes location accuracy radius\n" +
			"- Displays timestamp of when location was captured\n" +
			"- Location history may be available depending on platform\n\n" +
			"**Privacy Considerations:**\n" +
			"- Users may receive notification that device location was requested\n" +
			"- Location tracking should comply with organizational privacy policies\n" +
			"- Document legitimate business reasons for location requests\n" +
			"- Consider legal requirements in your jurisdiction\n\n" +
			"**Reference:** [Microsoft Graph API - Locate Device](https://learn.microsoft.com/en-us/graph/api/intune-devices-manageddevice-locatedevice?view=graph-rest-beta)",
		Attributes: map[string]schema.Attribute{
			"device_ids": schema.ListAttribute{
				ElementType: types.StringType,
				Required:    true,
				MarkdownDescription: "List of managed device IDs to locate. " +
					"Each ID must be a valid GUID format. Multiple devices can be located in a single action. " +
					"Example: `[\"12345678-1234-1234-1234-123456789abc\", \"87654321-4321-4321-4321-ba9876543210\"]`\n\n" +
					"**Important:** Devices must be online and have location services enabled to respond to the locate request. " +
					"Location data will be displayed in the Microsoft Intune admin center once the device reports its position.",
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
				MarkdownDescription: "Whether to validate that devices exist and support location services before attempting to locate them. " +
					"Disabling this can speed up planning but may result in runtime errors for non-existent or unsupported devices. " +
					"Default: `true`.",
			},
			"timeouts": commonschema.ActionTimeouts(ctx),
		},
	}
}
