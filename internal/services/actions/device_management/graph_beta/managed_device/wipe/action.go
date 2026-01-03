package graphBetaWipeManagedDevice

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
	ActionName = "microsoft365_graph_beta_device_management_managed_device_wipe"
	InvokeTimeout = 60
)

var (
	_ action.Action                   = &WipeManagedDeviceAction{}
	_ action.ActionWithConfigure      = &WipeManagedDeviceAction{}
	_ action.ActionWithValidateConfig = &WipeManagedDeviceAction{}
)

func NewWipeManagedDeviceAction() action.Action {
	return &WipeManagedDeviceAction{
		ReadPermissions: []string{
			"DeviceManagementConfiguration.ReadWrite.All",
			"DeviceManagementManagedDevices.ReadWrite.All",
			"DeviceManagementManagedDevices.PrivilegedOperations.All",
		},
		WritePermissions: []string{
			"DeviceManagementConfiguration.ReadWrite.All",
			"DeviceManagementManagedDevices.ReadWrite.All",
			"DeviceManagementManagedDevices.PrivilegedOperations.All",
		},
	}
}

type WipeManagedDeviceAction struct {
	client           *msgraphbetasdk.GraphServiceClient
	ReadPermissions  []string
	WritePermissions []string
}

func (a *WipeManagedDeviceAction) Metadata(ctx context.Context, req action.MetadataRequest, resp *action.MetadataResponse) {
	resp.TypeName = ActionName
}

func (a *WipeManagedDeviceAction) Configure(ctx context.Context, req action.ConfigureRequest, resp *action.ConfigureResponse) {
	a.client = client.SetGraphBetaClientForAction(ctx, req, resp, ActionName)
}

func (a *WipeManagedDeviceAction) Schema(ctx context.Context, req action.SchemaRequest, resp *action.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Wipes managed devices from Microsoft Intune using the " +
			"`/deviceManagement/managedDevices/{managedDeviceId}/wipe` endpoint. " +
			"This action performs a factory reset, removing all data (company and personal) from the device. " +
			"The device is returned to its out-of-box state and removed from Intune management. " +
			"This action supports wiping multiple devices in a single operation.\n\n" +
			"**Important Notes:**\n" +
			"- This action removes **ALL** data from the device unless `keep_user_data` is set to `true`\n" +
			"- For iOS/iPadOS devices, Activation Lock must be disabled or unlock code provided\n" +
			"- For Windows devices, you can use protected wipe to maintain UEFI-embedded licenses\n" +
			"- For Android devices, factory reset protection must be disabled\n" +
			"- This action cannot be reversed - all data will be permanently deleted\n\n" +
			"**Reference:** [Microsoft Graph API - Wipe Managed Device](https://learn.microsoft.com/en-us/graph/api/intune-devices-manageddevice-wipe?view=graph-rest-beta)",
		Attributes: map[string]schema.Attribute{
			"device_ids": schema.ListAttribute{
				ElementType: types.StringType,
				Required:    true,
				MarkdownDescription: "List of managed device IDs to wipe. " +
					"Each ID must be a valid GUID format. Multiple devices can be wiped in a single action. " +
					"Example: `[\"12345678-1234-1234-1234-123456789abc\", \"87654321-4321-4321-4321-ba9876543210\"]`",
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
			"keep_enrollment_data": schema.BoolAttribute{
				Optional: true,
				MarkdownDescription: "If `true`, maintains enrollment state data during wipe. " +
					"This allows the device to automatically re-enroll after being wiped. " +
					"Defaults to `false`.",
			},
			"keep_user_data": schema.BoolAttribute{
				Optional: true,
				MarkdownDescription: "If `true`, preserves user data during the wipe operation. " +
					"Only company data and managed apps are removed. " +
					"**Note:** Not supported on all device types. " +
					"Defaults to `false`.",
			},
			"macos_unlock_code": schema.StringAttribute{
				Optional: true,
				MarkdownDescription: "The 6-digit PIN required to unlock macOS devices with Activation Lock enabled. " +
					"Required for supervised macOS devices with Activation Lock. " +
					"Format: 6-digit numeric string.",
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^\d{6}$`),
						"must be a 6-digit numeric string",
					),
				},
			},
			"obliteration_behavior": schema.StringAttribute{
				Optional: true,
				MarkdownDescription: "Specifies the obliteration behavior for macOS 12+ devices with Apple M1 chip or Apple T2 Security Chip. " +
					"This controls fallback behavior when Erase All Content and Settings (EACS) cannot run.\n\n" +
					"Valid values:\n" +
					"- `default`: If EACS preflight fails, device responds with Error status and attempts to erase itself. " +
					"If EACS preflight succeeds but EACS fails, the device attempts to erase itself.\n" +
					"- `doNotObliterate`: If EACS preflight fails, device responds with Error status and doesn't attempt to erase. " +
					"If EACS preflight succeeds but EACS fails, the device doesn't attempt to erase itself.\n" +
					"- `obliterateWithWarning`: If EACS preflight fails, device responds with Acknowledged status and attempts to erase itself. " +
					"If EACS preflight succeeds but EACS fails, the device attempts to erase itself.\n" +
					"- `always`: The system doesn't attempt EACS. T2 and later devices always obliterate.\n\n" +
					"**Note:** This setting only applies to Mac computers with Apple M1 chip or Apple T2 Security Chip running macOS 12 or later. " +
					"It has no effect on machines prior to the T2 chip.\n\n" +
					"**Reference:** [obliterationBehavior enum](https://learn.microsoft.com/en-us/graph/api/resources/intune-devices-obliterationbehavior?view=graph-rest-beta)",
				Validators: []validator.String{
					stringvalidator.OneOf("default", "doNotObliterate", "obliterateWithWarning", "always"),
				},
			},
			"persist_esim_data_plan": schema.BoolAttribute{
				Optional: true,
				MarkdownDescription: "If `true`, preserves the eSIM data plan on the device during wipe. " +
					"Only applicable to devices with eSIM support. " +
					"Defaults to `false`.",
			},
			"use_protected_wipe": schema.BoolAttribute{
				Optional: true,
				MarkdownDescription: "If `true`, uses protected wipe for Windows 10/11 devices. " +
					"Protected wipe maintains UEFI-embedded product keys and recovery partition. " +
					"Only applicable to Windows devices. " +
					"Defaults to `false`.",
			},
			"timeouts": commonschema.ActionTimeouts(ctx),
		},
	}
}
