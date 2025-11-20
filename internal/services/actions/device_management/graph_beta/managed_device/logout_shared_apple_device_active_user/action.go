package graphBetaLogoutSharedAppleDeviceActiveUser

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
	ActionName = "microsoft365_graph_beta_device_management_managed_device_logout_shared_apple_device_active_user"
)

var (
	_ action.Action                   = &LogoutSharedAppleDeviceActiveUserAction{}
	_ action.ActionWithConfigure      = &LogoutSharedAppleDeviceActiveUserAction{}
	_ action.ActionWithValidateConfig = &LogoutSharedAppleDeviceActiveUserAction{}
)

func NewLogoutSharedAppleDeviceActiveUserAction() action.Action {
	return &LogoutSharedAppleDeviceActiveUserAction{
		ReadPermissions: []string{
			"DeviceManagementManagedDevices.PrivilegedOperations.All",
		},
		WritePermissions: []string{
			"DeviceManagementManagedDevices.PrivilegedOperations.All",
		},
	}
}

type LogoutSharedAppleDeviceActiveUserAction struct {
	client           *msgraphbetasdk.GraphServiceClient
	ReadPermissions  []string
	WritePermissions []string
}

func (a *LogoutSharedAppleDeviceActiveUserAction) Metadata(ctx context.Context, req action.MetadataRequest, resp *action.MetadataResponse) {
	resp.TypeName = ActionName
}

func (a *LogoutSharedAppleDeviceActiveUserAction) Configure(ctx context.Context, req action.ConfigureRequest, resp *action.ConfigureResponse) {
	a.client = client.SetGraphBetaClientForAction(ctx, req, resp, ActionName)
}

func (a *LogoutSharedAppleDeviceActiveUserAction) Schema(ctx context.Context, req action.SchemaRequest, resp *action.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Logs out the currently active user from Shared iPad devices using the " +
			"`/deviceManagement/managedDevices/{managedDeviceId}/logoutSharedAppleDeviceActiveUser` endpoint. " +
			"This action is specifically designed for iPads configured in Shared iPad mode, where multiple users " +
			"can use the same device while maintaining separate user environments.\n\n" +
			"**What is Shared iPad Mode?**\n" +
			"- Educational/enterprise feature for iPadOS\n" +
			"- Multiple users share single physical device\n" +
			"- Each user has separate data and settings\n" +
			"- Users log in with their Apple ID or Managed Apple ID\n" +
			"- Local caching of user data for offline access\n" +
			"- Requires supervised iPads enrolled via DEP/ABM\n\n" +
			"**What This Action Does:**\n" +
			"- Logs out currently active user from Shared iPad\n" +
			"- Returns device to login screen\n" +
			"- Preserves user data on device (cached locally)\n" +
			"- Allows next user to log in\n" +
			"- Does not remove user from device roster\n" +
			"- Does not delete user's cached data\n\n" +
			"**Platform Support:**\n" +
			"- **iPadOS**: Full support (Shared iPad mode only)\n" +
			"- **iOS**: Not supported (iPhones don't support Shared mode)\n" +
			"- **macOS**: Not supported\n" +
			"- **Windows**: Not supported\n" +
			"- **Android**: Not supported\n\n" +
			"**Common Use Cases:**\n" +
			"- Classroom management (switching students)\n" +
			"- End of class period user logout\n" +
			"- Preparing device for next user\n" +
			"- Remote user session management\n" +
			"- Enforcing session time limits\n" +
			"- Cart/lab device rotation\n" +
			"- Emergency user logout\n" +
			"- Troubleshooting user sessions\n\n" +
			"**Requirements:**\n" +
			"- iPad must be in Shared iPad mode\n" +
			"- Device must be supervised\n" +
			"- Must be enrolled via DEP/ABM\n" +
			"- User must be actively logged in\n" +
			"- Device must be online\n\n" +
			"**Important Notes:**\n" +
			"- Only affects Shared iPad devices\n" +
			"- Regular (non-shared) iPads: action has no effect\n" +
			"- User data remains cached on device\n" +
			"- User can log back in immediately\n" +
			"- Unsaved work may be lost\n" +
			"- Active apps will close\n\n" +
			"**Reference:** [Microsoft Graph API - Logout Shared Apple Device Active User](https://learn.microsoft.com/en-us/graph/api/intune-devices-manageddevice-logoutsharedappledeviceactiveuser?view=graph-rest-beta)",
		Attributes: map[string]schema.Attribute{
			"device_ids": schema.ListAttribute{
				ElementType: types.StringType,
				Required:    true,
				MarkdownDescription: "List of Shared iPad device IDs to log out the active user from. " +
					"Each ID must be a valid GUID format. Multiple devices can be processed in a single action. " +
					"Example: `[\"12345678-1234-1234-1234-123456789abc\", \"87654321-4321-4321-4321-ba9876543210\"]`\n\n" +
					"**Important:** This action only works on iPads configured in Shared iPad mode. " +
					"The action will fail or have no effect on regular (non-shared) iPads, iPhones, or other device types. " +
					"Ensure all device IDs refer to Shared iPad devices before executing this action.",
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
