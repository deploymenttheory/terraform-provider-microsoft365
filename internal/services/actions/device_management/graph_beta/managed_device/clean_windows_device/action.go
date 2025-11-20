package graphBetaCleanWindowsManagedDevice

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
	ActionName = "microsoft365_graph_beta_device_management_managed_device_clean_windows_device"
)

var (
	_ action.Action                   = &CleanWindowsManagedDeviceAction{}
	_ action.ActionWithConfigure      = &CleanWindowsManagedDeviceAction{}
	_ action.ActionWithValidateConfig = &CleanWindowsManagedDeviceAction{}
)

func NewCleanWindowsManagedDeviceAction() action.Action {
	return &CleanWindowsManagedDeviceAction{
		ReadPermissions: []string{
			"DeviceManagementManagedDevices.PrivilegedOperations.All",
		},
		WritePermissions: []string{
			"DeviceManagementManagedDevices.PrivilegedOperations.All",
		},
	}
}

type CleanWindowsManagedDeviceAction struct {
	client           *msgraphbetasdk.GraphServiceClient
	ReadPermissions  []string
	WritePermissions []string
}

func (a *CleanWindowsManagedDeviceAction) Metadata(ctx context.Context, req action.MetadataRequest, resp *action.MetadataResponse) {
	resp.TypeName = ActionName
}

func (a *CleanWindowsManagedDeviceAction) Configure(ctx context.Context, req action.ConfigureRequest, resp *action.ConfigureResponse) {
	a.client = client.SetGraphBetaClientForAction(ctx, req, resp, ActionName)
}

