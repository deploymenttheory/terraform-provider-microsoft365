package graphBetaDeviceEnrollmentConfiguration

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common"
	planmodifiers "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/plan_modifiers"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/schema"
	commonschemagraphbeta "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/schema/graph_beta/device_management"
	customValidator "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/validators"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
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
	ResourceName  = "graph_beta_device_management_device_enrollment_configuration"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180
)

var (
	// Basic resource interface (CRUD operations)
	_ resource.Resource = &DeviceEnrollmentConfigurationResource{}

	// Allows the resource to be configured with the provider client
	_ resource.ResourceWithConfigure = &DeviceEnrollmentConfigurationResource{}

	// Enables import functionality
	_ resource.ResourceWithImportState = &DeviceEnrollmentConfigurationResource{}

	// Enables plan modification/diff suppression
	_ resource.ResourceWithModifyPlan = &DeviceEnrollmentConfigurationResource{}
)

func NewDeviceEnrollmentConfigurationResource() resource.Resource {
	return &DeviceEnrollmentConfigurationResource{
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

type DeviceEnrollmentConfigurationResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ProviderTypeName string
	TypeName         string
	ReadPermissions  []string
	WritePermissions []string
	ResourcePath     string
}

// Metadata returns the resource type name.
func (r *DeviceEnrollmentConfigurationResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + ResourceName
}

// Configure sets the client for the resource.
func (r *DeviceEnrollmentConfigurationResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = common.SetGraphBetaClientForResource(ctx, req, resp, r.TypeName)
}

// ImportState imports the resource state.
func (r *DeviceEnrollmentConfigurationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Schema defines the schema for the resource.
func (r *DeviceEnrollmentConfigurationResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a Device Enrollment Configuration in Microsoft Intune.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					planmodifiers.UseStateForUnknownString(),
				},
				MarkdownDescription: "The Identifier of the entity.",
			},
			"display_name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The display name of the device enrollment configuration.",
			},
			"description": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The description of the device enrollment configuration.",
			},
			"priority": schema.Int32Attribute{
				Optional:            true,
				MarkdownDescription: "Priority is used when a user exists in multiple groups that are assigned enrollment configuration. Users are subject only to the configuration with the lowest priority value.",
			},
			"created_date_time": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The date time that the device enrollment configuration was created.",
			},
			"last_modified_date_time": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The date time that the device enrollment configuration was last modified.",
			},
			"version": schema.Int32Attribute{
				Computed:            true,
				MarkdownDescription: "The version of the device enrollment configuration.",
			},
			"role_scope_tag_ids": schema.SetAttribute{
				ElementType:         types.StringType,
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Set of scope tag IDs for this Settings Catalog template profile.",
				PlanModifiers: []planmodifier.Set{
					planmodifiers.DefaultSetValue(
						[]attr.Value{types.StringValue("0")},
					),
				},
			},
			"device_enrollment_configuration_type": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The type of the device enrollment configuration. Possible values are: unknown, limit, platformRestrictions, windowsHelloForBusiness, defaultLimit, defaultPlatformRestrictions, defaultWindowsHelloForBusiness, defaultWindows10EnrollmentCompletionPageConfiguration, windows10EnrollmentCompletionPageConfiguration, deviceComanagementAuthorityConfiguration, singlePlatformRestriction, enrollmentNotificationsConfiguration.",
				Validators: []validator.String{
					stringvalidator.OneOf(
						"limit",
						"platformRestrictions",
						"windowsHelloForBusiness",
						"defaultLimit",
						"defaultPlatformRestrictions",
						"defaultWindowsHelloForBusiness",
						"defaultWindows10EnrollmentCompletionPageConfiguration",
						"windows10EnrollmentCompletionPageConfiguration",
						"deviceComanagementAuthorityConfiguration",
						"singlePlatformRestriction",
						"enrollmentNotificationsConfiguration",
					),
				},
			},
			"timeouts": commonschema.Timeouts(ctx),
		},
		Blocks: map[string]schema.Block{
			"platform_restriction": schema.SingleNestedBlock{
				MarkdownDescription: "Single platform enrollment restriction configuration.",
				Attributes: map[string]schema.Attribute{
					"platform_type": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "The platform type this restriction applies to. Possible values are: `allPlatforms`, `ios`, `windows`, `windowsPhone`, `android`, `androidForWork`, `mac`, `linux`.",
						Validators: []validator.String{
							stringvalidator.OneOf(
								"allPlatforms",
								"ios",
								"windows",
								"windowsPhone",
								"android",
								"androidForWork",
								"mac",
								"linux",
							),
						},
					},
					"restriction": schema.SingleNestedAttribute{
						Optional:            true,
						MarkdownDescription: "The platform restriction settings.",
						Attributes: map[string]schema.Attribute{
							"platform_blocked": schema.BoolAttribute{
								Optional:            true,
								MarkdownDescription: "Block the platform from enrolling.",
							},
							"personal_device_enrollment_blocked": schema.BoolAttribute{
								Optional:            true,
								MarkdownDescription: "Block personally owned devices from enrolling.",
							},
							"os_minimum_version": schema.StringAttribute{
								Optional:            true,
								MarkdownDescription: "Minimum version of the platform.",
							},
							"os_maximum_version": schema.StringAttribute{
								Optional:            true,
								MarkdownDescription: "Maximum version of the platform.",
							},
							"blocked_manufacturers": schema.SetAttribute{
								ElementType:         types.StringType,
								Optional:            true,
								MarkdownDescription: "Collection of blocked manufacturers.",
							},
							"blocked_skus": schema.SetAttribute{
								ElementType:         types.StringType,
								Optional:            true,
								MarkdownDescription: "Collection of blocked SKUs.",
							},
						},
					},
				},
			},

			"enrollment_notifications": schema.SingleNestedBlock{
				MarkdownDescription: "Settings for enrollment notifications sent to end users during device enrollment.",
				Attributes: map[string]schema.Attribute{
					"platform_type": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "Platform type of the Enrollment Notification. Possible values are: `allPlatforms`, `ios`, `windows`, `windowsPhone`, `android`, `androidForWork`, `mac`, `linux`, `unknownFutureValue`.",
						Validators: []validator.String{
							stringvalidator.OneOf(
								"allPlatforms",
								"ios",
								"windows",
								"windowsPhone",
								"android",
								"androidForWork",
								"mac",
								"linux",
							),
						},
					},
					"template_type": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "Template type of the Enrollment Notification. Possible values are: `email`, `push`, `unknownFutureValue`.",
						Validators: []validator.String{
							stringvalidator.OneOf(
								"email",
								"push",
							),
						},
					},
					"notification_message_template_id": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "Notification Message Template ID in UUID/GUID format.",
					},
					"notification_templates": schema.SetAttribute{
						ElementType:         types.StringType,
						Optional:            true,
						MarkdownDescription: "The list of notification templates.",
					},
					"branding_options": schema.ListAttribute{
						ElementType:         types.StringType,
						Optional:            true,
						MarkdownDescription: "Branding Options for the Enrollment Notification. This is a bitmask that can include multiple values: `none`, `includeCompanyLogo`, `includeCompanyName`, `includeContactInformation`, `includeCompanyPortalLink`, `includeDeviceDetails`, `unknownFutureValue`.",
						Validators: []validator.List{
							customValidator.StringListAllowedValues(
								"none",
								"includeCompanyLogo",
								"includeCompanyName",
								"includeContactInformation",
								"includeCompanyPortalLink",
								"includeDeviceDetails",
								"unknownFutureValue",
							),
						},
					},
					"default_locale": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "Default Locale for the Enrollment Notification.",
					},
					"include_company_portal_link": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: "When TRUE, includes a link to the Company Portal in the notification.",
					},
					"send_push_notification": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: "When TRUE, sends a push notification along with email notification.",
					},
					"notification_title": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "The title of the notification message.",
					},
					"notification_body": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "The body text of the notification message.",
					},
					"notification_sender": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "The sender name for the notification message.",
					},
				},
			},
			"device_comanagement_authority": schema.SingleNestedBlock{
				MarkdownDescription: "Settings for configuring the device co-management authority between Intune and Configuration Manager.",
				Attributes: map[string]schema.Attribute{
					"managed_device_authority": schema.Int32Attribute{
						Optional:            true,
						MarkdownDescription: "Co-Management Authority configuration for managed devices. Defines how workloads are split between Intune and Configuration Manager.",
					},
					"install_configuration_manager_agent": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: "When TRUE, the Configuration Manager agent will be installed during device enrollment. When FALSE, the agent will not be installed.",
					},
					"configuration_manager_agent_command_line_argument": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "Command line arguments to pass to the Configuration Manager agent during installation.",
					},
				},
			},
			"windows10_enrollment_completion_page": schema.SingleNestedBlock{
				MarkdownDescription: "Windows 10 enrollment completion page settings which specify the information shown to users during device enrollment.",
				Attributes: map[string]schema.Attribute{
					"show_installation_progress": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: "When TRUE, shows installation progress to user. When false, hides installation progress. The default is false.",
					},
					"block_device_setup_retry_by_user": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: "When TRUE, blocks user from retrying the setup on installation failure. When false, user is allowed to retry. The default is false.",
					},
					"allow_device_reset_on_install_failure": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: "When TRUE, allows device reset on installation failure. When false, reset is blocked. The default is false.",
					},
					"allow_log_collection_on_install_failure": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: "When TRUE, allows log collection on installation failure. When false, log collection is not allowed. The default is false.",
					},
					"custom_error_message": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "The custom error message to show upon installation failure. Max length is 10000.",
					},
					"install_progress_timeout_in_minutes": schema.Int32Attribute{
						Optional:            true,
						MarkdownDescription: "The installation progress timeout in minutes. Default is 60 minutes.",
					},
					"allow_device_use_on_install_failure": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: "When TRUE, allows the user to continue using the device on installation failure. When false, blocks the user on installation failure. The default is false.",
					},
					"selected_mobile_app_ids": schema.SetAttribute{
						ElementType:         types.StringType,
						Optional:            true,
						MarkdownDescription: "Selected applications to track the installation status. It is in the form of an array of GUIDs.",
					},
					"allow_non_blocking_app_installation": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: "When TRUE, ESP (Enrollment Status Page) installs all required apps targeted during technician phase and ignores any failures for non-blocking apps. When FALSE, ESP fails on any error during app install. The default is false.",
					},
					"install_quality_updates": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: "Allows quality updates installation during OOBE.",
					},
					"track_install_progress_for_autopilot_only": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: "When TRUE, installation progress is tracked for only Autopilot enrollment scenarios. When false, other scenarios are tracked as well. The default is false.",
					},
					"disable_user_status_tracking_after_first_user": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: "When TRUE, disables showing installation progress for first user post enrollment. When false, enables progress. The default is false.",
					},
				},
			},

			"default_windows10_enrollment_completion_page": schema.SingleNestedBlock{
				MarkdownDescription: "Default settings for Windows 10 enrollment completion page which specify the information shown to users during device enrollment.",
				Attributes: map[string]schema.Attribute{
					"allow_devices_for_users": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: "When TRUE, allows devices for users in this profile. When FALSE, prevents devices for users in this profile.",
					},
					"show_installation_progress": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: "When TRUE, shows installation progress to user. When false, hides installation progress. The default is false.",
					},
					"allow_device_reset_on_install_failure": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: "When TRUE, allows device reset on installation failure. When false, reset is blocked. The default is false.",
					},
					"allow_log_collection_on_install_failure": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: "When TRUE, allows log collection on installation failure. When false, log collection is not allowed. The default is false.",
					},
					"custom_error_message": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "The custom error message to show upon installation failure. Max length is 10000.",
					},
					"install_progress_timeout_in_minutes": schema.Int32Attribute{
						Optional:            true,
						MarkdownDescription: "The installation progress timeout in minutes. Default is 60 minutes.",
					},
					"selected_mobile_app_ids": schema.SetAttribute{
						ElementType:         types.StringType,
						Optional:            true,
						MarkdownDescription: "Selected applications to track the installation status. It is in the form of an array of GUIDs.",
					},
					"track_install_progress_for_autopilot_only": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: "When TRUE, installation progress is tracked for only Autopilot enrollment scenarios. When false, other scenarios are tracked as well. The default is false.",
					},
					"disable_user_status_tracking_after_first_user": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: "When TRUE, disables showing installation progress for first user post enrollment. When false, enables progress. The default is false.",
					},
				},
			},

			"windows_hello_for_business": schema.SingleNestedBlock{
				MarkdownDescription: "Settings for Windows Hello for Business authentication on Windows devices.",
				Attributes: map[string]schema.Attribute{
					"state": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "Controls whether to allow the device to be configured for Windows Hello for Business. Possible values are: `notConfigured`, `enabled`, `disabled`.",
						Validators: []validator.String{
							stringvalidator.OneOf("notConfigured", "enabled", "disabled"),
						},
					},
					"pin_minimum_length": schema.Int32Attribute{
						Optional:            true,
						MarkdownDescription: "Controls the minimum number of characters required for the Windows Hello for Business PIN. This value must be between 4 and 127, inclusive.",
						Validators: []validator.Int32{
							int32validator.Between(4, 127),
						},
					},
					"pin_maximum_length": schema.Int32Attribute{
						Optional:            true,
						MarkdownDescription: "Controls the maximum number of characters allowed for the Windows Hello for Business PIN. This value must be between 4 and 127, inclusive.",
						Validators: []validator.Int32{
							int32validator.Between(4, 127),
						},
					},
					"pin_uppercase_characters_usage": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "Controls the ability to use uppercase letters in the Windows Hello for Business PIN. Possible values are: `allowed`, `required`, `disallowed`.",
						Validators: []validator.String{
							stringvalidator.OneOf("allowed", "required", "disallowed"),
						},
					},
					"pin_lowercase_characters_usage": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "Controls the ability to use lowercase letters in the Windows Hello for Business PIN. Possible values are: `allowed`, `required`, `disallowed`.",
						Validators: []validator.String{
							stringvalidator.OneOf("allowed", "required", "disallowed"),
						},
					},
					"pin_special_characters_usage": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "Controls the ability to use special characters in the Windows Hello for Business PIN. Possible values are: `allowed`, `required`, `disallowed`.",
						Validators: []validator.String{
							stringvalidator.OneOf("allowed", "required", "disallowed"),
						},
					},
					"security_device_required": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: "Controls whether to require a Trusted Platform Module (TPM) for provisioning Windows Hello for Business. If set to False, all devices can provision Windows Hello for Business even if there is not a usable TPM.",
					},
					"unlock_with_biometrics_enabled": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: "Controls the use of biometric gestures, such as face and fingerprint, as an alternative to the Windows Hello for Business PIN. If set to False, biometric gestures are not allowed.",
					},
					"remote_passport_enabled": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: "Controls the use of Remote Windows Hello for Business. Remote Windows Hello for Business provides the ability for a portable, registered device to be usable as a companion for desktop authentication.",
					},
					"pin_previous_block_count": schema.Int32Attribute{
						Optional:            true,
						MarkdownDescription: "Controls the ability to prevent users from using past PINs. This must be set between 0 and 50, inclusive.",
						Validators: []validator.Int32{
							int32validator.Between(0, 50),
						},
					},
					"pin_expiration_in_days": schema.Int32Attribute{
						Optional:            true,
						MarkdownDescription: "Controls the period of time (in days) that a PIN can be used before the system requires the user to change it. This must be set between 0 and 730, inclusive.",
						Validators: []validator.Int32{
							int32validator.Between(0, 730),
						},
					},
					"enhanced_biometrics_state": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "Controls the ability to use the anti-spoofing features for facial recognition on devices which support it. Possible values are: `notConfigured`, `enabled`, `disabled`.",
						Validators: []validator.String{
							stringvalidator.OneOf("notConfigured", "enabled", "disabled"),
						},
					},
					"security_key_for_sign_in": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "Security key for Sign In provides the capacity for remotely turning ON/OFF Windows Hello Security Key. Possible values are: `notConfigured`, `enabled`, `disabled`.",
						Validators: []validator.String{
							stringvalidator.OneOf("notConfigured", "enabled", "disabled"),
						},
					},
					"enhanced_sign_in_security": schema.Int32Attribute{
						Optional:            true,
						MarkdownDescription: "Setting to configure Enhanced sign-in security. Default is Not Configured.",
					},
				},
			},
			"device_enrollment_limit": schema.SingleNestedBlock{
				MarkdownDescription: "Settings that limit the number of devices a user can enroll.",
				Attributes: map[string]schema.Attribute{
					"limit": schema.Int32Attribute{
						Optional:            true,
						MarkdownDescription: "The maximum number of devices that a user can enroll. Maximum of 15.",
						Validators: []validator.Int32{
							int32validator.AtMost(15),
						},
					},
				},
			},
			"assignment": commonschemagraphbeta.WindowsUpdateAssignments(),
		},
	}
}
