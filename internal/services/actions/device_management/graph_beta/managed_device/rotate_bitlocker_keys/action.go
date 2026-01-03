package graphBetaRotateBitLockerKeys

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
	ActionName = "microsoft365_graph_beta_device_management_managed_device_rotate_bitlocker_keys"
	InvokeTimeout = 60
)

var (
	_ action.Action                   = &RotateBitLockerKeysAction{}
	_ action.ActionWithConfigure      = &RotateBitLockerKeysAction{}
	_ action.ActionWithValidateConfig = &RotateBitLockerKeysAction{}
)

func NewRotateBitLockerKeysAction() action.Action {
	return &RotateBitLockerKeysAction{
		ReadPermissions: []string{
			"DeviceManagementConfiguration.ReadWrite.All",
			"DeviceManagementManagedDevices.ReadWrite.All",
		},
		WritePermissions: []string{
			"DeviceManagementConfiguration.ReadWrite.All",
			"DeviceManagementManagedDevices.ReadWrite.All",
		},
	}
}

type RotateBitLockerKeysAction struct {
	client           *msgraphbetasdk.GraphServiceClient
	ReadPermissions  []string
	WritePermissions []string
}

func (a *RotateBitLockerKeysAction) Metadata(ctx context.Context, req action.MetadataRequest, resp *action.MetadataResponse) {
	resp.TypeName = ActionName
}

func (a *RotateBitLockerKeysAction) Configure(ctx context.Context, req action.ConfigureRequest, resp *action.ConfigureResponse) {
	a.client = client.SetGraphBetaClientForAction(ctx, req, resp, ActionName)
}

func (a *RotateBitLockerKeysAction) Schema(ctx context.Context, req action.SchemaRequest, resp *action.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Rotates BitLocker encryption recovery keys on Windows devices using the " +
			"`/deviceManagement/managedDevices/{managedDeviceId}/rotateBitLockerKeys` and " +
			"`/deviceManagement/comanagedDevices/{managedDeviceId}/rotateBitLockerKeys` endpoints. " +
			"This action generates new BitLocker recovery keys and escrows them to Intune, invalidating the previous recovery keys.\n\n" +
			"**What This Action Does:**\n" +
			"- Generates new BitLocker recovery password\n" +
			"- Escrows new recovery key to Intune/Azure AD\n" +
			"- Invalidates previous recovery keys\n" +
			"- Updates key protector on device\n" +
			"- Maintains encryption state (no re-encryption)\n" +
			"- Audits key rotation event\n\n" +
			"**When to Use:**\n" +
			"- Security incident or breach response\n" +
			"- Recovery key compromised or exposed\n" +
			"- Compliance policy requirement\n" +
			"- Regular security maintenance schedule\n" +
			"- Device ownership transfer\n" +
			"- Administrative access change\n" +
			"- Proactive security hardening\n\n" +
			"**Platform Support:**\n" +
			"- **Windows 10**: Pro, Enterprise, Education (v1703+)\n" +
			"- **Windows 11**: All editions with BitLocker\n" +
			"- **Other platforms**: Not applicable (no BitLocker)\n\n" +
			"**Important Considerations:**\n" +
			"- Only rotates recovery keys, not encryption keys\n" +
			"- Device must be online and connected\n" +
			"- BitLocker must be enabled and configured\n" +
			"- Previous recovery keys become invalid\n" +
			"- New keys escrowed to Azure AD/Intune\n" +
			"- No user interaction required\n" +
			"- No device restart needed\n\n" +
			"**Reference:** [Microsoft Graph API - Rotate BitLocker Keys](https://learn.microsoft.com/en-us/graph/api/intune-devices-manageddevice-rotatebitlockerkeys?view=graph-rest-beta)",
		Attributes: map[string]schema.Attribute{
			"managed_device_ids": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
				MarkdownDescription: "List of managed device IDs to rotate BitLocker keys on. These are Windows devices fully managed by Intune only. " +
					"Each ID must be a valid GUID format. BitLocker recovery keys will be rotated on these devices. " +
					"Example: `[\"12345678-1234-1234-1234-123456789abc\", \"87654321-4321-4321-4321-ba9876543210\"]`\n\n" +
					"**Note:** At least one of `managed_device_ids` or `comanaged_device_ids` must be provided. " +
					"You can provide both to rotate keys on different types of devices in one action.",
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
				MarkdownDescription: "List of co-managed device IDs to rotate BitLocker keys on. These are Windows devices managed by both Intune and " +
					"Configuration Manager (SCCM). Each ID must be a valid GUID format. " +
					"Example: `[\"12345678-1234-1234-1234-123456789abc\"]`\n\n" +
					"**Note:** Co-managed devices can have BitLocker keys escrowed to both Intune and Configuration Manager. " +
					"At least one of `managed_device_ids` or `comanaged_device_ids` must be provided.",
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
