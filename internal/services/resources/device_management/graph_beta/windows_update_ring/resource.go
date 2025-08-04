package graphBetaWindowsUpdateRing

import (
	"context"
	"regexp"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	planmodifiers "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/plan_modifiers"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
	commonschemagraphbeta "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema/graph_beta/device_management"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/validators"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

const (
	ResourceName  = "graph_beta_device_management_windows_update_ring"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180
)

var (
	// Basic resource interface (CRUD operations)
	_ resource.Resource = &WindowsUpdateRingResource{}

	// Allows the resource to be configured with the provider client
	_ resource.ResourceWithConfigure = &WindowsUpdateRingResource{}

	// Enables import functionality
	_ resource.ResourceWithImportState = &WindowsUpdateRingResource{}

	// Enables plan modification/diff suppression
	_ resource.ResourceWithModifyPlan = &WindowsUpdateRingResource{}
)

func NewWindowsUpdateRingResource() resource.Resource {
	return &WindowsUpdateRingResource{
		ReadPermissions: []string{
			"DeviceManagementConfiguration.Read.All",
		},
		WritePermissions: []string{
			"DeviceManagementConfiguration.ReadWrite.All",
		},
		ResourcePath: "/deviceManagement/deviceConfigurations",
	}
}

type WindowsUpdateRingResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ProviderTypeName string
	TypeName         string
	ReadPermissions  []string
	WritePermissions []string
	ResourcePath     string
}

// Metadata returns the resource type name.
func (r *WindowsUpdateRingResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	r.ProviderTypeName = req.ProviderTypeName
	r.TypeName = ResourceName
	resp.TypeName = r.FullTypeName()
}

// FullTypeName returns the full resource type name in the format "providername_resourcename".
func (r *WindowsUpdateRingResource) FullTypeName() string {
	return r.ProviderTypeName + "_" + r.TypeName
}

// Configure sets the client for the resource.
func (r *WindowsUpdateRingResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = client.SetGraphBetaClientForResource(ctx, req, resp, constants.PROVIDER_NAME+"_"+ResourceName)
}

