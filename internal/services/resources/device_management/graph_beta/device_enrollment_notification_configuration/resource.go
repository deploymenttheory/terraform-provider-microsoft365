package graphBetaDeviceEnrollmentNotificationConfiguration

import (
	"context"
	"regexp"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	planmodifiers "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/plan_modifiers"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/validators"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

const (
	ResourceName  = "graph_beta_device_management_device_enrollment_notification_configuration"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180
)

var (
	// Basic resource interface (CRUD operations)
	_ resource.Resource = &DeviceEnrollmentNotificationConfigurationResource{}

	// Allows the resource to be configured with the provider client
	_ resource.ResourceWithConfigure = &DeviceEnrollmentNotificationConfigurationResource{}

	// Enables import functionality
	_ resource.ResourceWithImportState = &DeviceEnrollmentNotificationConfigurationResource{}

	// Enables plan modification/diff suppression
	_ resource.ResourceWithModifyPlan = &DeviceEnrollmentNotificationConfigurationResource{}
)

func NewDeviceEnrollmentNotificationConfigurationResource() resource.Resource {
	return &DeviceEnrollmentNotificationConfigurationResource{
		ReadPermissions: []string{
			"DeviceManagementServiceConfig.Read.All",
			"DeviceManagementConfiguration.Read.All",
		},
		WritePermissions: []string{
			"DeviceManagementServiceConfig.ReadWrite.All",
			"DeviceManagementConfiguration.ReadWrite.All",
		},
		ResourcePath: "/deviceManagement/deviceEnrollmentConfigurations",
	}
}

type DeviceEnrollmentNotificationConfigurationResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ProviderTypeName string
	TypeName         string
	ReadPermissions  []string
	WritePermissions []string
	ResourcePath     string
}

// Metadata returns the resource type name.
func (r *DeviceEnrollmentNotificationConfigurationResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	r.ProviderTypeName = req.ProviderTypeName
	r.TypeName = ResourceName
	resp.TypeName = r.FullTypeName()
}

// FullTypeName returns the full type name of the resource for logging purposes.
func (r *DeviceEnrollmentNotificationConfigurationResource) FullTypeName() string {
	return r.ProviderTypeName + "_" + r.TypeName
}

// Configure sets the client for the resource.
func (r *DeviceEnrollmentNotificationConfigurationResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = client.SetGraphBetaClientForResource(ctx, req, resp, constants.PROVIDER_NAME+"_"+ResourceName)
}

