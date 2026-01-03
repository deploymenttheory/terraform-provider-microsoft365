package graphBetaManagedDevice

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

const (
	DataSourceName = "microsoft365_graph_beta_device_management_managed_device"
	ReadTimeout    = 180
)

var (
	_ datasource.DataSource              = &ManagedDeviceDataSource{}
	_ datasource.DataSourceWithConfigure = &ManagedDeviceDataSource{}
)

func NewManagedDeviceDataSource() datasource.DataSource {
	return &ManagedDeviceDataSource{
		ReadPermissions: []string{
			"DeviceManagementManagedDevices.Read.All",
		},
	}
}

type ManagedDeviceDataSource struct {
	client *msgraphbetasdk.GraphServiceClient

	ReadPermissions []string
}

func (d *ManagedDeviceDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = DataSourceName
}

func (d *ManagedDeviceDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.client = client.SetGraphBetaClientForDataSource(ctx, req, resp, DataSourceName)
}

func (d *ManagedDeviceDataSource) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Retrieves Managed Devices from Microsoft Intune using the `/deviceManagement/managedDevices` endpoint. This data source enables querying managed devices with advanced filtering capabilities.",
		Attributes: map[string]schema.Attribute{
			"filter_type": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Type of filter to apply. Valid values are: `all`, `id`, `device_name`, `serial_number`, `user_id`, `odata`.",
				Validators: []validator.String{
					stringvalidator.OneOf("all", "id", "device_name", "serial_number", "user_id", "odata"),
				},
			},
			"filter_value": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Value to filter by. Not required when filter_type is 'all' or 'odata'.",
			},
			"odata_filter": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "OData $filter parameter for filtering results. Only used when filter_type is 'odata'. Example: operatingSystem eq 'Windows'.",
			},
			"odata_top": schema.Int32Attribute{
				Optional:            true,
				MarkdownDescription: "OData $top parameter to limit the number of results. Only used when filter_type is 'odata'.",
			},
			"odata_skip": schema.Int32Attribute{
				Optional:            true,
				MarkdownDescription: "OData $skip parameter for pagination. Only used when filter_type is 'odata'.",
			},
			"odata_select": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "OData $select parameter to specify which fields to include. Only used when filter_type is 'odata'.",
			},
			"odata_orderby": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "OData $orderby parameter to sort results. Only used when filter_type is 'odata'. Example: deviceName.",
			},
			"odata_count": schema.BoolAttribute{
				Optional:            true,
				MarkdownDescription: "OData $count parameter to include count of total results. Only used when filter_type is 'odata'.",
			},
			"odata_search": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "OData $search parameter for full-text search. Only used when filter_type is 'odata'.",
			},
			"odata_expand": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "OData $expand parameter to include related entities. Only used when filter_type is 'odata'.",
			},
			"items": schema.ListNestedAttribute{
				Computed:            true,
				MarkdownDescription: "The list of managed devices that match the filter criteria.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The unique identifier for the managed device.",
						},
						"user_id": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The unique identifier for the user associated with the device.",
						},
						"device_name": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The name of the device as displayed in Intune.",
						},
						"hardware_information": schema.SingleNestedAttribute{
							Computed:            true,
							MarkdownDescription: "Hardware information for the device.",
							Attributes: map[string]schema.Attribute{
								"serial_number": schema.StringAttribute{
									Computed: true, MarkdownDescription: "Device serial number."},
								"total_storage_space": schema.Int64Attribute{
									Computed: true, MarkdownDescription: "Total storage space on the device in bytes."},
								"free_storage_space": schema.Int64Attribute{
									Computed: true, MarkdownDescription: "Free storage space on the device in bytes."},
								"imei": schema.StringAttribute{
									Computed: true, MarkdownDescription: "International Mobile Equipment Identity (IMEI) of the device."},
								"meid": schema.StringAttribute{
									Computed: true, MarkdownDescription: "Mobile Equipment Identifier (MEID) of the device."},
								"manufacturer": schema.StringAttribute{
									Computed: true, MarkdownDescription: "Device manufacturer."},
								"model": schema.StringAttribute{
									Computed: true, MarkdownDescription: "Device model."},
								"phone_number": schema.StringAttribute{
									Computed: true, MarkdownDescription: "Phone number associated with the device."},
								"subscriber_carrier": schema.StringAttribute{
									Computed: true, MarkdownDescription: "Mobile carrier for the device's SIM card."},
								"cellular_technology": schema.StringAttribute{
									Computed: true, MarkdownDescription: "Cellular technology used by the device (e.g., LTE, 5G)."},
								"wifi_mac": schema.StringAttribute{
									Computed: true, MarkdownDescription: "Wi-Fi MAC address of the device."},
								"operating_system_language": schema.StringAttribute{
									Computed: true, MarkdownDescription: "Language of the device's operating system."},
								"is_supervised": schema.BoolAttribute{
									Computed: true, MarkdownDescription: "Whether the device is supervised (Apple devices only)."},
								"is_encrypted": schema.BoolAttribute{
									Computed: true, MarkdownDescription: "Whether the device storage is encrypted."},
								"battery_serial_number": schema.StringAttribute{
									Computed: true, MarkdownDescription: "Serial number of the device's battery."},
								"battery_health_percentage": schema.Int64Attribute{
									Computed: true, MarkdownDescription: "Battery health as a percentage."},
								"battery_charge_cycles": schema.Int64Attribute{
									Computed: true, MarkdownDescription: "Number of battery charge cycles."},
								"is_shared_device": schema.BoolAttribute{
									Computed: true, MarkdownDescription: "Whether the device is a shared device."},
								"shared_device_cached_users": schema.ListNestedAttribute{
									Computed:            true,
									MarkdownDescription: "List of users cached on a shared device.",
									NestedObject: schema.NestedAttributeObject{
										Attributes: map[string]schema.Attribute{
											"user_principal_name": schema.StringAttribute{
												Computed: true, MarkdownDescription: "User principal name of the cached user."},
											"data_to_sync": schema.BoolAttribute{
												Computed: true, MarkdownDescription: "Whether there is data to sync for the user."},
											"data_quota": schema.Int64Attribute{
												Computed: true, MarkdownDescription: "Data quota for the user in MB."},
											"data_used": schema.Int64Attribute{
												Computed: true, MarkdownDescription: "Data used by the user in MB."},
										},
									},
								},
								"tpm_specification_version": schema.StringAttribute{
									Computed: true, MarkdownDescription: "TPM specification version."},
								"operating_system_edition": schema.StringAttribute{
									Computed: true, MarkdownDescription: "Edition of the device's operating system."},
								"device_full_qualified_domain_name": schema.StringAttribute{
									Computed: true, MarkdownDescription: "Fully qualified domain name of the device."},
								"device_guard_virtualization_based_security_hardware_requirement_state": schema.StringAttribute{
									Computed: true, MarkdownDescription: "Device Guard VBS hardware requirement state."},
								"device_guard_virtualization_based_security_state": schema.StringAttribute{
									Computed: true, MarkdownDescription: "Device Guard VBS state."},
								"device_guard_local_system_authority_credential_guard_state": schema.StringAttribute{
									Computed: true, MarkdownDescription: "Device Guard LSA Credential Guard state."},
								"os_build_number": schema.StringAttribute{
									Computed: true, MarkdownDescription: "Operating system build number."},
								"operating_system_product_type": schema.Int64Attribute{
									Computed: true, MarkdownDescription: "Product type of the operating system."},
								"ip_address_v4": schema.StringAttribute{
									Computed: true, MarkdownDescription: "IPv4 address of the device."},
								"subnet_address": schema.StringAttribute{
									Computed: true, MarkdownDescription: "Subnet address of the device."},
								"esim_identifier": schema.StringAttribute{
									Computed: true, MarkdownDescription: "eSIM identifier for the device."},
								"system_management_bios_version": schema.StringAttribute{
									Computed: true, MarkdownDescription: "System Management BIOS version."},
								"tpm_manufacturer": schema.StringAttribute{
									Computed: true, MarkdownDescription: "TPM manufacturer."},
								"tpm_version": schema.StringAttribute{
									Computed: true, MarkdownDescription: "TPM version."},
								"wired_ipv4_addresses": schema.ListAttribute{
									ElementType: types.StringType, Computed: true, MarkdownDescription: "List of wired IPv4 addresses for the device."},
								"battery_level_percentage": schema.Float64Attribute{
									Computed: true, MarkdownDescription: "Battery level as a percentage."},
								"resident_users_count": schema.Int64Attribute{
									Computed: true, MarkdownDescription: "Number of resident users on the device."},
								"product_name": schema.StringAttribute{
									Computed: true, MarkdownDescription: "Product name of the device."},
								"device_licensing_status": schema.StringAttribute{
									Computed: true, MarkdownDescription: "Device licensing status."},
								"device_licensing_last_error_code": schema.Int64Attribute{
									Computed: true, MarkdownDescription: "Last error code for device licensing."},
								"device_licensing_last_error_description": schema.StringAttribute{
									Computed: true, MarkdownDescription: "Last error description for device licensing."},
							},
						},
						"owner_type": schema.StringAttribute{
							Computed: true, MarkdownDescription: "Owner type of the device (e.g., company, personal)."},
						"managed_device_owner_type": schema.StringAttribute{
							Computed: true, MarkdownDescription: "Managed device owner type (e.g., company, personal)."},
						"device_action_results": schema.ListNestedAttribute{
							Computed:            true,
							MarkdownDescription: "List of device action results for the device.",
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"action_name": schema.StringAttribute{
										Computed: true, MarkdownDescription: "Name of the action performed on the device."},
									"action_state": schema.StringAttribute{
										Computed: true, MarkdownDescription: "State of the action (e.g., pending, completed)."},
									"start_date_time": schema.StringAttribute{
										Computed: true, MarkdownDescription: "Start time of the action."},
									"last_updated_date_time": schema.StringAttribute{
										Computed: true, MarkdownDescription: "Last update time of the action."},
								},
							},
						},
						"management_state": schema.StringAttribute{
							Computed: true, MarkdownDescription: "Management state of the device (e.g., retirePending, managed)."},
						"enrolled_date_time": schema.StringAttribute{
							Computed: true, MarkdownDescription: "Date and time when the device was enrolled."},
						"last_sync_date_time": schema.StringAttribute{
							Computed: true, MarkdownDescription: "Last time the device synced with Intune."},
						"chassis_type": schema.StringAttribute{
							Computed: true, MarkdownDescription: "Chassis type of the device (e.g., desktop, laptop)."},
						"operating_system": schema.StringAttribute{
							Computed: true, MarkdownDescription: "Operating system of the device."},
						"device_type": schema.StringAttribute{
							Computed: true, MarkdownDescription: "Type of the device (e.g., windowsRT, windows)."},
						"compliance_state": schema.StringAttribute{
							Computed: true, MarkdownDescription: "Compliance state of the device (e.g., compliant, noncompliant)."},
						"jail_broken": schema.StringAttribute{
							Computed: true, MarkdownDescription: "Indicates if the device is jailbroken (for iOS devices)."},
						"management_agent": schema.StringAttribute{
							Computed: true, MarkdownDescription: "Management agent used for the device (e.g., mdm, eas)."},
						"os_version": schema.StringAttribute{
							Computed: true, MarkdownDescription: "Operating system version of the device."},
						"eas_activated": schema.BoolAttribute{
							Computed: true, MarkdownDescription: "Whether Exchange ActiveSync is activated on the device."},
						"eas_device_id": schema.StringAttribute{
							Computed: true, MarkdownDescription: "Exchange ActiveSync device ID."},
						"eas_activation_date_time": schema.StringAttribute{
							Computed: true, MarkdownDescription: "Date and time when Exchange ActiveSync was activated."},
						"aad_registered": schema.BoolAttribute{
							Computed: true, MarkdownDescription: "Whether the device is Azure AD registered."},
						"azure_ad_registered": schema.BoolAttribute{
							Computed: true, MarkdownDescription: "Whether the device is Azure AD registered (legacy field)."},
						"device_enrollment_type": schema.StringAttribute{
							Computed: true, MarkdownDescription: "Type of device enrollment (e.g., userEnrollment, deviceEnrollmentManager)."},
						"lost_mode_state": schema.StringAttribute{
							Computed: true, MarkdownDescription: "State of lost mode on the device (e.g., enabled, disabled)."},
						"activation_lock_bypass_code": schema.StringAttribute{
							Computed: true, MarkdownDescription: "Activation lock bypass code for the device."},
						"email_address": schema.StringAttribute{
							Computed: true, MarkdownDescription: "Email address associated with the device."},
						"azure_active_directory_device_id": schema.StringAttribute{
							Computed: true, MarkdownDescription: "Azure Active Directory device ID."},
						"azure_ad_device_id": schema.StringAttribute{
							Computed: true, MarkdownDescription: "Azure AD device ID (legacy field)."},
						"device_registration_state": schema.StringAttribute{
							Computed: true, MarkdownDescription: "Registration state of the device."},
						"device_category_display_name": schema.StringAttribute{
							Computed: true, MarkdownDescription: "Display name of the device category."},
						"is_supervised": schema.BoolAttribute{
							Computed: true, MarkdownDescription: "Whether the device is supervised (Apple devices only)."},
						"exchange_last_successful_sync_date_time": schema.StringAttribute{
							Computed: true, MarkdownDescription: "Last successful Exchange sync date and time."},
						"exchange_access_state": schema.StringAttribute{
							Computed: true, MarkdownDescription: "Exchange access state for the device."},
						"exchange_access_state_reason": schema.StringAttribute{
							Computed: true, MarkdownDescription: "Reason for the Exchange access state."},
						"remote_assistance_session_url": schema.StringAttribute{
							Computed: true, MarkdownDescription: "URL for the remote assistance session."},
						"remote_assistance_session_error_details": schema.StringAttribute{
							Computed: true, MarkdownDescription: "Error details for the remote assistance session."},
						"is_encrypted": schema.BoolAttribute{
							Computed: true, MarkdownDescription: "Whether the device storage is encrypted."},
						"user_principal_name": schema.StringAttribute{
							Computed: true, MarkdownDescription: "User principal name associated with the device."},
						"model": schema.StringAttribute{
							Computed: true, MarkdownDescription: "Device model."},
						"manufacturer": schema.StringAttribute{
							Computed: true, MarkdownDescription: "Device manufacturer."},
						"imei": schema.StringAttribute{
							Computed: true, MarkdownDescription: "International Mobile Equipment Identity (IMEI) of the device."},
						"compliance_grace_period_expiration_date_time": schema.StringAttribute{
							Computed: true, MarkdownDescription: "Expiration date and time for the compliance grace period."},
						"serial_number": schema.StringAttribute{
							Computed: true, MarkdownDescription: "Device serial number."},
						"phone_number": schema.StringAttribute{
							Computed: true, MarkdownDescription: "Phone number associated with the device."},
						"android_security_patch_level": schema.StringAttribute{
							Computed: true, MarkdownDescription: "Android security patch level on the device."},
						"user_display_name": schema.StringAttribute{
							Computed: true, MarkdownDescription: "Display name of the user associated with the device."},
						"configuration_manager_client_enabled_features": schema.SingleNestedAttribute{
							Computed:            true,
							MarkdownDescription: "Configuration Manager client enabled features.",
							Attributes: map[string]schema.Attribute{
								"inventory":                   schema.BoolAttribute{Computed: true, MarkdownDescription: "Whether inventory is enabled."},
								"modern_apps":                 schema.BoolAttribute{Computed: true, MarkdownDescription: "Whether modern apps are enabled."},
								"resource_access":             schema.BoolAttribute{Computed: true, MarkdownDescription: "Whether resource access is enabled."},
								"device_configuration":        schema.BoolAttribute{Computed: true, MarkdownDescription: "Whether device configuration is enabled."},
								"compliance_policy":           schema.BoolAttribute{Computed: true, MarkdownDescription: "Whether compliance policy is enabled."},
								"windows_update_for_business": schema.BoolAttribute{Computed: true, MarkdownDescription: "Whether Windows Update for Business is enabled."},
								"endpoint_protection":         schema.BoolAttribute{Computed: true, MarkdownDescription: "Whether endpoint protection is enabled."},
								"office_apps":                 schema.BoolAttribute{Computed: true, MarkdownDescription: "Whether Office apps are enabled."},
							},
						},
						"wi_fi_mac_address": schema.StringAttribute{Computed: true, MarkdownDescription: "Wi-Fi MAC address of the device."},
						"device_health_attestation_state": schema.SingleNestedAttribute{
							Computed:            true,
							MarkdownDescription: "Device health attestation state.",
							Attributes: map[string]schema.Attribute{
								"last_update_date_time":                         schema.StringAttribute{Computed: true, MarkdownDescription: "Last update date and time for device health attestation."},
								"content_namespace_url":                         schema.StringAttribute{Computed: true, MarkdownDescription: "Content namespace URL for device health attestation."},
								"device_health_attestation_status":              schema.StringAttribute{Computed: true, MarkdownDescription: "Device health attestation status."},
								"content_version":                               schema.StringAttribute{Computed: true, MarkdownDescription: "Content version for device health attestation."},
								"issued_date_time":                              schema.StringAttribute{Computed: true, MarkdownDescription: "Issued date and time for device health attestation."},
								"attestation_identity_key":                      schema.StringAttribute{Computed: true, MarkdownDescription: "Attestation identity key."},
								"reset_count":                                   schema.Int64Attribute{Computed: true, MarkdownDescription: "Reset count for device health attestation."},
								"restart_count":                                 schema.Int64Attribute{Computed: true, MarkdownDescription: "Restart count for device health attestation."},
								"data_excution_policy":                          schema.StringAttribute{Computed: true, MarkdownDescription: "Data execution policy for device health attestation."},
								"bit_locker_status":                             schema.StringAttribute{Computed: true, MarkdownDescription: "BitLocker status for device health attestation."},
								"boot_manager_version":                          schema.StringAttribute{Computed: true, MarkdownDescription: "Boot manager version for device health attestation."},
								"code_integrity_check_version":                  schema.StringAttribute{Computed: true, MarkdownDescription: "Code integrity check version for device health attestation."},
								"secure_boot":                                   schema.StringAttribute{Computed: true, MarkdownDescription: "Secure boot status for device health attestation."},
								"boot_debugging":                                schema.StringAttribute{Computed: true, MarkdownDescription: "Boot debugging status for device health attestation."},
								"operating_system_kernel_debugging":             schema.StringAttribute{Computed: true, MarkdownDescription: "Operating system kernel debugging status for device health attestation."},
								"code_integrity":                                schema.StringAttribute{Computed: true, MarkdownDescription: "Code integrity status for device health attestation."},
								"test_signing":                                  schema.StringAttribute{Computed: true, MarkdownDescription: "Test signing status for device health attestation."},
								"safe_mode":                                     schema.StringAttribute{Computed: true, MarkdownDescription: "Safe mode status for device health attestation."},
								"windows_pe":                                    schema.StringAttribute{Computed: true, MarkdownDescription: "Windows PE status for device health attestation."},
								"early_launch_anti_malware_driver_protection":   schema.StringAttribute{Computed: true, MarkdownDescription: "Early launch anti-malware driver protection status."},
								"virtual_secure_mode":                           schema.StringAttribute{Computed: true, MarkdownDescription: "Virtual secure mode status for device health attestation."},
								"pcr_hash_algorithm":                            schema.StringAttribute{Computed: true, MarkdownDescription: "PCR hash algorithm for device health attestation."},
								"boot_app_security_version":                     schema.StringAttribute{Computed: true, MarkdownDescription: "Boot app security version for device health attestation."},
								"boot_manager_security_version":                 schema.StringAttribute{Computed: true, MarkdownDescription: "Boot manager security version for device health attestation."},
								"tpm_version":                                   schema.StringAttribute{Computed: true, MarkdownDescription: "TPM version for device health attestation."},
								"pcr0":                                          schema.StringAttribute{Computed: true, MarkdownDescription: "PCR0 value for device health attestation."},
								"secure_boot_configuration_policy_finger_print": schema.StringAttribute{Computed: true, MarkdownDescription: "Secure boot configuration policy fingerprint."},
								"code_integrity_policy":                         schema.StringAttribute{Computed: true, MarkdownDescription: "Code integrity policy for device health attestation."},
								"boot_revision_list_info":                       schema.StringAttribute{Computed: true, MarkdownDescription: "Boot revision list info for device health attestation."},
								"operating_system_rev_list_info":                schema.StringAttribute{Computed: true, MarkdownDescription: "Operating system revision list info for device health attestation."},
								"health_status_mismatch_info":                   schema.StringAttribute{Computed: true, MarkdownDescription: "Health status mismatch info for device health attestation."},
								"health_attestation_supported_status":           schema.StringAttribute{Computed: true, MarkdownDescription: "Health attestation supported status."},
								"memory_integrity_protection":                   schema.StringAttribute{Computed: true, MarkdownDescription: "Memory integrity protection status."},
								"memory_access_protection":                      schema.StringAttribute{Computed: true, MarkdownDescription: "Memory access protection status."},
								"virtualization_based_security":                 schema.StringAttribute{Computed: true, MarkdownDescription: "Virtualization based security status."},
								"firmware_protection":                           schema.StringAttribute{Computed: true, MarkdownDescription: "Firmware protection status."},
								"system_management_mode":                        schema.StringAttribute{Computed: true, MarkdownDescription: "System management mode status."},
								"secured_core_pc":                               schema.StringAttribute{Computed: true, MarkdownDescription: "Secured core PC status."},
							},
						},
						"subscriber_carrier":            schema.StringAttribute{Computed: true, MarkdownDescription: "Mobile carrier for the device's SIM card."},
						"meid":                          schema.StringAttribute{Computed: true, MarkdownDescription: "Mobile Equipment Identifier (MEID) of the device."},
						"total_storage_space_in_bytes":  schema.Int64Attribute{Computed: true, MarkdownDescription: "Total storage space in bytes."},
						"free_storage_space_in_bytes":   schema.Int64Attribute{Computed: true, MarkdownDescription: "Free storage space in bytes."},
						"managed_device_name":           schema.StringAttribute{Computed: true, MarkdownDescription: "Managed device name."},
						"partner_reported_threat_state": schema.StringAttribute{Computed: true, MarkdownDescription: "Partner reported threat state."},
						"retire_after_date_time":        schema.StringAttribute{Computed: true, MarkdownDescription: "Date and time after which the device will be retired."},
						"users_logged_on": schema.ListNestedAttribute{
							Computed:            true,
							MarkdownDescription: "List of users currently logged on to the device.",
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"user_id":               schema.StringAttribute{Computed: true, MarkdownDescription: "User ID of the logged on user."},
									"last_log_on_date_time": schema.StringAttribute{Computed: true, MarkdownDescription: "Last logon date and time for the user."},
								},
							},
						},
						"prefer_mdm_over_group_policy_applied_date_time": schema.StringAttribute{Computed: true, MarkdownDescription: "Date and time when MDM was preferred over group policy."},
						"autopilot_enrolled":                             schema.BoolAttribute{Computed: true, MarkdownDescription: "Whether the device is enrolled in Autopilot."},
						"require_user_enrollment_approval":               schema.BoolAttribute{Computed: true, MarkdownDescription: "Whether user enrollment approval is required."},
						"management_certificate_expiration_date":         schema.StringAttribute{Computed: true, MarkdownDescription: "Expiration date of the management certificate."},
						"iccid":                                          schema.StringAttribute{Computed: true, MarkdownDescription: "Integrated Circuit Card Identifier (ICCID) for the device's SIM card."},
						"udid":                                           schema.StringAttribute{Computed: true, MarkdownDescription: "Unique Device Identifier (UDID) for the device."},
						"role_scope_tag_ids":                             schema.ListAttribute{ElementType: types.StringType, Computed: true, MarkdownDescription: "List of role scope tag IDs assigned to the device."},
						"windows_active_malware_count":                   schema.Int64Attribute{Computed: true, MarkdownDescription: "Count of active malware instances on the device."},
						"windows_remediated_malware_count":               schema.Int64Attribute{Computed: true, MarkdownDescription: "Count of remediated malware instances on the device."},
						"notes":                                          schema.StringAttribute{Computed: true, MarkdownDescription: "Notes associated with the device."},
						"configuration_manager_client_health_state": schema.SingleNestedAttribute{
							Computed:            true,
							MarkdownDescription: "Configuration Manager client health state.",
							Attributes: map[string]schema.Attribute{
								"state":               schema.StringAttribute{Computed: true, MarkdownDescription: "Health state of the Configuration Manager client."},
								"error_code":          schema.Int64Attribute{Computed: true, MarkdownDescription: "Error code for the Configuration Manager client health state."},
								"last_sync_date_time": schema.StringAttribute{Computed: true, MarkdownDescription: "Last sync date and time for the Configuration Manager client."},
							},
						},
						"configuration_manager_client_information": schema.SingleNestedAttribute{
							Computed:            true,
							MarkdownDescription: "Configuration Manager client information.",
							Attributes: map[string]schema.Attribute{
								"client_identifier": schema.StringAttribute{Computed: true, MarkdownDescription: "Client identifier for the Configuration Manager client."},
								"is_blocked":        schema.BoolAttribute{Computed: true, MarkdownDescription: "Whether the Configuration Manager client is blocked."},
								"client_version":    schema.StringAttribute{Computed: true, MarkdownDescription: "Version of the Configuration Manager client."},
							},
						},
						"ethernet_mac_address":     schema.StringAttribute{Computed: true, MarkdownDescription: "Ethernet MAC address of the device."},
						"physical_memory_in_bytes": schema.Int64Attribute{Computed: true, MarkdownDescription: "Physical memory in bytes on the device."},
						"processor_architecture":   schema.StringAttribute{Computed: true, MarkdownDescription: "Processor architecture of the device (e.g., x86, x64)."},
						"specification_version":    schema.StringAttribute{Computed: true, MarkdownDescription: "Specification version of the device."},
						"join_type":                schema.StringAttribute{Computed: true, MarkdownDescription: "Join type of the device (e.g., azureADJoined)."},
						"sku_family":               schema.StringAttribute{Computed: true, MarkdownDescription: "SKU family of the device."},
						"security_patch_level":     schema.StringAttribute{Computed: true, MarkdownDescription: "Security patch level of the device."},
						"sku_number":               schema.Int32Attribute{Computed: true, MarkdownDescription: "SKU number of the device."},
						"management_features":      schema.StringAttribute{Computed: true, MarkdownDescription: "Management features enabled on the device."},
						"chrome_os_device_info": schema.ListNestedAttribute{
							Computed:            true,
							MarkdownDescription: "List of Chrome OS device information properties.",
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"name":       schema.StringAttribute{Computed: true, MarkdownDescription: "Name of the Chrome OS device property."},
									"value":      schema.StringAttribute{Computed: true, MarkdownDescription: "Value of the Chrome OS device property."},
									"value_type": schema.StringAttribute{Computed: true, MarkdownDescription: "Type of the value for the Chrome OS device property."},
									"updatable":  schema.BoolAttribute{Computed: true, MarkdownDescription: "Whether the Chrome OS device property is updatable."},
								},
							},
						},
						"enrollment_profile_name":                         schema.StringAttribute{Computed: true, MarkdownDescription: "Enrollment profile name for the device."},
						"bootstrap_token_escrowed":                        schema.BoolAttribute{Computed: true, MarkdownDescription: "Whether the bootstrap token is escrowed for the device."},
						"device_firmware_configuration_interface_managed": schema.BoolAttribute{Computed: true, MarkdownDescription: "Whether the device firmware configuration interface is managed."},
						"device_identity_attestation_detail": schema.SingleNestedAttribute{
							Computed:            true,
							MarkdownDescription: "Device identity attestation detail.",
							Attributes: map[string]schema.Attribute{
								"device_identity_attestation_status": schema.StringAttribute{Computed: true, MarkdownDescription: "Device identity attestation status."},
							},
						},
					},
				},
			},
			"timeouts": commonschema.ResourceTimeouts(ctx),
		},
	}
}
