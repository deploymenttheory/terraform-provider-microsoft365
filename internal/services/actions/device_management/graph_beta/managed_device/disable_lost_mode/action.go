package graphBetaDisableLostModeManagedDevice

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
	ActionName = "microsoft365_graph_beta_device_management_managed_device_disable_lost_mode"
	InvokeTimeout = 60
)

var (
	_ action.Action                   = &DisableLostModeManagedDeviceAction{}
	_ action.ActionWithConfigure      = &DisableLostModeManagedDeviceAction{}
	_ action.ActionWithValidateConfig = &DisableLostModeManagedDeviceAction{}
)

func NewDisableLostModeManagedDeviceAction() action.Action {
	return &DisableLostModeManagedDeviceAction{
		ReadPermissions: []string{
			"DeviceManagementManagedDevices.PrivilegedOperations.All",
		},
		WritePermissions: []string{
			"DeviceManagementManagedDevices.PrivilegedOperations.All",
		},
	}
}

type DisableLostModeManagedDeviceAction struct {
	client           *msgraphbetasdk.GraphServiceClient
	ReadPermissions  []string
	WritePermissions []string
}

func (a *DisableLostModeManagedDeviceAction) Metadata(ctx context.Context, req action.MetadataRequest, resp *action.MetadataResponse) {
	resp.TypeName = ActionName
}

func (a *DisableLostModeManagedDeviceAction) Configure(ctx context.Context, req action.ConfigureRequest, resp *action.ConfigureResponse) {
	a.client = client.SetGraphBetaClientForAction(ctx, req, resp, ActionName)
}

func (a *DisableLostModeManagedDeviceAction) Schema(ctx context.Context, req action.SchemaRequest, resp *action.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Disables lost mode on iOS/iPadOS managed devices using the " +
			"`/deviceManagement/managedDevices/{managedDeviceId}/disableLostMode` and " +
			"`/deviceManagement/comanagedDevices/{managedDeviceId}/disableLostMode` endpoints. " +
			"This action removes the device from lost mode, allowing normal device operation to resume. " +
			"Lost mode is a feature that helps locate and secure lost iOS/iPadOS devices by locking them " +
			"and displaying a custom message with contact information on the lock screen. " +
			"This action supports disabling lost mode on multiple devices in a single operation.\n\n" +
			"**Important Notes:**\n" +
			"- Only applicable to iOS and iPadOS devices (iOS 9.3+)\n" +
			"- Device must currently be in lost mode\n" +
			"- Device must be supervised\n" +
			"- Requires device to be online to receive command\n" +
			"- Once disabled, device returns to normal operation\n" +
			"- The custom lock screen message is removed\n\n" +
			"**Use Cases:**\n" +
			"- Device has been recovered and needs to be returned to service\n" +
			"- Lost mode was enabled in error\n" +
			"- Device location has been confirmed and no longer needs tracking\n" +
			"- User has regained possession of their device\n\n" +
			"**Platform Support:**\n" +
			"- **iOS/iPadOS**: Fully supported (iOS 9.3+, supervised devices only)\n" +
			"- **Other Platforms**: Not applicable - lost mode is iOS/iPadOS only\n\n" +
			"**Reference:** [Microsoft Graph API - Disable Lost Mode](https://learn.microsoft.com/en-us/graph/api/intune-devices-manageddevice-disablelostmode?view=graph-rest-beta)",
		Attributes: map[string]schema.Attribute{
			"managed_device_ids": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
				MarkdownDescription: "List of managed device IDs to disable lost mode for. These are iOS/iPadOS devices fully managed by Intune only. " +
					"Each ID must be a valid GUID format. Example: `[\"12345678-1234-1234-1234-123456789abc\", \"87654321-4321-4321-4321-ba9876543210\"]`\n\n" +
					"**Note:** At least one of `managed_device_ids` or `comanaged_device_ids` must be provided. " +
					"You can provide both to disable lost mode on different types of devices in one action.",
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
				MarkdownDescription: "List of co-managed device IDs to disable lost mode for. These are iOS/iPadOS devices managed by both Intune and " +
					"Configuration Manager (SCCM). Each ID must be a valid GUID format. " +
					"Example: `[\"12345678-1234-1234-1234-123456789abc\"]`\n\n" +
					"**Note:** At least one of `managed_device_ids` or `comanaged_device_ids` must be provided.",
				Validators: []validator.List{
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
