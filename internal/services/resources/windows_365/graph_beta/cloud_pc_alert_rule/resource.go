package graphBetaCloudPcAlertRule

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	planmodifiers "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/plan_modifiers"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

const (
	ResourceName  = "microsoft365_graph_beta_windows_365_cloud_pc_alert_rule"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180
)

var (
	// Basic resource interface (CRUD operations)
	_ resource.Resource = &CloudPcAlertRuleResource{}

	// Allows the resource to be configured with the provider client
	_ resource.ResourceWithConfigure = &CloudPcAlertRuleResource{}

	// Enables import functionality
	_ resource.ResourceWithImportState = &CloudPcAlertRuleResource{}

	// Enables plan modification/diff suppression
	_ resource.ResourceWithModifyPlan = &CloudPcAlertRuleResource{}
)

func NewCloudPcAlertRuleResource() resource.Resource {
	return &CloudPcAlertRuleResource{
		ReadPermissions:  []string{"DeviceManagementConfiguration.Read.All"},
		WritePermissions: []string{"DeviceManagementConfiguration.ReadWrite.All"},
	}
}

type CloudPcAlertRuleResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ReadPermissions  []string
	WritePermissions []string
}

func (r *CloudPcAlertRuleResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = ResourceName
}

func (r *CloudPcAlertRuleResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = client.SetGraphBetaClientForResource(ctx, req, resp, ResourceName)
}

