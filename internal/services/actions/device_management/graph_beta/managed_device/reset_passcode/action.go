package graphBetaResetManagedDevicePasscode

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
	ActionName = "graph_beta_device_management_managed_device_reset_passcode"
)

var (
	_ action.Action                   = &ResetManagedDevicePasscodeAction{}
	_ action.ActionWithConfigure      = &ResetManagedDevicePasscodeAction{}
	_ action.ActionWithValidateConfig = &ResetManagedDevicePasscodeAction{}
)

func NewResetManagedDevicePasscodeAction() action.Action {
	return &ResetManagedDevicePasscodeAction{
		ReadPermissions: []string{
			"DeviceManagementManagedDevices.PrivilegedOperations.All",
		},
		WritePermissions: []string{
			"DeviceManagementManagedDevices.PrivilegedOperations.All",
		},
	}
}

type ResetManagedDevicePasscodeAction struct {
	client           *msgraphbetasdk.GraphServiceClient
	ProviderTypeName string
	TypeName         string
	ReadPermissions  []string
	WritePermissions []string
}

func (a *ResetManagedDevicePasscodeAction) Metadata(ctx context.Context, req action.MetadataRequest, resp *action.MetadataResponse) {
	a.ProviderTypeName = req.ProviderTypeName
	a.TypeName = ActionName
	resp.TypeName = a.FullTypeName()
}

func (a *ResetManagedDevicePasscodeAction) FullTypeName() string {
	return a.ProviderTypeName + "_" + ActionName
}

func (a *ResetManagedDevicePasscodeAction) Configure(ctx context.Context, req action.ConfigureRequest, resp *action.ConfigureResponse) {
	a.client = client.SetGraphBetaClientForAction(ctx, req, resp, constants.PROVIDER_NAME+"_"+ActionName)
}

func (a *ResetManagedDevicePasscodeAction) Schema(ctx context.Context, req action.SchemaRequest, resp *action.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Resets the passcode on managed devices using the " +
			"`/deviceManagement/managedDevices/{managedDeviceId}/resetPasscode` endpoint. " +
			"This action removes the current device passcode/password and generates a new temporary passcode. " +
			"The new passcode is displayed to the administrator and must be communicated to the device user. " +
			"This action supports resetting passcodes on multiple devices in a single operation.\n\n" +
			"**Important Notes:**\n" +
			"- The device must be online and able to receive the command\n" +
			"- On iOS/iPadOS devices, the device must be supervised\n" +
			"- On Android devices, this removes the passcode requirement temporarily\n" +
			"- On Windows devices, the functionality varies by Windows version\n" +
			"- The new passcode is a temporary system-generated code that should be changed by the user\n" +
			"- This action requires the device to be enrolled and actively managed by Intune\n\n" +
			"**Use Cases:**\n" +
			"- User forgot device passcode and cannot unlock device\n" +
			"- Device locked after too many failed passcode attempts\n" +
			"- Administrative access needed for troubleshooting\n" +
			"- Security incident requiring immediate access to device\n\n" +
			"**Reference:** [Microsoft Graph API - Reset Passcode](https://learn.microsoft.com/en-us/graph/api/intune-devices-manageddevice-resetpasscode?view=graph-rest-beta)",
		Attributes: map[string]schema.Attribute{
			"device_ids": schema.ListAttribute{
				ElementType: types.StringType,
				Required:    true,
				MarkdownDescription: "List of managed device IDs to reset passcodes for. " +
					"Each ID must be a valid GUID format. Multiple device passcodes can be reset in a single action. " +
					"Example: `[\"12345678-1234-1234-1234-123456789abc\", \"87654321-4321-4321-4321-ba9876543210\"]`\n\n" +
					"**Important:** The new temporary passcode for each device will be displayed in Intune and must be " +
					"communicated to the device user. Users should change this temporary passcode immediately after regaining access.",
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
			"timeouts": commonschema.Timeouts(ctx),
		},
	}
}
