package graphBetaAndroidEnrollmentNotifications

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	planmodifiers "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/plan_modifiers"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
	commonschemagraphbeta "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema/graph_beta/device_management"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

const (
	ResourceName  = "graph_beta_device_management_android_enrollment_notifications"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180
)

var (
	// Basic resource interface (CRUD operations)
	_ resource.Resource = &AndroidEnrollmentNotificationsResource{}

	// Allows the resource to be configured with the provider client
	_ resource.ResourceWithConfigure = &AndroidEnrollmentNotificationsResource{}

	// Enables import functionality
	_ resource.ResourceWithImportState = &AndroidEnrollmentNotificationsResource{}

	// Enables plan modification/diff suppression
	_ resource.ResourceWithModifyPlan = &AndroidEnrollmentNotificationsResource{}
)

func NewAndroidEnrollmentNotificationsResource() resource.Resource {
	return &AndroidEnrollmentNotificationsResource{
		ReadPermissions: []string{
			"DeviceManagementConfiguration.Read.All",
		},
		WritePermissions: []string{
			"DeviceManagementConfiguration.ReadWrite.All",
		},
	}
}

type AndroidEnrollmentNotificationsResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ProviderTypeName string
	TypeName         string
	ReadPermissions  []string
	WritePermissions []string
}

// GetID returns the ID of the resource from the state.
func (r *AndroidEnrollmentNotificationsResource) GetID() string {
	return ResourceName
}

// Metadata returns the resource type name.
func (r *AndroidEnrollmentNotificationsResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	r.ProviderTypeName = req.ProviderTypeName
	r.TypeName = ResourceName
	resp.TypeName = r.FullTypeName()
}

// FullTypeName returns the full resource type name in the format "providername_resourcename".
func (r *AndroidEnrollmentNotificationsResource) FullTypeName() string {
	return r.ProviderTypeName + "_" + r.TypeName
}

// Configure sets the client for the resource.
func (r *AndroidEnrollmentNotificationsResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = client.SetGraphBetaClientForResource(ctx, req, resp, constants.PROVIDER_NAME+"_"+ResourceName)
}

// ImportState imports the resource state.
func (r *AndroidEnrollmentNotificationsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Schema defines the schema for the resource.
func (r *AndroidEnrollmentNotificationsResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages Android Enterprise notification configurations in Microsoft Intune. " +
			"This resource creates device enrollment notification configurations for Android for Work platform.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The unique identifier for the Android Enterprise notification configuration.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"display_name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The display name for the Android Enterprise notification configuration.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"description": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The description for the Android Enterprise notification configuration.",
			},
			"platform_type": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The platform type for the notification configuration. Must be either 'androidForWork' for Android Enterprise or 'android' for Android device andministrator.",
				Validators: []validator.String{
					stringvalidator.OneOf("androidForWork", "android"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"default_locale": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "The default locale for the notification configuration (e.g., 'en-US').",
				Default:             stringdefault.StaticString("en-US"),
			},
			"role_scope_tag_ids": schema.SetAttribute{
				ElementType:         types.StringType,
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Set of scope tag IDs for this Android Enterprise Notification configuration.",
				PlanModifiers: []planmodifier.Set{
					planmodifiers.DefaultSetValue(
						[]attr.Value{types.StringValue("0")},
					),
				},
			},
			"branding_options": schema.SetAttribute{
				ElementType:         types.StringType,
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "The branding options for the message template. Possible values are: none, includeCompanyLogo, includeCompanyName, includeContactInformation, includeCompanyPortalLink, includeDeviceDetails",
				Validators: []validator.Set{
					setvalidator.ValueStringsAre(stringvalidator.OneOf(
						"none",
						"includeCompanyLogo",
						"includeCompanyName",
						"includeContactInformation",
						"includeCompanyPortalLink",
						"includeDeviceDetails",
					)),
				},
				Default: setdefault.StaticValue(types.SetValueMust(types.StringType, []attr.Value{types.StringValue("none")})),
			},
			"notification_templates": schema.SetAttribute{
				ElementType:         types.StringType,
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "The notification template types for this configuration. Can be 'email', 'push', or both. Defaults to ['email', 'push'].",
				Validators: []validator.Set{
					setvalidator.ValueStringsAre(
						stringvalidator.OneOf("email", "push"),
					),
					setvalidator.SizeAtLeast(1),
					setvalidator.SizeAtMost(2),
				},
				PlanModifiers: []planmodifier.Set{
					planmodifiers.DefaultSetValue(
						[]attr.Value{types.StringValue("email"), types.StringValue("push")},
					),
				},
			},
			"priority": schema.Int32Attribute{
				Computed:            true,
				MarkdownDescription: "The priority of the notification configuration.",
			},
			"device_enrollment_configuration_type": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The type of device enrollment configuration.",
			},
			"localized_notification_messages": schema.SetNestedAttribute{
				Optional:            true,
				MarkdownDescription: "The localized notification messages for the configuration.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"locale": schema.StringAttribute{
							Required:            true,
							MarkdownDescription: "The locale for the notification message (e.g., 'en-us'). Must be in lowercase format.",
							PlanModifiers: []planmodifier.String{
								planmodifiers.EnsureLowerCaseString(),
							},
						},
						"subject": schema.StringAttribute{
							Required:            true,
							MarkdownDescription: "The subject of the notification message.",
						},
						"message_template": schema.StringAttribute{
							Required:            true,
							MarkdownDescription: "The template content of the notification message.",
						},
						"is_default": schema.BoolAttribute{
							Optional:            true,
							MarkdownDescription: "Whether this is the default notification message.",
						},
						"template_type": schema.StringAttribute{
							Required:            true,
							MarkdownDescription: "The type of template (email or push).",
							Validators: []validator.String{
								stringvalidator.OneOf("email", "push"),
							},
						},
					},
				},
			},
			"assignments": commonschemagraphbeta.AndroidNotificationAssignmentsSchema(),
			"timeouts":    commonschema.Timeouts(ctx),
		},
	}
}
