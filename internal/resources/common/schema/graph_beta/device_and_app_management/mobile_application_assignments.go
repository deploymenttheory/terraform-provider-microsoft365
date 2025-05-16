package schema

import (
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	planmodifiers "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/plan_modifiers"
)

func MobileAppAssignmentSchema() schema.ListNestedAttribute {
	return schema.ListNestedAttribute{
		Required: true,
		PlanModifiers: []planmodifier.List{
			listplanmodifier.UseStateForUnknown(),
			planmodifiers.MobileAppAssignmentsListModifier(),
		},
		NestedObject: schema.NestedAttributeObject{
			Attributes: map[string]schema.Attribute{
				"id": schema.StringAttribute{
					MarkdownDescription: "The ID of the app assignment associated with the Intune application.",
					Computed:            true,
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.UseStateForUnknown(),
					},
				},
				"intent": schema.StringAttribute{
					Required: true,
					MarkdownDescription: "The Intune app install intent defined by the admin. Possible values are:\n\n" +
						"- **available**: App is available for users to install\n" +
						"- **required**: App is required and will be automatically installed\n" +
						"- **uninstall**: App will be uninstalled\n" +
						"- **availableWithoutEnrollment**: App is available without Intune device enrollment",
					Validators: []validator.String{
						stringvalidator.OneOf(
							"available",
							"required",
							"uninstall",
							"availableWithoutEnrollment",
						),
					},
				},
				"source": schema.StringAttribute{
					MarkdownDescription: "The resource type which is the source for the assignment. Possible values are: direct, policySets.",
					Required:            true,
				},
				"source_id": schema.StringAttribute{
					MarkdownDescription: "The identifier of the source of the assignment.",
					Optional:            true,
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.UseStateForUnknown(),
					},
				},
				"target": schema.SingleNestedAttribute{
					Required: true,
					Attributes: map[string]schema.Attribute{
						"target_type": schema.StringAttribute{
							Required: true,
							MarkdownDescription: "The target group type for the application assignment. Possible values are:\n\n" +
								"- **allDevices**: Target all devices in the tenant\n" +
								"- **allLicensedUsers**: Target all licensed users in the tenant\n" +
								"- **androidFotaDeployment**: Target Android FOTA deployment\n" +
								"- **configurationManagerCollection**: Target System Centre Configuration Manager collection\n" +
								"- **exclusionGroupAssignment**: Target a specific Entra ID group for exclusion\n" +
								"- **groupAssignment**: Target a specific Entra ID group",
							Validators: []validator.String{
								stringvalidator.OneOf(
									"allDevices",
									"allLicensedUsers",
									"androidFotaDeployment",
									"configurationManagerCollection",
									"exclusionGroupAssignment",
									"groupAssignment",
								),
							},
						},
						"group_id": schema.StringAttribute{
							MarkdownDescription: "The entra ID group ID for the application assignment target. Required when target_type is 'groupAssignment', 'exclusionGroupAssignment', or 'androidFotaDeployment'.",
							Optional:            true,
							Validators: []validator.String{
								stringvalidator.RegexMatches(
									regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`),
									"Must be a valid GUID format (xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx)",
								),
							},
						},
						"collection_id": schema.StringAttribute{
							MarkdownDescription: "The SCCM group collection ID for the application assignment target. Default collections start with 'SMS', while custom collections start with your site code (e.g., 'MEM').",
							Optional:            true,
							Validators: []validator.String{
								// Validator for SCCM collection ID format
								stringvalidator.RegexMatches(
									regexp.MustCompile(`^[A-Za-z]{2,8}[0-9A-Za-z]{8}$`),
									"Must be a valid SCCM collection ID format. Default collections start with 'SMS' followed by an alphanumeric ID. Custom collections start with your site code (e.g., 'MEM') followed by an alphanumeric ID.",
								),
							},
						},
						"device_and_app_management_assignment_filter_id": schema.StringAttribute{
							MarkdownDescription: "The Id of the scope filter applied to the target assignment.",
							Optional:            true,
							Validators: []validator.String{
								stringvalidator.RegexMatches(
									regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`),
									"Must be a valid GUID format (xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx)",
								),
							},
						},
						"device_and_app_management_assignment_filter_type": schema.StringAttribute{
							Optional: true,
							Computed: true,
							Default:  stringdefault.StaticString("none"),
							MarkdownDescription: "The type of scope filter for the target assignment. Defaults to 'none'. Possible values are:\n\n" +
								"- **include**: Only include devices or users matching the filter\n" +
								"- **exclude**: Exclude devices or users matching the filter\n" +
								"- **none**: No assignment filter applied",
							Validators: []validator.String{
								stringvalidator.OneOf(
									"include",
									"exclude",
									"none",
								),
							},
						},
					},
				},
				"settings": schema.SingleNestedAttribute{
					Optional: true,
					Attributes: map[string]schema.Attribute{
						"android_managed_store": schema.SingleNestedAttribute{
							Optional: true,
							Attributes: map[string]schema.Attribute{
								"android_managed_store_app_track_ids": schema.ListAttribute{
									ElementType:         types.StringType,
									Optional:            true,
									MarkdownDescription: "The track IDs to enable for this app assignment.",
									Computed:            true,
									PlanModifiers: []planmodifier.List{
										listplanmodifier.UseStateForUnknown(),
									},
								},
								"auto_update_mode": schema.StringAttribute{
									Optional: true,
									MarkdownDescription: "The prioritization of automatic updates for this app assignment. Possible values are:\n\n" +
										"- **default**: Default auto-update mode\n" +
										"- **postponed**: Updates are postponed\n" +
										"- **priority**: Updates are prioritized\n" +
										"- **unknownFutureValue**: Reserved for future use",
									Validators: []validator.String{
										stringvalidator.OneOf(
											"default",
											"postponed",
											"priority",
											"unknownFutureValue",
										),
									},
								},
							},
						},
						"ios_lob": schema.SingleNestedAttribute{
							Optional: true,
							Attributes: map[string]schema.Attribute{
								"is_removable": schema.BoolAttribute{
									Optional:            true,
									Computed:            true,
									MarkdownDescription: "When TRUE, indicates that the app can be uninstalled by the user. When FALSE, indicates that the app cannot be uninstalled by the user. By default, this property is set to TRUE.",
									Default:             booldefault.StaticBool(true),
								},
								"prevent_managed_app_backup": schema.BoolAttribute{
									Optional:            true,
									Computed:            true,
									MarkdownDescription: "When TRUE, indicates that the app should not be backed up to iCloud. When FALSE, indicates that the app may be backed up to iCloud. By default, this property is set to FALSE.",
									Default:             booldefault.StaticBool(false),
								},
								"uninstall_on_device_removal": schema.BoolAttribute{
									Optional:            true,
									Computed:            true,
									MarkdownDescription: "When TRUE, indicates that the app should be uninstalled when the device is removed from Intune. When FALSE, indicates that the app will not be uninstalled when the device is removed from Intune. By default, this property is set to TRUE.",
									Default:             booldefault.StaticBool(true),
								},
								"vpn_configuration_id": schema.StringAttribute{
									Optional:            true,
									MarkdownDescription: "This is the unique identifier (Id) of the VPN Configuration to apply to the app.",
								},
							},
						},
						"ios_store": schema.SingleNestedAttribute{
							Optional: true,
							Attributes: map[string]schema.Attribute{
								"is_removable": schema.BoolAttribute{
									Optional:            true,
									Computed:            true,
									MarkdownDescription: "When TRUE, indicates that the app can be uninstalled by the user. When FALSE, indicates that the app cannot be uninstalled by the user. By default, this property is set to TRUE.",
									Default:             booldefault.StaticBool(true),
								},
								"prevent_managed_app_backup": schema.BoolAttribute{
									Optional:            true,
									Computed:            true,
									MarkdownDescription: "When TRUE, indicates that the app should not be backed up to iCloud. When FALSE, indicates that the app may be backed up to iCloud. By default, this property is set to FALSE.",
									Default:             booldefault.StaticBool(false),
								},
								"uninstall_on_device_removal": schema.BoolAttribute{
									Optional:            true,
									Computed:            true,
									MarkdownDescription: "When TRUE, indicates that the app should be uninstalled when the device is removed from Intune. When FALSE, indicates that the app will not be uninstalled when the device is removed from Intune. By default, this property is set to TRUE.",
									Default:             booldefault.StaticBool(true),
								},
								"vpn_configuration_id": schema.StringAttribute{
									Optional:            true,
									MarkdownDescription: "This is the unique identifier (Id) of the VPN Configuration to apply to the app.",
								},
							},
						},
						"ios_vpp": schema.SingleNestedAttribute{
							Optional: true,
							Attributes: map[string]schema.Attribute{
								"is_removable": schema.BoolAttribute{
									Optional:            true,
									Computed:            true,
									MarkdownDescription: "Whether or not the app can be removed by the user. By default, this property is set to FALSE.",
									Default:             booldefault.StaticBool(false),
								},
								"prevent_auto_app_update": schema.BoolAttribute{
									Optional:            true,
									Computed:            true,
									MarkdownDescription: "When TRUE, indicates that the app should not be automatically updated with the latest version from Apple app store. When FALSE, indicates that the app may be auto updated. By default, this property is set to FALSE.",
									Default:             booldefault.StaticBool(false),
								},
								"prevent_managed_app_backup": schema.BoolAttribute{
									Optional:            true,
									Computed:            true,
									MarkdownDescription: "When TRUE, indicates that the app should not be backed up to iCloud. When FALSE, indicates that the app may be backed up to iCloud. By default, this property is set to FALSE.",
									Default:             booldefault.StaticBool(false),
								},
								"uninstall_on_device_removal": schema.BoolAttribute{
									Optional:            true,
									Computed:            true,
									MarkdownDescription: "Whether or not to uninstall the app when device is removed from Intune. By default, this property is set to FALSE.",
									Default:             booldefault.StaticBool(false),
								},
								"use_device_licensing": schema.BoolAttribute{
									Optional:            true,
									Computed:            true,
									MarkdownDescription: "Whether or not to use device licensing. By default, this property is set to FALSE.",
									Default:             booldefault.StaticBool(false),
								},
								"vpn_configuration_id": schema.StringAttribute{
									Optional:            true,
									MarkdownDescription: "The VPN Configuration Id to apply for this app.",
								},
							},
						},
						"macos_lob": schema.SingleNestedAttribute{
							Optional: true,
							Attributes: map[string]schema.Attribute{
								"uninstall_on_device_removal": schema.BoolAttribute{
									Optional:            true,
									Computed:            true,
									MarkdownDescription: "When TRUE, the macOS LOB app will be uninstalled when the device is removed from Intune management. When FALSE, the macOS LOB app will not be uninstalled when the device is removed from management. By default, this property is set to FALSE.",
									Default:             booldefault.StaticBool(false),
								},
							},
						},
						"macos_vpp": schema.SingleNestedAttribute{
							Optional: true,
							Attributes: map[string]schema.Attribute{
								"prevent_auto_app_update": schema.BoolAttribute{
									Optional:            true,
									Computed:            true,
									MarkdownDescription: "When TRUE, indicates that the app should not be automatically updated with the latest version from Apple app store. When FALSE, indicates that the app may be auto updated. By default, this property is set to null which internally is treated as FALSE.",
									Default:             booldefault.StaticBool(false),
								},
								"prevent_managed_app_backup": schema.BoolAttribute{
									Optional:            true,
									Computed:            true,
									MarkdownDescription: "When TRUE, indicates that the app should not be backed up to iCloud. When FALSE, indicates that the app may be backed up to iCloud. By default, this property is set to null which internally is treated as FALSE.",
									Default:             booldefault.StaticBool(false),
								},
								"uninstall_on_device_removal": schema.BoolAttribute{
									Optional:            true,
									Computed:            true,
									MarkdownDescription: "When TRUE, the macOS VPP app will be uninstalled when the device is removed from Intune management. When FALSE, the macOS VPP app will not be uninstalled when the device is removed from management. By default, this property is set to FALSE.",
									Default:             booldefault.StaticBool(false),
								},
								"use_device_licensing": schema.BoolAttribute{
									Optional:            true,
									Computed:            true,
									MarkdownDescription: "When TRUE indicates that the macOS VPP app should use device-based licensing. When FALSE indicates that the macOS VPP app should use user-based licensing. By default, this property is set to FALSE.",
									Default:             booldefault.StaticBool(false),
								},
							},
						},
						"microsoft_store_for_business": schema.SingleNestedAttribute{
							Optional: true,
							Attributes: map[string]schema.Attribute{
								"use_device_context": schema.BoolAttribute{
									MarkdownDescription: "When TRUE, indicates that device execution context will be used for the Microsoft Store for Business mobile app. " +
										"When FALSE, indicates that user context will be used for the Microsoft Store for Business mobile app. " +
										"By default, this property is set to FALSE. Once this property has been set to TRUE it cannot be changed.",
									Optional: true,
									Computed: true,
									Default:  booldefault.StaticBool(false),
								},
							},
						},
						"win32_catalog": schema.SingleNestedAttribute{
							Optional: true,
							Attributes: map[string]schema.Attribute{
								"auto_update_settings": schema.SingleNestedAttribute{
									MarkdownDescription: "The auto-update settings to apply for this app assignment.",
									Optional:            true,
									Attributes: map[string]schema.Attribute{
										"auto_update_superseded_apps_state": schema.StringAttribute{
											Optional: true,
											MarkdownDescription: "The auto-update superseded apps setting for the app assignment. " +
												"Default value is notConfigured. Possible values are:\n\n" +
												"- **notConfigured**: Auto-update is not configured\n" +
												"- **enabled**: Auto-update is enabled\n" +
												"- **unknownFutureValue**: Reserved for future use",
											Validators: []validator.String{
												stringvalidator.OneOf(
													"notConfigured",
													"enabled",
													"unknownFutureValue",
												),
											},
										},
									},
								},
								"delivery_optimization_priority": schema.StringAttribute{
									Optional: true,
									MarkdownDescription: "The delivery optimization priority for this app assignment. This setting is not " +
										"supported in National Cloud environments. Possible values are:\n\n" +
										"- **notConfigured**: Not configured or background normal delivery optimization priority\n" +
										"- **foreground**: Foreground delivery optimization priority",
									Validators: []validator.String{
										stringvalidator.OneOf(
											"notConfigured",
											"foreground",
										),
									},
								},
								"install_time_settings": schema.SingleNestedAttribute{
									MarkdownDescription: "The install time settings to apply for this app assignment.",
									Optional:            true,
									Attributes: map[string]schema.Attribute{
										"use_local_time": schema.BoolAttribute{
											Optional:            true,
											MarkdownDescription: "Whether the local device time or UTC time should be used when determining the available and deadline times.",
										},
										"deadline_date_time": schema.StringAttribute{
											Optional:            true,
											MarkdownDescription: "The time at which the app should be installed.",
										},
										"start_date_time": schema.StringAttribute{
											Optional:            true,
											MarkdownDescription: "The time at which the app should be available for installation.",
										},
									},
								},
								"notifications": schema.StringAttribute{
									Optional: true,
									MarkdownDescription: "The notification status for this app assignment. Possible values are:\n\n" +
										"- **showAll**: Show all notifications\n" +
										"- **showReboot**: Show only reboot notifications\n" +
										"- **hideAll**: Hide all notifications",
									Validators: []validator.String{
										stringvalidator.OneOf(
											"showAll",
											"showReboot",
											"hideAll",
										),
									},
								},
								"restart_settings": schema.SingleNestedAttribute{
									MarkdownDescription: "The reboot settings to apply for this app assignment.",
									Optional:            true,
									Attributes: map[string]schema.Attribute{
										"grace_period_in_minutes": schema.Int32Attribute{
											Optional:            true,
											MarkdownDescription: "The number of minutes to wait before restarting the device after an app installation.",
										},
										"countdown_display_before_restart_in_minutes": schema.Int32Attribute{
											Optional:            true,
											MarkdownDescription: "The number of minutes before the restart time to display the countdown dialog for pending restarts.",
										},
										"restart_notification_snooze_duration_in_minutes": schema.Int32Attribute{
											Optional:            true,
											MarkdownDescription: "The number of minutes to snooze the restart notification dialog when the snooze button is selected.",
										},
									},
								},
							},
						},
						"win32_lob": schema.SingleNestedAttribute{
							Optional: true,
							Attributes: map[string]schema.Attribute{
								"auto_update_settings": schema.SingleNestedAttribute{
									MarkdownDescription: "The auto-update settings to apply for this app assignment.",
									Optional:            true,
									Attributes: map[string]schema.Attribute{
										"auto_update_superseded_apps_state": schema.StringAttribute{
											Optional: true,
											MarkdownDescription: "The auto-update superseded apps setting for the app assignment. " +
												"Default value is notConfigured. Possible values are:\n\n" +
												"- **notConfigured**: Auto-update is not configured\n" +
												"- **enabled**: Auto-update is enabled\n" +
												"- **unknownFutureValue**: Reserved for future use",
											Validators: []validator.String{
												stringvalidator.OneOf(
													"notConfigured",
													"enabled",
													"unknownFutureValue",
												),
											},
										},
									},
								},
								"delivery_optimization_priority": schema.StringAttribute{
									Optional: true,
									MarkdownDescription: "The delivery optimization priority for this app assignment. This setting is not" +
										"supported in National Cloud environments. Possible values are: notConfigured, foreground." +
										"- **notConfigured**: Not configured or background normal delivery optimization priority.\n" +
										"- **foreground**: Foreground delivery optimization priority.",
									Validators: []validator.String{
										stringvalidator.OneOf(
											"notConfigured",
											"foreground",
										),
									},
								},
								"install_time_settings": schema.SingleNestedAttribute{
									MarkdownDescription: "The install time settings to apply for this app assignment.",
									Optional:            true,
									Attributes: map[string]schema.Attribute{
										"use_local_time": schema.BoolAttribute{
											Optional:            true,
											MarkdownDescription: "Whether the local device time or UTC time should be used when determining the available and deadline times.",
										},
										"deadline_date_time": schema.StringAttribute{
											Optional:            true,
											MarkdownDescription: "The time at which the app should be installed.",
										},
										"start_date_time": schema.StringAttribute{
											Optional:            true,
											MarkdownDescription: "The time at which the app should be available for installation.",
										},
									},
								},
								"notifications": schema.StringAttribute{
									Optional: true,
									MarkdownDescription: "The notification status for this app assignment. Possible values are:\n\n" +
										"- **showAll**: Show all notifications\n" +
										"- **showReboot**: Show only reboot notifications\n" +
										"- **hideAll**: Hide all notifications",
									Validators: []validator.String{
										stringvalidator.OneOf(
											"showAll",
											"showReboot",
											"hideAll",
										),
									},
								},
								"restart_settings": schema.SingleNestedAttribute{
									MarkdownDescription: "The reboot settings to apply for this app assignment.",
									Optional:            true,
									Attributes: map[string]schema.Attribute{
										"grace_period_in_minutes": schema.Int32Attribute{
											Optional:            true,
											MarkdownDescription: "The number of minutes to wait before restarting the device after an app installation.",
										},
										"countdown_display_before_restart_in_minutes": schema.Int32Attribute{
											Optional:            true,
											MarkdownDescription: "The number of minutes before the restart time to display the countdown dialog for pending restarts.",
										},
										"restart_notification_snooze_duration_in_minutes": schema.Int32Attribute{
											Optional:            true,
											MarkdownDescription: "The number of minutes to snooze the restart notification dialog when the snooze button is selected.",
										},
									},
								},
							},
						},
						"windows_app_x": schema.SingleNestedAttribute{
							Optional: true,
							Attributes: map[string]schema.Attribute{
								"use_device_context": schema.BoolAttribute{
									MarkdownDescription: "When TRUE, indicates that device execution context will be used for the AppX mobile app. When FALSE, indicates that user context will be used for the AppX mobile app. By default, this property is set to FALSE. Once this property has been set to TRUE it cannot be changed.",
									Optional:            true,
									Computed:            true,
									Default:             booldefault.StaticBool(false),
								},
							},
						},
						"windows_universal_app_x": schema.SingleNestedAttribute{
							Optional: true,
							Attributes: map[string]schema.Attribute{
								"use_device_context": schema.BoolAttribute{
									MarkdownDescription: "If true, uses device execution context for Windows Universal AppX mobile app. Device-context install is not allowed when this type of app is targeted with Available intent. Defaults to false.",
									Optional:            true,
									Computed:            true,
									Default:             booldefault.StaticBool(false),
								},
							},
						},
						"win_get": schema.SingleNestedAttribute{
							Optional: true,
							Attributes: map[string]schema.Attribute{
								"install_time_settings": schema.SingleNestedAttribute{
									Optional: true,
									Attributes: map[string]schema.Attribute{
										"use_local_time": schema.BoolAttribute{
											MarkdownDescription: "Whether the local device time or UTC time should be used when determining the deadline times.",
											Optional:            true,
										},
										"deadline_date_time": schema.StringAttribute{
											MarkdownDescription: "The time at which the app should be installed.",
											Optional:            true,
										},
									},
								},
								"notifications": schema.StringAttribute{
									MarkdownDescription: "The notification settings for the assignment. The supported values are 'showAll', 'showReboot', 'hideAll'.",
									Optional:            true,
									Validators: []validator.String{
										stringvalidator.OneOf(
											"showAll",
											"showReboot",
											"hideAll",
										),
									},
								},
								"restart_settings": schema.SingleNestedAttribute{
									Optional: true,
									Attributes: map[string]schema.Attribute{
										"grace_period_in_minutes": schema.Int32Attribute{
											MarkdownDescription: "The number of minutes to wait before restarting the device after an app installation.",
											Optional:            true,
										},
										"countdown_display_before_restart_in_minutes": schema.Int32Attribute{
											MarkdownDescription: "The number of minutes before the restart time to display the countdown dialog for pending restarts.",
											Optional:            true,
										},
										"restart_notification_snooze_duration_in_minutes": schema.Int32Attribute{
											MarkdownDescription: "The number of minutes to snooze the restart notification dialog when the snooze button is selected.",
											Optional:            true,
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}
