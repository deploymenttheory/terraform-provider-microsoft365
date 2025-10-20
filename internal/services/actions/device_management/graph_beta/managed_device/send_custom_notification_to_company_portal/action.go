package graphBetaSendCustomNotificationToCompanyPortal

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
	ActionName = "graph_beta_device_management_managed_device_send_custom_notification_to_company_portal"
)

var (
	_ action.Action                   = &SendCustomNotificationToCompanyPortalAction{}
	_ action.ActionWithConfigure      = &SendCustomNotificationToCompanyPortalAction{}
	_ action.ActionWithValidateConfig = &SendCustomNotificationToCompanyPortalAction{}
)

func NewSendCustomNotificationToCompanyPortalAction() action.Action {
	return &SendCustomNotificationToCompanyPortalAction{
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

type SendCustomNotificationToCompanyPortalAction struct {
	client           *msgraphbetasdk.GraphServiceClient
	ProviderTypeName string
	TypeName         string
	ReadPermissions  []string
	WritePermissions []string
}

func (a *SendCustomNotificationToCompanyPortalAction) Metadata(ctx context.Context, req action.MetadataRequest, resp *action.MetadataResponse) {
	a.ProviderTypeName = req.ProviderTypeName
	a.TypeName = ActionName
	resp.TypeName = a.FullTypeName()
}

func (a *SendCustomNotificationToCompanyPortalAction) FullTypeName() string {
	return a.ProviderTypeName + "_" + ActionName
}

func (a *SendCustomNotificationToCompanyPortalAction) Configure(ctx context.Context, req action.ConfigureRequest, resp *action.ConfigureResponse) {
	a.client = client.SetGraphBetaClientForAction(ctx, req, resp, constants.PROVIDER_NAME+"_"+ActionName)
}

func (a *SendCustomNotificationToCompanyPortalAction) Schema(ctx context.Context, req action.SchemaRequest, resp *action.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Sends custom notifications to the Company Portal app on managed devices using the " +
			"`/deviceManagement/managedDevices/{managedDeviceId}/sendCustomNotificationToCompanyPortal` and " +
			"`/deviceManagement/comanagedDevices/{managedDeviceId}/sendCustomNotificationToCompanyPortal` endpoints. " +
			"This action enables IT administrators to send targeted messages to end users through the Company Portal app.\n\n" +
			"**What This Action Does:**\n" +
			"- Sends push notification to Company Portal app\n" +
			"- Displays custom message title and body\n" +
			"- Targets specific devices or users\n" +
			"- Supports customized messages per device\n" +
			"- Provides in-app notification visibility\n" +
			"- Enables two-way communication channel\n\n" +
			"**When to Use:**\n" +
			"- Compliance reminders and deadlines\n" +
			"- Security alert communications\n" +
			"- Policy update notifications\n" +
			"- Maintenance window announcements\n" +
			"- Action required messages\n" +
			"- User guidance and instructions\n" +
			"- Incident response communications\n\n" +
			"**Platform Support:**\n" +
			"- **Windows**: Company Portal app required\n" +
			"- **iOS/iPadOS**: Company Portal app required\n" +
			"- **Android**: Company Portal app required\n" +
			"- **macOS**: Company Portal app required\n\n" +
			"**Important Considerations:**\n" +
			"- Company Portal app must be installed\n" +
			"- Device must be enrolled in Intune\n" +
			"- User must be signed into Company Portal\n" +
			"- Device must have network connectivity\n" +
			"- Notifications appear in Company Portal app\n" +
			"- Consider user time zones for timing\n\n" +
			"**Reference:** [Microsoft Graph API - Send Custom Notification To Company Portal](https://learn.microsoft.com/en-us/graph/api/intune-devices-manageddevice-sendcustomnotificationtocompanyportal?view=graph-rest-beta)",
		Attributes: map[string]schema.Attribute{
			"timeouts": commonschema.Timeouts(ctx),
		},
		Blocks: map[string]schema.Block{
			"managed_devices": schema.ListNestedBlock{
				MarkdownDescription: "List of managed devices to send custom notifications to. These are devices fully managed by Intune only. " +
					"Each entry specifies a device ID and the custom notification title and body for that device.\n\n" +
					"**Note:** At least one of `managed_devices` or `comanaged_devices` must be provided. " +
					"You can provide both to send notifications to different types of devices in one action.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"device_id": schema.StringAttribute{
							Required: true,
							MarkdownDescription: "The unique identifier (GUID) of the managed device to send the notification to. " +
								"Example: `\"12345678-1234-1234-1234-123456789abc\"`",
							Validators: []validator.String{
								stringvalidator.RegexMatches(
									regexp.MustCompile(constants.GuidRegex),
									"device_id must be a valid GUID format (e.g., 12345678-1234-1234-1234-123456789abc)",
								),
							},
						},
						"notification_title": schema.StringAttribute{
							Required: true,
							MarkdownDescription: "The title of the custom notification to display in the Company Portal app. " +
								"Should be concise and descriptive. Maximum recommended length: 50-60 characters. " +
								"Example: `\"Action Required: Update Your Password\"`",
							Validators: []validator.String{
								stringvalidator.LengthAtLeast(1),
								stringvalidator.LengthAtMost(250),
							},
						},
						"notification_body": schema.StringAttribute{
							Required: true,
							MarkdownDescription: "The body/content of the custom notification to display in the Company Portal app. " +
								"Should provide clear instructions or information to the user. Maximum recommended length: 200-300 characters. " +
								"Example: `\"Your device password will expire in 3 days. Please update it to maintain access to company resources.\"`",
							Validators: []validator.String{
								stringvalidator.LengthAtLeast(1),
								stringvalidator.LengthAtMost(1000),
							},
						},
					},
				},
			},
			"comanaged_devices": schema.ListNestedBlock{
				MarkdownDescription: "List of co-managed devices to send custom notifications to. These are devices managed by both Intune and " +
					"Configuration Manager (SCCM). Each entry specifies a device ID and the custom notification title and body.\n\n" +
					"**Note:** At least one of `managed_devices` or `comanaged_devices` must be provided.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"device_id": schema.StringAttribute{
							Required: true,
							MarkdownDescription: "The unique identifier (GUID) of the co-managed device to send the notification to. " +
								"Example: `\"12345678-1234-1234-1234-123456789abc\"`",
							Validators: []validator.String{
								stringvalidator.RegexMatches(
									regexp.MustCompile(constants.GuidRegex),
									"device_id must be a valid GUID format (e.g., 12345678-1234-1234-1234-123456789abc)",
								),
							},
						},
						"notification_title": schema.StringAttribute{
							Required: true,
							MarkdownDescription: "The title of the custom notification to display in the Company Portal app. " +
								"Should be concise and descriptive. Maximum recommended length: 50-60 characters. " +
								"Example: `\"Compliance Alert\"`",
							Validators: []validator.String{
								stringvalidator.LengthAtLeast(1),
								stringvalidator.LengthAtMost(250),
							},
						},
						"notification_body": schema.StringAttribute{
							Required: true,
							MarkdownDescription: "The body/content of the custom notification to display in the Company Portal app. " +
								"Should provide clear instructions or information to the user. Maximum recommended length: 200-300 characters. " +
								"Example: `\"Your device is not compliant with corporate security policies. Please contact IT support.\"`",
							Validators: []validator.String{
								stringvalidator.LengthAtLeast(1),
								stringvalidator.LengthAtMost(1000),
							},
						},
					},
				},
			},
		},
	}
}
