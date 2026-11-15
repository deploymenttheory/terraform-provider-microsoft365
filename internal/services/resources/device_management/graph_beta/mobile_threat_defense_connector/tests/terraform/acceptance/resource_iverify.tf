# iVerify partner connector, created directly with a mix of settings.
# The acceptance tests only use partners that have no connector in the test tenant.
# Android and iOS MAM stay off because only one partner per platform may be enabled
# for MAM evaluation and the tenant's Defender connector already uses them.
resource "microsoft365_graph_beta_device_management_mobile_threat_defense_connector" "iverify" {
  id = "fc06ff9b-e2f9-4221-9a68-465684ee8662"

  android_enabled                                            = true
  ios_enabled                                                = true
  mac_enabled                                                = false
  windows_enabled                                            = true
  android_mobile_application_management_enabled              = false
  ios_mobile_application_management_enabled                  = false
  windows_mobile_application_management_enabled              = false
  partner_unsupported_os_version_blocked                     = true
  partner_unresponsiveness_threshold_in_days                 = 30
  allow_partner_to_collect_ios_application_metadata          = true
  allow_partner_to_collect_ios_personal_application_metadata = false
  allow_partner_to_collect_ios_certificate_metadata          = true
  allow_partner_to_collect_ios_personal_certificate_metadata = false
  grant_mobile_threat_defense_partner_role                   = false
}