func (a *CleanWindowsManagedDeviceAction) Schema(ctx context.Context, req action.SchemaRequest, resp *action.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Performs a clean operation on Windows managed and co-managed devices using the " +
			"`/deviceManagement/managedDevices/{managedDeviceId}/cleanWindowsDevice` and " +
			"`/deviceManagement/comanagedDevices/{managedDeviceId}/cleanWindowsDevice` endpoints. " +
			"This action provides a lighter-weight alternative to full device wipe, allowing IT administrators " +
			"to remove applications and settings while optionally preserving user data on each device independently.\n\n" +
			"**What Clean Windows Device Does:**\n" +
			"- Removes installed applications (except inbox Windows apps)\n" +
			"- Removes user profiles (unless `keep_user_data` is true for that device)\n" +
			"- Removes device configuration settings\n" +
			"- Removes company policies and profiles\n" +
			"- Can preserve user data per-device if specified\n" +
			"- Device remains enrolled in Intune\n" +
			"- Less destructive than full wipe\n\n" +
			"**Platform Support:**\n" +
			"- **Windows**: Full support (Windows 10/11)\n" +
			"- **Other platforms**: Not supported (Windows-only action)\n\n" +
			"**Clean vs Wipe vs Retire:**\n" +
			"- **Clean**: Removes apps/settings, optionally keeps user data, device stays enrolled\n" +
			"- **Wipe**: Factory reset, removes all data, device must re-enroll\n" +
			"- **Retire**: Removes company data only, preserves personal data\n\n" +
			"**Common Use Cases:**\n" +
			"- Device refresh without full rebuild\n" +
			"- Removing malware/unwanted applications\n" +
			"- Preparing device for new user (keeping OS)\n" +
			"- Troubleshooting device issues\n" +
			"- Compliance remediation\n" +
			"- Software bloat removal\n" +
			"- Maintaining device enrollment\n\n" +
			"**Important Considerations:**\n" +
			"- Device must be online to receive command\n" +
			"- User will lose unsaved work\n" +
			"- Installed applications will be removed\n" +
			"- Process may take several minutes\n" +
			"- Device remains in Intune (no re-enrollment needed)\n" +
			"- Each device can have different `keep_user_data` setting\n\n" +
			"**Reference:** [Microsoft Graph API - Clean Windows Device](https://learn.microsoft.com/en-us/graph/api/intune-devices-manageddevice-cleanwindowsdevice?view=graph-rest-beta)",
		Attributes: map[string]schema.Attribute{
			"ignore_partial_failures": schema.BoolAttribute{
				Optional: true,
				MarkdownDescription: "If set to `true`, the action will succeed even if some devices fail clean operation. " +
					"Failed devices will be reported as warnings instead of errors. " +
					"Default: `false` (action fails if any device fails).",
			},
			"validate_device_exists": schema.BoolAttribute{
				Optional: true,
				MarkdownDescription: "Whether to validate that devices exist and are Windows devices before attempting clean. " +
					"Disabling this can speed up planning but may result in runtime errors for non-existent or non-Windows devices. " +
					"Default: `true`.",
			},
			"timeouts": commonschema.Timeouts(ctx),
		},
		Blocks: map[string]schema.Block{
			"managed_devices": schema.ListNestedBlock{
				MarkdownDescription: "List of managed Windows devices to clean. These are devices fully managed by Intune only. " +
					"Each entry specifies a device ID and whether to preserve user data.\n\n" +
					"**Examples:**\n" +
					"```hcl\n" +
					"managed_devices = [\n" +
					"  {\n" +
					"    device_id       = \"12345678-1234-1234-1234-123456789abc\"\n" +
					"    keep_user_data  = false\n" +
					"  },\n" +
					"  {\n" +
					"    device_id       = \"87654321-4321-4321-4321-987654321cba\"\n" +
					"    keep_user_data  = true\n" +
					"  }\n" +
					"]\n" +
					"```\n\n" +
					"**Note:** At least one of `managed_devices` or `comanaged_devices` must be provided.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"device_id": schema.StringAttribute{
							Required: true,
							MarkdownDescription: "The unique identifier (GUID) of the managed Windows device to clean. " +
								"Device must be Windows 10 or Windows 11. " +
								"Example: `\"12345678-1234-1234-1234-123456789abc\"`",
							Validators: []validator.String{
								stringvalidator.RegexMatches(
									regexp.MustCompile(constants.GuidRegex),
									"device_id must be a valid GUID format (e.g., 12345678-1234-1234-1234-123456789abc)",
								),
							},
						},
						"keep_user_data": schema.BoolAttribute{
							Required: true,
							MarkdownDescription: "Determines whether user data should be preserved for this device during the clean operation. " +
								"**Required field** - must be explicitly set to `true` or `false`.\n\n" +
								"**When `false`:**\n" +
								"- User profiles removed\n" +
								"- User data deleted\n" +
								"- Applications removed\n" +
								"- Settings reset\n\n" +
								"**When `true`:**\n" +
								"- User profiles preserved\n" +
								"- User data kept (documents, desktop, etc.)\n" +
								"- Applications still removed\n" +
								"- Settings still reset",
						},
					},
				},
			},
			"comanaged_devices": schema.ListNestedBlock{
				MarkdownDescription: "List of co-managed Windows devices to clean. These are devices managed by both Intune and " +
					"Configuration Manager (SCCM). Each entry specifies a device ID and whether to preserve user data.\n\n" +
					"**Examples:**\n" +
					"```hcl\n" +
					"comanaged_devices = [\n" +
					"  {\n" +
					"    device_id       = \"aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee\"\n" +
					"    keep_user_data  = false\n" +
					"  }\n" +
					"]\n" +
					"```\n\n" +
					"**Note:** At least one of `managed_devices` or `comanaged_devices` must be provided.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"device_id": schema.StringAttribute{
							Required: true,
							MarkdownDescription: "The unique identifier (GUID) of the co-managed Windows device to clean. " +
								"Example: `\"12345678-1234-1234-1234-123456789abc\"`",
							Validators: []validator.String{
								stringvalidator.RegexMatches(
									regexp.MustCompile(constants.GuidRegex),
									"device_id must be a valid GUID format (e.g., 12345678-1234-1234-1234-123456789abc)",
								),
							},
						},
						"keep_user_data": schema.BoolAttribute{
							Required: true,
							MarkdownDescription: "Determines whether user data should be preserved for this device during the clean operation. " +
								"**Required field** - must be explicitly set to `true` or `false`.\n\n" +
								"**When `false`:**\n" +
								"- User profiles removed\n" +
								"- User data deleted\n" +
								"- Applications removed\n" +
								"- Settings reset\n\n" +
								"**When `true`:**\n" +
								"- User profiles preserved\n" +
								"- User data kept (documents, desktop, etc.)\n" +
								"- Applications still removed\n" +
								"- Settings still reset",
						},
					},
				},
			},
		},
	}
}
