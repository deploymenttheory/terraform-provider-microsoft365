package graphBetaDeprovisionManagedDevice

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
	ActionName = "microsoft365_graph_beta_device_management_managed_device_deprovision"
	InvokeTimeout = 60
)

var (
	_ action.Action                   = &DeprovisionManagedDeviceAction{}
	_ action.ActionWithConfigure      = &DeprovisionManagedDeviceAction{}
	_ action.ActionWithValidateConfig = &DeprovisionManagedDeviceAction{}
)

func NewDeprovisionManagedDeviceAction() action.Action {
	return &DeprovisionManagedDeviceAction{
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

type DeprovisionManagedDeviceAction struct {
	client           *msgraphbetasdk.GraphServiceClient
	ReadPermissions  []string
	WritePermissions []string
}

func (a *DeprovisionManagedDeviceAction) Metadata(ctx context.Context, req action.MetadataRequest, resp *action.MetadataResponse) {
	resp.TypeName = ActionName
}

func (a *DeprovisionManagedDeviceAction) Configure(ctx context.Context, req action.ConfigureRequest, resp *action.ConfigureResponse) {
	a.client = client.SetGraphBetaClientForAction(ctx, req, resp, ActionName)
}

func (a *DeprovisionManagedDeviceAction) Schema(ctx context.Context, req action.SchemaRequest, resp *action.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Deprovisions Windows managed devices from Intune management using the " +
			"`/deviceManagement/managedDevices/{managedDeviceId}/deprovision` and " +
			"`/deviceManagement/comanagedDevices/{managedDeviceId}/deprovision` endpoints. " +
			"This action removes management capabilities from a device while allowing it to remain enrolled. " +
			"Deprovisioning is less destructive than wiping or retiring a device, as it only removes management " +
			"policies and profiles without deleting user data or removing the device entirely. This is useful when " +
			"transitioning devices between management solutions or preparing devices for different management scenarios.\n\n" +
			"**Important Notes:**\n" +
			"- Device remains enrolled in Intune after deprovisioning\n" +
			"- Management policies and profiles are removed\n" +
			"- User data is preserved on the device\n" +
			"- Less disruptive than wipe or retire actions\n" +
			"- Requires a reason to be specified for auditing\n" +
			"- Primarily used for Windows devices\n\n" +
			"**Use Cases:**\n" +
			"- Transitioning devices to different management authority\n" +
			"- Preparing devices for repurposing\n" +
			"- Removing management overhead without data loss\n" +
			"- Troubleshooting management issues\n" +
			"- Moving from co-management to different configuration\n\n" +
			"**Platform Support:**\n" +
			"- **Windows**: Primary platform for deprovisioning\n" +
			"- **Other Platforms**: Limited or no support\n\n" +
			"**Reference:** [Microsoft Graph API - Deprovision](https://learn.microsoft.com/en-us/graph/api/intune-devices-manageddevice-deprovision?view=graph-rest-beta)",
		Attributes: map[string]schema.Attribute{
			"managed_devices": schema.ListNestedAttribute{
				Optional: true,
				MarkdownDescription: "List of managed devices to deprovision. These are Windows devices " +
					"fully managed by Intune only.\n\n" +
					"**Note:** At least one of `managed_devices` or `comanaged_devices` must be provided.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"device_id": schema.StringAttribute{
							Required: true,
							MarkdownDescription: "The unique identifier (GUID) of the managed device to deprovision.\n\n" +
								"Example: `\"12345678-1234-1234-1234-123456789abc\"`",
							Validators: []validator.String{
								stringvalidator.RegexMatches(
									regexp.MustCompile(constants.GuidRegex),
									"device_id must be a valid GUID format (e.g., 12345678-1234-1234-1234-123456789abc)",
								),
							},
						},
						"deprovision_reason": schema.StringAttribute{
							Required: true,
							MarkdownDescription: "The reason for deprovisioning this device. This is required for auditing and tracking purposes.\n\n" +
								"Examples:\n" +
								"- `\"Transitioning to new management solution\"`\n" +
								"- `\"Device repurposing\"`\n" +
								"- `\"Troubleshooting management issues\"`\n" +
								"- `\"User requested management removal\"`\n" +
								"- `\"Moving to co-management\"`",
							Validators: []validator.String{
								stringvalidator.LengthAtLeast(1),
							},
						},
					},
				},
			},
			"comanaged_devices": schema.ListNestedAttribute{
				Optional: true,
				MarkdownDescription: "List of co-managed devices to deprovision. These are Windows devices " +
					"managed by both Intune and Configuration Manager (SCCM).\n\n" +
					"**Note:** At least one of `managed_devices` or `comanaged_devices` must be provided.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"device_id": schema.StringAttribute{
							Required: true,
							MarkdownDescription: "The unique identifier (GUID) of the co-managed device to deprovision.\n\n" +
								"Example: `\"abcdef12-3456-7890-abcd-ef1234567890\"`",
							Validators: []validator.String{
								stringvalidator.RegexMatches(
									regexp.MustCompile(constants.GuidRegex),
									"device_id must be a valid GUID format (e.g., 12345678-1234-1234-1234-123456789abc)",
								),
							},
						},
						"deprovision_reason": schema.StringAttribute{
							Required: true,
							MarkdownDescription: "The reason for deprovisioning this device. This is required for auditing and tracking purposes.\n\n" +
								"Examples:\n" +
								"- `\"Transitioning from co-management\"`\n" +
								"- `\"Device repurposing\"`\n" +
								"- `\"Management authority change\"`",
							Validators: []validator.String{
								stringvalidator.LengthAtLeast(1),
							},
						},
					},
				},
			},
			"timeouts": commonschema.ActionTimeouts(ctx),
		},
	}
}
