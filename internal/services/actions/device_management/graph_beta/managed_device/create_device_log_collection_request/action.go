package graphBetaCreateDeviceLogCollectionRequestManagedDevice

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
	ActionName = "microsoft365_graph_beta_device_management_managed_device_create_device_log_collection_request"
)

var (
	_ action.Action                   = &CreateDeviceLogCollectionRequestManagedDeviceAction{}
	_ action.ActionWithConfigure      = &CreateDeviceLogCollectionRequestManagedDeviceAction{}
	_ action.ActionWithValidateConfig = &CreateDeviceLogCollectionRequestManagedDeviceAction{}
)

func NewCreateDeviceLogCollectionRequestManagedDeviceAction() action.Action {
	return &CreateDeviceLogCollectionRequestManagedDeviceAction{
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

type CreateDeviceLogCollectionRequestManagedDeviceAction struct {
	client           *msgraphbetasdk.GraphServiceClient
	ReadPermissions  []string
	WritePermissions []string
}

func (a *CreateDeviceLogCollectionRequestManagedDeviceAction) Metadata(ctx context.Context, req action.MetadataRequest, resp *action.MetadataResponse) {
	resp.TypeName = ActionName
}

func (a *CreateDeviceLogCollectionRequestManagedDeviceAction) Configure(ctx context.Context, req action.ConfigureRequest, resp *action.ConfigureResponse) {
	a.client = client.SetGraphBetaClientForAction(ctx, req, resp, ActionName)
}

func (a *CreateDeviceLogCollectionRequestManagedDeviceAction) Schema(ctx context.Context, req action.SchemaRequest, resp *action.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Creates a device log collection request for Windows managed devices using the " +
			"`/deviceManagement/managedDevices/{managedDeviceId}/createDeviceLogCollectionRequest` and " +
			"`/deviceManagement/comanagedDevices/{managedDeviceId}/createDeviceLogCollectionRequest` endpoints. " +
			"This action initiates the collection of diagnostic logs from Windows devices, which are essential for " +
			"troubleshooting device issues, analyzing compliance problems, and supporting technical investigations. " +
			"The collected logs are uploaded to Intune and can be downloaded for analysis. This action is critical " +
			"for IT support teams when diagnosing device-specific problems or investigating security incidents.\n\n" +
			"**Important Notes:**\n" +
			"- Only applicable to Windows devices (Windows 10/11)\n" +
			"- Device must be online to receive collection request\n" +
			"- Log collection runs on the device and uploads results\n" +
			"- Logs are available in Intune portal after collection completes\n" +
			"- Collection includes system logs, event logs, and diagnostic data\n" +
			"- Log files have expiration dates for security\n\n" +
			"**Use Cases:**\n" +
			"- Troubleshooting device configuration issues\n" +
			"- Investigating compliance failures or policy problems\n" +
			"- Supporting help desk tickets requiring detailed diagnostics\n" +
			"- Analyzing application deployment failures\n" +
			"- Security incident investigation and forensics\n" +
			"- Proactive monitoring and preventive maintenance\n\n" +
			"**Platform Support:**\n" +
			"- **Windows**: Fully supported (Windows 10 version 1709 or later, Windows 11)\n" +
			"- **Other Platforms**: Not supported (macOS, iOS/iPadOS, Android use different logging mechanisms)\n\n" +
			"**Reference:** [Microsoft Graph API - Create Device Log Collection Request](https://learn.microsoft.com/en-us/graph/api/intune-devices-manageddevice-createdevicelogcollectionrequest?view=graph-rest-beta)",
		Attributes: map[string]schema.Attribute{
			"timeouts": commonschema.Timeouts(ctx),
		},
		Blocks: map[string]schema.Block{
			"managed_devices": schema.ListNestedBlock{
				MarkdownDescription: "List of managed devices to collect logs from. These are Windows devices " +
					"fully managed by Intune only. Each device can have its own template type configuration.\n\n" +
					"**Note:** At least one of `managed_devices` or `comanaged_devices` must be provided.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"device_id": schema.StringAttribute{
							Required: true,
							MarkdownDescription: "The unique identifier (GUID) of the managed device to collect logs from. " +
								"This must be a Windows device running Windows 10 version 1709 or later, or Windows 11.\n\n" +
								"Example: `\"12345678-1234-1234-1234-123456789abc\"`",
							Validators: []validator.String{
								stringvalidator.RegexMatches(
									regexp.MustCompile(constants.GuidRegex),
									"device_id must be a valid GUID format (e.g., 12345678-1234-1234-1234-123456789abc)",
								),
							},
						},
						"template_type": schema.StringAttribute{
							Optional: true,
							MarkdownDescription: "The template type for the log collection. Determines the scope and type of logs collected.\n\n" +
								"Valid values:\n" +
								"- `\"predefined\"` (default): Uses the standard predefined log collection template that includes common system and diagnostic logs\n" +
								"- `\"unknownFutureValue\"`: Reserved for future expansion\n\n" +
								"If not specified, defaults to `\"predefined\"`.",
							Validators: []validator.String{
								stringvalidator.OneOf("predefined", "unknownFutureValue"),
							},
						},
					},
				},
			},
			"comanaged_devices": schema.ListNestedBlock{
				MarkdownDescription: "List of co-managed devices to collect logs from. These are Windows devices " +
					"managed by both Intune and Configuration Manager (SCCM).\n\n" +
					"**Note:** At least one of `managed_devices` or `comanaged_devices` must be provided.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"device_id": schema.StringAttribute{
							Required: true,
							MarkdownDescription: "The unique identifier (GUID) of the co-managed device to collect logs from. " +
								"This must be a Windows device running Windows 10 version 1709 or later, or Windows 11.\n\n" +
								"Example: `\"abcdef12-3456-7890-abcd-ef1234567890\"`",
							Validators: []validator.String{
								stringvalidator.RegexMatches(
									regexp.MustCompile(constants.GuidRegex),
									"device_id must be a valid GUID format (e.g., 12345678-1234-1234-1234-123456789abc)",
								),
							},
						},
						"template_type": schema.StringAttribute{
							Optional: true,
							MarkdownDescription: "The template type for the log collection. Determines the scope and type of logs collected.\n\n" +
								"Valid values:\n" +
								"- `\"predefined\"` (default): Uses the standard predefined log collection template\n" +
								"- `\"unknownFutureValue\"`: Reserved for future expansion\n\n" +
								"If not specified, defaults to `\"predefined\"`.",
							Validators: []validator.String{
								stringvalidator.OneOf("predefined", "unknownFutureValue"),
							},
						},
					},
				},
			},
		},
	}
}
