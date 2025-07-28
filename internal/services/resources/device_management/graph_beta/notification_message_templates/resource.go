package graphBetaNotificationMessageTemplates

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	planmodifiers "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/plan_modifiers"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

const (
	ResourceName  = "graph_beta_device_management_notification_message_template"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180
)

var (
	// Basic resource interface (CRUD operations)
	_ resource.Resource = &NotificationMessageTemplateResource{}

	// Allows the resource to be configured with the provider client
	_ resource.ResourceWithConfigure = &NotificationMessageTemplateResource{}

	// Enables import functionality
	_ resource.ResourceWithImportState = &NotificationMessageTemplateResource{}

	// Enables plan modification/diff suppression
	_ resource.ResourceWithModifyPlan = &NotificationMessageTemplateResource{}
)

func NewNotificationMessageTemplateResource() resource.Resource {
	return &NotificationMessageTemplateResource{
		ReadPermissions: []string{
			"DeviceManagementServiceConfig.Read.All",
		},
		WritePermissions: []string{
			"DeviceManagementServiceConfig.ReadWrite.All",
		},
		ResourcePath: "/deviceManagement/notificationMessageTemplates",
	}
}

type NotificationMessageTemplateResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ProviderTypeName string
	TypeName         string
	ReadPermissions  []string
	WritePermissions []string
	ResourcePath     string
}

// Metadata returns the resource type name.
func (r *NotificationMessageTemplateResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	r.ProviderTypeName = req.ProviderTypeName
	r.TypeName = ResourceName
	resp.TypeName = r.FullTypeName()
}

func (r *NotificationMessageTemplateResource) FullTypeName() string {
	return r.ProviderTypeName + "_" + r.TypeName
}

func (r *NotificationMessageTemplateResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = client.SetGraphBetaClientForResource(ctx, req, resp, constants.PROVIDER_NAME+"_"+ResourceName)
}
func (r *NotificationMessageTemplateResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}


func (r *NotificationMessageTemplateResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages an Intune notification message template for compliance notifications",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Manages notification message templates in Microsoft Intune using the `/deviceManagement/notificationMessageTemplates` endpoint. Notification message templates define the content and branding of compliance notifications sent to users.",
				PlanModifiers: []planmodifier.String{
					planmodifiers.UseStateForUnknownString(),
				},
			},
			"display_name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Display name for the notification message template",
			},
			"description": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Description of the notification message template",
			},
			"default_locale": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The default locale to fallback onto when the requested locale is not available",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"branding_options": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "The branding options for the message template. Possible values are: none, includeCompanyLogo, includeCompanyName, includeContactInformation, includeCompanyPortalLink, includeDeviceDetails",
				Validators: []validator.String{
					stringvalidator.OneOf("none", "includeCompanyLogo", "includeCompanyName", "includeContactInformation", "includeCompanyPortalLink", "includeDeviceDetails"),
				},
				PlanModifiers: []planmodifier.String{
					planmodifiers.DefaultValueString("none"),
				},
			},
			"role_scope_tag_ids": schema.SetAttribute{
				ElementType:         types.StringType,
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "List of scope tag IDs for this notification message template",
				PlanModifiers: []planmodifier.Set{
					planmodifiers.DefaultSetValue(
						[]attr.Value{types.StringValue("0")},
					),
				},
			},
			"last_modified_date_time": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "DateTime the notification message template was last modified",
			},
			"localized_notification_messages": schema.SetNestedAttribute{
				Optional:            true,
				MarkdownDescription: "The list of localized notification messages for this template",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "Unique identifier for the localized notification message",
						},
						"locale": schema.StringAttribute{
							Required:            true,
							MarkdownDescription: "The locale for the notification message (e.g., en-US, es-ES)",
							Validators: []validator.String{
								stringvalidator.LengthAtLeast(2),
							},
						},
						"subject": schema.StringAttribute{
							Required:            true,
							MarkdownDescription: "The subject of the notification message",
							Validators: []validator.String{
								stringvalidator.LengthAtLeast(1),
							},
						},
						"message_template": schema.StringAttribute{
							Required:            true,
							MarkdownDescription: "The message template text that can include tokens like {DeviceName}, {UserName}, etc.",
							Validators: []validator.String{
								stringvalidator.LengthAtLeast(1),
							},
						},
						"is_default": schema.BoolAttribute{
							Optional:            true,
							Computed:            true,
							MarkdownDescription: "Indicates if this is the default message for the template",
						},
						"last_modified_date_time": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "DateTime the localized message was last modified",
						},
					},
				},
			},
			"timeouts": commonschema.Timeouts(ctx),
		},
	}
}
