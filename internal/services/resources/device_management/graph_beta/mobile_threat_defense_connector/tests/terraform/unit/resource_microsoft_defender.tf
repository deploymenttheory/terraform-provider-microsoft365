resource "microsoft365_graph_beta_device_management_mobile_threat_defense_connector" "microsoft_defender" {
  id = "fc780465-2017-40d4-a0c5-307022471b92"

  android_enabled                                   = true
  ios_enabled                                       = true
  windows_enabled                                   = true
  partner_unsupported_os_version_blocked            = true
  partner_unresponsiveness_threshold_in_days        = 5
  allow_partner_to_collect_ios_application_metadata = true
  microsoft_defender_for_endpoint_attach_enabled    = true
}
