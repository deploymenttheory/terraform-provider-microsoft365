package graphBetaPauseConfigurationRefreshManagedDevice

import (
	"context"
	"regexp"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/action/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

const (
	ActionName = "graph_beta_device_management_managed_device_pause_configuration_refresh"
)

var (
	_ action.Action                   = &PauseConfigurationRefreshManagedDeviceAction{}
	_ action.ActionWithConfigure      = &PauseConfigurationRefreshManagedDeviceAction{}
	_ action.ActionWithValidateConfig = &PauseConfigurationRefreshManagedDeviceAction{}
)

func NewPauseConfigurationRefreshManagedDeviceAction() action.Action {
	return &PauseConfigurationRefreshManagedDeviceAction{
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

type PauseConfigurationRefreshManagedDeviceAction struct {
	client           *msgraphbetasdk.GraphServiceClient
	ProviderTypeName string
	TypeName         string
	ReadPermissions  []string
	WritePermissions []string
}

func (a *PauseConfigurationRefreshManagedDeviceAction) Metadata(ctx context.Context, req action.MetadataRequest, resp *action.MetadataResponse) {
	a.ProviderTypeName = req.ProviderTypeName
	a.TypeName = ActionName
	resp.TypeName = a.FullTypeName()
}

func (a *PauseConfigurationRefreshManagedDeviceAction) FullTypeName() string {
	return a.ProviderTypeName + "_" + ActionName
}

func (a *PauseConfigurationRefreshManagedDeviceAction) Configure(ctx context.Context, req action.ConfigureRequest, resp *action.ConfigureResponse) {
	a.client = client.SetGraphBetaClientForAction(ctx, req, resp, constants.PROVIDER_NAME+"_"+ActionName)
}

func (a *PauseConfigurationRefreshManagedDeviceAction) Schema(ctx context.Context, req action.SchemaRequest, resp *action.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Initiates a command to pause configuration refresh on managed Windows devices using the " +
			"`/deviceManagement/managedDevices/{managedDeviceId}/pauseConfigurationRefresh` and " +
			"`/deviceManagement/comanagedDevices/{managedDeviceId}/pauseConfigurationRefresh` endpoints. " +
			"This action temporarily prevents devices from receiving and applying new configuration policies from Intune, " +
			"which is useful during maintenance windows, troubleshooting, or when you need to prevent policy changes from " +
			"being applied to specific devices for a defined period.\n\n" +
			"**Important Notes:**\n" +
			"- Only works on Windows 10/11 devices\n" +
			"- Configuration refresh automatically resumes after the pause period expires\n" +
			"- Maximum pause period is typically 24 hours (1440 minutes)\n" +
			"- Does not affect existing applied policies, only prevents new policy updates\n" +
			"- Device can still check in and report status\n" +
			"- Critical security updates may still be applied\n" +
			"- User can still manually sync from Company Portal\n\n" +
			"**Use Cases:**\n" +
			"- Maintenance windows for critical applications\n" +
			"- Troubleshooting policy conflicts\n" +
			"- Testing policy changes in staging\n" +
			"- Preventing policy updates during business-critical operations\n" +
			"- Temporary freeze during incident response\n" +
			"- User acceptance testing (UAT) phases\n\n" +
			"**Platform Support:**\n" +
			"- **Windows**: Windows 10/11\n" +
			"- **Other Platforms**: Not supported\n\n" +
			"**Reference:** [Microsoft Graph API - Pause Configuration Refresh](https://learn.microsoft.com/en-us/graph/api/intune-devices-manageddevice-pauseconfigurationrefresh?view=graph-rest-beta)",
		Attributes: map[string]schema.Attribute{
			"timeouts": commonschema.Timeouts(ctx),
		},
		Blocks: map[string]schema.Block{
			"managed_devices": schema.ListNestedBlock{
				MarkdownDescription: "List of managed devices to pause configuration refresh for. " +
					"Each device can have a different pause duration based on specific requirements.\n\n" +
					"**Note:** At least one of `managed_devices` or `comanaged_devices` must be provided.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"device_id": schema.StringAttribute{
							Required: true,
							MarkdownDescription: "The unique identifier (GUID) of the managed device to pause configuration refresh for. " +
								"This should be a Windows 10 or Windows 11 device managed by Intune.\n\n" +
								"Example: `\"12345678-1234-1234-1234-123456789abc\"`",
							Validators: []validator.String{
								stringvalidator.RegexMatches(
									regexp.MustCompile(constants.GuidRegex),
									"device_id must be a valid GUID format (e.g., 12345678-1234-1234-1234-123456789abc)",
								),
							},
						},
						"pause_time_period_in_minutes": schema.Int64Attribute{
							Required: true,
							MarkdownDescription: "The duration in minutes to pause configuration refresh for this device. " +
								"Configuration refresh will automatically resume after this period expires.\n\n" +
								"**Valid Range:** 1 to 1440 minutes (1 minute to 24 hours)\n\n" +
								"**Common Values:**\n" +
								"- `60` - 1 hour (short maintenance)\n" +
								"- `120` - 2 hours (application updates)\n" +
								"- `240` - 4 hours (extended maintenance)\n" +
								"- `480` - 8 hours (business day)\n" +
								"- `1440` - 24 hours (full day)\n\n" +
								"Example: `120` (2 hours)",
							Validators: []validator.Int64{
								int64validator.Between(1, 1440),
							},
						},
					},
				},
			},
			"comanaged_devices": schema.ListNestedBlock{
				MarkdownDescription: "List of co-managed devices to pause configuration refresh for. " +
					"These are devices managed by both Intune and Configuration Manager (SCCM).\n\n" +
					"**Note:** At least one of `managed_devices` or `comanaged_devices` must be provided.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"device_id": schema.StringAttribute{
							Required: true,
							MarkdownDescription: "The unique identifier (GUID) of the co-managed device to pause configuration refresh for.\n\n" +
								"Example: `\"abcdef12-3456-7890-abcd-ef1234567890\"`",
							Validators: []validator.String{
								stringvalidator.RegexMatches(
									regexp.MustCompile(constants.GuidRegex),
									"device_id must be a valid GUID format (e.g., 12345678-1234-1234-1234-123456789abc)",
								),
							},
						},
						"pause_time_period_in_minutes": schema.Int64Attribute{
							Required: true,
							MarkdownDescription: "The duration in minutes to pause configuration refresh for this device.\n\n" +
								"**Valid Range:** 1 to 1440 minutes (1 minute to 24 hours)",
							Validators: []validator.Int64{
								int64validator.Between(1, 1440),
							},
						},
					},
				},
			},
		},
	}
}

