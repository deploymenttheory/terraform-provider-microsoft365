resource "microsoft365_graph_beta_device_management_mobile_threat_defense_connector" "test" {
  id = "8d0ed095-8191-4bd3-8a41-953b22d51ff7"

  android_enabled                                            = true
  ios_enabled                                                = true
  mac_enabled                                                = true
  windows_enabled                                            = true
  android_mobile_application_management_enabled              = true
  ios_mobile_application_management_enabled                  = true
  windows_mobile_application_management_enabled              = true
  partner_unsupported_os_version_blocked                     = true
  partner_unresponsiveness_threshold_in_days                 = 14
  allow_partner_to_collect_ios_application_metadata          = true
  allow_partner_to_collect_ios_personal_application_metadata = true
  allow_partner_to_collect_ios_certificate_metadata          = true
  allow_partner_to_collect_ios_personal_certificate_metadata = true
  grant_mobile_threat_defense_partner_role                   = true
}
