package graphBetaBypassActivationLockManagedDevice

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
	ActionName    = "microsoft365_graph_beta_device_management_managed_device_bypass_activation_lock"
	InvokeTimeout = 60
)

var (
	_ action.Action                   = &BypassActivationLockManagedDeviceAction{}
	_ action.ActionWithConfigure      = &BypassActivationLockManagedDeviceAction{}
	_ action.ActionWithValidateConfig = &BypassActivationLockManagedDeviceAction{}
)

func NewBypassActivationLockManagedDeviceAction() action.Action {
	return &BypassActivationLockManagedDeviceAction{}
}

type BypassActivationLockManagedDeviceAction struct {
	client *msgraphbetasdk.GraphServiceClient
}

func (a *BypassActivationLockManagedDeviceAction) Metadata(ctx context.Context, req action.MetadataRequest, resp *action.MetadataResponse) {
	resp.TypeName = ActionName
}

func (a *BypassActivationLockManagedDeviceAction) Configure(ctx context.Context, req action.ConfigureRequest, resp *action.ConfigureResponse) {
	a.client = client.SetGraphBetaClientForAction(ctx, req, resp, ActionName)
}

func (a *BypassActivationLockManagedDeviceAction) Schema(ctx context.Context, req action.SchemaRequest, resp *action.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Bypasses Activation Lock on iOS, iPadOS, and macOS devices using the " +
			"`/deviceManagement/managedDevices/{managedDeviceId}/bypassActivationLock` endpoint. " +
			"Activation Lock is an Apple security feature that prevents unauthorized use of a device after it has been erased. " +
			"When Find My iPhone/iPad/Mac is enabled and a device is erased, Activation Lock requires the original Apple ID " +
			"and password before the device can be reactivated. This action generates a bypass code that allows IT administrators " +
			"to reactivate managed devices without the user's Apple ID credentials.\n\n" +
			"**What is Activation Lock?**\n" +
			"- Security feature built into iOS, iPadOS, and macOS\n" +
			"- Automatically enabled when Find My iPhone/iPad/Mac is turned on\n" +
			"- Prevents device reactivation after factory reset without Apple ID credentials\n" +
			"- Helps prevent theft and unauthorized device reuse\n" +
			"- Links device to specific Apple ID\n\n" +
			"**Important Notes:**\n" +
			"- Device must be supervised (iOS/iPadOS) or enrolled via DEP/ABM (macOS)\n" +
			"- Activation Lock must currently be enabled on the device\n" +
			"- Generates a bypass code stored in Intune for future use\n" +
			"- Bypass code can be retrieved from device properties in Intune portal\n" +
			"- Code can be used during device setup to bypass Activation Lock screen\n" +
			"- Does not disable Find My iPhone/iPad/Mac, only provides bypass capability\n" +
			"- Bypass code remains valid until Activation Lock is disabled by user\n\n" +
			"**Use Cases:**\n" +
			"- Wiping and reassigning corporate devices to new employees\n" +
			"- Recovering devices from departing employees who forgot to disable Find My\n" +
			"- Preparing devices for return to vendor or recycling\n" +
			"- Enabling IT to factory reset and redeploy devices\n" +
			"- Handling devices with lost or forgotten Apple ID credentials\n" +
			"- Bulk device preparation and provisioning\n\n" +
			"**Platform Support:**\n" +
			"- **iOS**: Supported (iOS 7.1+, supervised devices only)\n" +
			"- **iPadOS**: Supported (supervised devices only)\n" +
			"- **macOS**: Supported (macOS 10.11+, DEP/ABM enrolled devices)\n" +
			"- **Other Platforms**: Not supported (Activation Lock is Apple-only feature)\n\n" +
			"**How to Use Bypass Code:**\n" +
			"1. Issue bypass command via this action\n" +
			"2. Retrieve bypass code from Intune portal (device properties)\n" +
			"3. Erase/wipe the device (using wipe action or manually)\n" +
			"4. When device shows Activation Lock screen during setup\n" +
			"5. Enter bypass code in password field\n" +
			"6. Device will bypass Activation Lock and complete setup\n\n" +
			"**Workflow Example:**\n" +
			"```\n" +
			"1. Employee leaves organization with device in lost mode\n" +
			"2. IT issues bypass activation lock command (this action)\n" +
			"3. IT retrieves bypass code from Intune portal\n" +
			"4. IT wipes device (removes all data)\n" +
			"5. During setup, device shows Activation Lock screen\n" +
			"6. IT enters bypass code to unlock device\n" +
			"7. Device can now be re-enrolled and assigned to new user\n" +
			"```\n\n" +
			"**Security Considerations:**\n" +
			"- Bypass code should be treated as sensitive credential\n" +
			"- Only authorized IT staff should have access to bypass codes\n" +
			"- Document usage for compliance and audit purposes\n" +
			"- Consider implementing approval workflow for bypass requests\n" +
			"- Verify device ownership before issuing bypass\n" +
			"- Bypass does not affect device security after reactivation\n\n" +
			"**Limitations:**\n" +
			"- Cannot bypass Activation Lock on personal (non-supervised) iOS/iPadOS devices\n" +
			"- Cannot bypass Activation Lock on macOS devices not enrolled via DEP/ABM\n" +
			"- Bypass code only works for device it was generated for\n" +
			"- Must have Activation Lock bypass code retrieved before device is erased\n" +
			"- Some older device models may not support this feature\n\n" +
			"**Best Practices:**\n" +
			"- Issue bypass command before wiping device when possible\n" +
			"- Store bypass codes securely in password manager or secure vault\n" +
			"- Document which devices have bypass codes generated\n" +
			"- Include activation lock bypass in device offboarding procedures\n" +
			"- Test bypass process in controlled environment first\n" +
			"- Verify device supervision status before attempting bypass\n" +
			"- Consider enabling automatic bypass code escrow during enrollment\n\n" +
			"**Reference:** [Microsoft Graph API - Bypass Activation Lock](https://learn.microsoft.com/en-us/graph/api/intune-devices-manageddevice-bypassactivationlock?view=graph-rest-beta)",
		Attributes: map[string]schema.Attribute{
			"device_ids": schema.ListAttribute{
				ElementType: types.StringType,
				Required:    true,
				MarkdownDescription: "List of managed device IDs to generate Activation Lock bypass codes for. " +
					"Each ID must be a valid GUID format. Multiple devices can have bypass codes generated in a single action. " +
					"Example: `[\"12345678-1234-1234-1234-123456789abc\", \"87654321-4321-4321-4321-ba9876543210\"]`\n\n" +
					"**Important:** Devices must be supervised iOS/iPadOS devices or DEP/ABM enrolled macOS devices with Activation Lock enabled. " +
					"The bypass code will be stored in Intune and can be retrieved from device properties in the admin portal. " +
					"This code is required to reactivate the device after a factory reset if Activation Lock is enabled.",
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
				MarkdownDescription: "If set to `true`, the action will succeed even if some devices fail Activation Lock bypass. " +
					"Failed devices will be reported as warnings instead of errors. " +
					"Default: `false` (action fails if any device fails).",
			},
			"validate_device_exists": schema.BoolAttribute{
				Optional: true,
				MarkdownDescription: "Whether to validate that devices exist and support Activation Lock before attempting bypass. " +
					"Disabling this can speed up planning but may result in runtime errors for non-existent or unsupported devices. " +
					"Default: `true`.",
			},
			"timeouts": commonschema.ActionTimeouts(ctx),
		},
	}
}
