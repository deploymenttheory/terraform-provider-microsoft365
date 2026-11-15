// REF: https://learn.microsoft.com/en-us/graph/api/resources/intune-onboarding-mobilethreatdefenseconnector?view=graph-rest-beta
package graphBetaMobileThreatDefenseConnector

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// MobileThreatDefenseConnectorResourceModel represents the schema for the MobileThreatDefenseConnector resource
type MobileThreatDefenseConnectorResourceModel struct {
	ID                                                  types.String   `tfsdk:"id"`
	AndroidEnabled                                      types.Bool     `tfsdk:"android_enabled"`
	IosEnabled                                          types.Bool     `tfsdk:"ios_enabled"`
	MacEnabled                                          types.Bool     `tfsdk:"mac_enabled"`
	WindowsEnabled                                      types.Bool     `tfsdk:"windows_enabled"`
	AndroidMobileApplicationManagementEnabled           types.Bool     `tfsdk:"android_mobile_application_management_enabled"`
	IosMobileApplicationManagementEnabled               types.Bool     `tfsdk:"ios_mobile_application_management_enabled"`
	WindowsMobileApplicationManagementEnabled           types.Bool     `tfsdk:"windows_mobile_application_management_enabled"`
	PartnerUnsupportedOsVersionBlocked                  types.Bool     `tfsdk:"partner_unsupported_os_version_blocked"`
	PartnerUnresponsivenessThresholdInDays              types.Int32    `tfsdk:"partner_unresponsiveness_threshold_in_days"`
	AllowPartnerToCollectIosApplicationMetadata         types.Bool     `tfsdk:"allow_partner_to_collect_ios_application_metadata"`
	AllowPartnerToCollectIosPersonalApplicationMetadata types.Bool     `tfsdk:"allow_partner_to_collect_ios_personal_application_metadata"`
	AllowPartnerToCollectIosCertificateMetadata         types.Bool     `tfsdk:"allow_partner_to_collect_ios_certificate_metadata"`
	AllowPartnerToCollectIosPersonalCertificateMetadata types.Bool     `tfsdk:"allow_partner_to_collect_ios_personal_certificate_metadata"`
	MicrosoftDefenderForEndpointAttachEnabled           types.Bool     `tfsdk:"microsoft_defender_for_endpoint_attach_enabled"`
	GrantMobileThreatDefensePartnerRole                 types.Bool     `tfsdk:"grant_mobile_threat_defense_partner_role"`
	LaunchMobileThreatDefensePartnerOnSetupEnabled      types.Bool     `tfsdk:"launch_mobile_threat_defense_partner_on_setup_enabled"`
	AndroidDeviceBlockedOnMissingPartnerData            types.Bool     `tfsdk:"android_device_blocked_on_missing_partner_data"`
	IosDeviceBlockedOnMissingPartnerData                types.Bool     `tfsdk:"ios_device_blocked_on_missing_partner_data"`
	MacDeviceBlockedOnMissingPartnerData                types.Bool     `tfsdk:"mac_device_blocked_on_missing_partner_data"`
	WindowsDeviceBlockedOnMissingPartnerData            types.Bool     `tfsdk:"windows_device_blocked_on_missing_partner_data"`
	PartnerState                                        types.String   `tfsdk:"partner_state"`
	LastHeartbeatDateTime                               types.String   `tfsdk:"last_heartbeat_date_time"`
	PartnerDisplayName                                  types.String   `tfsdk:"partner_display_name"`
	PartnerSupportsAndroidMam                           types.Bool     `tfsdk:"partner_supports_android_mam"`
	PartnerSupportsIosMam                               types.Bool     `tfsdk:"partner_supports_ios_mam"`
	PartnerSupportsWindowsMam                           types.Bool     `tfsdk:"partner_supports_windows_mam"`
	PartnerSupportsIosCertificateSync                   types.Bool     `tfsdk:"partner_supports_ios_certificate_sync"`
	PartnerMinimumSupportedAndroidVersion               types.String   `tfsdk:"partner_minimum_supported_android_version"`
	PartnerMinimumSupportedIosVersion                   types.String   `tfsdk:"partner_minimum_supported_ios_version"`
	PartnerMinimumSupportedMacVersion                   types.String   `tfsdk:"partner_minimum_supported_mac_version"`
	PartnerMinimumSupportedWindowsVersion               types.String   `tfsdk:"partner_minimum_supported_windows_version"`
	Timeouts                                            timeouts.Value `tfsdk:"timeouts"`
}
