package graphBetaRotateFileVaultKeyManagedDevice

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
	ActionName = "microsoft365_graph_beta_device_management_managed_device_rotate_file_vault_key"
)

var (
	_ action.Action                   = &RotateFileVaultKeyManagedDeviceAction{}
	_ action.ActionWithConfigure      = &RotateFileVaultKeyManagedDeviceAction{}
	_ action.ActionWithValidateConfig = &RotateFileVaultKeyManagedDeviceAction{}
)

func NewRotateFileVaultKeyManagedDeviceAction() action.Action {
	return &RotateFileVaultKeyManagedDeviceAction{
		ReadPermissions: []string{
			"DeviceManagementManagedDevices.PrivilegedOperations.All",
		},
		WritePermissions: []string{
			"DeviceManagementManagedDevices.PrivilegedOperations.All",
		},
	}
}

type RotateFileVaultKeyManagedDeviceAction struct {
	client           *msgraphbetasdk.GraphServiceClient
	ReadPermissions  []string
	WritePermissions []string
}

func (a *RotateFileVaultKeyManagedDeviceAction) Metadata(ctx context.Context, req action.MetadataRequest, resp *action.MetadataResponse) {
	resp.TypeName = ActionName
}

func (a *RotateFileVaultKeyManagedDeviceAction) Configure(ctx context.Context, req action.ConfigureRequest, resp *action.ConfigureResponse) {
	a.client = client.SetGraphBetaClientForAction(ctx, req, resp, ActionName)
}

func (a *RotateFileVaultKeyManagedDeviceAction) Schema(ctx context.Context, req action.SchemaRequest, resp *action.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Rotates the FileVault recovery key for macOS managed devices using the " +
			"`/deviceManagement/managedDevices/{managedDeviceId}/rotateFileVaultKey` and " +
			"`/deviceManagement/comanagedDevices/{managedDeviceId}/rotateFileVaultKey` endpoints. " +
			"This action generates a new FileVault recovery key and escrows it with Intune, ensuring that administrators " +
			"can recover encrypted macOS devices if users forget their passwords or lose access. Regular key rotation is " +
			"a security best practice that limits the window of exposure if a key is compromised. This action supports " +
			"rotating keys on multiple devices in a single operation.\n\n" +
			"**Important Notes:**\n" +
			"- Only applicable to macOS devices with FileVault enabled\n" +
			"- Generates a new personal recovery key\n" +
			"- New key is escrowed with Intune automatically\n" +
			"- Previous recovery key becomes invalid\n" +
			"- Device must be online to receive rotation command\n" +
			"- User does not need to be logged in\n" +
			"- No user interaction required for rotation\n\n" +
			"**Use Cases:**\n" +
			"- Regular security key rotation compliance\n" +
			"- After potential key compromise or exposure\n" +
			"- When changing device ownership or assignment\n" +
			"- As part of security incident response\n" +
			"- Periodic rotation per security policy\n" +
			"- Before device reassignment to new users\n\n" +
			"**Platform Support:**\n" +
			"- **macOS**: Fully supported on devices with FileVault enabled\n" +
			"- **Other Platforms**: Not applicable (FileVault is macOS-only)\n\n" +
			"**Reference:** [Microsoft Graph API - Rotate FileVault Key](https://learn.microsoft.com/en-us/graph/api/intune-devices-manageddevice-rotatefilevaultkey?view=graph-rest-beta)",
		Attributes: map[string]schema.Attribute{
			"managed_device_ids": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
				MarkdownDescription: "List of managed device IDs (GUIDs) to rotate FileVault keys for. These are macOS devices " +
					"fully managed by Intune only. Each device must have FileVault encryption enabled.\n\n" +
					"**Note:** At least one of `managed_device_ids` or `comanaged_device_ids` must be provided. " +
					"You can provide both to rotate keys on different types of devices in one action.\n\n" +
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
				MarkdownDescription: "List of co-managed device IDs (GUIDs) to rotate FileVault keys for. These are macOS devices " +
					"managed by both Intune and Configuration Manager (SCCM).\n\n" +
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
			"timeouts": commonschema.Timeouts(ctx),
		},
	}
}