// ImportState imports the resource state.
func (r *DeviceEnrollmentNotificationConfigurationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *DeviceEnrollmentNotificationConfigurationResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages device enrollment notification configurations using the `/deviceManagement/deviceEnrollmentConfigurations` endpoint. Enrollment notification configurations are used to send notifications during device enrollment processes.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					planmodifiers.UseStateForUnknownString(),
				},
				MarkdownDescription: "The unique identifier of the device enrollment notification configuration.",
			},
			"display_name": schema.StringAttribute{
				MarkdownDescription: "The display name of the device enrollment configuration.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "The description of the device enrollment configuration.",
				Optional:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"priority": schema.Int32Attribute{
				MarkdownDescription: "Priority is used when a user exists in multiple groups that are assigned enrollment configuration. Users are subject only to the configuration with the lowest priority value.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Int32{
					int32planmodifier.UseStateForUnknown(),
					int32planmodifier.RequiresReplace(),
				},
				Validators: []validator.Int32{
					int32validator.AtLeast(0),
				},
			},
			"created_date_time": schema.StringAttribute{
				MarkdownDescription: "Created date time in UTC of the device enrollment configuration. This property is read-only.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplace(),
				},
			},
			"last_modified_date_time": schema.StringAttribute{
				MarkdownDescription: "Last modified date time in UTC of the device enrollment configuration. This property is read-only.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"version": schema.Int32Attribute{
				MarkdownDescription: "The version of the device enrollment configuration. This property is read-only.",
				Computed:            true,
				PlanModifiers: []planmodifier.Int32{
					int32planmodifier.UseStateForUnknown(),
					int32planmodifier.RequiresReplace(),
				},
			},
			"role_scope_tag_ids": schema.SetAttribute{
				ElementType:         types.StringType,
				MarkdownDescription: "Optional role scope tags for the enrollment restrictions.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Set{
					planmodifiers.DefaultSetValue(
						[]attr.Value{types.StringValue("0")},
					),
					setplanmodifier.RequiresReplace(),
				},
			},
			"device_enrollment_configuration_type": schema.StringAttribute{
				MarkdownDescription: "Support for Enrollment Configuration Type. This property is read-only.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplace(),
				},
			},
			"platform_type": schema.StringAttribute{
				MarkdownDescription: "Platform type of the Enrollment Notification. Possible values are: `allPlatforms`, `ios`, `windows`, `windowsPhone`, `android`, `androidForWork`, `mac`, `linux`, `unknownFutureValue`.",
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("allPlatforms", "ios", "windows", "windowsPhone", "android", "androidForWork", "mac", "linux", "unknownFutureValue"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"template_types": schema.SetAttribute{
				ElementType:         types.StringType,
				MarkdownDescription: "Template types of the Enrollment Notification. Possible values are: `email`, `push`.",
				Required:            true,
				Validators: []validator.Set{
					setvalidator.ValueStringsAre(stringvalidator.OneOf("email", "push")),
				},
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.RequiresReplace(),
				},
			},
			"notification_message_template_id": schema.StringAttribute{
				MarkdownDescription: "Notification Message Template Id.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"notification_templates": schema.SetAttribute{
				ElementType:         types.StringType,
				MarkdownDescription: "The list of notification data.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.RequiresReplace(),
				},
			},
			"branding_options": schema.SetAttribute{
				ElementType:         types.StringType,
				MarkdownDescription: "Branding Options for the Email footer and header. Possible values are: `none`, `includeCompanyLogo`, `includeCompanyName`, `includeContactInformation`, `includeCompanyPortalLink`, `includeDeviceDetails`, `unknownFutureValue`.",
				Optional:            true,
				Computed:            true,
				Validators: []validator.Set{
					setvalidator.ValueStringsAre(stringvalidator.OneOf("none", "includeCompanyLogo", "includeCompanyName", "includeContactInformation", "includeCompanyPortalLink", "includeDeviceDetails", "unknownFutureValue")),
				},
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.RequiresReplace(),
				},
			},
			"push_localized_message": schema.SingleNestedAttribute{
				MarkdownDescription: "Configuration for push notification localized message. Changes to this block will trigger replacement of the entire resource.",
				Optional:            true,
				Attributes: map[string]schema.Attribute{
					"locale": schema.StringAttribute{
						MarkdownDescription: "The Locale for which this message is destined (e.g., 'en-us', 'es-es').",
						Optional:            true,
						Computed:            true,
						Default:             stringdefault.StaticString("en-us"),
						Validators: []validator.String{
							stringvalidator.RegexMatches(
								regexp.MustCompile(`^[a-z]{2}(-[a-z]{2})?$`),
								"must be a valid locale format (e.g., 'en', 'en-us', 'es-es')",
							),
						},
					},
					"subject": schema.StringAttribute{
						MarkdownDescription: "The Message Template Subject.",
						Required:            true,
					},
					"message_template": schema.StringAttribute{
						MarkdownDescription: "The Message Template content.",
						Required:            true,
					},
					"is_default": schema.BoolAttribute{
						MarkdownDescription: "Flag to indicate whether or not this is the default locale for language fallback.",
						Optional:            true,
						Computed:            true,
						Default:             booldefault.StaticBool(true),
					},
				},
				Validators: []validator.Object{
					validators.BlockRequiresSetValue("template_types", "push", "push_localized_message"),
				},
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.RequiresReplace(),
				},
			},
			"email_localized_message": schema.SingleNestedAttribute{
				MarkdownDescription: "Configuration for email notification localized message. Changes to this block will trigger replacement of the entire resource.",
				Optional:            true,
				Attributes: map[string]schema.Attribute{
					"locale": schema.StringAttribute{
						MarkdownDescription: "The Locale for which this message is destined (e.g., 'en-us', 'es-es').",
						Optional:            true,
						Computed:            true,
						Default:             stringdefault.StaticString("en-us"),
						Validators: []validator.String{
							stringvalidator.RegexMatches(
								regexp.MustCompile(`^[a-z]{2}(-[a-z]{2})?$`),
								"must be a valid locale format (e.g., 'en', 'en-us', 'es-es')",
							),
						},
					},
					"subject": schema.StringAttribute{
						MarkdownDescription: "The Message Template Subject.",
						Required:            true,
						Validators: []validator.String{
							validators.RequiredWhenSetContains("template_types", "email"),
						},
					},
					"message_template": schema.StringAttribute{
						MarkdownDescription: "The Message Template content.",
						Required:            true,
						Validators: []validator.String{
							validators.RequiredWhenSetContains("template_types", "email"),
						},
					},
					"is_default": schema.BoolAttribute{
						MarkdownDescription: "Flag to indicate whether or not this is the default locale for language fallback.",
						Optional:            true,
						Computed:            true,
						Default:             booldefault.StaticBool(true),
					},
				},
				Validators: []validator.Object{
					validators.BlockRequiresSetValue("template_types", "email", "email_localized_message"),
				},
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.RequiresReplace(),
				},
			},
			"assignments": schema.SetNestedAttribute{
				MarkdownDescription: "The list of assignments for the device enrollment configuration. This will overwrite any existing assignments.",
				Optional:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"target": schema.SingleNestedAttribute{
							Required:            true,
							MarkdownDescription: "The target for the assignment.",
							Attributes: map[string]schema.Attribute{
								"target_type": schema.StringAttribute{
									Description: "The type of target for the assignment. Possible values are `group`, or `allLicensedUsers`.",
									Required:    true,
									Validators: []validator.String{
										stringvalidator.OneOf("group", "allLicensedUsers"),
									},
								},
								"group_id": schema.StringAttribute{
									Description: "The ID of the group to be targeted. This is required when `target_type` is `group`.",
									Optional:    true,
									Validators: []validator.String{
										validators.RequiredWhenEquals("target_type", types.StringValue("group")),
									},
								},
								"device_and_app_management_assignment_filter_id": schema.StringAttribute{
									Description: "The ID of the filter to be applied to the assignment. Filters are not supported for `allDevices` and `allLicensedUsers` assignment targets.",
									Optional:    true,
									Computed:    true,
									Default:     stringdefault.StaticString("00000000-0000-0000-0000-000000000000"),
								},
								"device_and_app_management_assignment_filter_type": schema.StringAttribute{
									Description: "The type of filter to be applied (`include` or `exclude`).",
									Optional:    true,
									Computed:    true,
									Default:     stringdefault.StaticString("none"),
									Validators: []validator.String{
										stringvalidator.OneOf("include", "exclude", "none"),
									},
								},
							},
						},
					},
				},
			},
			"timeouts": commonschema.Timeouts(ctx),
		},
	}
}
