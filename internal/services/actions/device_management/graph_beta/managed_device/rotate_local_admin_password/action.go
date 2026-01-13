package graphBetaRotateLocalAdminPasswordManagedDevice

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
	ActionName    = "microsoft365_graph_beta_device_management_managed_device_rotate_local_admin_password"
	InvokeTimeout = 60
)

var (
	_ action.Action                   = &RotateLocalAdminPasswordManagedDeviceAction{}
	_ action.ActionWithConfigure      = &RotateLocalAdminPasswordManagedDeviceAction{}
	_ action.ActionWithValidateConfig = &RotateLocalAdminPasswordManagedDeviceAction{}
)

func NewRotateLocalAdminPasswordManagedDeviceAction() action.Action {
	return &RotateLocalAdminPasswordManagedDeviceAction{
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

type RotateLocalAdminPasswordManagedDeviceAction struct {
	client           *msgraphbetasdk.GraphServiceClient
	ReadPermissions  []string
	WritePermissions []string
}

func (a *RotateLocalAdminPasswordManagedDeviceAction) Metadata(ctx context.Context, req action.MetadataRequest, resp *action.MetadataResponse) {
	resp.TypeName = ActionName
}

func (a *RotateLocalAdminPasswordManagedDeviceAction) Configure(ctx context.Context, req action.ConfigureRequest, resp *action.ConfigureResponse) {
	a.client = client.SetGraphBetaClientForAction(ctx, req, resp, ActionName)
}

func (a *RotateLocalAdminPasswordManagedDeviceAction) Schema(ctx context.Context, req action.SchemaRequest, resp *action.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Initiates manual rotation of the local administrator password on managed Windows devices using the " +
			"`/deviceManagement/managedDevices/{managedDeviceId}/rotateLocalAdminPassword` and " +
			"`/deviceManagement/comanagedDevices/{managedDeviceId}/rotateLocalAdminPassword` endpoints. " +
			"This action works with Windows Local Administrator Password Solution (LAPS) to generate and rotate local admin " +
			"passwords on Windows devices. The new password is automatically generated, stored securely in Azure AD or Intune, " +
			"and can be retrieved by authorized administrators. This enhances security by ensuring regular password rotation and " +
			"centralized password management for local administrator accounts.\n\n" +
			"**Important Notes:**\n" +
			"- Only works on Windows 10/11 devices with Windows LAPS enabled\n" +
			"- Requires Windows LAPS policy configured in Intune\n" +
			"- New password automatically generated and stored in Azure AD/Intune\n" +
			"- Password retrievable by authorized administrators\n" +
			"- Does not affect device operation or require restart\n" +
			"- Critical for security compliance and privileged access management\n" +
			"- Should be performed regularly or after admin account compromise\n\n" +
			"**Use Cases:**\n" +
			"- Regular security password rotation (quarterly/semi-annually)\n" +
			"- Post-incident password reset after compromise\n" +
			"- Compliance requirements for privileged account management\n" +
			"- Onboarding/offboarding administrator access\n" +
			"- Audit preparation and security validation\n" +
			"- Zero Trust privileged access implementation\n\n" +
			"**Platform Support:**\n" +
			"- **Windows**: Windows 10/11 with Windows LAPS enabled\n" +
			"- **Other Platforms**: Not supported (Windows LAPS-specific)\n\n" +
			"**Reference:** [Microsoft Graph API - Rotate Local Admin Password](https://learn.microsoft.com/en-us/graph/api/intune-devices-manageddevice-rotatelocaladminpassword?view=graph-rest-beta)",
		Attributes: map[string]schema.Attribute{
			"managed_device_ids": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
				MarkdownDescription: "List of managed device IDs (GUIDs) to rotate local administrator passwords for. " +
					"These are devices fully managed by Intune with Windows LAPS enabled.\n\n" +
					"**Note:** At least one of `managed_device_ids` or `comanaged_device_ids` must be provided. " +
					"You can provide both to rotate passwords on different types of devices in one action.\n\n" +
					"**Important:** Devices must have Windows LAPS policy configured and enabled. The new password will be " +
					"automatically generated and stored securely in Azure AD or Intune.\n\n" +
					"Example: `[\"12345678-1234-1234-1234-123456789abc\", \"87654321-4321-4321-4321-ba9876543210\"]`",
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
				MarkdownDescription: "List of co-managed device IDs (GUIDs) to rotate local administrator passwords for. " +
					"These are devices managed by both Intune and Configuration Manager (SCCM) with Windows LAPS enabled.\n\n" +
					"**Note:** At least one of `managed_device_ids` or `comanaged_device_ids` must be provided.\n\n" +
					"Example: `[\"abcdef12-3456-7890-abcd-ef1234567890\"]`",
				Validators: []validator.List{
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
				MarkdownDescription: "Whether to validate that devices exist and are Windows devices before attempting to rotate local admin passwords. " +
					"Disabling this can speed up planning but may result in runtime errors for non-existent or non-Windows devices. " +
					"Default: `true`.",
			},
			"timeouts": commonschema.ActionTimeouts(ctx),
		},
	}
}
