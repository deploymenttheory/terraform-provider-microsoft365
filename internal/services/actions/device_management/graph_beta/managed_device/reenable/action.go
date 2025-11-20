package graphBetaReenableManagedDevice

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
	ActionName = "microsoft365_graph_beta_device_management_managed_device_reenable"
)

var (
	_ action.Action                   = &ReenableManagedDeviceAction{}
	_ action.ActionWithConfigure      = &ReenableManagedDeviceAction{}
	_ action.ActionWithValidateConfig = &ReenableManagedDeviceAction{}
)

func NewReenableManagedDeviceAction() action.Action {
	return &ReenableManagedDeviceAction{
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

type ReenableManagedDeviceAction struct {
	client           *msgraphbetasdk.GraphServiceClient
	ReadPermissions  []string
	WritePermissions []string
}

func (a *ReenableManagedDeviceAction) Metadata(ctx context.Context, req action.MetadataRequest, resp *action.MetadataResponse) {
	resp.TypeName = ActionName
}

func (a *ReenableManagedDeviceAction) Configure(ctx context.Context, req action.ConfigureRequest, resp *action.ConfigureResponse) {
	a.client = client.SetGraphBetaClientForAction(ctx, req, resp, ActionName)
}

func (a *ReenableManagedDeviceAction) Schema(ctx context.Context, req action.SchemaRequest, resp *action.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Re-enables previously disabled managed devices in Intune using the " +
			"`/deviceManagement/managedDevices/{managedDeviceId}/reenable` and " +
			"`/deviceManagement/comanagedDevices/{managedDeviceId}/reenable` endpoints. " +
			"This action restores a disabled device's ability to interact with Intune services, allowing it to sync and " +
			"receive policy updates again. Re-enabling is the counterpart to the disable action and restores full management " +
			"capabilities to devices that were temporarily suspended. This is useful after resolving security incidents, " +
			"compliance violations, or completing investigations that required temporary device suspension.\n\n" +
			"**Important Notes:**\n" +
			"- Only works on previously disabled devices\n" +
			"- Restores sync capability with Intune\n" +
			"- Re-enables policy application\n" +
			"- Maintains existing enrollment\n" +
			"- Reverses the disable action\n" +
			"- All platforms supported\n\n" +
			"**Use Cases:**\n" +
			"- Restoring devices after security investigation completion\n" +
			"- Re-enabling compliant devices after violations resolved\n" +
			"- Ending temporary quarantine period\n" +
			"- Resuming management after troubleshooting\n" +
			"- Restoring devices after policy fixes\n" +
			"- Completing incident response procedures\n\n" +
			"**Platform Support:**\n" +
			"- **All Platforms**: Windows, macOS, iOS/iPadOS, Android\n\n" +
			"**Reference:** [Microsoft Graph API - Reenable](https://learn.microsoft.com/en-us/graph/api/intune-devices-manageddevice-reenable?view=graph-rest-beta)",
		Attributes: map[string]schema.Attribute{
			"managed_device_ids": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
				MarkdownDescription: "List of managed device IDs (GUIDs) to re-enable. These are devices " +
					"fully managed by Intune that were previously disabled.\n\n" +
					"**Note:** At least one of `managed_device_ids` or `comanaged_device_ids` must be provided. " +
					"You can provide both to re-enable different types of devices in one action.\n\n" +
					"**Important:** Re-enabled devices will be able to sync with Intune and receive policy updates again.\n\n" +
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
				MarkdownDescription: "List of co-managed device IDs (GUIDs) to re-enable. These are devices " +
					"managed by both Intune and Configuration Manager (SCCM) that were previously disabled.\n\n" +
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
