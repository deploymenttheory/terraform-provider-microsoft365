package graphBetaDisableManagedDevice

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
	ActionName = "microsoft365_graph_beta_device_management_managed_device_disable"
)

var (
	_ action.Action                   = &DisableManagedDeviceAction{}
	_ action.ActionWithConfigure      = &DisableManagedDeviceAction{}
	_ action.ActionWithValidateConfig = &DisableManagedDeviceAction{}
)

func NewDisableManagedDeviceAction() action.Action {
	return &DisableManagedDeviceAction{
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

type DisableManagedDeviceAction struct {
	client           *msgraphbetasdk.GraphServiceClient
	ReadPermissions  []string
	WritePermissions []string
}

func (a *DisableManagedDeviceAction) Metadata(ctx context.Context, req action.MetadataRequest, resp *action.MetadataResponse) {
	resp.TypeName = ActionName
}

func (a *DisableManagedDeviceAction) Configure(ctx context.Context, req action.ConfigureRequest, resp *action.ConfigureResponse) {
	a.client = client.SetGraphBetaClientForAction(ctx, req, resp, ActionName)
}

func (a *DisableManagedDeviceAction) Schema(ctx context.Context, req action.SchemaRequest, resp *action.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Disables managed devices from Intune management using the " +
			"`/deviceManagement/managedDevices/{managedDeviceId}/disable` and " +
			"`/deviceManagement/comanagedDevices/{managedDeviceId}/disable` endpoints. " +
			"This action disables a device's ability to interact with Intune services while maintaining its enrollment record. " +
			"Disabled devices cannot receive policies, sync with Intune, or perform managed operations until re-enabled. " +
			"This is useful for temporarily suspending device management without fully removing the device from Intune, " +
			"such as during investigations, compliance violations, or security incidents.\n\n" +
			"**Important Notes:**\n" +
			"- Device remains enrolled but cannot sync or receive policies\n" +
			"- Management operations are suspended\n" +
			"- Device can be re-enabled later\n" +
			"- Less permanent than retire or wipe\n" +
			"- Useful for temporary suspensions\n" +
			"- Security and compliance enforcement\n\n" +
			"**Use Cases:**\n" +
			"- Security incident response (suspected compromise)\n" +
			"- Compliance violations requiring device suspension\n" +
			"- Temporary device quarantine\n" +
			"- Investigation of device issues\n" +
			"- Preventing policy application during troubleshooting\n" +
			"- Temporary management suspension\n\n" +
			"**Platform Support:**\n" +
			"- **All Platforms**: Windows, macOS, iOS/iPadOS, Android\n\n" +
			"**Reference:** [Microsoft Graph API - Disable](https://learn.microsoft.com/en-us/graph/api/intune-devices-manageddevice-disable?view=graph-rest-beta)",
		Attributes: map[string]schema.Attribute{
			"managed_device_ids": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
				MarkdownDescription: "List of managed device IDs (GUIDs) to disable. These are devices " +
					"fully managed by Intune only.\n\n" +
					"**Note:** At least one of `managed_device_ids` or `comanaged_device_ids` must be provided. " +
					"You can provide both to disable different types of devices in one action.\n\n" +
					"**Important:** Disabled devices will not be able to sync with Intune or receive policy updates " +
					"until they are re-enabled.\n\n" +
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
				MarkdownDescription: "List of co-managed device IDs (GUIDs) to disable. These are devices " +
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
