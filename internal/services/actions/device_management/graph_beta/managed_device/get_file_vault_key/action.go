package graphBetaGetFileVaultKeyManagedDevice

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
	ActionName = "microsoft365_graph_beta_device_management_managed_device_get_file_vault_key"
	InvokeTimeout = 60
)

var (
	_ action.Action                   = &GetFileVaultKeyManagedDeviceAction{}
	_ action.ActionWithConfigure      = &GetFileVaultKeyManagedDeviceAction{}
	_ action.ActionWithValidateConfig = &GetFileVaultKeyManagedDeviceAction{}
)

func NewGetFileVaultKeyManagedDeviceAction() action.Action {
	return &GetFileVaultKeyManagedDeviceAction{
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

type GetFileVaultKeyManagedDeviceAction struct {
	client           *msgraphbetasdk.GraphServiceClient
	ReadPermissions  []string
	WritePermissions []string
}

func (a *GetFileVaultKeyManagedDeviceAction) Metadata(ctx context.Context, req action.MetadataRequest, resp *action.MetadataResponse) {
	resp.TypeName = ActionName
}

func (a *GetFileVaultKeyManagedDeviceAction) Configure(ctx context.Context, req action.ConfigureRequest, resp *action.ConfigureResponse) {
	a.client = client.SetGraphBetaClientForAction(ctx, req, resp, ActionName)
}

func (a *GetFileVaultKeyManagedDeviceAction) Schema(ctx context.Context, req action.SchemaRequest, resp *action.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Retrieves the FileVault recovery key for macOS managed devices using the " +
			"`/deviceManagement/managedDevices/{managedDeviceId}/getFileVaultKey` and " +
			"`/deviceManagement/comanagedDevices/{managedDeviceId}/getFileVaultKey` endpoints. " +
			"This action allows administrators to retrieve escrowed FileVault recovery keys for device recovery purposes. " +
			"The recovery key is displayed in the action output and can be used to unlock an encrypted macOS device when " +
			"a user has forgotten their password or is otherwise unable to access the device. This is a critical capability " +
			"for IT support and device recovery scenarios.\n\n" +
			"**Important Security Notes:**\n" +
			"- Recovery keys are highly sensitive credentials\n" +
			"- Keys grant full access to encrypted device data\n" +
			"- Access to keys should be audited and restricted\n" +
			"- Only retrieve keys when necessary for device recovery\n" +
			"- Keys are displayed in plain text in action output\n" +
			"- Ensure proper security controls on Terraform state\n" +
			"- Consider security implications before using in automation\n\n" +
			"**Use Cases:**\n" +
			"- Emergency device recovery when user cannot log in\n" +
			"- Unlocking devices for departing employees\n" +
			"- Technical support scenarios requiring device access\n" +
			"- Disaster recovery and business continuity\n" +
			"- Device repurposing or reassignment preparation\n\n" +
			"**Platform Support:**\n" +
			"- **macOS**: Fully supported on devices with FileVault enabled and keys escrowed\n" +
			"- **Other Platforms**: Not applicable (FileVault is macOS-only)\n\n" +
			"**Reference:** [Microsoft Graph API - Get FileVault Key](https://learn.microsoft.com/en-us/graph/api/intune-devices-manageddevice-getfilevaultkey?view=graph-rest-beta)",
		Attributes: map[string]schema.Attribute{
			"managed_device_ids": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
				MarkdownDescription: "List of managed device IDs (GUIDs) to retrieve FileVault keys for. These are macOS devices " +
					"fully managed by Intune only. Each device must have FileVault encryption enabled with key escrowed to Intune.\n\n" +
					"**Note:** At least one of `managed_device_ids` or `comanaged_device_ids` must be provided. " +
					"You can provide both to retrieve keys from different types of devices in one action.\n\n" +
					"**Security Warning:** Retrieved keys will be displayed in action output and may be stored in Terraform state.\n\n" +
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
				MarkdownDescription: "List of co-managed device IDs (GUIDs) to retrieve FileVault keys for. These are macOS devices " +
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
			"ignore_partial_failures": schema.BoolAttribute{
				Optional: true,
				MarkdownDescription: "If set to `true`, the action will succeed even if some operations fail. " +
					"Failed operations will be reported as warnings instead of errors. " +
					"Default: `false` (action fails if any operation fails).",
			},
			"validate_device_exists": schema.BoolAttribute{
				Optional: true,
				MarkdownDescription: "Whether to validate that devices exist and are macOS devices with FileVault enabled before attempting to retrieve keys. " +
					"Disabling this can speed up planning but may result in runtime errors for non-existent or unsupported devices. " +
					"Default: `true`.",
			},
			"timeouts": commonschema.ActionTimeouts(ctx),
		},
	}
}
