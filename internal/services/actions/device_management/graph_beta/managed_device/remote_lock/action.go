package graphBetaRemoteLockManagedDevice

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
	ActionName = "microsoft365_graph_beta_device_management_managed_device_remote_lock"
	InvokeTimeout = 60
)

var (
	_ action.Action                   = &RemoteLockManagedDeviceAction{}
	_ action.ActionWithConfigure      = &RemoteLockManagedDeviceAction{}
	_ action.ActionWithValidateConfig = &RemoteLockManagedDeviceAction{}
)

func NewRemoteLockManagedDeviceAction() action.Action {
	return &RemoteLockManagedDeviceAction{
		ReadPermissions: []string{
			"DeviceManagementManagedDevices.PrivilegedOperations.All",
		},
		WritePermissions: []string{
			"DeviceManagementManagedDevices.PrivilegedOperations.All",
		},
	}
}

type RemoteLockManagedDeviceAction struct {
	client           *msgraphbetasdk.GraphServiceClient
	ReadPermissions  []string
	WritePermissions []string
}

func (a *RemoteLockManagedDeviceAction) Metadata(ctx context.Context, req action.MetadataRequest, resp *action.MetadataResponse) {
	resp.TypeName = ActionName
}

func (a *RemoteLockManagedDeviceAction) Configure(ctx context.Context, req action.ConfigureRequest, resp *action.ConfigureResponse) {
	a.client = client.SetGraphBetaClientForAction(ctx, req, resp, ActionName)
}

func (a *RemoteLockManagedDeviceAction) Schema(ctx context.Context, req action.SchemaRequest, resp *action.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Remotely locks managed devices using the " +
			"`/deviceManagement/managedDevices/{managedDeviceId}/remoteLock` endpoint. " +
			"This action immediately locks the device screen, requiring the user to enter their passcode to unlock it. " +
			"This is useful for securing lost or stolen devices, or for security compliance scenarios. " +
			"This action supports remotely locking multiple devices in a single operation.\n\n" +
			"**Important Notes:**\n" +
			"- The device must be online and able to receive the command\n" +
			"- The device will lock immediately when it receives the command\n" +
			"- The user's existing passcode remains unchanged\n" +
			"- The user will need to enter their passcode to unlock the device\n" +
			"- For lost/stolen devices, consider using remote lock before more drastic measures\n" +
			"- This action does not remove any data from the device\n\n" +
			"**Use Cases:**\n" +
			"- Lost or stolen device - immediate security action\n" +
			"- Security incident - prevent unauthorized access\n" +
			"- Compliance enforcement - ensure device is secured\n" +
			"- Unattended device in public location\n" +
			"- User reported potential device compromise\n\n" +
			"**Platform Support:**\n" +
			"- **iOS/iPadOS**: Fully supported (iOS 9.0+)\n" +
			"- **Android**: Supported on Android Enterprise devices (work profile and fully managed)\n" +
			"- **Windows**: Supported on Windows 10/11 devices\n" +
			"- **macOS**: Supported on managed Mac computers\n\n" +
			"**Reference:** [Microsoft Graph API - Remote Lock](https://learn.microsoft.com/en-us/graph/api/intune-devices-manageddevice-remotelock?view=graph-rest-beta)",
		Attributes: map[string]schema.Attribute{
			"device_ids": schema.ListAttribute{
				ElementType: types.StringType,
				Required:    true,
				MarkdownDescription: "List of managed device IDs to remotely lock. " +
					"Each ID must be a valid GUID format. Multiple devices can be locked in a single action. " +
					"Example: `[\"12345678-1234-1234-1234-123456789abc\", \"87654321-4321-4321-4321-ba9876543210\"]`\n\n" +
					"**Important:** Devices will lock immediately when they receive the command. " +
					"Ensure you have authorization to lock these devices. For lost or stolen devices, " +
					"this provides an immediate security measure while you determine next steps (locate, wipe, etc.).",
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
			"timeouts": commonschema.ActionTimeouts(ctx),
		},
	}
}
