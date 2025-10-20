package graphBetaTriggerConfigurationManagerActionManagedDevice

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
	ActionName = "graph_beta_device_management_managed_device_trigger_configuration_manager_action"
)

var (
	_ action.Action                   = &TriggerConfigurationManagerActionManagedDeviceAction{}
	_ action.ActionWithConfigure      = &TriggerConfigurationManagerActionManagedDeviceAction{}
	_ action.ActionWithValidateConfig = &TriggerConfigurationManagerActionManagedDeviceAction{}
)

func NewTriggerConfigurationManagerActionManagedDeviceAction() action.Action {
	return &TriggerConfigurationManagerActionManagedDeviceAction{
		ReadPermissions: []string{
			"DeviceManagementManagedDevices.PrivilegedOperations.All",
		},
		WritePermissions: []string{
			"DeviceManagementManagedDevices.PrivilegedOperations.All",
		},
	}
}

type TriggerConfigurationManagerActionManagedDeviceAction struct {
	client           *msgraphbetasdk.GraphServiceClient
	ProviderTypeName string
	TypeName         string
	ReadPermissions  []string
	WritePermissions []string
}

func (a *TriggerConfigurationManagerActionManagedDeviceAction) Metadata(ctx context.Context, req action.MetadataRequest, resp *action.MetadataResponse) {
	a.ProviderTypeName = req.ProviderTypeName
	a.TypeName = ActionName
	resp.TypeName = a.FullTypeName()
}

func (a *TriggerConfigurationManagerActionManagedDeviceAction) FullTypeName() string {
	return a.ProviderTypeName + "_" + ActionName
}

func (a *TriggerConfigurationManagerActionManagedDeviceAction) Configure(ctx context.Context, req action.ConfigureRequest, resp *action.ConfigureResponse) {
	a.client = client.SetGraphBetaClientForAction(ctx, req, resp, constants.PROVIDER_NAME+"_"+ActionName)
}

