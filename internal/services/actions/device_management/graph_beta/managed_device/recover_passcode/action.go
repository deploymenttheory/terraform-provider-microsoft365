package graphBetaRecoverManagedDevicePasscode

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
	ActionName = "microsoft365_graph_beta_device_management_managed_device_recover_passcode"
	InvokeTimeout = 60
)

var (
	_ action.Action                   = &RecoverManagedDevicePasscodeAction{}
	_ action.ActionWithConfigure      = &RecoverManagedDevicePasscodeAction{}
	_ action.ActionWithValidateConfig = &RecoverManagedDevicePasscodeAction{}
)

func NewRecoverManagedDevicePasscodeAction() action.Action {
	return &RecoverManagedDevicePasscodeAction{
		ReadPermissions: []string{
			"DeviceManagementManagedDevices.PrivilegedOperations.All",
		},
		WritePermissions: []string{
			"DeviceManagementManagedDevices.PrivilegedOperations.All",
		},
	}
}

type RecoverManagedDevicePasscodeAction struct {
	client           *msgraphbetasdk.GraphServiceClient
	ReadPermissions  []string
	WritePermissions []string
}

func (a *RecoverManagedDevicePasscodeAction) Metadata(ctx context.Context, req action.MetadataRequest, resp *action.MetadataResponse) {
	resp.TypeName = ActionName
}

func (a *RecoverManagedDevicePasscodeAction) Configure(ctx context.Context, req action.ConfigureRequest, resp *action.ConfigureResponse) {
	a.client = client.SetGraphBetaClientForAction(ctx, req, resp, ActionName)
}

func (a *RecoverManagedDevicePasscodeAction) Schema(ctx context.Context, req action.SchemaRequest, resp *action.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Recovers passcodes for managed devices using the " +
			"`/deviceManagement/managedDevices/{managedDeviceId}/recoverPasscode` endpoint. " +
			"This action retrieves existing passcodes that are escrowed in Intune, which is different from " +
			"reset passcode that generates new temporary passcodes. Recover passcode is primarily used for " +
			"iOS/iPadOS devices where passcodes may be escrowed during enrollment or management.\n\n" +
			"**Important Notes:**\n" +
			"- Retrieves existing escrowed passcode from Intune\n" +
			"- Different from reset passcode (which creates new passcode)\n" +
			"- Passcode must have been previously escrowed\n" +
			"- Primarily for iOS/iPadOS supervised devices\n" +
			"- Retrieved passcode displayed in Intune portal\n" +
			"- May not be available for all device types\n\n" +
			"**Use Cases:**\n" +
			"- User forgot their device passcode (iOS/iPadOS)\n" +
			"- Supervised device lockout recovery\n" +
			"- Administrative access to escrowed passcodes\n" +
			"- Device recovery without factory reset\n" +
			"- Emergency access to locked devices\n" +
			"- Help desk support for locked devices\n\n" +
			"**Platform Support:**\n" +
			"- **iOS/iPadOS**: Supported (supervised devices with passcode escrow)\n" +
			"- **macOS**: Limited (may work with specific configurations)\n" +
			"- **Windows**: Not typically supported for passcode recovery\n" +
			"- **Android**: Not typically supported for passcode recovery\n\n" +
			"**Passcode Escrow:**\n" +
			"- Passcodes must be escrowed during device enrollment\n" +
			"- Not all devices escrow passcodes automatically\n" +
			"- Supervised iOS/iPadOS devices typically escrow passcodes\n" +
			"- Check device enrollment configuration for escrow settings\n" +
			"- Recovery may fail if passcode not escrowed\n\n" +
			"**Recover vs Reset Passcode:**\n" +
			"- **Recover**: Retrieves existing escrowed passcode (no change to device)\n" +
			"- **Reset**: Generates new temporary passcode (device must be unlocked and reset)\n" +
			"- Use recover first if passcode is escrowed\n" +
			"- Use reset if recover fails or passcode not escrowed\n\n" +
			"**Reference:** [Microsoft Graph API - Recover Passcode](https://learn.microsoft.com/en-us/graph/api/intune-devices-manageddevice-recoverpasscode?view=graph-rest-beta)",
		Attributes: map[string]schema.Attribute{
			"device_ids": schema.ListAttribute{
				ElementType: types.StringType,
				Required:    true,
				MarkdownDescription: "List of managed device IDs to recover passcodes for. " +
					"Each ID must be a valid GUID format. Multiple devices can have passcodes recovered in a single action. " +
					"Example: `[\"12345678-1234-1234-1234-123456789abc\", \"87654321-4321-4321-4321-ba9876543210\"]`\n\n" +
					"**Important:** This action retrieves existing escrowed passcodes. If a passcode was not escrowed " +
					"during device enrollment, the recovery will fail. Check device properties in Intune to verify " +
					"passcode escrow status before attempting recovery.",
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