// ImportState imports the resource state.
func (r *WindowsUpdateRingResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Schema defines the schema for the resource.
func (r *WindowsUpdateRingResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages Windows Update for Business configuration policies using the `/deviceManagement/deviceConfigurations` endpoint. This resource controls Windows Update settings including feature update deferrals, quality update schedules, driver management, and restart behaviors for managed Windows 10/11 devices.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					planmodifiers.UseStateForUnknownString(),
				},
				MarkdownDescription: "Key of the entity. Inherited from deviceConfiguration.",
			},
			"display_name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Admin provided name of the device configuration. Inherited from deviceConfiguration.",
				Validators: []validator.String{
					stringvalidator.LengthAtMost(255),
				},
			},
			"description": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Admin provided description of the Device Configuration. Inherited from deviceConfiguration.",
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
			"microsoft_update_service_allowed": schema.BoolAttribute{
				Required:            true,
				MarkdownDescription: "When TRUE, allows Microsoft Update Service. When FALSE, does not allow Microsoft Update Service. Returned by default. Query parameters are not supported.",
			},
			"drivers_excluded": schema.BoolAttribute{
				Required:            true,
				MarkdownDescription: "When TRUE, excludes Windows update Drivers. When FALSE, does not exclude Windows update Drivers. Returned by default. Query parameters are not supported.",
			},
			"quality_updates_deferral_period_in_days": schema.Int32Attribute{
				Required:            true,
				MarkdownDescription: "Defer Quality Updates by these many days with valid range from 0 to 30 days. Returned by default. Query parameters are not supported.",
			},
			"feature_updates_deferral_period_in_days": schema.Int32Attribute{
				Required:            true,
				MarkdownDescription: "Defer Feature Updates by these many days with valid range from 0 to 30 days. Returned by default. Query parameters are not supported.",
			},
			"allow_windows11_upgrade": schema.BoolAttribute{
				Required:            true,
				MarkdownDescription: "When TRUE, allows eligible Windows 10 devices to latest Windows 11 release. When FALSE, implies the device stays on the existing operating system. Returned by default. Query parameters are not supported.",
			},

			"skip_checks_before_restart": schema.BoolAttribute{
				Required:            true,
				MarkdownDescription: "When TRUE, skips all checks before restart: Battery level = 40%, User presence, Display Needed, Presentation mode, Full screen mode, phone call state, game mode etc. When FALSE, does not skip all checks before restart. Returned by default. Query parameters are not supported.",
			},
			"business_ready_updates_only": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString("userDefined"),
				MarkdownDescription: "Enable pre-release builds if you want devices to be on a Windows Insider channel." +
					"Enabling pre-release builds will cause devices to reboot. Determines which update branch devices will " +
					"receive their updates from. Possible values are: UserDefined, All, BusinessReadyOnly, WindowsInsiderBuildFast, " +
					"WindowsInsiderBuildSlow, WindowsInsiderBuildRelease." +
					"UserDefined equates to 'Not configured' in the gui." +
					"all equates to 'Not configured' in the gui." +
					"windowsInsiderBuildRelease equates to 'Windows Insider - Release Preview' in the gui." +
					"windowsInsiderBuildSlow equates to 'Beta Channel' in the gui." +
					"windowsInsiderBuildFast equates to ' Dev Channel' in the gui.",
				Validators: []validator.String{
					stringvalidator.OneOf(
						"userDefined",
						"all",
						"businessReadyOnly",
						"windowsInsiderBuildFast",
						"windowsInsiderBuildSlow",
						"windowsInsiderBuildRelease",
					),
				},
			},
			"automatic_update_mode": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The Automatic Update Mode. Possible values are: UserDefined, NotifyDownload, AutoInstallAtMaintenanceTime, AutoInstallAndRebootAtMaintenanceTime, AutoInstallAndRebootAtScheduledTime, AutoInstallAndRebootWithoutEndUserControl, WindowsDefault. UserDefined is the default value, no intent. Returned by default. Query parameters are not supported. Possible values are: userDefined, notifyDownload, autoInstallAtMaintenanceTime, autoInstallAndRebootAtMaintenanceTime, autoInstallAndRebootAtScheduledTime, autoInstallAndRebootWithoutEndUserControl.",
				Validators: []validator.String{
					stringvalidator.OneOf(
						"userDefined",                               // reset to default - no other fields should be set
						"notifyDownload",                            // notify download - no other fields should be set
						"autoInstallAtMaintenanceTime",              // auto install at maintenance time - requires active_hours_start and active_hours_end to be set
						"autoInstallAndRebootAtMaintenanceTime",     // auto install and reboot at maintenance time - requires active_hours_start and active_hours_end to be set
						"autoInstallAndRebootAtScheduledTime",       // auto install and reboot at scheduled time - requires active_hours_start and active_hours_end to be set and update_weeks to be set
						"autoInstallAndRebootWithoutEndUserControl", // auto install and reboot without end user control - no other fields should be set
					),
				},
			},
			"delivery_optimization_mode": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "The Delivery Optimization Mode. Possible values are: UserDefined, HttpOnly, HttpWithPeeringNat, HttpWithPeeringPrivateGroup, HttpWithInternetPeering, SimpleDownload, BypassMode. UserDefined allows the user to set. Returned by default. Query parameters are not supported. Possible values are: userDefined, httpOnly, httpWithPeeringNat, httpWithPeeringPrivateGroup, httpWithInternetPeering, simpleDownload, bypassMode.",
				Validators: []validator.String{
					stringvalidator.OneOf(
						"userDefined",
						"httpOnly",
						"httpWithPeeringNat",
						"httpWithPeeringPrivateGroup",
						"httpWithInternetPeering",
						"simpleDownload",
						"bypassMode",
					),
				},
			},
			"prerelease_features": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "The Pre-Release Features. Possible values are: UserDefined, SettingsOnly, SettingsAndExperimentations, NotAllowed. UserDefined is the default value, no intent. Returned by default. Query parameters are not supported. Possible values are: userDefined, settingsOnly, settingsAndExperimentations, notAllowed.",
				Validators: []validator.String{
					stringvalidator.OneOf(
						"userDefined",
						"settingsOnly",
						"settingsAndExperimentations",
						"notAllowed",
					),
				},
			},
			"update_weeks": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Schedule the update installation on the weeks of the month. Possible values are: UserDefined, FirstWeek, SecondWeek, ThirdWeek, FourthWeek, EveryWeek. Returned by default. Query parameters are not supported. Possible values are: userDefined, firstWeek, secondWeek, thirdWeek, fourthWeek, everyWeek, unknownFutureValue.",
				Validators: []validator.String{
					stringvalidator.OneOf(
						"userDefined",
						"firstWeek",
						"secondWeek",
						"thirdWeek",
						"fourthWeek",
						"everyWeek",
						"unknownFutureValue",
					),
					validators.RequiredWith("automatic_update_mode", "autoInstallAndRebootAtScheduledTime"),
				},
			},
			"active_hours_start": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Active Hours Start. Part of the Installation Schedule.",
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(constants.TimeFormatHHMMSSRegex),
						"must be in format HH:MM:SS",
					),
					validators.RequiredWith("automatic_update_mode", "autoInstallAtMaintenanceTime"),
					validators.RequiredWith("automatic_update_mode", "autoInstallAndRebootAtMaintenanceTime"),
					validators.RequiredWith("automatic_update_mode", "autoInstallAndRebootAtScheduledTime"),
				},
			},
			"active_hours_end": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Active Hours End. Part of the Installation Schedule.",
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(constants.TimeFormatHHMMSSRegex),
						"must be in format HH:MM:SS",
					),
					validators.RequiredWith("automatic_update_mode", "autoInstallAtMaintenanceTime"),
					validators.RequiredWith("automatic_update_mode", "autoInstallAndRebootAtMaintenanceTime"),
					validators.RequiredWith("automatic_update_mode", "autoInstallAndRebootAtScheduledTime"),
				},
			},
			"user_pause_access": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("enabled"),
				MarkdownDescription: "Specifies whether to enable end user's access to pause software updates. Possible values are: NotConfigured, Enabled, Disabled. Returned by default. Query parameters are not supported. Possible values are: notConfigured, enabled, disabled.",
				Validators: []validator.String{
					stringvalidator.OneOf("notConfigured", "enabled", "disabled"),
				},
			},
			"user_windows_update_scan_access": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("enabled"),
				MarkdownDescription: "Specifies whether to disable user's access to scan Windows Update. Possible values are: NotConfigured, Enabled, Disabled. Returned by default. Query parameters are not supported. Possible values are: notConfigured, enabled, disabled.",
				Validators: []validator.String{
					stringvalidator.OneOf("notConfigured", "enabled", "disabled"),
				},
			},
			"update_notification_level": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("defaultNotifications"),
				MarkdownDescription: "Specifies what Windows Update notifications users see. Possible values are: NotConfigured, DefaultNotifications, RestartWarningsOnly, DisableAllNotifications. Returned by default. Query parameters are not supported. Possible values are: notConfigured, defaultNotifications, restartWarningsOnly, disableAllNotifications, unknownFutureValue.",
				Validators: []validator.String{
					stringvalidator.OneOf(
						"notConfigured",
						"defaultNotifications",
						"restartWarningsOnly",
						"disableAllNotifications",
					),
				},
			},
			"feature_updates_rollback_window_in_days": schema.Int32Attribute{
				Required:            true,
				MarkdownDescription: "The number of days after a Feature Update for which a rollback is valid with valid range from 2 to 60 days. Returned by default. Query parameters are not supported.",
			},
			"deadline_settings": schema.SingleNestedAttribute{
				Optional:            true,
				MarkdownDescription: "Settings for update installation deadlines and reboot behavior.",
				Attributes: map[string]schema.Attribute{
					"deadline_for_feature_updates_in_days": schema.Int32Attribute{
						Required:            true,
						MarkdownDescription: "Number of days before feature updates are installed automatically with valid range from 0 to 30 days. Returned by default. Query parameters are not supported.",
						Validators: []validator.Int32{
							int32validator.Between(0, 30),
						},
					},
					"deadline_for_quality_updates_in_days": schema.Int32Attribute{
						Required:            true,
						MarkdownDescription: "Number of days before quality updates are installed automatically with valid range from 0 to 30 days. Returned by default. Query parameters are not supported.",
						Validators: []validator.Int32{
							int32validator.Between(0, 30),
						},
					},
					"deadline_grace_period_in_days": schema.Int32Attribute{
						Required:            true,
						MarkdownDescription: "Number of days after deadline until restarts occur automatically with valid range from 0 to 7 days. Returned by default. Query parameters are not supported.",
						Validators: []validator.Int32{
							int32validator.Between(0, 7),
						},
					},
					"postpone_reboot_until_after_deadline": schema.BoolAttribute{
						Required:            true,
						MarkdownDescription: "When TRUE the device should wait until deadline for rebooting outside of active hours. When FALSE the device should not wait until deadline for rebooting outside of active hours. Returned by default. Query parameters are not supported.",
					},
				},
			},
			"engaged_restart_deadline_in_days": schema.Int32Attribute{
				Optional:            true,
				MarkdownDescription: "Deadline in days before automatically scheduling and executing a pending restart outside of active hours, with valid range from 2 to 30 days. Returned by default. Query parameters are not supported.",
			},
			"engaged_restart_snooze_schedule_in_days": schema.Int32Attribute{
				Optional:            true,
				MarkdownDescription: "Number of days a user can snooze Engaged Restart reminder notifications with valid range from 1 to 3 days. Returned by default. Query parameters are not supported.",
			},
			"engaged_restart_transition_schedule_in_days": schema.Int32Attribute{
				Optional:            true,
				MarkdownDescription: "Number of days before transitioning from Auto Restarts scheduled outside of active hours to Engaged Restart, which requires the user to schedule, with valid range from 0 to 30 days. Returned by default. Query parameters are not supported.",
			},
			"auto_restart_notification_dismissal": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("notConfigured"),
				MarkdownDescription: "Specify the method by which the auto-restart required notification is dismissed. Possible values are: NotConfigured, Automatic, User. Returned by default. Query parameters are not supported. Possible values are: notConfigured, automatic, user, unknownFutureValue.",
				Validators: []validator.String{
					stringvalidator.OneOf("notConfigured", "automatic", "user"),
				},
			},
			"schedule_restart_warning_in_hours": schema.Int32Attribute{
				Optional:            true,
				MarkdownDescription: "Specify the period for auto-restart warning reminder notifications. Supported values: 2, 4, 8, 12 or 24 (hours). Returned by default. Query parameters are not supported.",
				Validators: []validator.Int32{
					int32validator.OneOf(2, 4, 8, 12, 24),
				},
			},
			"schedule_imminent_restart_warning_in_minutes": schema.Int32Attribute{
				Optional:            true,
				MarkdownDescription: "Specify the period for auto-restart imminent warning notifications. Supported values: 15, 30 or 60 (minutes). Returned by default. Query parameters are not supported.",
				Validators: []validator.Int32{
					int32validator.OneOf(15, 30, 60),
				},
			},
			"engaged_restart_snooze_schedule_for_feature_updates_in_days": schema.Int32Attribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Number of days a user can snooze Engaged Restart reminder notifications for feature updates.",
			},
			"engaged_restart_transition_schedule_for_feature_updates_in_days": schema.Int32Attribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Number of days before transitioning from Auto Restarts scheduled outside of active hours to Engaged Restart for feature updates.",
			},
			"quality_updates_paused": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "When TRUE, assigned devices are paused from receiving quality updates for up to 35 days from the time you pause the ring. When FALSE, does not pause Quality Updates. Returned by default. Query parameters are not supported.",
			},
			"feature_updates_paused": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "When TRUE, assigned devices are paused from receiving feature updates for up to 35 days from the time you pause the ring. When FALSE, does not pause Feature Updates. Returned by default. Query parameters are not supported.s",
			},
			"feature_updates_pause_expiry_date_time": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The date and time when feature updates pause expires. This value is in ISO 8601 format, in UTC time.",
			},
			"feature_updates_rollback_start_date_time": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The date and time when feature updates rollback started. This value is in ISO 8601 format, in UTC time.",
			},
			"feature_updates_pause_start_date": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The date when feature updates are paused. This value is in ISO 8601 format, in UTC time.",
			},
			"quality_updates_pause_expiry_date_time": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The date and time when quality updates pause expires. This value is in ISO 8601 format, in UTC time.",
			},
			"quality_updates_pause_start_date": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The date when quality updates are paused. This value is in ISO 8601 format, in UTC time.",
			},
			"quality_updates_rollback_start_date_time": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The date and time when quality updates rollback started. This value is in ISO 8601 format, in UTC time.",
			},
			"assignments": commonschemagraphbeta.DeviceConfigurationWithAllGroupAssignmentsAndFilterSchema(),
			"timeouts":    commonschema.Timeouts(ctx),
		},
	}
}
