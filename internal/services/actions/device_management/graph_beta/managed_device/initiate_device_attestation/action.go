package graphBetaInitiateDeviceAttestationManagedDevice

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
	ActionName = "graph_beta_device_management_managed_device_initiate_device_attestation"
)

var (
	_ action.Action                   = &InitiateDeviceAttestationManagedDeviceAction{}
	_ action.ActionWithConfigure      = &InitiateDeviceAttestationManagedDeviceAction{}
	_ action.ActionWithValidateConfig = &InitiateDeviceAttestationManagedDeviceAction{}
)

func NewInitiateDeviceAttestationManagedDeviceAction() action.Action {
	return &InitiateDeviceAttestationManagedDeviceAction{
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

type InitiateDeviceAttestationManagedDeviceAction struct {
	client           *msgraphbetasdk.GraphServiceClient
	ProviderTypeName string
	TypeName         string
	ReadPermissions  []string
	WritePermissions []string
}

func (a *InitiateDeviceAttestationManagedDeviceAction) Metadata(ctx context.Context, req action.MetadataRequest, resp *action.MetadataResponse) {
	a.ProviderTypeName = req.ProviderTypeName
	a.TypeName = ActionName
	resp.TypeName = a.FullTypeName()
}

func (a *InitiateDeviceAttestationManagedDeviceAction) FullTypeName() string {
	return a.ProviderTypeName + "_" + ActionName
}

func (a *InitiateDeviceAttestationManagedDeviceAction) Configure(ctx context.Context, req action.ConfigureRequest, resp *action.ConfigureResponse) {
	a.client = client.SetGraphBetaClientForAction(ctx, req, resp, constants.PROVIDER_NAME+"_"+ActionName)
}

func (a *InitiateDeviceAttestationManagedDeviceAction) Schema(ctx context.Context, req action.SchemaRequest, resp *action.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Initiates device attestation on managed Windows devices using the " +
			"`/deviceManagement/managedDevices/{managedDeviceId}/initiateDeviceAttestation` and " +
			"`/deviceManagement/comanagedDevices/{managedDeviceId}/initiateDeviceAttestation` endpoints. " +
			"Device attestation is a security feature that uses the Trusted Platform Module (TPM) to cryptographically " +
			"verify the device's boot integrity, security configuration, and overall health status. This attestation " +
			"process creates a trusted baseline that can be used for conditional access, compliance policies, and zero-trust " +
			"security models. The TPM provides hardware-rooted proof that the device has not been tampered with and is in a known good state.\n\n" +
			"**Important Notes:**\n" +
			"- Only works on Windows devices with TPM 1.2 or 2.0\n" +
			"- Performs cryptographic verification of device health\n" +
			"- Creates attestation report for compliance validation\n" +
			"- Does not affect device operation or user access\n" +
			"- Results stored in Intune for policy enforcement\n" +
			"- Critical for Zero Trust security architecture\n" +
			"- Should be performed periodically for compliance\n\n" +
			"**Use Cases:**\n" +
			"- Conditional access policy enforcement\n" +
			"- Compliance validation for security standards\n" +
			"- Zero Trust security model implementation\n" +
			"- Periodic device health verification\n" +
			"- Pre-deployment security validation\n" +
			"- Post-incident device integrity checks\n\n" +
			"**Platform Support:**\n" +
			"- **Windows**: Devices with TPM 1.2/2.0 and secure boot\n" +
			"- **Other Platforms**: Not supported (Windows-specific feature)\n\n" +
			"**Reference:** [Microsoft Graph API - Initiate Device Attestation](https://learn.microsoft.com/en-us/graph/api/intune-devices-manageddevice-initiatedeviceattestation?view=graph-rest-beta)",
		Attributes: map[string]schema.Attribute{
			"managed_device_ids": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
				MarkdownDescription: "List of managed device IDs (GUIDs) to initiate device attestation for. " +
					"These are devices fully managed by Intune.\n\n" +
					"**Note:** At least one of `managed_device_ids` or `comanaged_device_ids` must be provided. " +
					"You can provide both to initiate attestation on different types of devices in one action.\n\n" +
					"**Important:** This action uses the TPM to cryptographically verify device health and security state. " +
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
				MarkdownDescription: "List of co-managed device IDs (GUIDs) to initiate device attestation for. " +
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