func (r *CloudPcAlertRuleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *CloudPcAlertRuleResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Windows 365 Cloud PC Alert Rule using the Microsoft Graph Beta API.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The unique identifier for the alert rule. Read-only.",
				PlanModifiers: []planmodifier.String{
					planmodifiers.UseStateForUnknownString(),
				},
			},
			"alert_rule_template": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The rule template of the alert event. Possible values: cloudPcProvisionScenario, cloudPcImageUploadScenario, cloudPcOnPremiseNetworkConnectionCheckScenario, unknownFutureValue, cloudPcInGracePeriodScenario, cloudPcFrontlineInsufficientLicensesScenario, cloudPcInaccessibleScenario, cloudPcFrontlineConcurrencyScenario.",
				Validators: []validator.String{
					stringvalidator.OneOf(
						"cloudPcProvisionScenario",
						"cloudPcImageUploadScenario",
						"cloudPcOnPremiseNetworkConnectionCheckScenario",
						"unknownFutureValue",
						"cloudPcInGracePeriodScenario",
						"cloudPcFrontlineInsufficientLicensesScenario",
						"cloudPcInaccessibleScenario",
						"cloudPcFrontlineConcurrencyScenario",
					),
				},
			},
			"description": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The rule description.",
			},
			"display_name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The display name of the rule.",
			},
			"enabled": schema.BoolAttribute{
				Optional:            true,
				MarkdownDescription: "Whether the rule is enabled.",
			},
			"is_system_rule": schema.BoolAttribute{
				Required:            true,
				MarkdownDescription: "Indicates whether the rule is a system rule. Read-only.",
			},
			"notification_channels": schema.ListNestedAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "The notification channels of the rule.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"notification_channel_type": schema.StringAttribute{
							Optional:            true,
							MarkdownDescription: "The type of the notification channel. Possible values: portal, email, phoneCall, sms, unknownFutureValue.",
							Validators: []validator.String{
								stringvalidator.OneOf("portal", "email", "phoneCall", "sms", "unknownFutureValue"),
							},
						},
						"notification_receivers": schema.ListNestedAttribute{
							Optional:            true,
							MarkdownDescription: "Notification receivers for the channel.",
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"contact_information": schema.StringAttribute{
										Optional:            true,
										MarkdownDescription: "Contact information for the receiver.",
									},
									"locale": schema.StringAttribute{
										Optional:            true,
										MarkdownDescription: "Locale for the receiver.",
									},
								},
							},
						},
					},
				},
			},
			"severity": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The severity of the rule. Possible values: unknown, informational, warning, critical, unknownFutureValue.",
				Validators: []validator.String{
					stringvalidator.OneOf("unknown", "informational", "warning", "critical", "unknownFutureValue"),
				},
			},
			"threshold": schema.SingleNestedAttribute{
				Optional:            true,
				MarkdownDescription: "The conditions that determine when to send alerts. Deprecated. Use conditions instead.",
				Attributes: map[string]schema.Attribute{
					"aggregation": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "Aggregation method. Possible values: count, percentage, affectedCloudPcCount, affectedCloudPcPercentage, unknownFutureValue.",
						Validators: []validator.String{
							stringvalidator.OneOf("count", "percentage", "affectedCloudPcCount", "affectedCloudPcPercentage", "unknownFutureValue"),
						},
					},
					"operator": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "Operator. Possible values: greaterOrEqual, equal, greater, less, lessOrEqual, notEqual, unknownFutureValue.",
						Validators: []validator.String{
							stringvalidator.OneOf("greaterOrEqual", "equal", "greater", "less", "lessOrEqual", "notEqual", "unknownFutureValue"),
						},
					},
					"target": schema.Int32Attribute{
						Optional:            true,
						MarkdownDescription: "The target threshold value.",
					},
				},
			},
			"conditions": schema.ListNestedAttribute{
				Optional:            true,
				MarkdownDescription: "The conditions that determine when to send alerts.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"relationship_type": schema.StringAttribute{
							Optional:            true,
							MarkdownDescription: "The relationship type. Possible values: and, or, unknownFutureValue.",
							Validators: []validator.String{
								stringvalidator.OneOf("and", "or", "unknownFutureValue"),
							},
						},
						"condition_category": schema.StringAttribute{
							Optional:            true,
							MarkdownDescription: "The property that the rule condition monitors. Possible values: provisionFailures, imageUploadFailures, azureNetworkConnectionCheckFailures, cloudPcInGracePeriod, frontlineInsufficientLicenses, cloudPcConnectionErrors, cloudPcHostHealthCheckFailures, cloudPcZoneOutage, unknownFutureValue, frontlineBufferUsageDuration, frontlineBufferUsageThreshold.",
							Validators: []validator.String{
								stringvalidator.OneOf(
									"provisionFailures",
									"imageUploadFailures",
									"azureNetworkConnectionCheckFailures",
									"cloudPcInGracePeriod",
									"frontlineInsufficientLicenses",
									"cloudPcConnectionErrors",
									"cloudPcHostHealthCheckFailures",
									"cloudPcZoneOutage",
									"unknownFutureValue",
									"frontlineBufferUsageDuration",
									"frontlineBufferUsageThreshold",
								),
							},
						},
						"aggregation": schema.StringAttribute{
							Optional:            true,
							MarkdownDescription: "The built-in aggregation method for the rule condition. Possible values: count, percentage, affectedCloudPcCount, affectedCloudPcPercentage, unknownFutureValue, durationInMinutes.",
							Validators: []validator.String{
								stringvalidator.OneOf(
									"count",
									"percentage",
									"affectedCloudPcCount",
									"affectedCloudPcPercentage",
									"unknownFutureValue",
									"durationInMinutes",
								),
							},
						},
						"operator": schema.StringAttribute{
							Optional:            true,
							MarkdownDescription: "The built-in operator for the rule condition. Possible values: greaterOrEqual, equal, greater, less, lessOrEqual, notEqual, unknownFutureValue.",
							Validators: []validator.String{
								stringvalidator.OneOf(
									"greaterOrEqual",
									"equal",
									"greater",
									"less",
									"lessOrEqual",
									"notEqual",
									"unknownFutureValue",
								),
							},
						},
						"threshold_value": schema.StringAttribute{
							Optional:            true,
							MarkdownDescription: "The threshold value of the alert condition. The threshold value can be a number in string form or string like 'WestUS'.",
						},
					},
				},
			},
			"timeouts": commonschema.ResourceTimeouts(ctx),
		},
	}
}