func (a *TriggerConfigurationManagerActionManagedDeviceAction) Schema(ctx context.Context, req action.SchemaRequest, resp *action.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Triggers Configuration Manager client actions on Windows managed and co-managed devices using the " +
			"`/deviceManagement/managedDevices/{managedDeviceId}/triggerConfigurationManagerAction` and " +
			"`/deviceManagement/comanagedDevices/{managedDeviceId}/triggerConfigurationManagerAction` endpoints. " +
			"This action allows administrators to remotely invoke specific Configuration Manager (SCCM) operations on devices " +
			"that have the Configuration Manager client installed. This is particularly useful for co-managed devices where " +
			"Intune and Configuration Manager work together to manage devices. Actions include policy refresh, application " +
			"evaluation, antivirus scans, and more.\n\n" +
			"**Important Notes:**\n" +
			"- Requires Configuration Manager client installed on device\n" +
			"- Primarily used for co-managed devices (Intune + Configuration Manager)\n" +
			"- Device must be online to receive the action trigger\n" +
			"- Different actions available for different management scenarios\n" +
			"- Actions execute on the Configuration Manager client side\n\n" +
			"**Use Cases:**\n" +
			"- Force policy refresh after configuration changes\n" +
			"- Trigger application deployment evaluation\n" +
			"- Initiate antivirus scans remotely\n" +
			"- Wake up clients for scheduled operations\n" +
			"- Update Windows Defender signatures\n" +
			"- Synchronize device state with Configuration Manager\n\n" +
			"**Platform Support:**\n" +
			"- **Windows**: Fully supported (devices with Configuration Manager client)\n" +
			"- **Other Platforms**: Not supported (Configuration Manager is Windows-only)\n\n" +
			"**Reference:** [Microsoft Graph API - Trigger Configuration Manager Action](https://learn.microsoft.com/en-us/graph/api/intune-devices-manageddevice-triggerconfigurationmanageraction?view=graph-rest-beta)",
		Attributes: map[string]schema.Attribute{
			"timeouts": commonschema.Timeouts(ctx),
		},
		Blocks: map[string]schema.Block{
			"managed_devices": schema.ListNestedBlock{
				MarkdownDescription: "List of managed devices to trigger Configuration Manager actions on. These are Windows devices " +
					"fully managed by Intune that also have the Configuration Manager client installed.\n\n" +
					"**Note:** At least one of `managed_devices` or `comanaged_devices` must be provided.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"device_id": schema.StringAttribute{
							Required: true,
							MarkdownDescription: "The unique identifier (GUID) of the managed device to trigger the action on. " +
								"This must be a Windows device with Configuration Manager client installed.\n\n" +
								"Example: `\"12345678-1234-1234-1234-123456789abc\"`",
							Validators: []validator.String{
								stringvalidator.RegexMatches(
									regexp.MustCompile(constants.GuidRegex),
									"device_id must be a valid GUID format (e.g., 12345678-1234-1234-1234-123456789abc)",
								),
							},
						},
						"action": schema.StringAttribute{
							Required: true,
							MarkdownDescription: "The Configuration Manager action to trigger on this device.\n\n" +
								"Valid values:\n" +
								"- `\"refreshMachinePolicy\"`: Refresh the device's machine-level policies from Configuration Manager\n" +
								"- `\"refreshUserPolicy\"`: Refresh the current user's policies from Configuration Manager\n" +
								"- `\"wakeUpClient\"`: Wake up the Configuration Manager client for immediate activity\n" +
								"- `\"appEvaluation\"`: Trigger application deployment evaluation cycle\n" +
								"- `\"quickScan\"`: Initiate a quick antivirus scan using Windows Defender\n" +
								"- `\"fullScan\"`: Initiate a full antivirus scan using Windows Defender\n" +
								"- `\"windowsDefenderUpdateSignatures\"`: Update Windows Defender antivirus signatures\n\n" +
								"Example: `\"refreshMachinePolicy\"`",
							Validators: []validator.String{
								stringvalidator.OneOf(
									"refreshMachinePolicy",
									"refreshUserPolicy",
									"wakeUpClient",
									"appEvaluation",
									"quickScan",
									"fullScan",
									"windowsDefenderUpdateSignatures",
								),
							},
						},
					},
				},
			},
			"comanaged_devices": schema.ListNestedBlock{
				MarkdownDescription: "List of co-managed devices to trigger Configuration Manager actions on. These are Windows devices " +
					"managed by both Intune and Configuration Manager (SCCM). This is the most common scenario for this action.\n\n" +
					"**Note:** At least one of `managed_devices` or `comanaged_devices` must be provided.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"device_id": schema.StringAttribute{
							Required: true,
							MarkdownDescription: "The unique identifier (GUID) of the co-managed device to trigger the action on. " +
								"This must be a Windows device with Configuration Manager client installed.\n\n" +
								"Example: `\"abcdef12-3456-7890-abcd-ef1234567890\"`",
							Validators: []validator.String{
								stringvalidator.RegexMatches(
									regexp.MustCompile(constants.GuidRegex),
									"device_id must be a valid GUID format (e.g., 12345678-1234-1234-1234-123456789abc)",
								),
							},
						},
						"action": schema.StringAttribute{
							Required: true,
							MarkdownDescription: "The Configuration Manager action to trigger on this device.\n\n" +
								"Valid values:\n" +
								"- `\"refreshMachinePolicy\"`: Refresh the device's machine-level policies\n" +
								"- `\"refreshUserPolicy\"`: Refresh the current user's policies\n" +
								"- `\"wakeUpClient\"`: Wake up the Configuration Manager client\n" +
								"- `\"appEvaluation\"`: Trigger application deployment evaluation\n" +
								"- `\"quickScan\"`: Initiate a quick antivirus scan\n" +
								"- `\"fullScan\"`: Initiate a full antivirus scan\n" +
								"- `\"windowsDefenderUpdateSignatures\"`: Update Windows Defender signatures\n\n" +
								"Example: `\"appEvaluation\"`",
							Validators: []validator.String{
								stringvalidator.OneOf(
									"refreshMachinePolicy",
									"refreshUserPolicy",
									"wakeUpClient",
									"appEvaluation",
									"quickScan",
									"fullScan",
									"windowsDefenderUpdateSignatures",
								),
							},
						},
					},
				},
			},
		},
	}
}
