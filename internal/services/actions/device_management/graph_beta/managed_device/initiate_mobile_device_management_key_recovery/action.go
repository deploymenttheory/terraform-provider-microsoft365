package graphBetaInitiateMobileDeviceManagementKeyRecoveryManagedDevice

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
	ActionName = "graph_beta_device_management_managed_device_initiate_mobile_device_management_key_recovery"
)

var (
	_ action.Action                   = &InitiateMobileDeviceManagementKeyRecoveryManagedDeviceAction{}
	_ action.ActionWithConfigure      = &InitiateMobileDeviceManagementKeyRecoveryManagedDeviceAction{}
	_ action.ActionWithValidateConfig = &InitiateMobileDeviceManagementKeyRecoveryManagedDeviceAction{}
)

func NewInitiateMobileDeviceManagementKeyRecoveryManagedDeviceAction() action.Action {
	return &InitiateMobileDeviceManagementKeyRecoveryManagedDeviceAction{
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

type InitiateMobileDeviceManagementKeyRecoveryManagedDeviceAction struct {
	client           *msgraphbetasdk.GraphServiceClient
	ProviderTypeName string
	TypeName         string
	ReadPermissions  []string
	WritePermissions []string
}

func (a *InitiateMobileDeviceManagementKeyRecoveryManagedDeviceAction) Metadata(ctx context.Context, req action.MetadataRequest, resp *action.MetadataResponse) {
	a.ProviderTypeName = req.ProviderTypeName
	a.TypeName = ActionName
	resp.TypeName = a.FullTypeName()
}

func (a *InitiateMobileDeviceManagementKeyRecoveryManagedDeviceAction) FullTypeName() string {
	return a.ProviderTypeName + "_" + ActionName
}

func (a *InitiateMobileDeviceManagementKeyRecoveryManagedDeviceAction) Configure(ctx context.Context, req action.ConfigureRequest, resp *action.ConfigureResponse) {
	a.client = client.SetGraphBetaClientForAction(ctx, req, resp, constants.PROVIDER_NAME+"_"+ActionName)
}

func (a *InitiateMobileDeviceManagementKeyRecoveryManagedDeviceAction) Schema(ctx context.Context, req action.SchemaRequest, resp *action.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Initiates Mobile Device Management (MDM) key recovery and TPM attestation on managed Windows devices using the " +
			"`/deviceManagement/managedDevices/{managedDeviceId}/initiateMobileDeviceManagementKeyRecovery` and " +
			"`/deviceManagement/comanagedDevices/{managedDeviceId}/initiateMobileDeviceManagementKeyRecovery` endpoints. " +
			"This action performs BitLocker recovery key escrow and Trusted Platform Module (TPM) attestation to ensure recovery keys " +
			"are properly stored in Azure AD and the device's TPM is healthy. This is critical for security compliance, data recovery " +
			"scenarios, and ensuring encrypted devices can be recovered if users forget passwords or encounter hardware issues.\n\n" +
			"**Important Notes:**\n" +
			"- Only works on Windows devices with BitLocker and TPM enabled\n" +
			"- Escrows BitLocker recovery keys to Azure AD\n" +
			"- Performs TPM health attestation\n" +
			"- Does not encrypt/decrypt the device\n" +
			"- Does not affect device operation or user access\n" +
			"- Essential for compliance and disaster recovery\n" +
			"- Should be run periodically for key rotation\n\n" +
			"**Use Cases:**\n" +
			"- Ensuring BitLocker recovery keys are escrowed\n" +
			"- Compliance auditing for encryption key management\n" +
			"- Verifying TPM attestation and health\n" +
			"- Periodic key rotation and refresh\n" +
			"- Pre-deployment validation for new devices\n" +
			"- Recovery preparation for critical devices\n\n" +
			"**Platform Support:**\n" +
			"- **Windows**: Devices with BitLocker and TPM 1.2/2.0\n" +
			"- **Other Platforms**: Not supported (Windows-specific feature)\n\n" +
			"**Reference:** [Microsoft Graph API - Initiate Mobile Device Management Key Recovery](https://learn.microsoft.com/en-us/graph/api/intune-devices-manageddevice-initiatemobiledevicemanagementkeyrecovery?view=graph-rest-beta)",
		Attributes: map[string]schema.Attribute{
			"managed_device_ids": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
				MarkdownDescription: "List of managed device IDs (GUIDs) to initiate MDM key recovery and TPM attestation for. " +
					"These are devices fully managed by Intune.\n\n" +
					"**Note:** At least one of `managed_device_ids` or `comanaged_device_ids` must be provided. " +
					"You can provide both to initiate key recovery on different types of devices in one action.\n\n" +
					"**Important:** This action escrows BitLocker recovery keys to Azure AD and performs TPM attestation. " +
					"It does not affect device operation or user access.\n\n" +
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
				MarkdownDescription: "List of co-managed device IDs (GUIDs) to initiate MDM key recovery and TPM attestation for. " +
					"These are devices managed by both Intune and Configuration Manager (SCCM).\n\n" +
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

